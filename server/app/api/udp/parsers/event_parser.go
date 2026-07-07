package parsers

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// FastestLap mirrors the C struct FastestLap.
type FastestLap struct {
	VehicleIdx uint8   // Vehicle index of car achieving fastest lap
	LapTime    float32 // Lap time is in seconds
}

// Retirement mirrors the C struct Retirement.
type Retirement struct {
	VehicleIdx uint8 // Vehicle index of car retiring
	Reason     uint8 // Reason - 0 = invalid, 1 = retired, 2 = finished, 3 = terminal damage, 4 = inactive, 5 = not enough laps completed, 6 = black flagged, 7 = red flagged, 8 = mechanical failure, 9 = session skipped, 10 = session simulated
}

// DRSDisabled mirrors the C struct DRSDisabled.
type DRSDisabled struct {
	Reason uint8 // 0 = Wet track, 1 = Safety car deployed, 2 = Red flag, 3 = Min lap not reached
}

// TeamMateInPits mirrors the C struct TeamMateInPits.
type TeamMateInPits struct {
	VehicleIdx uint8 // Vehicle index of team mate
}

// RaceWinner mirrors the C struct RaceWinner.
type RaceWinner struct {
	VehicleIdx uint8 // Vehicle index of the race winner
}

// Penalty mirrors the C struct Penalty.
type Penalty struct {
	PenaltyType      uint8 // Penalty type – see Appendices
	InfringementType uint8 // Infringement type – see Appendices
	VehicleIdx       uint8 // Vehicle index of the car the penalty is applied to
	OtherVehicleIdx  uint8 // Vehicle index of the other car involved
	Time             uint8 // Time gained, or time spent doing action in seconds
	LapNum           uint8 // Lap the penalty occurred on
	PlacesGained     uint8 // Number of places gained by this
}

// SpeedTrap mirrors the C struct SpeedTrap.
type SpeedTrap struct {
	VehicleIdx                 uint8   // Vehicle index of the vehicle triggering speed trap
	Speed                      float32 // Top speed achieved in kilometres per hour
	IsOverallFastestInSession  uint8   // Overall fastest speed in session = 1, otherwise 0
	IsDriverFastestInSession   uint8   // Fastest speed for driver in session = 1, otherwise 0
	FastestVehicleIdxInSession uint8   // Vehicle index of the vehicle that is the fastest in this session
	FastestSpeedInSession      float32 // Speed of the vehicle that is the fastest in this session
}

// StartLights mirrors the C struct StartLIghts.
type StartLights struct {
	NumLights uint8 // Number of lights showing
}

// DriveThroughPenaltyServed mirrors the C struct DriveThroughPenaltyServed.
type DriveThroughPenaltyServed struct {
	VehicleIdx uint8 // Vehicle index of the vehicle serving drive through
}

// StopGoPenaltyServed mirrors the C struct StopGoPenaltyServed.
type StopGoPenaltyServed struct {
	VehicleIdx uint8   // Vehicle index of the vehicle serving stop go
	StopTime   float32 // Time spent serving stop go in seconds
}

// Flashback mirrors the C struct Flashback.
type Flashback struct {
	FlashbackFrameIdentifier uint32  // Frame identifier flashed back to
	FlashbackSessionTime     float32 // Session time flashed back to
}

// Buttons mirrors the C struct Buttons.
type Buttons struct {
	ButtonStatus uint32 // Bit flags specifying which buttons are being pressed currently - see appendices
}

// Overtake mirrors the C struct Overtake.
type Overtake struct {
	OvertakingVehicleIdx     uint8 // Vehicle index of the vehicle overtaking
	BeingOvertakenVehicleIdx uint8 // Vehicle index of the vehicle being overtaken
}

// SafetyCarEvent mirrors the C struct SafetyCar.
type SafetyCarEvent struct {
	SafetyCarType uint8 // 0 = No Safety Car, 1 = Full Safety Car, 2 = Virtual Safety Car, 3 = Formation Lap Safety Car
	EventType     uint8 // 0 = Deployed, 1 = Returning, 2 = Returned, 3 = Resume Race
}

// Collision mirrors the C struct Collision.
type Collision struct {
	Vehicle1Idx uint8 // Vehicle index of the first vehicle involved in the collision
	Vehicle2Idx uint8 // Vehicle index of the second vehicle involved in the collision
}

// EventDataDetails mirrors the C union EventDataDetails.
// Note: C unions overlay all members on the same memory; Go has no native union type.
// This struct lays out each variant sequentially and does NOT binary.Read correctly
// as-is for a true union — you'll need custom decode logic to interpret the raw bytes
// according to m_eventStringCode before populating the relevant field.
type EventDataDetails struct {
	FastestLap                FastestLap
	Retirement                Retirement
	DRSDisabled               DRSDisabled
	TeamMateInPits            TeamMateInPits
	RaceWinner                RaceWinner
	Penalty                   Penalty
	SpeedTrap                 SpeedTrap
	StartLights               StartLights
	DriveThroughPenaltyServed DriveThroughPenaltyServed
	StopGoPenaltyServed       StopGoPenaltyServed
	Flashback                 Flashback
	Buttons                   Buttons
	Overtake                  Overtake
	SafetyCar                 SafetyCarEvent
	Collision                 Collision
}

// PacketEventData mirrors the C struct PacketEventData.
type PacketEventData struct {
	Header          PacketHeader
	EventStringCode [4]uint8
	EventDetails    EventDataDetails
}

func ParseEventPacket(payload []byte) (*PacketEventData, error) {
	reader := bytes.NewReader(payload)

	var header PacketHeader
	if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
		return nil, fmt.Errorf("decode event header: %w", err)
	}

	var code [4]uint8
	if err := binary.Read(reader, binary.LittleEndian, &code); err != nil {
		return nil, fmt.Errorf("decode event code: %w", err)
	}

	packet := &PacketEventData{
		Header:          header,
		EventStringCode: code,
	}

	codeStr := bytesToEventCode(code)

	switch codeStr {
	case "FTLP":
		var d FastestLap
		if err := binary.Read(reader, binary.LittleEndian, &d); err != nil {
			return nil, fmt.Errorf("decode FastestLap: %w", err)
		}
		packet.EventDetails.FastestLap = d
	case "RTMT":
		var d Retirement
		if err := binary.Read(reader, binary.LittleEndian, &d); err != nil {
			return nil, fmt.Errorf("decode Retirement: %w", err)
		}
		packet.EventDetails.Retirement = d
	case "BUTN":
		var d Buttons
		if err := binary.Read(reader, binary.LittleEndian, &d); err != nil {
			return nil, fmt.Errorf("decode Buttons: %w", err)
		}
		packet.EventDetails.Buttons = d
	case "SSTA", "SEND", "DRSE", "CHQF", "RDFL", "LGOT":
		// no payload
	default:
		return nil, fmt.Errorf("unknown event code: %s", codeStr)
	}

	return packet, nil
}
func bytesToEventCode(code [4]uint8) string {
	n := 0
	for n < len(code) && code[n] != 0 {
		n++
	}
	return string(code[:n])
}
