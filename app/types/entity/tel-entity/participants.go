package telentity

import (
	"errors"
)

/*
This is a list of participants in the race. If the vehicle is controlled by AI, then the name will be the
driver name. If this is a multiplayer game, the names will be the Steam Id on PC, or the LAN name if
appropriate.
N.B. on Xbox, the names will always be the driver name, on PlayStation the name will be the LAN
name if playing a LAN game, otherwise it will be the driver name.
The array should be indexed by vehicle index.
Frequency: Every 5 second
*/
type LiveryColour struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

type ParticipantData struct {
	// AIControlled    bool
	DriverId uint8 // see appendix, 255 if network human
	// NetworkId       uint8 // unique identifier for network players
	TeamId        uint8 // see appendix
	MyTeam        bool
	RaceNumber    uint8
	Nationality   uint8  // see appendix
	Name          string // UTF-8, truncated with … (U+2026) if too long
	YourTelemetry string
	// ShowOnlineNames bool
	// TechLevel       uint16 // F1 World tech level
	// Platform        string
	// NumColours      uint8 // Number of colours valid for this car
	// LiveryColours   [4]LiveryColour
}

type PacketParticipantsData struct {
	Header        UDPHeader
	NumActiveCars uint8 // Number of active cars - should match number of cars on HUD
	Participants  [22]ParticipantData
}

func (p *PacketParticipantsData) Me() (ParticipantData, int, error) {
	idx := int(p.Header.PlayerCarIndex)

	if idx >= len(p.Participants) || idx >= int(p.NumActiveCars) {
		return ParticipantData{}, -1, errors.New("telentity: invalid player car index")
	}

	return p.Participants[idx], idx, nil
}
