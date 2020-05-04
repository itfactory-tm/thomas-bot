package discordha

import (
	"context"
	"fmt"
	"math/rand"

	"go.etcd.io/etcd/clientv3"
)

// LockVoice locks a voice channel ID, returns true if successful
func (h *HA) LockVoice(channelID string) (bool, error) {
	return h.lockKey(fmt.Sprintf("voice-%s", channelID), false)
}

// UnlockVoice unlocks a voice channel ID
func (h *HA) UnlockVoice(channelID string) error {
	return h.unlockKey(fmt.Sprintf("voice-%s", channelID))
}

// SendVoiceCommand sends a string command to the instance handling the voice channel
// These can be received using WatchVoiceCommands
func (h *HA) SendVoiceCommand(channelID, command string) error {
	grant, err := h.etcd.Grant(context.TODO(), int64(30))
	if err != nil {
		return err
	}
	_, err = h.etcd.Put(context.TODO(), fmt.Sprintf("/voice/command/%s/%d", channelID, rand.Intn(9999999)), command, clientv3.WithLease(grant.ID))
	return err
}

// WatchVoiceCommands gives a channel with commands transmitted by SendVoiceCommand
func (h *HA) WatchVoiceCommands(ctx context.Context, channelID string) chan string {
	out := make(chan string)
	w := h.etcd.Watch(ctx, fmt.Sprintf("/voice/command/%s/", channelID), clientv3.WithPrefix())
	go func() {
		for wresp := range w {
			if wresp.Canceled {
				close(out)
				break
			}
			for _, ev := range wresp.Events {
				if ev.IsCreate() {
					out <- string(ev.Kv.Value)
				}
			}
		}
	}()

	return out
}
