package monitor

import (
	"dubmer-bono/app/utility"
	"fmt"
	"time"
)

func (s *Service) InformOfPosition() {
	lapdata, _ := s.GetLapData()
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
		return
	}

	s.PushToRadio(fmt.Sprintf(template, position), 1, time.Minute)
}
func (s *Service) InformOfGapToFront() {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return
	}
	if lapdata.CarPosition == 1 {
		// leading the race, no car ahead to have a gap to
		return
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
		return
	}
	s.PushToRadio(fmt.Sprintf(template, gap.GapString()), 1, time.Minute)
}

func (s *Service) InformOfGapToLeader() {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return
	}
	if lapdata.CarPosition == 1 {
		return
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
		return
	}
	s.PushToRadio(fmt.Sprintf(template, gap), 1, time.Minute)
}

func (s *Service) InformOfCurrentLap() {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return
	}
	if lapdata.CurrentLapNum == 0 {
		// not yet in a valid lap (e.g. still on formation/pre-session)
		return
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
		return
	}
	s.PushToRadio(fmt.Sprintf(template, lap), 1, time.Minute)
}

func (s *Service) InformOfLastLapTime() {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return
	}
	if lapdata.LastLapTime == 0 {
		// no completed lap yet this session
		return
	}
	if lapdata.CurrentLapInvalid {
		// don't announce a time tied to an invalidated lap as if it were a clean reference
		return
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
		return
	}
	s.PushToRadio(fmt.Sprintf(template, lapTime), 1, time.Minute)
}

func (s *Service) InformOfPitStatus() {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return
	}
	if lapdata.PitStatus == "" || lapdata.PitStatus == "NONE" {
		// not in/approaching the pits, nothing worth radioing
		return
	}
	status := lapdata.PitStatus

	templates := []string{
		"Pit status %s",
		"You're %s",
		"Currently %s on strategy",
	}
	template, err := utility.RandomString(templates)
	if err != nil {
		s.logger.Errorf("failed to select pit status phrase: %v", err)
		return
	}
	s.PushToRadio(fmt.Sprintf(template, status), 1, time.Minute)
}

func (s *Service) InformOfNumPitStops() {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return
	}
	if lapdata.NumPitStops == 0 {
		// no stops made yet, nothing to report
		return
	}
	stops := lapdata.NumPitStops

	templates := []string{
		"You've made %d pit stops so far",
		"%d stops on this strategy",
		"Pit stop count %d",
	}
	template, err := utility.RandomString(templates)
	if err != nil {
		s.logger.Errorf("failed to select pit stop count phrase: %v", err)
		return
	}
	s.PushToRadio(fmt.Sprintf(template, stops), 1, time.Minute)
}

func (s *Service) InformOfSector() {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return
	}
	if lapdata.Sector == "" {
		return
	}
	sector := lapdata.Sector

	templates := []string{
		"Entering %s",
		"Now in %s",
		"%s coming up",
	}
	template, err := utility.RandomString(templates)
	if err != nil {
		s.logger.Errorf("failed to select sector phrase: %v", err)
		return
	}
	s.PushToRadio(fmt.Sprintf(template, sector), 1, time.Minute)
}

func (s *Service) InformOfGridPosition() {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return
	}
	if lapdata.GridPosition == 0 {
		// not a meaningful grid slot (e.g. pit lane start / data not populated)
		return
	}
	if lapdata.GridPosition == lapdata.CarPosition {
		// unchanged from the start, not interesting to call out
		return
	}
	grid := lapdata.GridPosition

	templates := []string{
		"You started P %d",
		"Grid position was P %d",
		"Started the race from P %d",
	}
	template, err := utility.RandomString(templates)
	if err != nil {
		s.logger.Errorf("failed to select grid position phrase: %v", err)
		return
	}
	s.PushToRadio(fmt.Sprintf(template, grid), 1, time.Minute)
}

func (s *Service) InformOfSpeedTrap() {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return
	}
	if lapdata.SpeedTrapFastestLap == 255 || lapdata.SpeedTrapFastestSpeed == 0 {
		// 255 = not set per the game's spec, or zero-value meaning no reading yet
		return
	}
	speed := lapdata.SpeedTrapFastestSpeed

	templates := []string{
		"Fastest speed trap so far %.0f kilometers per hour",
		"Top speed recorded %.0f kph",
		"Speed trap best %.0f kph",
	}
	template, err := utility.RandomString(templates)
	if err != nil {
		s.logger.Errorf("failed to select speed trap phrase: %v", err)
		return
	}
	s.PushToRadio(fmt.Sprintf(template, speed), 1, time.Minute)
}

func (s *Service) InformOfTotalWarnings() {
	lapdata, err := s.GetLapData()
	if err != nil {
		s.logger.Errorf("failed to get lap data: %v", err)
		return
	}
	if lapdata.TotalWarnings == 0 {
		// no warnings to report, skip rather than announce zero
		return
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
		return
	}
	s.PushToRadio(fmt.Sprintf(template, warnings), 1, time.Minute)
}
