package store

import (
	"context"

	"github.com/harshitrajsinha/van-server-devops/models"
)

type EngineStoreInterface interface {
	GetEngineById(ctx context.Context, id string) (interface{}, error)
	GetAllEngine(ctx context.Context) (interface{}, error)
	CreateEngine(ctx context.Context, engineReq *models.Engine) (int64, error)
	UpdateEngine(ctx context.Context, id string, engineReq *models.Engine) (int64, error)
	DeleteEngine(ctx context.Context, id string) (int64, error)
}

type VanStoreInterface interface {
	GetVanById(ctx context.Context, id string) (interface{}, error)
	GetAllVan(ctx context.Context) (interface{}, error)
	CreateVan(ctx context.Context, vanReq *models.Van) (int64, error)
	UpdateVan(ctx context.Context, id string, vanReq *models.Van) (int64, error)
	DeleteVan(ctx context.Context, id string) (int64, error)
}
