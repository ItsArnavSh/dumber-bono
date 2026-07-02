package telentity

/*
This packet details car statuses for all the cars in the race.
Frequency: Rate as specified in menus
Size: 1239 bytes
Version: 1
*/

type CarStatusData struct {
	TractionControl         string
	AntiLockBrakes          bool
	FuelMix                 string
	FrontBrakeBias          uint8 // Front brake bias (percentage)
	PitLimiterStatus        bool
	FuelInTank              float32 // Current fuel mass
	FuelCapacity            float32
	FuelRemainingLaps       float32 // Fuel remaining in terms of laps (value on MFD)
	MaxRPM                  uint16
	IdleRPM                 uint16
	MaxGears                uint8
	DRSAllowed              bool
	DRSActivationDistance   uint16 // 0 = DRS not available, non-zero = available in X metres
	ActualTyreCompound      string
	VisualTyreCompound      string
	TyresAgeLaps            uint8 // Age in laps of the current set of tyres
	VehicleFiaFlags         string
	EnginePowerICE          float32 // Engine power output of ICE (W)
	EnginePowerMGUK         float32 // Engine power output of MGU-K (W)
	ERSStoreEnergy          float32 // ERS energy store in Joules
	ERSDeployMode           string
	ERSHarvestedThisLapMGUK float32
	ERSHarvestedThisLapMGUH float32
	ERSDeployedThisLap      float32
	NetworkPaused           bool
}

type PacketCarStatusData struct {
	Header        UDPHeader
	CarStatusData [22]CarStatusData
}
