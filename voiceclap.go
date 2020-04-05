package main

import (
	"fmt"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/mixer"

	"github.com/bwmarrin/discordgo"
)

// TODO: automate these
const itfDiscord = "687565213943332875"
const audioChannel = "688370622228725848"

var audioConnected = false
var voiceQueueChan = make(chan string)

func connectVoice(dg *discordgo.Session) {
	audioConnected = true
	dgv, err := dg.ChannelVoiceJoin(itfDiscord, audioChannel, false, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	encoder := mixer.NewEncoder()
	encoder.VC = dgv
	go encoder.Run()

	doneChan := make(chan struct{})
	go func() {
		var i uint64
		for {
			select {
			case f := <-voiceQueueChan:
				go encoder.Queue(uint64(i), f)
				i++
			case <-doneChan:
				return
			}
		}
	}()

	time.Sleep(5 * time.Second) // waiting for first audio
	for !encoder.HasFinishedAll() {
		time.Sleep(5 * time.Second)
	}

	// Close connections once all are played
	doneChan <- struct{}{}
	dgv.Disconnect()
	dgv.Close()
	encoder.Stop()
	audioConnected = false
}
