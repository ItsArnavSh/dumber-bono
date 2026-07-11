package monitor

import (
	"dubmer-bono/app/types/entity"
	telentity "dubmer-bono/app/types/entity/tel-entity"
	"strconv"
)

func (s *Service) MyCarIndex() int {
	var ID int
	s.repo.Cache.Get(string(entity.GAMESESSION), string(entity.SESSIONID), &ID)
	return ID
}

func (s *Service) GetLapData() (telentity.LapData, error) {
	var lapdata telentity.LapData
	if err := s.repo.Cache.Get(string(entity.LAPDATA), strconv.Itoa(s.MyCarIndex()), &lapdata); err != nil {
		return telentity.LapData{}, err
	}
	return lapdata, nil
}

func (s *Service) GetSessionData() (telentity.PacketSessionData, error) {
	var sessionData telentity.PacketSessionData
	if err := s.repo.Cache.Get(string(entity.SESSIONDATA), string(entity.SESSIONDATA), &sessionData); err != nil {
		return telentity.PacketSessionData{}, err
	}
	return sessionData, nil
}

func (s *Service) GetMySession() (telentity.PacketSessionData, error) {
	return s.GetSessionData()
}

func (s *Service) GetCarSetupData() (telentity.CarSetupData, error) {
	var carSetup telentity.CarSetupData
	if err := s.repo.Cache.Get(string(entity.CARSETUP), strconv.Itoa(s.MyCarIndex()), &carSetup); err != nil {
		return telentity.CarSetupData{}, err
	}
	return carSetup, nil
}

func (s *Service) GetCarStatusData() (telentity.CarStatusData, error) {
	var carStatus telentity.CarStatusData
	if err := s.repo.Cache.Get(string(entity.CARSTATUS), strconv.Itoa(s.MyCarIndex()), &carStatus); err != nil {
		return telentity.CarStatusData{}, err
	}
	return carStatus, nil
}

func (s *Service) GetLobbyInfoData() (telentity.PacketLobbyInfoData, error) {
	var lobbyInfo telentity.PacketLobbyInfoData
	if err := s.repo.Cache.Get(string(entity.LOBBYINFO), string(entity.LOBBYINFO), &lobbyInfo); err != nil {
		return telentity.PacketLobbyInfoData{}, err
	}
	return lobbyInfo, nil
}

func (s *Service) GetCarDamageData() (telentity.CarDamageData, error) {
	var carDamage telentity.CarDamageData
	if err := s.repo.Cache.Get(string(entity.CARDAMAGE), strconv.Itoa(s.MyCarIndex()), &carDamage); err != nil {
		return telentity.CarDamageData{}, err
	}
	return carDamage, nil
}

func (s *Service) GetSessionHistoryData() (telentity.SessionHistoryPacket, error) {
	var sessionHistory telentity.SessionHistoryPacket
	if err := s.repo.Cache.Get(string(entity.SESSIONHISTORY), string(entity.SESSIONHISTORY), &sessionHistory); err != nil {
		return telentity.SessionHistoryPacket{}, err
	}
	return sessionHistory, nil
}

func (s *Service) GetTyreSetData() (telentity.TyreSetPacket, error) {
	var tyreSet telentity.TyreSetPacket
	if err := s.repo.Cache.Get(string(entity.TYRESET), string(entity.TYRESET), &tyreSet); err != nil {
		return telentity.TyreSetPacket{}, err
	}
	return tyreSet, nil
}

func (s *Service) GetParticipantData() (telentity.ParticipantData, error) {
	var participant telentity.ParticipantData
	if err := s.repo.Cache.Get(string(entity.PARTICIPANT), strconv.Itoa(s.MyCarIndex()), &participant); err != nil {
		return telentity.ParticipantData{}, err
	}
	return participant, nil
}
