package ingestion

import (
	"context"
	"dubmer-bono/app/api/udp/parsers"
	"dubmer-bono/app/service"
	"dubmer-bono/app/types"
	"dubmer-bono/app/types/entity"
	telentity "dubmer-bono/app/types/entity/tel-entity"
	"strconv"

	"go.uber.org/zap"
)

//The Ingestion Service Saves the Data in the relevant DBs so that the other processes can query into it

type Service struct {
	repo   *service.Repository
	logger *zap.SugaredLogger
	acc    TeleAccumulator
}

var _ types.Ingestion = &Service{}

func NewService(ctx context.Context, logger *zap.SugaredLogger, root string, repo *service.Repository) (types.Ingestion, error) {
	serv := &Service{repo: repo, logger: logger}
	serv.acc = TeleAccumulator{
		signal_func: serv.SignalPush,
	}
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

func (s *Service) IngestSessionPacket(payload telentity.PacketSessionData) {}
func (s *Service) IngestLapPacket(payload telentity.LapDataPacket) {
	ctx := context.Background()
	s.acc.UpsertLapData(ctx, payload)
	kv := make(map[string]any, len(payload.LapData))

	for i := range payload.LapData {
		kv[strconv.Itoa(i)] = payload.LapData[i]
	}

	if err := s.repo.Cache.BulkSet(string(entity.LAPDATA), kv); err != nil {
		s.logger.Error(err.Error())
		return
	}
}
func (s *Service) IngestEventPacket(payload telentity.PacketEventData) {}
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
	for i := range payload.NumActiveCars {
		s.repo.Cache.Set(string(entity.PARTICIPANT), string(i), payload.Participants[i])
	}
}
func (s *Service) IngestCarSetupPacket(payload telentity.PacketCarSetupData) {

}
func (s *Service) IngestMotionPacket(payload telentity.MotionPacket) {
	ctx := context.Background()
	s.acc.UpsertMotion(ctx, payload)
}
func (s *Service) IngestTelemetryPacket(payload telentity.PacketCarTelemetryData) {
	ctx := context.Background()
	s.acc.UpsertTelemetry(ctx, payload)
}
func (s *Service) IngestCarStatusPacket(payload telentity.PacketCarStatusData) {}
func (s *Service) IngestLobbyInfoPacket(payload telentity.PacketLobbyInfoData) {}
func (s *Service) IngestCarDamagePacket(payload telentity.PacketCarDamageData) {
	//We will save this data for
}
func (s *Service) IngestSessionHistoryPacket(payload telentity.SessionHistoryPacket) {}
func (s *Service) IngestTyreSetPacket(payload telentity.TyreSetPacket)               {}
