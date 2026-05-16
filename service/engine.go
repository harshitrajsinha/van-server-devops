package service

import (
	"context"

	"github.com/harshitrajsinha/van-server-devops/models"
	"github.com/harshitrajsinha/van-server-devops/store"
)

type EngineService struct {
	store store.EngineStoreInterface
}

func NewEngineService(store store.EngineStoreInterface) *EngineService {
	return &EngineService{
		store: store,
	}
}

func (s *EngineService) GetEngineByID(ctx context.Context, id string) (interface{}, error) {
	engine, err := s.store.GetEngineById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &engine, nil
}

func (s *EngineService) GetAllEngine(ctx context.Context) (interface{}, error) {
	engine, err := s.store.GetAllEngine(ctx)
	if err != nil {
		return nil, err
	}
	return &engine, nil
}

func (s *EngineService) CreateEngine(ctx context.Context, engineReq *models.Engine) (int64, error) {

	createdEngine, err := s.store.CreateEngine(ctx, engineReq)
	if err != nil {
		return -1, err
	}

	return createdEngine, nil
}

func (s *EngineService) UpdateEngine(ctx context.Context, id string, engineReq *models.Engine) (int64, error) {

	updatedEngine, err := s.store.UpdateEngine(ctx, id, engineReq)
	if err != nil {
		return -1, err
	}
	return updatedEngine, nil
}

func (s *EngineService) DeleteEngine(ctx context.Context, id string) (int64, error) {

	deletedEngine, err := s.store.DeleteEngine(ctx, id)
	if err != nil {
		return -1, err
	}
	return deletedEngine, nil
}
