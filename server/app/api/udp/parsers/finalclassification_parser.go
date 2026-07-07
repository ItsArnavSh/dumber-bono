package parsers

// FinalClassificationData mirrors the C struct FinalClassificationData.
type FinalClassificationData struct {
	Position          uint8    // Finishing position
	NumLaps           uint8    // Number of laps completed
	GridPosition      uint8    // Grid position of the car
	Points            uint8    // Number of points scored
	NumPitStops       uint8    // Number of pit stops made
	ResultStatus      uint8    // Result status - 0 = invalid, 1 = inactive, 2 = active, 3 = finished, 4 = didnotfinish, 5 = disqualified, 6 = not classified, 7 = retired
	ResultReason      uint8    // Result reason - 0 = invalid, 1 = retired, 2 = finished, 3 = terminal damage, 4 = inactive, 5 = not enough laps completed, 6 = black flagged, 7 = red flagged, 8 = mechanical failure, 9 = session skipped, 10 = session simulated
	BestLapTimeInMS   uint32   // Best lap time of the session in milliseconds
	TotalRaceTime     float64  // Total race time in seconds without penalties
	PenaltiesTime     uint8    // Total penalties accumulated in seconds
	NumPenalties      uint8    // Number of penalties applied to this driver
	NumTyreStints     uint8    // Number of tyres stints up to maximum
	TyreStintsActual  [8]uint8 // Actual tyres used by this driver
	TyreStintsVisual  [8]uint8 // Visual tyres used by this driver
	TyreStintsEndLaps [8]uint8 // The lap number stints end on
}

// PacketFinalClassificationData mirrors the C struct PacketFinalClassificationData.
type PacketFinalClassificationData struct {
	Header             PacketHeader // Header
	NumCars            uint8        // Number of cars in the final classification
	ClassificationData [22]FinalClassificationData
}
