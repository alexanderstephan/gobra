package sound

import (
	"flag"
	"io"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/oto/v2"
)

var (
	context         *oto.Context
	sampleRate      = flag.Int("samplerate", 44100, "sample rate")
	channelNum      = flag.Int("channelnum", 3, "number of channels")
	bitDepthInBytes = flag.Int("bitdepthinbytes", 2, "bit depth in bytes")
)

const FreqA = 300 // This is the frequency of the generated sound.

type sineWave struct {
	freq   float64
	length int64
	pos    int64

	remaining []byte
}

func newSineWave(freq float64, duration time.Duration) *sineWave {
	l := int64(*channelNum) * int64(*bitDepthInBytes) * int64(*sampleRate) * int64(duration) / int64(time.Second)
	l = l / 4 * 4
	return &sineWave{
		freq:   freq,
		length: l,
	}
}

// Read reads data from the sineWave and fills the provided buffer.
func (s *sineWave) Read(buf []byte) (int, error) {
	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		s.remaining = s.remaining[n:]
		return n, nil
	}

	if s.pos == s.length {
		return 0, io.EOF
	}

	eof := false
	if s.pos+int64(len(buf)) > s.length {
		buf = buf[:s.length-s.pos]
		eof = true
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	length := float64(*sampleRate) / float64(s.freq)

	num := (*bitDepthInBytes) * (*channelNum)
	p := s.pos / int64(num)
	switch *bitDepthInBytes {
	case 1:
		for i := 0; i < len(buf)/num; i++ {
			const max = 128
			b := int(math.Sin(2*math.Pi*float64(p)/length) * 0.3 * max)
			for ch := 0; ch < *channelNum; ch++ {
				buf[num*i+ch] = byte(b + 128)
			}
			p++
		}
	case 2:
		for i := 0; i < len(buf)/num; i++ {
			const max = 32767
			b := int16(math.Sin(2*math.Pi*float64(p)/length) * 0.3 * max)
			for ch := 0; ch < *channelNum; ch++ {
				buf[num*i+2*ch] = byte(b)
				buf[num*i+1+2*ch] = byte(b >> 8)
			}
			p++
		}
	}

	s.pos += int64(len(buf))

	n := len(buf)
	if origBuf != nil {
		n = copy(origBuf, buf)
		s.remaining = buf[n:]
	}

	if eof {
		return n, io.EOF
	}
	return n, nil
}

// Play plays a sound at a given frequency for a duration in milliseconds.
func Play(freq float64, duration time.Duration) oto.Player {
	s := newSineWave(freq, duration)
	p := context.NewPlayer(s)
	p.Play()
	return p
}

// InitSound initializes oto so sound can be played.
func InitSound() {
	c, ready, err := oto.NewContext(*sampleRate, *channelNum, *bitDepthInBytes)
	if err != nil {
		log.Fatal(err)
	}
	<-ready
	context = c
}

// PlayFoodSound plays a high pitched sound that should be played when the snake eats food.
func PlayFoodSound() {
	Play(FreqA, 250*time.Millisecond)
}
