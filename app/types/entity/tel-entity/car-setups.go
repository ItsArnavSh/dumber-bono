package telentity

/*
This packet details the car setups for each vehicle in the session. Note that in multiplayer games, other
player cars will appear as blank, you will only be able to see your own car setup, regardless of the
“Your Telemetry” setting. Spectators will also not be able to see any car setups.
Frequency: 2 per second
Size: 1133 bytes
Version: 1
*/

//Note: Not using in V1

type CarSetupData struct {
	FrontWing              uint8   // Front wing aero
	RearWing               uint8   // Rear wing aero
	OnThrottle             uint8   // Differential adjustment on throttle (percentage)
	OffThrottle            uint8   // Differential adjustment off throttle (percentage)
	FrontCamber            float32 // Front camber angle (suspension geometry)
	RearCamber             float32 // Rear camber angle (suspension geometry)
	FrontToe               float32 // Front toe angle (suspension geometry)
	RearToe                float32 // Rear toe angle (suspension geometry)
	FrontSuspension        uint8
	RearSuspension         uint8
	FrontAntiRollBar       uint8
	RearAntiRollBar        uint8
	FrontSuspensionHeight  uint8   // Front ride height
	RearSuspensionHeight   uint8   // Rear ride height
	BrakePressure          uint8   // Brake pressure (percentage)
	BrakeBias              uint8   // Brake bias (percentage)
	EngineBraking          uint8   // Engine braking (percentage)
	RearLeftTyrePressure   float32 // PSI
	RearRightTyrePressure  float32 // PSI
	FrontLeftTyrePressure  float32 // PSI
	FrontRightTyrePressure float32 // PSI
	Ballast                uint8
	FuelLoad               float32
}

type PacketCarSetupData struct {
	Header             UDPHeader
	CarSetups          [22]CarSetupData
	NextFrontWingValue float32 // Value of front wing after next pit stop - player only
}
