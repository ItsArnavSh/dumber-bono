package parsers

// TyreSetData mirrors the C struct TyreSetData.
type TyreSetData struct {
	ActualTyreCompound uint8 // Actual tyre compound used
	VisualTyreCompound uint8 // Visual tyre compound used
	Wear               uint8 // Tyre wear (percentage)
	Available          uint8 // Whether this set is currently available
	RecommendedSession uint8 // Recommended session for tyre set, see appendix
	LifeSpan           uint8 // Laps left in this tyre set
	UsableLife         uint8 // Max number of laps recommended for this compound
	LapDeltaTime       int16 // Lap delta time in milliseconds compared to fitted set
	Fitted             uint8 // Whether the set is fitted or not
}

// PacketTyreSetsData mirrors the C struct PacketTyreSetsData.
type PacketTyreSetsData struct {
	Header      PacketHeader    // Header
	CarIdx      uint8           // Index of the car this data relates to
	TyreSetData [20]TyreSetData // 13 (dry) + 7 (wet)
	FittedIdx   uint8           // Index into array of fitted tyre
}
