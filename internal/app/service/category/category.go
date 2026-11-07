package scategory

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/Pak3n/catalog-service/internal/app/entity"
	"github.com/Pak3n/catalog-service/internal/app/repository"
	"github.com/Pak3n/catalog-service/internal/app/service"
)

type srv struct {
	repoCategory repository.Category
	repoProduct  repository.Product
}

func NewService(repoCategory repository.Category, repoProduct repository.Product) service.Category {
	return &srv{
		repoCategory: repoCategory,
		repoProduct:  repoProduct,
	}
}

func (s *srv) Create(ctx context.Context, req entity.RequestCategoryCreate) (entity.Category, error) {
	category := entity.Category{
		Name: req.Name,
	}

	err := s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {
		return s.repoCategory.Create(ctx, category)
	})
	if err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (s *srv) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Category, error) {
	return s.repoCategory.GetByGUIDs(ctx, guids)
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestCategoryUpdate) (entity.Category, error) {
	category := entity.Category{
		GUID: guid,
		Name: req.Name,
	}

	err := s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {
		return s.repoCategory.Update(ctx, category)
	})
	if err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	err := s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {
		cats, err := s.repoCategory.GetByGUIDs(ctx, []uuid.UUID{guid})
		if err != nil {
			return err
		}
		if len(cats) == 0 {
			return entity.ErrNotFound
		}

		products, err := s.repoProduct.List(ctx, nil, &guid, nil, nil)
		if err != nil {
			return err
		}
		if len(products) > 0 {
			return entity.ErrCategoryHasProducts
		}

		return s.repoCategory.Delete(ctx, guid)
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *srv) List(ctx context.Context) ([]entity.Category, error) {
	return s.repoCategory.List(ctx, nil)
}
