package radio

import (
	"dubmer-bono/app/types/entity"
	"dubmer-bono/app/utility"
	"fmt"
)

const maxPriority = 5

func (s *Service) GetMessageByMinPriority() (string, bool) {
	priority := maxPriority
	allowed_priority := s.getDriverPressure()
	for priority >= allowed_priority {
		s.maplock.Lock()
		vc, ok := s.prio_sorted_vc[priority]
		s.maplock.Unlock()
		if ok {
			rad_msg, ok := vc.Pop()
			//If muted only print top level messages, SC other events etc
			if s.muted && rad_msg.Priority != maxPriority {
				continue
			}
			if ok {
				switch p := rad_msg.Message.(type) {
				case entity.DirectMessage:
					return p.Text, true
				case entity.FunctionMessage:
					return p.Fn(), true
				}
			}
		}
		priority--
	}
	return "", false
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
