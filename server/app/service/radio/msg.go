package radio

// 5 means MAX PRIORITY and 1 is Min priority
// If We are on a 4 Turn, the radio will only get 4 priority basically
func (s *Service) GetMessageByMaxPriority(priority int) (string, bool) {
	for priority > 0 {
		vc := s.prio_sorted_vc[priority]
		rad_msg, ok := vc.Pop()
		if ok {
			return rad_msg.Message, true
		}
		priority--
	}
	return "", false
}
