package telentity

type LobbyInfoData struct {
	AiControlled bool
	TeamId       string
	Nationality  string
	Name         string
	CarNumber    uint8
}

type PacketLobbyInfoData struct {
	Header       UDPHeader
	NumPlayers   uint8
	LobbyPlayers [22]LobbyInfoData
}
