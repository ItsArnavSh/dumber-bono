package duckdb

import (
	"context"
	"dubmer-bono/app/types/entity"
	"fmt"
	"sync"
	"time"
)

type BatchedFrames struct {
	telemetry []entity.TelemetryFrame
	mu        sync.Mutex
}

func (d *DuckDB) BatchProcess(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			frames := d.batch.Swap()
			if len(frames) == 0 {
				continue
			}
			if err := d.writeTelemetrysBatch(ctx, frames); err != nil {
				fmt.Println(err)
				// batch is lost here
			}

		case <-ctx.Done():
			frames := d.batch.Swap()
			if len(frames) > 0 {
				_ = d.writeTelemetrysBatch(context.Background(), frames)
			}
			return
		}
	}
}
func (b *BatchedFrames) Swap() []entity.TelemetryFrame {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := b.telemetry
	b.telemetry = make([]entity.TelemetryFrame, 0, cap(out))
	return out
}
