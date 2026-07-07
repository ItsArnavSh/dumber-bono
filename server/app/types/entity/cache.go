package entity

type NameSpace string

const (
	GAMESESSION NameSpace = "gamesession"
	PARTICIPANT NameSpace = "participant"
	LAPDATA     NameSpace = "lapdata"
)

type Key string

const (
	PLAYERINDEX Key = "playerindex"
	SESSIONID   Key = "sessionid"
	MYCARID     Key = "mycarid"
	MYDRIVERID  Key = "mydriverid"
	MYTEAMID    Key = "myteamid"
	MYRACEINDEX Key = "myraceindex"
)
