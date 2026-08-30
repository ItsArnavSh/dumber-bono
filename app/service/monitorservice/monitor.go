package monitor

import (
	"context"
	"dubmer-bono/app/service"
	"dubmer-bono/app/types"
	"dubmer-bono/app/types/entity"
	"dubmer-bono/app/types/entity/consts"
	telentity "dubmer-bono/app/types/entity/tel-entity"
	"fmt"
	"strconv"
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
	s.shuffler = shuffler
	s.RegisterShufflerFunctions()
	return s, nil
}
func (s *Service) RegisterShufflerFunctions() {
	s.shuffler.Register(string(entity.SerFuncGapToFront), s.SendInformOfGapToFront)
	s.shuffler.Register(string(entity.SerFuncGapToLeader), s.SendInformOfGapToLeader)
	s.shuffler.Register(string(entity.SerFuncLastLapTime), s.SendInformOfLastLapTime)
	s.shuffler.Register(string(entity.SerFuncTotalWarnings), s.SendInformOfTotalWarnings)
}

func (s *Service) GetPressure() int {
	return s.driver_pressure
}

func (s *Service) PushToRadio(message entity.RadioPayload, priority int, expire_after time.Duration) {
	s.msg_chan <- entity.RadioMessage{
		Message:  message,
		Priority: priority,
		Expiry:   time.Now().Add(expire_after),
	}

}

// IsGameActive reports whether the game session was updated within the last 5 seconds.
func (s *Service) IsGameActive() bool {
	var latest_update time.Time
	if err := s.repo.Cache.Get(string(entity.GAMESESSION), string(entity.LASTUPDATED), &latest_update); err != nil {
		return false
	}
	return time.Since(latest_update) <= 5*time.Second
}

func (s *Service) RandomStatsMonitor(ctx context.Context) {
	wait := time.Second * 30

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		session, _ := s.GetSessionData()
		if s.IsGameActive() && session.TotalLaps >= 1 {
			s.shuffler.RunSubset(3) // run N, weighted toward ones not run recently
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
	}
}

type PlayerDetails struct {
	Driver_name, Team_name, Player_nation string
}

func (s *Service) getPlayerDetailsByID(id int) PlayerDetails {
	fmt.Println("called by carno ", id)
	var data telentity.ParticipantData
	if err := s.repo.Cache.Get(string(entity.PARTICIPANT), strconv.Itoa(id), &data); err != nil {
		s.logger.Warnf("no cached participant for car %d: %v", id, err)
	}
	fmt.Println(data)
	team_name := consts.TeamIDs[uint16(data.TeamId)]
	driver_name := consts.DriverIDs[uint16(data.DriverId)]
	nationality := consts.NationalityIDs[uint16(data.Nationality)]

	return PlayerDetails{
		Team_name:     team_name,
		Driver_name:   driver_name,
		Player_nation: nationality,
	}
}
func (s *Service) getCarPosition(id uint8) int {
	var data telentity.LapData
	if err := s.repo.Cache.Get(string(entity.LAPDATA), strconv.Itoa(int(id)), &data); err != nil {
		s.logger.Errorf("failed to get lap data for car %d: %v", id, err)
		return 0
	}
	return int(data.CarPosition)
}
