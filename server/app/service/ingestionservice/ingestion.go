package ingestion

import (
	"dubmer-bono/app/service/internal"
	"dubmer-bono/app/types"
)

//The Ingestion Service Saves the Data in the relevant DBs so that the other processes can query into it

type Service struct {
	repo *internal.Repository
}

var _ types.Ingestion = &Service{}

func NewService(root string) (types.Ingestion, error) {
	repo, err := internal.NewRepository(root)
	return Service{repo: repo}, err
}
