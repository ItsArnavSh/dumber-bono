package internal

import (
	"context"
	"dubmer-bono/app/service/internal/clickhouse"
	"dubmer-bono/app/service/internal/postgresql"
	"dubmer-bono/app/service/internal/qdrant"
	rdb "dubmer-bono/app/service/internal/redis"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	pb "github.com/qdrant/go-client/qdrant"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	db           *pgxpool.Pool
	redis        *redis.Client
	chdb         driver.Conn
	qdrantClient pb.CollectionsClient
	qdrantConn   *grpc.ClientConn
}

func NewRepository(ctx context.Context) (*Repository, error) {
	db, err := postgresql.NewPostgresPool(ctx)
	if err != nil {
		return nil, err
	}

	rdb, err := rdb.NewRedisClient(ctx)
	if err != nil {
		return nil, err
	}

	chdb, err := clickhouse.NewClickHouseConn(ctx)
	if err != nil {
		return nil, err
	}

	qdrantClient, qdrantConn, err := qdrant.NewQdrantClient(ctx)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:           db,
		redis:        rdb,
		chdb:         chdb,
		qdrantClient: qdrantClient,
		qdrantConn:   qdrantConn,
	}, nil
}
