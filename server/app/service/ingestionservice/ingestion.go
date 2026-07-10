package ingestion

import (
	"context"
	"dubmer-bono/app/api/udp/parsers"
	"dubmer-bono/app/service/internal"
	"dubmer-bono/app/types"
	"dubmer-bono/app/types/entity"
	telentity "dubmer-bono/app/types/entity/tel-entity"
	"strconv"

	"go.uber.org/zap"
)

//The Ingestion Service Saves the Data in the relevant DBs so that the other processes can query into it

type Service struct {
	repo   *internal.Repository
	logger *zap.Logger
	acc    TeleAccumulator
}

var _ types.Ingestion = &Service{}

func NewService(ctx context.Context, root string) (types.Ingestion, error) {
	repo, err := internal.NewRepository(ctx, root)
	if err != nil {
		return nil, err
	}
	serv := &Service{repo: repo}
	serv.acc = TeleAccumulator{
		signal_func: serv.SignalPush,
	}
	return serv, nil
}

func (s *Service) IngestHeader(payload *parsers.PacketHeader) {
	var session_id uint64
	err := s.repo.Cache.Get(string(entity.GAMESESSION), string(entity.SESSIONID), &session_id)
	s.logger.Error(err.Error())
	if payload.SessionUID != session_id {
		//TODO: Nuke Current Session
	}
	err = s.repo.Cache.Set(string(entity.GAMESESSION), string(entity.SESSIONID), payload.SessionUID)
	s.logger.Error(err.Error())
	err = s.repo.Cache.Set(string(entity.GAMESESSION), string(entity.PLAYERINDEX), payload.PlayerCarIndex)
	s.logger.Error(err.Error())
	//TODO: Add A SigKill Functionality
}
func (s *Service) IngestMotionPacket(payload telentity.MotionPacket) {
	//Ignoring For Now, Maybe could be useful in the future
}
func (s *Service) IngestSessionPacket(payload telentity.PacketSessionData) {}
func (s *Service) IngestLapPacket(payload telentity.LapDataPacket) {
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
func (s *Service) IngestTelemetryPacket(payload telentity.PacketCarTelemetryData) {}
func (s *Service) IngestCarStatusPacket(payload telentity.PacketCarStatusData)    {}
func (s *Service) IngestLobbyInfoPacket(payload telentity.PacketLobbyInfoData)    {}
func (s *Service) IngestCarDamagePacket(payload telentity.PacketCarDamageData) {
	//We will save this data for
}
func (s *Service) IngestSessionHistoryPacket(payload telentity.SessionHistoryPacket) {}
func (s *Service) IngestTyreSetPacket(payload telentity.TyreSetPacket)               {}
