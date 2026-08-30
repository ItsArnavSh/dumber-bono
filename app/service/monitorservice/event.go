package monitor

import (
	"dubmer-bono/app/types/entity"
	telentity "dubmer-bono/app/types/entity/tel-entity"
	"fmt"
	"time"
)

const (
	priorityCollision          = 5
	priorityRetirement         = 5
	prioritySafetyCarEvent     = 5
	priorityRaceWinner         = 4
	priorityPenalty            = 4
	priorityStopGoServed       = 2
	priorityDriveThroughServed = 2
	priorityStartLights        = 1
	priorityFastestLap         = 2
	priorityOvertake           = 2
	prioritySpeedTrap          = 1
	priorityTeamMateInPits     = 4
	priorityDRSDisabled        = 5
)

func (s *Service) EventMonitor(evt telentity.PacketEventData) {
	switch d := evt.Details.(type) {
	case telentity.FastestLap:
		s.handleFastestLap(d)
	case telentity.Retirement:
		s.handleRetirement(d)
	case telentity.DRSDisabled:
		s.handleDRSDisabled(d)
	case telentity.TeamMateInPits:
		s.handleTeamMateInPits(d)
	case telentity.RaceWinner:
		s.handleRaceWinner(d)
	case telentity.Penalty:
		s.handlePenalty(d)
	case telentity.DriveThroughPenaltyServed:
		s.handleDriveThroughPenaltyServed(d)
	case telentity.StopGoPenaltyServed:
		s.handleStopGoPenaltyServed(d)
	case telentity.Flashback:
		// no radio message, flashback is a rewind not an announceable event
	case telentity.Buttons:
		// no radio message, raw button state isn't announceable
	case telentity.Overtake:
		s.handleOvertake(d)
	case telentity.SafetyCarEvent:
		s.handleSafetyCarEvent(d)
	case telentity.Collision:
		s.handleCollision(d)
	case telentity.SpeedTrap:
		s.handleSpeedTrap(d)
	case telentity.StartLights:
		s.handleStartLights(d)
	default:
		//Ignore unhandled events
	}
}

func (s *Service) handleFastestLap(d telentity.FastestLap) {
	driver := s.getPlayerDetailsByID(int(d.VehicleIdx))
	msg := fmt.Sprintf("Fastest lap set by %s, %.3f seconds", driver.Driver_name, d.LapTime)
	s.PushToRadio(entity.DirectMessage{Text: msg}, priorityFastestLap, time.Minute)
}

func (s *Service) handleRetirement(d telentity.Retirement) {
	mycar := uint8(s.MyCarIndex())
	priority := 2
	if d.VehicleIdx == mycar {
		priority = 4
	}
	driver := s.getPlayerDetailsByID(int(d.VehicleIdx))
	msg := fmt.Sprintf("%s has retired, %s", driver.Driver_name, d.Reason)
	s.PushToRadio(entity.DirectMessage{Text: msg}, priority, time.Minute)
}

func (s *Service) handleDRSDisabled(d telentity.DRSDisabled) {
	msg := fmt.Sprintf("DRS disabled, %s", d.Reason)
	s.PushToRadio(entity.DirectMessage{Text: msg}, priorityDRSDisabled, time.Minute)
}

func (s *Service) handleTeamMateInPits(d telentity.TeamMateInPits) {
	driver := s.getPlayerDetailsByID(int(d.VehicleIdx))
	msg := fmt.Sprintf("Your team mate, %s, is in the pits", driver.Driver_name)
	s.PushToRadio(entity.DirectMessage{Text: msg}, priorityTeamMateInPits, time.Minute)
}

func (s *Service) handleRaceWinner(d telentity.RaceWinner) {
	driver := s.getPlayerDetailsByID(int(d.VehicleIdx))
	msg := fmt.Sprintf("%s has won the race", driver.Driver_name)
	s.PushToRadio(entity.DirectMessage{Text: msg}, priorityRaceWinner, time.Minute)
}

func (s *Service) handlePenalty(d telentity.Penalty) {
	mycar := uint8(s.MyCarIndex())
	priority := 2
	if d.VehicleIdx == mycar {
		priority = 4
	}
	driver := s.getPlayerDetailsByID(int(d.VehicleIdx))
	msg := fmt.Sprintf("Penalty for %s, %d seconds, %d places gained", driver.Driver_name, d.Time, d.PlacesGained)
	s.PushToRadio(entity.DirectMessage{Text: msg}, priority, time.Minute)
}

