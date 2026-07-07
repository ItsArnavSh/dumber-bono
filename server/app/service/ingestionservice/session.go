package ingestion

import (
	"context"
	"dubmer-bono/app/types/entity"
	telentity "dubmer-bono/app/types/entity/tel-entity"
)

func (s *Service) SessionHandler(ctx context.Context, session telentity.PacketSessionData) {
	s.repo.Cache.Set(ctx, string(entity.SESSION), "", "")
}
