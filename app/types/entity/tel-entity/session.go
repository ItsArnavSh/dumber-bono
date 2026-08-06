package telentity

/*
 The session packet includes details about the current session in progress.
 Frequency: 2 per second
 Size: 753 bytes
 Version: 1
*/

type MarshalZone struct {
	ZoneStart float32 //Fraction (0..1) of way through the lap the marshal zone starts
	ZoneFlag  string
}

type WeatherForecastSample struct {
	SessionType            string //Ref SessionTypes in f125-consts.go
	TimeOffset             uint8  // Time in minutes the forecast is for
	Weather                string
	TrackTemperature       int8 //In deg Cel
	TrackTemperatureChange string
	AirTemperature         int8 //In deg Cel
	AirTemperatureChange   string
	RainPercentage         uint8 //Chance of Rain
}

type PacketSessionData struct {
	Header           UDPHeader
	Weather          string
	TrackTemperature int8
	AirTempterature  int8
	TotalLaps        int8
	SessionType      string
	Track            string
	Formula          string
	SessionTimeLeft  uint16 // Time left in seconds
	SessionDuration  uint16
	PitSpeedLimit    uint8 // Speed limit in km/h
	GamePaused       bool  // Network Game Only
	// IsSpectating              bool
	// SpectatorCarIndex         uint8
	NumMarshalZones        uint8           // Number of marshal zones to follow
	MarshalZones           [21]MarshalZone // List of Marshal Zones
	SafetyCarStatus        string
	NetworkGame            string
	WeatherForecastSamples []WeatherForecastSample // fixed: array, not single struct

	//SeasonLinkIdentifier   uint32 // Persists across saves
	//WeekendLinkIdentifier  uint32 // Persists across saves
	//SessionLinkIdentifier  uint32 // Persists across saves
	PitStopWindowIdealLap  uint8 // Ideal lap to pit on for current strategy (player)
	PitStopWindowLatestLap uint8 // Latest lap to pit on for current strategy (player)
	PitStopRejoinPosition  uint8 // Predicted position to rejoin at (player)

	//SteeringAssist        bool
	//BrakingAssist         string
	//GearboxAssist         string
	//PitAssist             bool
	//PitReleaseAssist      bool
	//ERSAssist             bool
	//DRSAssist             bool
	//DynamicRacingLine     string
	//DynamicRacingLineType string
	GameMode      uint8  // Game mode id - see appendix
	RuleSet       uint8  // Ruleset - see appendix
	TimeOfDay     uint32 // Local time of day - minutes since midnight
	SessionLength string

	// SpeedUnitsLeadPlayer            string
	// TemperatureUnitsLeadPlayer      string
	// SpeedUnitsSecondaryPlayer       string
	// TemperatureUnitsSecondaryPlayer string

	// NumSafetyCarPeriods        uint8
	// NumVirtualSafetyCarPeriods uint8
	// NumRedFlagPeriods          uint8

	// EqualCarPerformance          bool
	// RecoveryMode                 string
	// FlashbackLimit               string
	// SurfaceType                  string
	// LowFuelMode                  string
	// RaceStarts                   string
	// TyreTemperature              string
	// PitLaneTyreSim               bool
	// CarDamage                    string
	// CarDamageRate                string
	// Collisions                   string
	// CollisionsOffForFirstLapOnly bool
	// MPUnsafePitRelease           bool // Multiplayer
	// MPOffForGriefing             bool // Multiplayer
	// CornerCuttingStringency      string
	// ParcFermeRules               bool
	// PitStopExperience            string
	// SafetyCar                    string
	// SafetyCarExperience          string
	// FormationLap                 bool
	// FormationLapExperience       string
	// RedFlags                     string
	// AffectsLicenceLevelSolo      bool
	// AffectsLicenceLevelMP        bool

	// NumSessionsInWeekend uint8
	// WeekendStructure     [12]uint8 // Session type ids for weekend structure - see appendix

	Sector2LapDistanceStart float32 // Distance in m around track where sector 2 starts
	Sector3LapDistanceStart float32 // Distance in m around track where sector 3 starts
}
