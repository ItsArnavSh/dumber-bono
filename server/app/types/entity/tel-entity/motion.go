package telentity

import "dubmer-bono/app/types/entity"

/*
 The motion packet gives physics data for all the cars being driven.
 N.B. For the normalised vectors below, to convert to float values divide by 32767.0f – 16-bit signed
 values are used to pack the data and on the assumption that direction values are always between -1.0f
 and 1.0f.
 Frequency: Rate as specified in menus
 Size: 1349 bytes
 Version: 1
*/

type CarMotion struct {
	WorldPosition   entity.Coordinates[float32] //World space position in meters
	WorldVelocity   entity.Coordinates[float32] //Velocity in that direction in m/s
	WorldForwardDir entity.Coordinates[int16]   //Worls Space Froward Direction Normalized
	WorldRightDir   entity.Coordinates[int16]
	GForce          entity.LatLon
	Orientation     entity.Orientation //Drone Stuff
}

type MotionPacket struct {
	header UDPHeader
	cars   [22]CarMotion //Data for all the cars
}
