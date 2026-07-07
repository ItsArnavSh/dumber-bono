package parsers

// PacketLapPositionsData mirrors the C struct PacketLapPositionsData.
type PacketLapPositionsData struct {
	Header                PacketHeader  // Header
	NumLaps               uint8         // Number of laps in the data
	LapStart              uint8         // Index of the lap where the data starts, 0 indexed
	PositionForVehicleIdx [50][22]uint8 // Array holding the position of the car in a given lap, 0 if no record
}
