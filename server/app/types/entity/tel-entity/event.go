package telentity

/*
This packet gives details of events that happen during the course of a session.
Frequency: When the event occurs
Size: 45 bytes
Version: 1
*/

// EventDataDetails is a C union in the source spec - Go has no union type,
// so each variant becomes its own struct below. Only the struct matching
// the current EventStringCode is valid to read; the rest are meaningless
// for a given event. Check EventStringCode (or use the EventCodeDescriptions
// map below) before deciding which field on PacketEventData to read.

// EventDataDetails is a C union in the source spec - Go has no union type,
// so each variant becomes its own struct below. Only the struct matching
// the current EventStringCode is valid to read; the rest are meaningless
// for a given event. Check EventStringCode (or use the EventCodeDescriptions
// map below) before deciding which field on PacketEventData to read.

type FastestLap struct {
	VehicleIdx uint8   // Vehicle index of car achieving fastest lap
	LapTime    float32 // Lap time in seconds
}

type Retirement struct {
	VehicleIdx uint8 // Vehicle index of car retiring
	Reason     string
}

type DRSDisabled struct {
	Reason string
}

type TeamMateInPits struct {
	VehicleIdx uint8 // Vehicle index of team mate
}

type RaceWinner struct {
	VehicleIdx uint8 // Vehicle index of the race winner
}

type Penalty struct {
	PenaltyType      uint8 // see Appendices
	InfringementType uint8 // see Appendices
	VehicleIdx       uint8 // Vehicle index of the car the penalty is applied to
	OtherVehicleIdx  uint8 // Vehicle index of the other car involved
	Time             uint8 // Time gained, or time spent doing action in seconds
	LapNum           uint8 // Lap the penalty occurred on
	PlacesGained     uint8 // Number of places gained by this
}

type SpeedTrap struct {
	VehicleIdx                 uint8   // Vehicle index of the vehicle triggering speed trap
	Speed                      float32 // Top speed achieved in km/h
	IsOverallFastestInSession  bool
	IsDriverFastestInSession   bool
	FastestVehicleIdxInSession uint8   // Vehicle index of the fastest vehicle in this session
	FastestSpeedInSession      float32 // Speed of the fastest vehicle in this session
}

type StartLights struct {
	NumLights uint8 // Number of lights showing
}

type DriveThroughPenaltyServed struct {
	VehicleIdx uint8 // Vehicle index of the vehicle serving drive through
}

type StopGoPenaltyServed struct {
	VehicleIdx uint8   // Vehicle index of the vehicle serving stop go
	StopTime   float32 // Time spent serving stop go in seconds
}

type Flashback struct {
	FlashbackFrameIdentifier uint32  // Frame identifier flashed back to
	FlashbackSessionTime     float32 // Session time flashed back to
}

type Buttons struct {
	ButtonStatus uint32 // Bit flags specifying which buttons are being pressed - see appendices
}

type Overtake struct {
	OvertakingVehicleIdx     uint8 // Vehicle index of the vehicle overtaking
	BeingOvertakenVehicleIdx uint8 // Vehicle index of the vehicle being overtaken
}

type SafetyCarEvent struct {
	SafetyCarType string
	EventType     string
}

type Collision struct {
	Vehicle1Idx uint8 // Vehicle index of the first vehicle involved in the collision
	Vehicle2Idx uint8 // Vehicle index of the second vehicle involved in the collision
}

// EventDetails is implemented by every event variant struct below.
// It exists purely as a marker so PacketEventData.Details can hold any
// one of them - there's no shared behaviour to call, since the structs
// don't share fields (Penalty has 7 fields, Collision has 2). A type
// switch on Details is how you recover the concrete type; see
// PacketEventData.EventStringCode / EventCodeDescriptions to know which
// one to expect before switching.
type EventDetails interface {
	isEventDetails()
}

func (FastestLap) isEventDetails()                {}
func (Retirement) isEventDetails()                {}
func (DRSDisabled) isEventDetails()               {}
func (TeamMateInPits) isEventDetails()            {}
func (RaceWinner) isEventDetails()                {}
func (Penalty) isEventDetails()                   {}
func (SpeedTrap) isEventDetails()                 {}
func (StartLights) isEventDetails()               {}
func (DriveThroughPenaltyServed) isEventDetails() {}
func (StopGoPenaltyServed) isEventDetails()       {}
func (Flashback) isEventDetails()                 {}
func (Buttons) isEventDetails()                   {}
func (Overtake) isEventDetails()                  {}
func (SafetyCarEvent) isEventDetails()            {}
func (Collision) isEventDetails()                 {}

// PacketEventData mirrors the union with a single interface-typed field.
// Details holds exactly one of the event structs above (as a value, not
// a pointer) - use EventStringCode to know which type to expect, then a
// type switch to recover it:
//
//	switch d := packet.Details.(type) {
//	case Retirement:
//	    // d.VehicleIdx, d.Reason
//	case SpeedTrap:
//	    // d.VehicleIdx, d.Speed, ...
//	}
type PacketEventData struct {
	Header          UDPHeader
	EventStringCode string // 4-char code, see EventCodeDescriptions below
	Details         EventDetails
}
