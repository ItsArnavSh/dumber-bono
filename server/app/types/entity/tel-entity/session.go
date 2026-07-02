package telentity

/*
 The session packet includes details about the current session in progress.
 Frequency: 2 per second
 Size: 753 bytes
 Version: 1
*/

type RaceZoneFlag string

const (
	GREEN  RaceZoneFlag = "Green"
	BLUE   RaceZoneFlag = "Blue"
	YELLOW RaceZoneFlag = "Yellow"
	//No provision for double yellow in the docs...Austria aah
)

type MarshalZone struct {
	ZoneStart float32 //Fraction (0..1) of way through the lap the marshal zone starts
	ZoneFlag  RaceZoneFlag
}