func (s *Service) handleSpeedTrap(d telentity.SpeedTrap) {
	driver := s.getPlayerDetailsByID(int(d.VehicleIdx))
	msg := fmt.Sprintf("Speed trap, %s clocked %.1f km/h", driver.Driver_name, d.Speed)
	s.PushToRadio(entity.DirectMessage{Text: msg}, prioritySpeedTrap, time.Minute)
}

func (s *Service) handleStartLights(d telentity.StartLights) {
	msg := fmt.Sprintf("%d lights showing", d.NumLights)
	s.PushToRadio(entity.DirectMessage{Text: msg}, priorityStartLights, time.Minute)
}

func (s *Service) handleDriveThroughPenaltyServed(d telentity.DriveThroughPenaltyServed) {
	mycar := uint8(s.MyCarIndex())
	priority := 2
	if d.VehicleIdx == mycar {
		priority = 4
	}
	driver := s.getPlayerDetailsByID(int(d.VehicleIdx))
	msg := fmt.Sprintf("%s has served their drive through penalty", driver.Driver_name)
	s.PushToRadio(entity.DirectMessage{Text: msg}, priority, time.Minute)
}

func (s *Service) handleStopGoPenaltyServed(d telentity.StopGoPenaltyServed) {
	mycar := uint8(s.MyCarIndex())
	priority := 2
	if d.VehicleIdx == mycar {
		priority = 4
	}
	driver := s.getPlayerDetailsByID(int(d.VehicleIdx))
	msg := fmt.Sprintf("%s has served their stop go penalty, %.1f seconds", driver.Driver_name, d.StopTime)
	s.PushToRadio(entity.DirectMessage{Text: msg}, priority, time.Minute)
}

func (s *Service) handleOvertake(d telentity.Overtake) {
	mycar := uint8(s.MyCarIndex())

	if d.OvertakingVehicleIdx == mycar || d.BeingOvertakenVehicleIdx == mycar {
		// you already know you overtook someone, just get your new position out fast
		s.PushToRadio(entity.FunctionMessage{Fn: s.InformOfPosition}, 3, 10*time.Second)
		return
	}

	overtakerPos := s.getCarPosition(d.OvertakingVehicleIdx)
	overtakenPos := s.getCarPosition(d.BeingOvertakenVehicleIdx)
	mypos := s.getCarPosition(mycar)

	nearMe := isNearPosition(mypos, overtakerPos) || isNearPosition(mypos, overtakenPos)
	topThree := overtakerPos <= 3 || overtakenPos <= 3

	if !nearMe && !topThree {
		return
	}

	overtaker := s.getPlayerDetailsByID(int(d.OvertakingVehicleIdx))
	overtaken := s.getPlayerDetailsByID(int(d.BeingOvertakenVehicleIdx))
	msg := fmt.Sprintf("%s has overtaken %s, now P%d", overtaker.Driver_name, overtaken.Driver_name, overtakerPos)
	s.PushToRadio(entity.DirectMessage{Text: msg}, 1, time.Minute)
}

func isNearPosition(mypos, pos int) bool {
	diff := mypos - pos
	if diff < 0 {
		diff = -diff
	}
	return diff <= 2
}
func (s *Service) handleSafetyCarEvent(d telentity.SafetyCarEvent) {
	msg := fmt.Sprintf("Safety car, %s, %s", d.SafetyCarType, d.EventType)
	s.PushToRadio(entity.DirectMessage{Text: msg}, prioritySafetyCarEvent, time.Minute)
}

func (s *Service) handleCollision(d telentity.Collision) {
	mycar := uint8(s.MyCarIndex())

	if d.Vehicle1Idx == mycar || d.Vehicle2Idx == mycar {
		// you already felt the hit, just inform
		return
	}

	car1 := s.getPlayerDetailsByID(int(d.Vehicle1Idx))
	car2 := s.getPlayerDetailsByID(int(d.Vehicle2Idx))
	lap, _ := s.GetLapData()
	if lap.CurrentLapNum != 1 {
		msg := fmt.Sprintf("Collision between %s and %s", car1.Driver_name, car2.Driver_name)

		s.PushToRadio(entity.DirectMessage{Text: msg}, 2, time.Minute)
	}
}
