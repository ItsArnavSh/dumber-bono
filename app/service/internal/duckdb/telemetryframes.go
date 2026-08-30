package duckdb

import (
	"context"
	"database/sql/driver"
	"dubmer-bono/app/types/entity"

	"github.com/marcboeker/go-duckdb"
)

func (d *DuckDB) InsertFrame(ctx context.Context, f entity.TelemetryFrame) {
	d.batch.mu.Lock()
	d.batch.telemetry = append(d.batch.telemetry, f)
	d.batch.mu.Unlock()
}
func (d *DuckDB) writeTelemetrysBatch(ctx context.Context, frames []entity.TelemetryFrame) error {
	conn, err := d.conn.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	return conn.Raw(func(driverConn any) error {
		appender, err := duckdb.NewAppenderFromConn(driverConn.(driver.Conn), "", "telemetry_frames")
		if err != nil {
			return err
		}
		defer func() {
			_ = appender.Close()
		}()

		for _, f := range frames {
			if err := appender.AppendRow(
				f.SessionID, f.CarNo, f.FrameTime,
				f.Speed, f.Throttle, f.Steer, f.Brake, f.Clutch, f.Gear, f.EngineRPM, f.DRS,
				f.PosX, f.PosY, f.PosZ,
				f.VelX, f.VelY, f.VelZ,
				f.FwdX, f.FwdY, f.FwdZ,
				f.GForceLat, f.GForceLon,
				f.Yaw, f.Pitch, f.Roll,
				f.CarPosition, f.DeltaToFrontMS, f.DeltaToLeaderMS, f.LapDistance,
			); err != nil {
				return err
			}
		}
		return nil
	})
}
