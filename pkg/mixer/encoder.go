package mixer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/go-audio/wav"

	"github.com/bwmarrin/discordgo"
	"github.com/hraban/opus"
)

// InputStream is the type used to control an input file stream
type InputStream struct {
	id uint64

	bufLock sync.Mutex
	buf     []int16
}

// Encoder is a discordgo voice encoder that can mix streams
type Encoder struct {
	queueLock    sync.Mutex
	inputStreams map[uint64]*InputStream
	stop         chan bool

	encoder *opus.Encoder

	VC *discordgo.VoiceConnection

	pcmbuf []int16
}

// NewInputStream gives a new InputStream for a given ID
func NewInputStream(id uint64) *InputStream {
	return &InputStream{
		id: id,
	}
}

// HandleFile handles an a WAV file decoding
func (is *InputStream) HandleFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		log(err)
		return
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log(err)
		return
	}

	d := wav.NewDecoder(bytes.NewReader(data))
	if !d.IsValidFile() {
		log("Not valid WAV")
		return
	}

	a, err := d.FullPCMBuffer()
	if err != nil {
		log(err)
	}
	is.bufLock.Lock()
	for _, pcm := range a.Data {
		is.buf = append(is.buf, int16(pcm))
	}
	is.bufLock.Unlock()
}

// Read gives back a buffer for the audio file, implements io.Reader
func (is *InputStream) Read(b []int16) (n int, err error) {
	is.bufLock.Lock()
	defer is.bufLock.Unlock()
	if len(is.buf) < 1 {
		return 0, nil
	}

	n = copy(b, is.buf)
	is.buf = is.buf[n:]
	return
}

// NewEncoder gives a new Encoder
func NewEncoder() *Encoder {
	enc, err := opus.NewEncoder(48000, 1, opus.AppAudio)
	if err != nil {
		panic("Failed creating encoder: " + err.Error())
	}
	return &Encoder{
		stop:         make(chan bool),
		inputStreams: make(map[uint64]*InputStream),
		encoder:      enc,
	}
}

// Stop stops the encoder
func (e *Encoder) Stop() {
	close(e.stop)
}

// Queue queues a WAV file input for a stream ID, different IDs live mix the audio, same IDs play on serial
func (e *Encoder) Queue(id uint64, path string) {

	st, ok := e.inputStreams[id]
	if !ok {
		st = NewInputStream(id)
		e.queueLock.Lock()
		e.inputStreams[id] = st
		e.queueLock.Unlock()
	}

	st.HandleFile(path)
}

// Run runs the encoder
func (e *Encoder) Run() {
	log("Encoder running")
	for {
		select {
		case <-e.stop:
			log("Encoder stopping")
			return
		default:
			e.processQueue()
		}
	}
}

func (e *Encoder) processQueue() {
	e.queueLock.Lock()
	mixedPCM := make([]int16, 48*20*1)
	for _, st := range e.inputStreams {
		userPCM := make([]int16, 48*20*1)
		n, _ := st.Read(userPCM)
		if n < 1 {
			continue
		}

		for i := 0; i < len(userPCM); i++ {

			// Mix it
			v := int32(mixedPCM[i] + userPCM[i])
			// Clip
			if v > 0x7fff {
				v = 0x7fff
			} else if v < -0x7fff {
				v = -0x7fff
			}
			mixedPCM[i] = int16(v)
		}
	}
	e.queueLock.Unlock()

	output := make([]byte, 0xfff)
	n, err := e.encoder.Encode(mixedPCM, output)
	if err != nil {
		log("Failed encode: ", err)
	}

	e.VC.OpusSend <- output[:n]
}

// HasFinishedAll returns true if all streams have no more audio to play
func (e *Encoder) HasFinishedAll() bool {
	for _, st := range e.inputStreams {
		if len(st.buf) > 0 {
			return false
		}
	}

	return true
}

func log(s ...interface{}) {
	fmt.Fprintln(os.Stderr, s...)
}
