package telentity

/*
This packet contains lap times and tyre usage for the session. This packet works slightly differently
to other packets. To reduce CPU and bandwidth, each packet relates to a specific vehicle and is
sent every 1/20 s, and the vehicle being sent is cycled through. Therefore in a 20 car race you
should receive an update for each vehicle at least once per second.
Note that at the end of the race, after the final classification packet has been sent, a final bulk update
of all the session histories for the vehicles in that session will be sent.
Frequency: 20 per second but cycling through cars
Size: 1460 bytes
Version: 1
*/

type LapHistoryData struct {
	LapTimeInMS            uint32 // Lap time in milliseconds
	Sector1TimeMSPart      uint16 // Sector 1 milliseconds part
	Sector1TimeMinutesPart uint8  // Sector 1 whole minute part
	Sector2TimeMSPart      uint16 // Sector 2 time milliseconds part
	Sector2TimeMinutesPart uint8  // Sector 2 whole minute part
	Sector3TimeMSPart      uint16 // Sector 3 time milliseconds part
	Sector3TimeMinutesPart uint8  // Sector 3 whole minute part
	LapValidBitFlags       uint8  // bitmask, see IsLapValid / IsSectorValid below
}

const (
	LapValidFlag     uint8 = 0x01
	Sector1ValidFlag uint8 = 0x02
	Sector2ValidFlag uint8 = 0x04
	Sector3ValidFlag uint8 = 0x08
)

func (l LapHistoryData) IsLapValid() bool {
	return l.LapValidBitFlags&LapValidFlag != 0
}

func (l LapHistoryData) IsSector1Valid() bool {
	return l.LapValidBitFlags&Sector1ValidFlag != 0
}

func (l LapHistoryData) IsSector2Valid() bool {
	return l.LapValidBitFlags&Sector2ValidFlag != 0
}

func (l LapHistoryData) IsSector3Valid() bool {
	return l.LapValidBitFlags&Sector3ValidFlag != 0
}

type TyreStintHistoryData struct {
	EndLap             uint8 // Lap the tyre usage ends on (255 = current tyre)
	TyreActualCompound string
	TyreVisualCompound string
}

type SessionHistoryPacket struct {
	Header                UDPHeader
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
