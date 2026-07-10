package monitor

import (
	"context"
	"dubmer-bono/app/types/entity"
	"fmt"
	"time"
)

func (s *Service) monitorPressure(ctx context.Context) {
	ticker := time.NewTicker(time.Second / 4)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var sessionid uint64
			var carno int
			if err := s.repo.Cache.Get(string(entity.GAMESESSION), string(entity.SESSIONID), &sessionid); err != nil {
				fmt.Println("session cache miss:", err)
				continue
			}
			if err := s.repo.Cache.Get(string(entity.GAMESESSION), string(entity.MYCARID), &carno); err != nil {
				fmt.Println("car id cache miss:", err)
				continue
			}
			fmt.Printf("querying with session=%d car=%d\n", sessionid, carno)
			pressure, err := s.repo.OLAP.GetPressureFactors(ctx, sessionid, uint8(carno))
			if err != nil {
				fmt.Println("pressure query error:", err)
				continue
			}
			fmt.Printf("%+v\n", pressure)
		case <-ctx.Done():
			return
		}
	}
}
