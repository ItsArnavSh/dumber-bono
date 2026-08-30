package ingestion

import (
	"context"
	"dubmer-bono/app/types/entity"
	telentity "dubmer-bono/app/types/entity/tel-entity"
	"testing"
	"time"
)

func TestLapTimestampToMS(t *testing.T) {
	tests := []struct {
		name string
		ts   entity.LapTimeStamp
		want uint32
	}{
		{"zero", entity.LapTimeStamp{}, 0},
		{"only ms", entity.LapTimeStamp{MS: 500}, 500},
		{"one minute", entity.LapTimeStamp{Minute: 1, MS: 0}, 60000},
		{"mixed", entity.LapTimeStamp{Minute: 1, MS: 25000}, 85000},
		{"two minutes", entity.LapTimeStamp{Minute: 2, MS: 999}, 120999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lapTimestampToMS(tt.ts); got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTeleAccumulatorUpsert(t *testing.T) {
	ctx := context.Background()
	acc := TeleAccumulator{signal_func: func(ctx context.Context) {}}

	acc.UpsertMotion(ctx, telentity.MotionPacket{})
	acc.UpsertTelemetry(ctx, telentity.PacketCarTelemetryData{})
	acc.UpsertLapData(ctx, telentity.LapDataPacket{})

	acc.mu.RLock()
	defer acc.mu.RUnlock()
	if acc.motion_update.IsZero() {
		t.Error("motion_update not set")
	}
	if acc.tele_update.IsZero() {
		t.Error("tele_update not set")
	}
	if acc.lap_update.IsZero() {
		t.Error("lap_update not set")
	}
}

func TestTeleAccumulatorUpsertPreservesPacket(t *testing.T) {
	ctx := context.Background()
	acc := TeleAccumulator{signal_func: func(ctx context.Context) {}}

	packet := telentity.PacketCarTelemetryData{}
	packet.CarTelemetryData[0].Speed = 300
	acc.UpsertTelemetry(ctx, packet)

	acc.mu.RLock()
	defer acc.mu.RUnlock()
	if acc.tele.CarTelemetryData[0].Speed != 300 {
		t.Errorf("Speed = %d, want 300", acc.tele.CarTelemetryData[0].Speed)
	}
}

func TestSignalPushSkipsWhenIncomplete(t *testing.T) {
	ctx := context.Background()
	s := &Service{
		acc: TeleAccumulator{
			signal_func: func(ctx context.Context) {},
		},
	}
	// Only motion and telemetry set, lap missing -> SignalPush should no-op
	s.acc.UpsertMotion(ctx, telentity.MotionPacket{})
	s.acc.UpsertTelemetry(ctx, telentity.PacketCarTelemetryData{})

	s.SignalPush(ctx)
	// No panic and no state change expected; just verify it returned
}

func TestSignalPushSkipsWhenStale(t *testing.T) {
	ctx := context.Background()
	s := &Service{
		acc: TeleAccumulator{
			signal_func: func(ctx context.Context) {},
		},
	}
	s.acc.UpsertMotion(ctx, telentity.MotionPacket{})
	s.acc.UpsertTelemetry(ctx, telentity.PacketCarTelemetryData{})
	s.acc.UpsertLapData(ctx, telentity.LapDataPacket{})

	// Simulate stale motion frame that exceeds tolerance
	s.acc.mu.Lock()
	s.acc.motion_update = s.acc.motion_update.Add(-time.Hour)
	s.acc.mu.Unlock()

	s.SignalPush(ctx)
	// Should no-op due to CheckFramesFresh failing; ensure no panic
}
