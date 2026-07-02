package consts

var RetirementReason = map[int16]string{
	0:  "invalid",
	1:  "retired",
	2:  "finished",
	3:  "terminal damage",
	4:  "inactive",
	5:  "not enough laps completed",
	6:  "black flagged",
	7:  "red flagged",
	8:  "mechanical failure",
	9:  "session skipped",
	10: "session simulated",
}

var DRSDisabledReason = map[int16]string{
	0: "Wet track",
	1: "Safety car deployed",
	2: "Red flag",
	3: "Min lap not reached",
}

var SafetyCarType = map[int16]string{
	0: "No Safety Car",
	1: "Full Safety Car",
	2: "Virtual Safety Car",
	3: "Formation Lap Safety Car",
}

var SafetyCarEventType = map[int16]string{
	0: "Deployed",
	1: "Returning",
	2: "Returned",
	3: "Resume Race",
}

// EventCodeDescriptions maps the 4-char m_eventStringCode to a
// human-readable description of the event.
var EventCodeDescriptions = map[string]string{
	"SSTA": "Session Started - sent when the session starts",
	"SEND": "Session Ended - sent when the session ends",
	"FTLP": "Fastest Lap - when a driver achieves the fastest lap",
	"RTMT": "Retirement - when a driver retires",
	"DRSE": "DRS enabled - race control have enabled DRS",
	"DRSD": "DRS disabled - race control have disabled DRS",
	"TMPT": "Team mate in pits - your team mate has entered the pits",
	"CHQF": "Chequered flag - the chequered flag has been waved",
	"RCWN": "Race Winner - the race winner is announced",
	"PENA": "Penalty Issued - a penalty has been issued, details in event",
	"SPTP": "Speed Trap Triggered - speed trap has been triggered by fastest speed",
	"STLG": "Start lights - number shown",
	"LGOT": "Lights out",
	"DTSV": "Drive through served - drive through penalty served",
	"SGSV": "Stop go served - stop go penalty served",
	"FLBK": "Flashback - flashback activated",
	"BUTN": "Button status - button status changed",
	"RDFL": "Red Flag - red flag shown",
	"OVTK": "Overtake - overtake occurred",
	"SCAR": "Safety Car - safety car event, details in event",
	"COLL": "Collision - collision between two vehicles has occurred",
}
