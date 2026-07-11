package radio

func (s *Service) GetMessageByMinPriority(allowed_priority int) (string, bool) {
	priority := 5
	for priority >= allowed_priority {
		vc := s.prio_sorted_vc[priority]
		rad_msg, ok := vc.Pop()
		if ok {
			return rad_msg.Message, true
		}
		priority--
	}
	return "", false
}

func (s *Service) AddMessage() {
	for msg := range s.msg_chan {
		queue := s.prio_sorted_vc[msg.Priority]
		queue.Push(msg)
	}
}
