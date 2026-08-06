package duckdb

import (
	"context"
	"dubmer-bono/app/types/entity"
	"strconv"
)

func (d *DuckDB) GetPressureFactors(ctx context.Context, sessionID uint64, carNo uint8) (entity.PilotPressurePhysicalFactors, error) {
	const query = `
		SELECT g_force_lat,g_force_lon, steer, brake,lap_distance
		FROM telemetry_frames
		WHERE session_id = CAST(? AS UBIGINT) AND car_no = ?
		ORDER BY frame_time DESC
		LIMIT 1
	`

	var pf entity.PilotPressurePhysicalFactors
	err := d.conn.QueryRowContext(ctx, query,
		strconv.FormatUint(sessionID, 10), carNo,
	).Scan(&pf.GForceLat, &pf.GFroceLon, &pf.Steer, &pf.Brake, &pf.LapDistance)
	if err != nil {
		return entity.PilotPressurePhysicalFactors{}, err
	}

	return pf, nil
}
