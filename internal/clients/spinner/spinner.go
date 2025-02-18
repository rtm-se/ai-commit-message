package spinner

import (
	"fmt"
	"sync"
	"time"
)

type Spinner struct {
	stop chan bool
	wg   *sync.WaitGroup
}

const spinAnimation = "▁▂▃▄▅▆▇█▇▆▅▄▃▁"

func NewSpinner() *Spinner {
	return &Spinner{
		stop: make(chan bool),
		wg:   &sync.WaitGroup{},
	}
}
func (s *Spinner) hideCursor() {
	fmt.Print("\033[?25l")
}
func (s *Spinner) showCursor() {
	fmt.Print("\033[?25h")
}

func (s *Spinner) clearLine() {
	fmt.Print("\010")
}
func (s *Spinner) putChar(c string) {
	fmt.Print(c, "\010")
}

func (s *Spinner) Stop() {
	defer func() {
		close(s.stop)
		s.stop = nil
	}()
	s.stop <- true
	s.wg.Wait()
	s.clearLine()
	s.showCursor()
}

// should be run in Gorutine
func (s *Spinner) Spin() {
	if s.stop == nil {
		s.stop = make(chan bool)
	}
	defer s.wg.Done()
	s.wg.Add(1)
	s.hideCursor()
	for {
		select {
		case <-s.stop:
			return
		default:
			for _, character := range spinAnimation {
				s.putChar(string(character))
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

}
