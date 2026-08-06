package consts

var PitStatus = map[int16]string{0: "none",
	1: "pitting",
	2: "in pit area",
}
var CurrentSector = map[int16]string{0: "sector1",
	1: "sector2",
	2: "sector3",
}
var Sector = map[int16]string{
	0: "sector1",
	1: "sector2",
	2: "sector3",
}

var DriverStatus = map[int16]string{
	0: "in garage",
	1: "flying lap",
	2: "in lap",
	3: "out lap",
	4: "on track",
}

var ResultStatus = map[int16]string{
	0: "invalid",
	1: "inactive",
	2: "active",
	3: "finished",
	4: "didnotfinish",
	5: "disqualified",
	6: "not classified",
	7: "retired",
}
