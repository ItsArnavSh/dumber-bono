package telentity

import "dubmer-bono/app/types/entity"

type CarDamageData struct {
	TyresWear            entity.Tyres[float32] // Tyre wear (percentage)
	TyresDamage          entity.Tyres[uint8]   // Tyre damage (percentage)
	BrakesDamage         entity.Tyres[uint8]   // Brakes damage (percentage)
	TyreBlisters         entity.Tyres[uint8]   // Tyre blisters value (percentage)
	FrontLeftWingDamage  uint8                 // Front left wing damage (percentage)
	FrontRightWingDamage uint8                 // Front right wing damage (percentage)
	RearWingDamage       uint8                 // Rear wing damage (percentage)
	FloorDamage          uint8                 // Floor damage (percentage)
	DiffuserDamage       uint8                 // Diffuser damage (percentage)
	SidepodDamage        uint8                 // Sidepod damage (percentage)
	DRSFault             bool                  // Indicator for DRS fault
	ERSFault             bool                  // Indicator for ERS fault
	GearBoxDamage        uint8                 // Gear box damage (percentage)
	EngineDamage         uint8                 // Engine damage (percentage)
	EngineMGUHWear       uint8                 // Engine wear MGU-H (percentage)
	EngineESWear         uint8                 // Engine wear ES (percentage)
	EngineCEWear         uint8                 // Engine wear CE (percentage)
	EngineICEWear        uint8                 // Engine wear ICE (percentage)
	EngineMGUKWear       uint8                 // Engine wear MGU-K (percentage)
	EngineTCWear         uint8                 // Engine wear TC (percentage)
	EngineBlown          bool                  // Engine blown
	EngineSeized         bool                  // Engine seized
}

type PacketCarDamageData struct {
	Header        UDPHeader
	CarDamageData [22]CarDamageData
}
