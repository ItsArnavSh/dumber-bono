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
		t.logger.Debug("SignalPush: waiting on motion or telemetry packet, skipping")
		return
	}

	avg, err := utility.CheckFrameGap(t.acc.motion_update, t.acc.tele_update, time.Millisecond*50)
	if err != nil {
		t.logger.Debugf("SignalPush: frame gap too large, dropping pair: %v", err)
		return
	}

	var sessionID uint64
	if err := t.repo.Cache.Get(string(entity.GAMESESSION), string(entity.SESSIONID), &sessionID); err != nil {
		t.logger.Debugf("SignalPush: session id not found in cache, skipping: %v", err)
		return
	}

	// snapshot before reset, since reading t.acc after zeroing it would give you empty data
	motionSnapshot := t.acc.motion
	teleSnapshot := t.acc.tele

	t.acc.motion = telentity.MotionPacket{}
	t.acc.tele = telentity.PacketCarTelemetryData{}

	t.logger.Debugf("SignalPush: writing 20 car frames, session=%d frameTime=%s (UTC)",
		sessionID, avg.UTC().Format(time.RFC3339Nano))

	for carIndex := range 20 {
		motion := motionSnapshot.Cars[carIndex]
		tele := teleSnapshot.CarTelemetryData[carIndex]

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
			PosX:      motion.WorldPosition.X,
			PosY:      motion.WorldPosition.Y,
			PosZ:      motion.WorldPosition.Z,
			VelX:      motion.WorldVelocity.X,
			VelY:      motion.WorldVelocity.Y,
			VelZ:      motion.WorldVelocity.Z,
			FwdX:      float32(motion.WorldForwardDir.X),
			FwdY:      float32(motion.WorldForwardDir.Y),
			FwdZ:      float32(motion.WorldForwardDir.Z),
			GForceLat: motion.GForce.Lateral,
			GForceLon: motion.GForce.Longitudinal,
			Yaw:       motion.Orientation.Yaw,
			Pitch:     motion.Orientation.Pitch,
			Roll:      motion.Orientation.Roll,
		})
	}
}
