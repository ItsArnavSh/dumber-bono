package monitor

import (
	"context"
	"dubmer-bono/app/service"
	"dubmer-bono/app/types"
	"dubmer-bono/app/types/entity"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

/*
 * The monitor service is concerned with primarily listening for events
 *  And based on that, generating text, and pusing to the Radio pipeline
 */

type Service struct {
	repo            *service.Repository
	logger          *zap.SugaredLogger
	driver_pressure int
	shuffler        *Shuffler
	msg_chan        chan<- entity.RadioMessage
}

var _ types.Monitor = &Service{}

func NewService(ctx context.Context, logger *zap.SugaredLogger, root string, repo *service.Repository, msg_chan chan<- entity.RadioMessage) (types.Monitor, error) {
	s := &Service{repo: repo, logger: logger, msg_chan: msg_chan}
	go s.monitorPressure(ctx)
	go s.RandomStatsMonitor(ctx)
	shuffler := NewShuffler()
	shuffler.Register("position", s.InformOfPosition)
	shuffler.Register("gap_to_front", s.InformOfGapToFront)
	shuffler.Register("gap_to_leader", s.InformOfGapToLeader)
	shuffler.Register("current_lap", s.InformOfCurrentLap)
	shuffler.Register("last_lap_time", s.InformOfLastLapTime)
	shuffler.Register("pit_status", s.InformOfPitStatus)
	shuffler.Register("num_pit_stops", s.InformOfNumPitStops)
	shuffler.Register("sector", s.InformOfSector)
	shuffler.Register("grid_position", s.InformOfGridPosition)
	shuffler.Register("speed_trap", s.InformOfSpeedTrap)
	shuffler.Register("total_warnings", s.InformOfTotalWarnings)
	s.shuffler = shuffler
	return s, nil
}

func (s *Service) GetPressure() int {
	return s.driver_pressure
}

func (s *Service) PushToRadio(message string, priority int, expire_after time.Duration) {
	s.msg_chan <- entity.RadioMessage{
		Message:  message,
		Priority: priority,
		Expiry:   time.Now().Add(expire_after),
	}

}

func (s *Service) RandomStatsMonitor(ctx context.Context) {
	intervals := []time.Duration{0, time.Second * 5, time.Second * 15, time.Second * 30}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		wait := intervals[rand.Intn(len(intervals))]
		if wait == 0 {
			continue
		}
		s.shuffler.RunSubset(3) // run 3 of the 11, weighted toward ones not run recently
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
	}
}
