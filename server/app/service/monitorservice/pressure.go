package monitor

import (
	"context"
	"dubmer-bono/app/types/entity"
	"fmt"
	"math"
	"time"
)

func (s *Service) monitorPressure(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
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
			pressure, err := s.repo.OLAP.GetPressureFactors(ctx, sessionid, uint8(carno))
			if err != nil {
				fmt.Println("pressure query error:", err)
				continue
			}
			s.AlterConfidence(pressure)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) AlterConfidence(pressure entity.PilotPressurePhysicalFactors) {
	// Normalizing all the units between 0 and 1
	steer := math.Abs(float64(pressure.Steer))                 // Direction doesnt matter here
	lat := min(math.Abs(float64(pressure.GForceLat))/5.0, 1.0) // F1 ~5-6G
	lon := min(math.Abs(float64(pressure.GFroceLon))/5.0, 1.0)
	brake := pressure.Brake // Already normalized
	gForce := (lat*lat + lon*lon) / 2
	score :=
		0.45*gForce +
			0.25*steer +
			0.30*float64(brake)

	//Pressure was under 0.01 for straights and went upto ~0.5 on crazier corners
	s.driver_pressure = int(math.Floor(min(score*10, 5)))

	//Now the possible values are 0 1 2 3 4 5
	//And we can signal based on them
	fmt.Printf("%f : %d\n", pressure.LapDistance, s.driver_pressure)
}
