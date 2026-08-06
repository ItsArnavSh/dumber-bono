package radio

import (
	"dubmer-bono/app/types/entity"
	"dubmer-bono/app/utility"
	"fmt"
	"io"
	"strings"
)

const maxPriority = 5

// GetMessageByMinPriority pops the highest-priority allowed message and
// returns it as an io.Reader, regardless of the underlying payload type:
//   - DirectMessage / FunctionMessage: wrapped in a strings.Reader over the
//     resolved text.
//   - IOPipe: its *io.PipeReader is returned directly, so callers can read
//     it incrementally as it streams in (e.g. from the LLM).
func (s *Service) GetMessageByMinPriority() (io.Reader, bool) {
	priority := maxPriority
	allowed_priority := s.getDriverPressure()
	for priority >= allowed_priority {
		s.maplock.Lock()
		vc, ok := s.prio_sorted_vc[priority]
		s.maplock.Unlock()
		if ok {
			rad_msg, ok := vc.Pop()
			// If muted only surface top-level messages, SC events etc.
			if s.muted && rad_msg.Priority != maxPriority {
				continue
			}
			if ok {
				switch p := rad_msg.Message.(type) {
				case entity.DirectMessage:
					return strings.NewReader(p.Text), true
				case entity.FunctionMessage:
					return strings.NewReader(p.Fn()), true
				case entity.IOPipe:
					return p.Pipe, true
				}
			}
		}
		priority--
	}
	return nil, false
}

func (s *Service) MessageChanListner() {
	for msg := range s.msg_chan {
		fmt.Println("received message")
		queue, ok := s.prio_sorted_vc[msg.Priority]
		if !ok {
			q := utility.NewQueue[entity.RadioMessage]()
			queue = &q
			s.prio_sorted_vc[msg.Priority] = queue
		}
		queue.Push(msg)
	}
}
