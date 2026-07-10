package ingestion

import (
	"context"
	"dubmer-bono/app/types/entity"
	telentity "dubmer-bono/app/types/entity/tel-entity"
	"dubmer-bono/app/utility"
	"time"
)

type TeleAccumulator struct {
	motion        telentity.MotionPacket
	motion_update time.Time

	tele        telentity.PacketCarTelemetryData
	tele_update time.Time

	lap        telentity.LapDataPacket
	lap_update time.Time

	signal_func func(context.Context)
}

func (t *TeleAccumulator) UpsertMotion(ctx context.Context, packet telentity.MotionPacket) {
	t.motion = packet
	t.motion_update = time.Now()
	t.signal_func(ctx)
}

func (t *TeleAccumulator) UpsertTelemetry(ctx context.Context, packet telentity.PacketCarTelemetryData) {
	t.tele = packet
	t.tele_update = time.Now()
}

func (t *TeleAccumulator) UpsertLapData(ctx context.Context, packet telentity.LapDataPacket) {
	t.lap = packet
	t.lap_update = time.Now()
}

func (t *Service) SignalPush(ctx context.Context) {
	if t.acc.motion == (telentity.MotionPacket{}) ||
		t.acc.tele == (telentity.PacketCarTelemetryData{}) ||
		t.acc.lap == (telentity.LapDataPacket{}) {
		return // haven't received all three packet types yet
	}

	avg, err := utility.CheckFramesFresh(time.Millisecond*50,
		t.acc.motion_update, t.acc.tele_update, t.acc.lap_update)
	if err != nil {
		// one of the three streams drifted too far from the others, don't trust this pairing
		return
	}

	var sessionID uint64
	if err := t.repo.Cache.Get(string(entity.GAMESESSION), string(entity.SESSIONID), &sessionID); err != nil {
		return
	}

	// snapshot before reset
	motionSnapshot := t.acc.motion
	teleSnapshot := t.acc.tele
	lapSnapshot := t.acc.lap

	t.acc.motion = telentity.MotionPacket{}
	t.acc.tele = telentity.PacketCarTelemetryData{}
	t.acc.lap = telentity.LapDataPacket{}

	for carIndex := range 20 {
		motion := motionSnapshot.Cars[carIndex]
		tele := teleSnapshot.CarTelemetryData[carIndex]
		lap := lapSnapshot.LapData[carIndex]

		t.repo.OLAP.InsertFrame(ctx, entity.TelemetryFrame{
			SessionID: sessionID,
			CarNo:     uint8(carIndex),
			FrameTime: avg,

			Speed:     float32(tele.Speed),
			Throttle:  tele.Throttle,
			Steer:     tele.Steer,
			Brake:     tele.Brake,
			Clutch:    float32(tele.Clutch) / 100.0,
			Gear:      tele.Gear,
			EngineRPM: tele.EngineRPM,
			DRS:       tele.DRS,

			PosX: motion.WorldPosition.X,
			PosY: motion.WorldPosition.Y,
			PosZ: motion.WorldPosition.Z,

			VelX: motion.WorldVelocity.X,
			VelY: motion.WorldVelocity.Y,
			VelZ: motion.WorldVelocity.Z,

			FwdX: float32(motion.WorldForwardDir.X),
			FwdY: float32(motion.WorldForwardDir.Y),
			FwdZ: float32(motion.WorldForwardDir.Z),

			GForceLat: motion.GForce.Lateral,
			GForceLon: motion.GForce.Longitudinal,

			Yaw:   motion.Orientation.Yaw,
			Pitch: motion.Orientation.Pitch,
			Roll:  motion.Orientation.Roll,

			CarPosition:     lap.CarPosition,
			DeltaToFrontMS:  lapTimestampToMS(lap.DeltaToFront),
			DeltaToLeaderMS: lapTimestampToMS(lap.DeltaToRaceLeader),
			LapDistance:     lap.LapDistance,
		})
	}
}

func lapTimestampToMS(ts entity.LapTimeStamp) uint32 {
	return uint32(ts.Minute)*60000 + uint32(ts.MS)
}
