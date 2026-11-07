package sproduct

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/Pak3n/catalog-service/internal/app/entity"
	"github.com/Pak3n/catalog-service/internal/app/repository"
	"github.com/Pak3n/catalog-service/internal/app/service"
)

type srv struct {
	repoProduct  repository.Product
	repoCategory repository.Category
}

func NewService(repoProduct repository.Product, repoCategory repository.Category) service.Product {
	return &srv{
		repoProduct:  repoProduct,
		repoCategory: repoCategory,
	}
}

func (s *srv) Create(ctx context.Context, req entity.RequestProductCreate) (entity.Product, error) {
	product := entity.Product{
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		CategoryGUID: req.CategoryGUID,
	}

	err := s.repoProduct.InsideTx(ctx, func(ctx context.Context) error {
		return s.repoProduct.Create(ctx, product)
	})
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (s *srv) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Product, error) {
	return s.repoProduct.GetByGUIDs(ctx, guids)
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error) {
	product := entity.Product{
		GUID:         guid,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		CategoryGUID: req.CategoryGUID,
	}

	err := s.repoProduct.InsideTx(ctx, func(ctx context.Context) error {
		return s.repoProduct.Update(ctx, product)
	})
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	err := s.repoProduct.InsideTx(ctx, func(ctx context.Context) error {
		products, err := s.repoProduct.GetByGUIDs(ctx, []uuid.UUID{guid})
		if err != nil {
			return err
		}
		if len(products) == 0 {
			return entity.ErrNotFound
		}

		return s.repoProduct.Delete(ctx, guid)
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *srv) List(ctx context.Context, req entity.RequestProductList) ([]entity.Product, error) {
	return s.repoProduct.List(ctx, nil, req.CategoryGUID, req.MinPrice, req.MaxPrice)
}
