package consts

var TempChange = map[int16]string{0: "up",
	1: "down",
	2: "no change",
}
var WeatherForecastTypes = map[uint16]string{
	0: "clear",
	1: "light cloud",
	2: "overcast",
	3: "light rain",
	4: "heavy rain",
	5: "storm",
}
var Formula = map[int16]string{0: "F1 Modern",
	1: "F1 Classic",
	2: "F2",
	3: "F1 Generic",
	4: "Beta",
	6: "Esports",
	8: "F1 World",
	9: "F1 Elimination",
}
var SafetyCarStatus = map[int16]string{0: "no safety car",
	1: "full",
	2: "virtual",
	3: "formation lap",
}
var NetworkGame = map[int16]string{0: "offline",
	1: "online",
}
var BrakingAssist = map[int16]string{0: "off",
	1: "low",
	2: "medium",
	3: "high",
}
var GearboxAssist = map[int16]string{1: "manual",
	2: "manual & suggested gear",
	3: "auto",
}
var DynamicRacingLine = map[int16]string{0: "off",
	1: "corners only",
	2: "full",
}
var DynamicRacingLineType = map[int16]string{0: "2D",
	1: "3D",
}

var SessionLength = map[int16]string{
	0: "None",
	2: "Very Short",
	3: "Short",
	4: "Medium",
	5: "Medium Long",
	6: "Long",
	7: "Full",
}

var SpeedUnits = map[int16]string{
	0: "MPH",
	1: "KPH",
}

var TemperatureUnits = map[int16]string{
	0: "Celsius",
	1: "Fahrenheit",
}

var RecoveryMode = map[int16]string{
	0: "None",
	1: "Flashbacks",
	2: "Auto-recovery",
}

var FlashbackLimit = map[int16]string{
	0: "Low",
	1: "Medium",
	2: "High",
	3: "Unlimited",
}

var SurfaceType = map[int16]string{
	0: "Simplified",
	1: "Realistic",
}

var LowFuelMode = map[int16]string{
	0: "Easy",
	1: "Hard",
}

var RaceStarts = map[int16]string{
	0: "Manual",
	1: "Assisted",
}

var TyreTemperature = map[int16]string{
	0: "Surface only",
	1: "Surface & Carcass",
}

var CarDamage = map[int16]string{
	0: "Off",
	1: "Reduced",
	2: "Standard",
	3: "Simulation",
}

var CarDamageRate = map[int16]string{
	0: "Reduced",
	1: "Standard",
	2: "Simulation",
}

var Collisions = map[int16]string{
	0: "Off",
	1: "Player-to-Player Off",
	2: "On",
}

var CornerCuttingStringency = map[int16]string{
	0: "Regular",
	1: "Strict",
}

var PitStopExperience = map[int16]string{
	0: "Automatic",
	1: "Broadcast",
	2: "Immersive",
}

var SafetyCar = map[int16]string{
	0: "Off",
	1: "Reduced",
	2: "Standard",
	3: "Increased",
}

var SafetyCarExperience = map[int16]string{
	0: "Broadcast",
	1: "Immersive",
}

var FormationLapExperience = map[int16]string{
	0: "Broadcast",
	1: "Immersive",
}

var RedFlags = map[int16]string{
	0: "Off",
	1: "Reduced",
	2: "Standard",
	3: "Increased",
}
