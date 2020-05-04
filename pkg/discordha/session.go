package discordha

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.etcd.io/etcd/clientv3"
)

// etcd states to store in value
const (
	statusNone     = "0"
	statusHandling = "1"
	statusOk       = "2"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// HA is a helper struct for high available discordgo using etcd
type HA struct {
	config     Config
	etcd       *clientv3.Client
	locksMutex sync.Mutex
	locks      map[string]clientv3.LeaseID
	bgContext  context.Context
}

// Config contains the configuration for HA
type Config struct {
	Session            *discordgo.Session
	HA                 bool
	LockUpdateInterval time.Duration
	LockTTL            time.Duration
	EtcdEndpoints      []string
	Context            context.Context
}

// New gives a HA instance for a given configuration
func New(c Config) (*HA, error) {
	if !c.HA {
		return &HA{
			config: c,
		}, nil
	}

	if c.LockUpdateInterval == 0 {
		c.LockUpdateInterval = time.Second * 3
	}
	if c.LockTTL == 0 {
		c.LockTTL = time.Second * 10
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   c.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	var s = &HA{
		config:    c,
		etcd:      client,
		locks:     map[string]clientv3.LeaseID{},
		bgContext: c.Context,
	}

	go s.lockUpdateLoop()
	go s.logLoop()

	return s, nil
}

func (h *HA) lockUpdateLoop() {
	for {
		time.Sleep(h.config.LockUpdateInterval)
		h.locksMutex.Lock()
		for _, lease := range h.locks {
			err := h.keepAlive(lease)
			if err != nil {
				log.Printf("Etcd keepalive error: %q\n", err)
			}
		}
		h.locksMutex.Unlock()
	}
}

func (h *HA) logLoop() {
	for {
		time.Sleep(time.Minute)
		h.locksMutex.Lock()
		log.Printf("I own %d locks\n", len(h.locks))
		h.locksMutex.Unlock()
	}
}

// Lock tries to acquire a lock on an event, it will return true if
// the instance that requests it may process the request.
func (h *HA) Lock(obj interface{}) (bool, error) {
	if !h.config.HA {
		// Non HA, development instance probably
		return true, nil
	}

	hash, err := h.getObjectHash(obj)
	if err != nil {
		log.Printf("Hash error:%q\n", err)
		return false, err
	}
	key := fmt.Sprintf("/locks/%s", hash)
	return h.lockKey(key, true)
}

func (h *HA) lockKey(key string, waitForFailure bool) (bool, error) {
	grant, err := h.etcd.Grant(h.bgContext, int64(h.config.LockTTL.Seconds()))
	if err != nil {
		return false, err
	}

	txn, err := h.etcd.Txn(h.bgContext).
		// txn value comparisons are lexical
		If(clientv3.Compare(clientv3.Value(key), ">", statusNone)).
		Else(clientv3.OpPut(key, statusHandling, clientv3.WithLease(grant.ID))).
		Commit()

	if err != nil {
		return false, err
	}

	if txn.Succeeded {
		// Lock exists!
		if !waitForFailure {
			return false, nil
		}
		ctx, cancel := context.WithCancel(h.bgContext)

		w := h.etcd.Watch(ctx, key)
		for wresp := range w {
			if wresp.Canceled {
				break
			}
			for _, ev := range wresp.Events {
				if ev.IsModify() && string(ev.Kv.Value) == statusOk {
					// other server succeeded!
					return false, nil
				}
				if ev.Type == clientv3.EventTypeDelete {
					return h.lockKey(key, waitForFailure) // re-lock!
				}
			}
		}
		cancel()
		return false, nil
	}

	h.locksMutex.Lock()
	h.locks[key] = grant.ID
	h.locksMutex.Unlock()

	return true, nil
}

// Unlock will release a lock on an event
func (h *HA) Unlock(obj interface{}) error {
	if !h.config.HA {
		// Non HA, development instance probably
		return nil
	}

	hash, err := h.getObjectHash(obj)
	if err != nil {
		log.Printf("Hash error:%q\n", err)
		return err
	}
	key := fmt.Sprintf("/locks/%s", hash)

	return h.unlockKey(key)
}

func (h *HA) unlockKey(key string) error {
	_, err := h.etcd.Put(h.bgContext, key, statusOk, clientv3.WithLease(h.locks[key]))
	if err != nil {
		log.Printf("Failed to set status OK: %q retrying", err)
		time.Sleep(5 * time.Second)
		return h.unlockKey(key)
	}
	time.Sleep(300 * time.Millisecond)

	_, err = h.etcd.Delete(h.bgContext, key)
	if err != nil {
		log.Printf("Failed to delete key: %q retrying", err)
		time.Sleep(5 * time.Second)
		return h.unlockKey(key)
	}

	h.locksMutex.Lock()
	delete(h.locks, key)
	h.locksMutex.Unlock()

	return nil
}

func (h *HA) keepAlive(leaseID clientv3.LeaseID) error {
	_, err := h.etcd.KeepAlive(h.bgContext, leaseID)
	return err
}

func (h *HA) getObjectHash(v interface{}) (string, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	hasher := sha256.New()
	hasher.Write(jsonData)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil)), nil
}
