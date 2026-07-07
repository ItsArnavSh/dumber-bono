package parsers

// LiveryColour mirrors the C struct LiveryColour - RGB value of a colour.
type LiveryColour struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// ParticipantData mirrors the C struct ParticipantData.
type ParticipantData struct {
	AiControlled    uint8           // Whether the vehicle is AI (1) or Human (0) controlled
	DriverId        uint8           // Driver id - see appendix, 255 if network human
	NetworkId       uint8           // Network id – unique identifier for network players
	TeamId          uint8           // Team id - see appendix
	MyTeam          uint8           // My team flag – 1 = My Team, 0 = otherwise
	RaceNumber      uint8           // Race number of the car
	Nationality     uint8           // Nationality of the driver
	Name            [32]byte        // Name of participant in UTF-8 format – null terminated. Will be truncated with … (U+2026) if too long
	YourTelemetry   uint8           // The player's UDP setting, 0 = restricted, 1 = public
	ShowOnlineNames uint8           // The player's show online names setting, 0 = off, 1 = on
	TechLevel       uint16          // F1 World tech level
	Platform        uint8           // 1 = Steam, 3 = PlayStation, 4 = Xbox, 6 = Origin, 255 = unknown
	NumColours      uint8           // Number of colours valid for this car
	LiveryColours   [4]LiveryColour // Colours for the car
}

// PacketParticipantsData mirrors the C struct PacketParticipantsData.
type PacketParticipantsData struct {
	Header        PacketHeader // Header
	NumActiveCars uint8        // Number of active cars in the data – should match number of cars on HUD
	Participants  [22]ParticipantData
}
