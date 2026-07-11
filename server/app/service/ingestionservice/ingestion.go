package ingestion

import (
	"context"
	"dubmer-bono/app/api/udp/parsers"
	"dubmer-bono/app/service"
	"dubmer-bono/app/types"
	"dubmer-bono/app/types/entity"
	telentity "dubmer-bono/app/types/entity/tel-entity"
	"dubmer-bono/app/utility"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// The Ingestion Service Saves the Data in the relevant DBs so that the other processes can query into it
type Service struct {
	repo   *service.Repository
	logger *zap.SugaredLogger
	acc    TeleAccumulator

	lapThrottler         utility.Throttler
	sessionThrottler     utility.Throttler
	eventThrottler       utility.Throttler
	participantThrottler utility.Throttler
	carSetupThrottler    utility.Throttler
	carStatusThrottler   utility.Throttler
	lobbyInfoThrottler   utility.Throttler
	carDamageThrottler   utility.Throttler
	sessionHistThrottler utility.Throttler
	tyreSetThrottler     utility.Throttler
}

var _ types.Ingestion = &Service{}

func NewService(ctx context.Context, logger *zap.SugaredLogger, root string, repo *service.Repository) (types.Ingestion, error) {
	serv := &Service{repo: repo, logger: logger}
	serv.acc = TeleAccumulator{
		signal_func: serv.SignalPush,
	}

	interval := time.Second * 1
	serv.lapThrottler = utility.Throttler{Interval: interval}
	serv.sessionThrottler = utility.Throttler{Interval: interval}
	serv.eventThrottler = utility.Throttler{Interval: interval}
	serv.participantThrottler = utility.Throttler{Interval: interval}
	serv.carSetupThrottler = utility.Throttler{Interval: interval}
	serv.carStatusThrottler = utility.Throttler{Interval: interval}
	serv.lobbyInfoThrottler = utility.Throttler{Interval: interval}
	serv.carDamageThrottler = utility.Throttler{Interval: interval}
	serv.sessionHistThrottler = utility.Throttler{Interval: interval}
	serv.tyreSetThrottler = utility.Throttler{Interval: interval}

	return serv, nil
}

func (s *Service) IngestHeader(payload *parsers.PacketHeader) error {
	var sessionID uint64
	err := s.repo.Cache.Get(string(entity.GAMESESSION), string(entity.SESSIONID), &sessionID)
	if err != nil {
		s.logger.Warnf("no cached session found: %v", err)
		// no session cached yet, treat this packet's session as new
		sessionID = 0
	}
	if payload.SessionUID != sessionID {
		// TODO: Nuke Current Session
	}
	if err := s.repo.Cache.Set(string(entity.GAMESESSION), string(entity.SESSIONID), payload.SessionUID); err != nil {
		s.logger.Errorf("failed to set session id in cache: %v", err)
		return err
	}
	if err := s.repo.Cache.Set(string(entity.GAMESESSION), string(entity.PLAYERINDEX), payload.PlayerCarIndex); err != nil {
		s.logger.Errorf("failed to set player index in cache: %v", err)
		return err
	}
	return nil
}

func (s *Service) IngestSessionPacket(payload telentity.PacketSessionData) {
	if !s.sessionThrottler.Allow() {
		return
	}
	if err := s.repo.Cache.Set(string(entity.SESSIONDATA), string(entity.SESSIONDATA), payload); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *Service) IngestLapPacket(payload telentity.LapDataPacket) {
	ctx := context.Background()
	s.acc.UpsertLapData(ctx, payload)

	if !s.lapThrottler.Allow() {
		return
	}
	kv := make(map[string]any, len(payload.LapData))
	for i := range payload.LapData {
		kv[strconv.Itoa(i)] = payload.LapData[i]
	}
	if err := s.repo.Cache.BulkSet(string(entity.LAPDATA), kv); err != nil {
		s.logger.Error(err.Error())
		return
	}
}

func (s *Service) IngestEventPacket(payload telentity.PacketEventData) {

}

func (s *Service) IngestParticipantPacket(payload telentity.PacketParticipantsData) {
	me, idx, err := payload.Me()
	if err != nil {
		// log and bail - can't cache anything meaningful without a valid player index
		return
	}
	s.repo.Cache.Set(string(entity.GAMESESSION), string(entity.MYCARID), idx)
	s.repo.Cache.Set(string(entity.GAMESESSION), string(entity.MYDRIVERID), me.DriverId)
	s.repo.Cache.Set(string(entity.GAMESESSION), string(entity.MYTEAMID), me.TeamId)
	s.repo.Cache.Set(string(entity.GAMESESSION), string(entity.MYRACEINDEX), me.RaceNumber)

	if !s.participantThrottler.Allow() {
		return
	}
	for i := range payload.NumActiveCars {
		s.repo.Cache.Set(string(entity.PARTICIPANT), string(i), payload.Participants[i])
	}
}

func (s *Service) IngestCarSetupPacket(payload telentity.PacketCarSetupData) {
	if !s.carSetupThrottler.Allow() {
		return
	}
	kv := make(map[string]any, len(payload.CarSetups))
	for i := range payload.CarSetups {
		kv[strconv.Itoa(i)] = payload.CarSetups[i]
	}
	if err := s.repo.Cache.BulkSet(string(entity.CARSETUP), kv); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *Service) IngestMotionPacket(payload telentity.MotionPacket) {
	ctx := context.Background()
	s.acc.UpsertMotion(ctx, payload)
}

func (s *Service) IngestTelemetryPacket(payload telentity.PacketCarTelemetryData) {
	ctx := context.Background()
	s.acc.UpsertTelemetry(ctx, payload)
}

func (s *Service) IngestCarStatusPacket(payload telentity.PacketCarStatusData) {
	if !s.carStatusThrottler.Allow() {
		return
	}
	kv := make(map[string]any, len(payload.CarStatusData))
	for i := range payload.CarStatusData {
		kv[strconv.Itoa(i)] = payload.CarStatusData[i]
	}
	if err := s.repo.Cache.BulkSet(string(entity.CARSTATUS), kv); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *Service) IngestLobbyInfoPacket(payload telentity.PacketLobbyInfoData) {
	if !s.lobbyInfoThrottler.Allow() {
		return
	}
	kv := make(map[string]any, payload.NumPlayers)
	for i := 0; i < int(payload.NumPlayers); i++ {
		kv[strconv.Itoa(i)] = payload.LobbyPlayers[i]
	}
	if err := s.repo.Cache.BulkSet(string(entity.LOBBYINFO), kv); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *Service) IngestCarDamagePacket(payload telentity.PacketCarDamageData) {
	//We will save this data for
	if !s.carDamageThrottler.Allow() {
		return
	}
	kv := make(map[string]any, len(payload.CarDamageData))
	for i := range payload.CarDamageData {
		kv[strconv.Itoa(i)] = payload.CarDamageData[i]
	}
	if err := s.repo.Cache.BulkSet(string(entity.CARDAMAGE), kv); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *Service) IngestSessionHistoryPacket(payload telentity.SessionHistoryPacket) {
	if !s.sessionHistThrottler.Allow() {
		return
	}
	if err := s.repo.Cache.Set(string(entity.SESSIONHISTORY), string(entity.SESSIONHISTORY), payload); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *Service) IngestTyreSetPacket(payload telentity.TyreSetPacket) {
	if !s.tyreSetThrottler.Allow() {
		return
	}
	if err := s.repo.Cache.Set(string(entity.TYRESET), string(entity.TYRESET), payload); err != nil {
		s.logger.Error(err.Error())
	}
}
