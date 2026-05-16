package service

import (
	"context"

	"github.com/harshitrajsinha/van-server-devops/models"
	"github.com/harshitrajsinha/van-server-devops/store"
)

type VanService struct {
	store store.VanStoreInterface
}

func NewVanService(store store.VanStoreInterface) *VanService {
	return &VanService{
		store: store,
	}
}

func (v *VanService) GetVanById(ctx context.Context, id string) (interface{}, error) {
	van, err := v.store.GetVanById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &van, nil
}

func (v *VanService) GetAllVan(ctx context.Context) (interface{}, error) {
	van, err := v.store.GetAllVan(ctx)
	if err != nil {
		return nil, err
	}
	return &van, nil
}

func (v *VanService) CreateVan(ctx context.Context, vanReq *models.Van) (int64, error) {

	createdVan, err := v.store.CreateVan(ctx, vanReq)
	if err != nil {
		return -1, err
	}

	return createdVan, nil
}

func (v *VanService) UpdateVan(ctx context.Context, id string, vanReq *models.Van) (int64, error) {

	updatedVan, err := v.store.UpdateVan(ctx, id, vanReq)
	if err != nil {
		return -1, err
	}
	return updatedVan, nil
}

func (v *VanService) DeleteVan(ctx context.Context, id string) (int64, error) {

	deletedVan, err := v.store.DeleteVan(ctx, id)
	if err != nil {
		return -1, err
	}
	return deletedVan, nil
}
