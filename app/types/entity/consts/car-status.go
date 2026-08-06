package consts

var TractionControl = map[int16]string{
	0: "off",
	1: "medium",
	2: "full",
}

var FuelMix = map[int16]string{
	0: "lean",
	1: "standard",
	2: "rich",
	3: "max",
}

// ActualTyreCompound covers F1 Modern, F1 Classic, and F2 - IDs overlap
// across formulas (e.g. 7/8 mean inter/wet in both F1 Modern and F2),
// so use in conjunction with the Formula field to disambiguate.
var ActualTyreCompound = map[int16]string{
	7:  "inter",
	8:  "wet",
	9:  "dry (F1 Classic)",
	10: "wet (F1 Classic)",
	11: "super soft (F2)",
	12: "soft (F2)",
	13: "medium (F2)",
	14: "hard (F2)",
	15: "wet (F2)",
	16: "C5",
	17: "C4",
	18: "C3",
	19: "C2",
	20: "C1",
	21: "C0",
	22: "C6",
}

// VisualTyreCompound - F2 '20 IDs overlap with F1 Modern IDs, use in
// conjunction with the Formula field to disambiguate.
var VisualTyreCompound = map[int16]string{
	7:  "inter",
	8:  "wet",
	15: "wet (F2 '20)",
	16: "soft",
	17: "medium",
	18: "hard",
	19: "super soft (F2 '20)",
	20: "soft (F2 '20)",
	21: "medium (F2 '20)",
	22: "hard (F2 '20)",
}

var VehicleFiaFlags = map[int16]string{
	-1: "invalid/unknown",
	0:  "none",
	1:  "green",
	2:  "blue",
	3:  "yellow",
}

var ERSDeployMode = map[int16]string{
	0: "none",
	1: "medium",
	2: "hotlap",
	3: "overtake",
}
