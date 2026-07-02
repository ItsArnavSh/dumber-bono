package telentity

import "dubmer-bono/app/types/entity"

/*
The lap data packet gives details of all the cars in the session.
Frequency: Rate as specified in menus
Size: 1285 bytes
Version: 1
*/

type LapData struct {
	LastLapTime            uint32 //in MS
	CurrentLapTime         uint32 //in MS
	SectorTimes            entity.SectorWise
	DeltaToFront           entity.LapTimeStamp
	DeltaToRaceLeader      entity.LapTimeStamp
	LapDistance            float32
	TotalDistanceInSession float32
	SafetyCarDelta         float32 //In seconds
	CarPosition            uint8
	CurrentLapNum          uint8
	PitStatus              string
	NumPitStops            uint8
	Sector                 string

	CurrentLapInvalid           bool
	Penalties                   uint8 // Accumulated time penalties in seconds to be added
	TotalWarnings               uint8 // Accumulated number of warnings issued
	CornerCuttingWarnings       uint8 // Accumulated number of corner cutting warnings issued
	NumUnservedDriveThroughPens uint8 // Num drive through pens left to serve
	NumUnservedStopGoPens       uint8 // Num stop go pens left to serve
	GridPosition                uint8 // Grid position the vehicle started the race in
	DriverStatus                string
	ResultStatus                string
	PitLaneTimerActive          bool
	PitLaneTimeInLaneInMS       uint16 // If active, the current time spent in the pit lane in ms
	PitStopTimerInMS            uint16 // Time of the actual pit stop in ms
	PitStopShouldServePen       bool
	SpeedTrapFastestSpeed       float32 // Fastest speed through speed trap for this car in kmph
	SpeedTrapFastestLap         uint8   // Lap no the fastest speed was achieved, 255 = not set
}

type LapDataPacket struct {
	header  UDPHeader
	LapData [22]LapData
}
