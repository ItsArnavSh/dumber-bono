package monitor

import (
	"dubmer-bono/app/utility"
	"fmt"
	"time"
)

func (s *Service) InformOfPosition() string {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return ""
	}
	position := lapdata.CarPosition
	templates := []string{
		"You are currently P %d",
		"Currently running P %d",
		"You're in P %d right now",
		"Position check, P %d",
		"Holding P %d at the moment",
	}
	template, err := utility.RandomString(templates)
	if err != nil {
		s.logger.Errorf("failed to select position phrase: %v", err)
		return ""
	}
	return fmt.Sprintf(template, position)
}

func (s *Service) InformOfGapToFront() string {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return ""
	}
	if lapdata.CarPosition == 1 {
		// leading the race, no car ahead to have a gap to
		return ""
	}
	gap := lapdata.DeltaToFront
	templates := []string{
		"Gap to car ahead %s",
		"You're %s behind the car in front",
		"Delta to front is %s",
	}
	template, err := utility.RandomString(templates)
	if err != nil {
		s.logger.Errorf("failed to select gap to front phrase: %v", err)
		return ""
	}
	return fmt.Sprintf(template, gap.GapString())
}

func (s *Service) InformOfGapToLeader() string {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return ""
	}
	if lapdata.CarPosition == 1 {
		return ""
	}
	gap := lapdata.DeltaToRaceLeader.GapString()
	templates := []string{
		"Gap to the leader %s",
		"You're %s off the race leader",
		"Delta to leader is %s",
	}
	template, err := utility.RandomString(templates)
	if err != nil {
		s.logger.Errorf("failed to select gap to leader phrase: %v", err)
		return ""
	}
	return fmt.Sprintf(template, gap)
}

func (s *Service) InformOfCurrentLap() string {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return ""
	}
	if lapdata.CurrentLapNum == 0 {
		// not yet in a valid lap (e.g. still on formation/pre-session)
		return ""
	}
	lap := lapdata.CurrentLapNum
	templates := []string{
		"Currently on lap %d",
		"Lap %d now",
		"This is lap %d",
	}
	template, err := utility.RandomString(templates)
	if err != nil {
		s.logger.Errorf("failed to select current lap phrase: %v", err)
		return ""
	}
	return fmt.Sprintf(template, lap)
}

func (s *Service) InformOfLastLapTime() string {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return ""
	}
	if lapdata.LastLapTime == 0 {
		// no completed lap yet this session
		return ""
	}
	if lapdata.CurrentLapInvalid {
		// don't announce a time tied to an invalidated lap as if it were a clean reference
		return ""
	}
	lapTime := time.Duration(lapdata.LastLapTime) * time.Millisecond
	templates := []string{
		"Last lap %s",
		"That lap was %s",
		"Previous lap time %s",
	}
	template, err := utility.RandomString(templates)
	if err != nil {
		s.logger.Errorf("failed to select last lap time phrase: %v", err)
		return ""
	}
	return fmt.Sprintf(template, lapTime)
}

func (s *Service) InformOfTotalWarnings() string {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return ""
	}
	if lapdata.TotalWarnings == 0 {
		// no warnings to report, skip rather than announce zero
		return ""
	}
	warnings := lapdata.TotalWarnings
	templates := []string{
		"You've got %d warnings so far",
		"Warning count is %d",
		"%d warnings on record",
	}
	template, err := utility.RandomString(templates)
	if err != nil {
		s.logger.Errorf("failed to select warnings phrase: %v", err)
		return ""
	}
	return fmt.Sprintf(template, warnings)
}
