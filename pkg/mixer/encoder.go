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

var (
	Silence = []byte{0xF8, 0xFF, 0xFE}
)

type InputStream struct {
	id      uint64
	decoder *opus.Decoder

	bufLock sync.Mutex
	buf     []int16
}

type Encoder struct {
	queueLock    sync.Mutex
	inputStreams map[uint64]*InputStream
	stop         chan bool

	encoder *opus.Encoder

	VC *discordgo.VoiceConnection

	pcmbuf []int16
}

func NewInputStream(id uint64) *InputStream {
	dec, err := opus.NewDecoder(48000, 2)
	if err != nil {
		panic("Failed creating decoder: " + err.Error())
	}

	return &InputStream{
		decoder: dec,
		id:      id,
	}
}

// Handles an incoming voice packet
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

func NewEncoder() *Encoder {
	enc, err := opus.NewEncoder(48000, 2, opus.AppAudio)
	if err != nil {
		panic("Failed creating encoder: " + err.Error())
	}
	return &Encoder{
		stop:         make(chan bool),
		inputStreams: make(map[uint64]*InputStream),
		encoder:      enc,
	}
}

func (e *Encoder) Stop() {
	close(e.stop)
}

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

func (e *Encoder) Run() {
	log("Encoder running")
	// ticker := time.NewTicker(time.Millisecond * 120)
	for {
		select {
		case <-e.stop:
			log("Encoder stopping")
			return
		default:
			e.ProcessQueue()
		}
	}
}

func (e *Encoder) ProcessQueue() {
	e.queueLock.Lock()
	mixedPCM := make([]int16, 48*20*2)
	for _, st := range e.inputStreams {
		userPCM := make([]int16, 48*20*2)
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
