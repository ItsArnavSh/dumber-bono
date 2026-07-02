package telentity

type TyreSetData struct {
	ActualTyreCompound string
	VisualTyreCompound string
	Wear               uint8 // Tyre wear (percentage)
	Available          bool  // Whether this set is currently available
	RecommendedSession uint8 // Recommended session for tyre set, see appendix
	LifeSpan           uint8 // Laps left in this tyre set
	UsableLife         uint8 // Max number of laps recommended for this compound
	LapDeltaTime       int16 // Lap delta time in milliseconds compared to fitted set
	Fitted             bool  // Whether the set is fitted or not
}

type TyreSetPacket struct {
	header     UDPHeader
	carid      uint8 //Index to which it relates to
	tyredata   [20]TyreSetData
	fittedtyre uint8 //Fitted tyre
}
