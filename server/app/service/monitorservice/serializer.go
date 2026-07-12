package monitor

import (
	"dubmer-bono/app/types/entity"
	"time"
)

func (s *Service) SendInformOfPosition() {
	s.PushToRadio(entity.FunctionMessage{Fn: s.InformOfPosition}, 1, time.Minute)
}

func (s *Service) SendInformOfGapToFront() {
	s.PushToRadio(entity.FunctionMessage{Fn: s.InformOfGapToFront}, 1, time.Minute)
}

func (s *Service) SendInformOfGapToLeader() {
	s.PushToRadio(entity.FunctionMessage{Fn: s.InformOfGapToLeader}, 1, time.Minute)
}

func (s *Service) SendInformOfCurrentLap() {
	s.PushToRadio(entity.FunctionMessage{Fn: s.InformOfCurrentLap}, 1, time.Minute)
}

func (s *Service) SendInformOfLastLapTime() {
	s.PushToRadio(entity.FunctionMessage{Fn: s.InformOfLastLapTime}, 1, time.Minute)
}

func (s *Service) SendInformOfTotalWarnings() {
	s.PushToRadio(entity.FunctionMessage{Fn: s.InformOfTotalWarnings}, 1, time.Minute)
}
