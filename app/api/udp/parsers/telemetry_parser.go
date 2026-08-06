package parsers

// CarTelemetryData mirrors the C struct CarTelemetryData.
type CarTelemetryData struct {
	Speed                   uint16     // Speed of car in kilometres per hour
	Throttle                float32    // Amount of throttle applied (0.0 to 1.0)
	Steer                   float32    // Steering (-1.0 (full lock left) to 1.0 (full lock right))
	Brake                   float32    // Amount of brake applied (0.0 to 1.0)
	Clutch                  uint8      // Amount of clutch applied (0 to 100)
	Gear                    int8       // Gear selected (1-8, N=0, R=-1)
	EngineRPM               uint16     // Engine RPM
	Drs                     uint8      // 0 = off, 1 = on
	RevLightsPercent        uint8      // Rev lights indicator (percentage)
	RevLightsBitValue       uint16     // Rev lights (bit 0 = leftmost LED, bit 14 = rightmost LED)
	BrakesTemperature       [4]uint16  // Brakes temperature (celsius)
	TyresSurfaceTemperature [4]uint8   // Tyres surface temperature (celsius)
	TyresInnerTemperature   [4]uint8   // Tyres inner temperature (celsius)
	EngineTemperature       uint16     // Engine temperature (celsius)
	TyresPressure           [4]float32 // Tyres pressure (PSI)
	SurfaceType             [4]uint8   // Driving surface, see appendices
}

// PacketCarTelemetryData mirrors the C struct PacketCarTelemetryData.
type PacketCarTelemetryData struct {
	Header                       PacketHeader // Header
	CarTelemetryData             [22]CarTelemetryData
	MfdPanelIndex                uint8 // Index of MFD panel open - 255 = MFD closed. Single player, race – 0 = Car setup, 1 = Pits, 2 = Damage, 3 = Engine, 4 = Temperatures. May vary depending on game mode
	MfdPanelIndexSecondaryPlayer uint8 // See above
	SuggestedGear                int8  // Suggested gear for the player (1-8), 0 if no gear suggested
}
