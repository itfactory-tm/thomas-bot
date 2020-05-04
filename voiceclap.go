package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/mixer"

	"github.com/bwmarrin/discordgo"
)

// TODO: automate these
const itfDiscord = "687565213943332875"
const audioChannel = "688370622228725848"

var audioConnected = false

func connectVoice(dg *discordgo.Session, connected chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if audioConnected {
		connected <- struct{}{}
		return
	}
	gotLock, err := ha.LockVoice(audioChannel)
	if err != nil {
		connected <- struct{}{}
		log.Println(err)
		return
	}
	if !gotLock {
		connected <- struct{}{}
		return
	}

	audioConnected = true
	voiceQueueChan := ha.WatchVoiceCommands(ctx, audioChannel)

	dgv, err := dg.ChannelVoiceJoin(itfDiscord, audioChannel, false, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	connected <- struct{}{}

	encoder := mixer.NewEncoder()
	encoder.VC = dgv
	go encoder.Run()

	doneChan := make(chan struct{})
	go func() {
		var i uint64
		for {
			select {
			case f := <-voiceQueueChan:
				log.Println(f)
				go encoder.Queue(uint64(i), f)
				i++
			case <-doneChan:
				ha.UnlockVoice(audioChannel)
				return
			}
		}
	}()

	time.Sleep(5 * time.Second) // waiting for first audio
	for !encoder.HasFinishedAll() {
		time.Sleep(5 * time.Second)
	}

	// Close connections once all are played

	dgv.Disconnect()
	dgv.Close()
	encoder.Stop()
	audioConnected = false
	ha.UnlockVoice(audioChannel)
	doneChan <- struct{}{}
}
