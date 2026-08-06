package entity

type NameSpace string

const (
	GAMESESSION    NameSpace = "gamesession"
	PARTICIPANT    NameSpace = "participant"
	LAPDATA        NameSpace = "lapdata"
	SESSIONDATA    NameSpace = "sessiondata"
	CARSETUP       NameSpace = "carsetup"
	CARSTATUS      NameSpace = "carstatus"
	LOBBYINFO      NameSpace = "lobbyinfo"
	CARDAMAGE      NameSpace = "cardamage"
	SESSIONHISTORY NameSpace = "lobbyinfo"
	TYRESET        NameSpace = "tyreset"
)

type Key string

const (
	PLAYERINDEX Key = "playerindex"
	SESSIONID   Key = "sessionid"
	MYCARID     Key = "mycarid"
	MYDRIVERID  Key = "mydriverid"
	MYTEAMID    Key = "myteamid"
	MYRACEINDEX Key = "myraceindex"
	LASTUPDATED Key = "lastupdated"
)
