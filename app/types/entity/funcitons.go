package entity

type SerFunc string

const (
	SerFuncPosition      SerFunc = "position"
	SerFuncGapToFront    SerFunc = "gap_to_front"
	SerFuncGapToLeader   SerFunc = "gap_to_leader"
	SerFuncCurrentLap    SerFunc = "current_lap"
	SerFuncLastLapTime   SerFunc = "last_lap_time"
	SerFuncTotalWarnings SerFunc = "total_warnings"
)
