package mappers

import (
	"errors"

	"dubmer-bono/app/api/udp/parsers"
	"dubmer-bono/app/types/entity/consts"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

func bytesToEventCode(code [4]uint8) string {
	n := 0
	for n < len(code) && code[n] != 0 {
		n++
	}
	return string(code[:n])
}

func MapToEventData(data *parsers.PacketEventData) (telentity.PacketEventData, error) {
	if data == nil {
		return telentity.PacketEventData{}, errors.New("mappers: nil PacketEventData")
	}

	header, err := MapToHeader(&data.Header)
	if err != nil {
		return telentity.PacketEventData{}, err
	}

	code := bytesToEventCode(data.EventStringCode)
	d := data.EventDetails

	var details telentity.EventDetails

	switch code {
	case "FTLP":
		details = telentity.FastestLap{
			VehicleIdx: d.FastestLap.VehicleIdx,
			LapTime:    d.FastestLap.LapTime,
		}
	case "RTMT":
		details = telentity.Retirement{
			VehicleIdx: d.Retirement.VehicleIdx,
			Reason:     consts.RetirementReason[int16(d.Retirement.Reason)],
		}
	case "TMPT":
		details = telentity.TeamMateInPits{
			VehicleIdx: d.TeamMateInPits.VehicleIdx,
		}
	case "RCWN":
		details = telentity.RaceWinner{
			VehicleIdx: d.RaceWinner.VehicleIdx,
		}
	case "PENA":
		details = telentity.Penalty{
			PenaltyType:      d.Penalty.PenaltyType,
			InfringementType: d.Penalty.InfringementType,
			VehicleIdx:       d.Penalty.VehicleIdx,
			OtherVehicleIdx:  d.Penalty.OtherVehicleIdx,
			Time:             d.Penalty.Time,
			LapNum:           d.Penalty.LapNum,
			PlacesGained:     d.Penalty.PlacesGained,
		}
	case "SPTP":
		details = telentity.SpeedTrap{
			VehicleIdx:                 d.SpeedTrap.VehicleIdx,
			Speed:                      d.SpeedTrap.Speed,
			IsOverallFastestInSession:  d.SpeedTrap.IsOverallFastestInSession != 0,
			IsDriverFastestInSession:   d.SpeedTrap.IsDriverFastestInSession != 0,
			FastestVehicleIdxInSession: d.SpeedTrap.FastestVehicleIdxInSession,
			FastestSpeedInSession:      d.SpeedTrap.FastestSpeedInSession,
		}
	case "STLG":
		details = telentity.StartLights{
			NumLights: d.StartLights.NumLights,
		}
	case "DTSV":
		details = telentity.DriveThroughPenaltyServed{
			VehicleIdx: d.DriveThroughPenaltyServed.VehicleIdx,
		}
	case "SGSV":
		details = telentity.StopGoPenaltyServed{
			VehicleIdx: d.StopGoPenaltyServed.VehicleIdx,
			StopTime:   d.StopGoPenaltyServed.StopTime,
		}
	case "FLBK":
		details = telentity.Flashback{
			FlashbackFrameIdentifier: d.Flashback.FlashbackFrameIdentifier,
			FlashbackSessionTime:     d.Flashback.FlashbackSessionTime,
		}
	case "BUTN":
		details = telentity.Buttons{
			ButtonStatus: d.Buttons.ButtonStatus,
		}
	case "OVTK":
		details = telentity.Overtake{
			OvertakingVehicleIdx:     d.Overtake.OvertakingVehicleIdx,
			BeingOvertakenVehicleIdx: d.Overtake.BeingOvertakenVehicleIdx,
		}
	case "SCAR":
		details = telentity.SafetyCarEvent{
			SafetyCarType: consts.SafetyCarType[int16(d.SafetyCar.SafetyCarType)],
			EventType:     consts.SafetyCarEventType[int16(d.SafetyCar.EventType)],
		}
	case "COLL":
		details = telentity.Collision{
			Vehicle1Idx: d.Collision.Vehicle1Idx,
			Vehicle2Idx: d.Collision.Vehicle2Idx,
		}
	case "DRSD":
		details = telentity.DRSDisabled{
			Reason: consts.DRSDisabledReason[int16(d.DRSDisabled.Reason)],
		}
	case "SSTA", "SEND", "DRSE", "CHQF", "RDFL", "LGOT":
		// No accompanying details struct for these events per spec.
		details = nil
	default:
		return telentity.PacketEventData{}, errors.New("mappers: unknown event string code: " + code)
	}

	return telentity.PacketEventData{
		Header:          header,
		EventStringCode: code,
		Details:         details,
	}, nil
}
