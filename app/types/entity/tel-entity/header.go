package telentity

type UDPHeader struct { //Header with every packet
	PacketFormat            uint16 //2025 for my version
	GameYear                uint8  //Last two digits of game year
	GameMajorVersion        uint8  // X.00
	GameMinorVersion        uint8  //1.XX
	PacketVersion           uint8
	PacketID                uint8
	SessionUID              uint64  //Unique identifier for the session
	SessionTime             float32 //Session Timestamp
	FrameIdentifier         uint32  //For ordering purposes, does not go back on flashbacks
	PlayerCarIndex          uint8   //Index of Player's car in the Array
	SecondaryPlayerCarIndex uint8   //In case of Splitscreen (Can be a future feature)
}
