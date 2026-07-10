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
	tele          telentity.PacketCarTelemetryData
	tele_update   time.Time
	signal_func   func(context.Context)
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

func (t *Service) SignalPush(ctx context.Context) {
	if t.acc.motion == (telentity.MotionPacket{}) || t.acc.tele == (telentity.PacketCarTelemetryData{}) {
		return // haven't received both packet types yet
	}

	avg, err := utility.CheckFrameGap(t.acc.motion_update, t.acc.tele_update, time.Millisecond*50)
	if err != nil {
		// motion and telemetry timestamps drifted too far apart, don't trust this pairing
		return
	}

	//Reset
	t.acc.motion = telentity.MotionPacket{}
	t.acc.tele = telentity.PacketCarTelemetryData{}
	var sessionID uint32
	if err := t.repo.Cache.Get(string(entity.SESSIONID), string(entity.SESSIONID), &sessionID); err != nil {
		return
	}

	for carIndex := range 20 {
		motion := t.acc.motion.Cars[carIndex]
		tele := t.acc.tele.CarTelemetryData[carIndex]

		t.repo.OLAP.InsertFrame(ctx, entity.TelemetryFrame{
			SessionID: sessionID,
			CarNo:     uint8(carIndex),
			FrameTime: avg,

			Speed:     float32(tele.Speed),
			Throttle:  tele.Throttle,
			Steer:     tele.Steer,
			Brake:     tele.Brake,
			Clutch:    float32(tele.Clutch) / 100.0, // 0-100 -> 0.0-1.0 to match entity
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
		})
	}
}
