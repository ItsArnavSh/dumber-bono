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
		vc, ok := s.prio_sorted_vc[priority]
		if ok {
			rad_msg, ok := vc.Pop()
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
