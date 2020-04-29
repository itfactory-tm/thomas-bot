package discordha

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

// HA is a helper struct for high available discordgo using etcd
type HA struct {
	config     Config
	etcd       *clientv3.Client
	locksMutex sync.Mutex
	locks      map[string]clientv3.LeaseID
}

// Config contains the configuration for HA
type Config struct {
	Session            *discordgo.Session
	HA                 bool
	LockUpdateInterval time.Duration
	LockTTL            time.Duration
	EtcdEndpoints      []string
}

// New gives a HA instance for a given configuration
func New(c Config) (*HA, error) {
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
		config: c,
		etcd:   client,
		locks:  map[string]clientv3.LeaseID{},
	}

	go s.lockUpdateLoop()

	return s, nil
}

func (h *HA) lockUpdateLoop() {
	for {
		time.Sleep(h.config.LockUpdateInterval)
		h.locksMutex.Lock()
		for _, lease := range h.locks {
			go func(l clientv3.LeaseID) {
				err := h.keepAlive(lease)
				if err != nil {
					log.Println(err)
				}
			}(lease)
		}
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
		return false, err
	}
	key := fmt.Sprintf("/locks/%s", hash)
	return h.lockKey(key)
}

func (h *HA) lockKey(key string) (bool, error) {
	grant, err := h.etcd.Grant(context.TODO(), int64(h.config.LockTTL.Seconds()))
	if err != nil {
		return false, err
	}

	txn, err := h.etcd.Txn(context.TODO()).
		// txn value comparisons are lexical
		If(clientv3.Compare(clientv3.Value(key), ">", statusNone)).
		Else(clientv3.OpPut(key, statusHandling, clientv3.WithLease(grant.ID))).
		Commit()

	if err != nil {
		return false, err
	}

	if txn.Succeeded {
		// Lock exists!
		// TODO: add watcher in case server expires
		ctx, cancel := context.WithCancel(context.TODO())

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
					return h.Lock(key) // re-lock!
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
		return err
	}
	key := fmt.Sprintf("/locks/%s", hash)

	return h.unlockKey(key)
}

func (h *HA) unlockKey(key string) error {
	_, err := h.etcd.Put(context.TODO(), key, statusOk, clientv3.WithLease(h.locks[key]))
	if err != nil {
		return err
	}
	time.Sleep(300 * time.Millisecond)

	h.locksMutex.Lock()
	delete(h.locks, key)
	h.locksMutex.Unlock()

	_, err = h.etcd.Delete(context.TODO(), key)
	return err
}

func (h *HA) keepAlive(leaseID clientv3.LeaseID) error {
	h.etcd.KeepAlive(context.TODO(), leaseID)
	return nil
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

var ErrorCacheKeyNotExist = errors.New("Cache key does not exist")

func (h *HA) CacheRead(cache, key string, want interface{}) (interface{}, error) {
	resp, err := h.etcd.Get(context.TODO(), fmt.Sprintf("/cache/%s/%s", cache, key))
	if err != nil {
		return nil, err
	}

	if resp.Count < 1 {
		return nil, ErrorCacheKeyNotExist
	}

	err = json.Unmarshal(resp.Kvs[0].Value, &want)
	if err != nil {
		return nil, err
	}

	return want, nil
}

func (h *HA) CacheWrite(cache, key string, data interface{}, ttl time.Duration) error {
	grant, err := h.etcd.Grant(context.TODO(), int64(ttl.Seconds()))
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = h.etcd.Put(context.TODO(), fmt.Sprintf("/cache/%s/%s", cache, key), string(jsonData), clientv3.WithLease(grant.ID))
	return err
}
