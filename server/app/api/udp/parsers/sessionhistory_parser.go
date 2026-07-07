package parsers

// LapHistoryData mirrors the C struct LapHistoryData.
type LapHistoryData struct {
	LapTimeInMS            uint32 // Lap time in milliseconds
	Sector1TimeMSPart      uint16 // Sector 1 milliseconds part
	Sector1TimeMinutesPart uint8  // Sector 1 whole minute part
	Sector2TimeMSPart      uint16 // Sector 2 time milliseconds part
	Sector2TimeMinutesPart uint8  // Sector 2 whole minute part
	Sector3TimeMSPart      uint16 // Sector 3 time milliseconds part
	Sector3TimeMinutesPart uint8  // Sector 3 whole minute part
	LapValidBitFlags       uint8  // 0x01 bit set-lap valid, 0x02 bit set-sector 1 valid, 0x04 bit set-sector 2 valid, 0x08 bit set-sector 3 valid
}

// TyreStintHistoryData mirrors the C struct TyreStintHistoryData.
type TyreStintHistoryData struct {
	EndLap             uint8 // Lap the tyre usage ends on (255 of current tyre)
	TyreActualCompound uint8 // Actual tyres used by this driver
	TyreVisualCompound uint8 // Visual tyres used by this driver
}

// PacketSessionHistoryData mirrors the C struct PacketSessionHistoryData.
type PacketSessionHistoryData struct {
	Header                PacketHeader        // Header
	CarIdx                uint8               // Index of the car this lap data relates to
	NumLaps               uint8               // Num laps in the data (including current partial lap)
	NumTyreStints         uint8               // Number of tyre stints in the data
	BestLapTimeLapNum     uint8               // Lap the best lap time was achieved on
	BestSector1LapNum     uint8               // Lap the best Sector 1 time was achieved on
	BestSector2LapNum     uint8               // Lap the best Sector 2 time was achieved on
	BestSector3LapNum     uint8               // Lap the best Sector 3 time was achieved on
	LapHistoryData        [100]LapHistoryData // 100 laps of data max
	TyreStintsHistoryData [8]TyreStintHistoryData
}
