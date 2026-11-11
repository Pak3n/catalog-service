package sproduct

import (
	"context"
	"time"

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
	var product entity.Product

	err := s.repoProduct.InsideTx(ctx, func(ctx context.Context) error {
		// 1. Проверяем, есть ли уже продукт с таким именем (внутри транзакции)
		existing, err := s.repoProduct.List(ctx, &req.Name, nil, nil, nil)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return entity.ErrAlreadyExists
		}

		// 2. Проверяем, существует ли категория, к которой привязываем продукт
		categories, err := s.repoCategory.GetByGUIDs(ctx, []uuid.UUID{req.CategoryGUID})
		if err != nil {
			return err
		}
		if len(categories) == 0 {
			return entity.ErrNotFound
		}

		// 3. Создаем сущность
		now := time.Now()
		product = entity.Product{
			GUID:         uuid.Must(uuid.NewV4()),
			Name:         req.Name,
			Description:  req.Description,
			Price:        req.Price,
			CategoryGUID: req.CategoryGUID,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		// 4. Сохраняем в БД
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
	var updatedProduct entity.Product

	err := s.repoProduct.InsideTx(ctx, func(ctx context.Context) error {
		// 1. Получаем текущий продукт
		products, err := s.repoProduct.GetByGUIDs(ctx, []uuid.UUID{guid})
		if err != nil {
			return err
		}
		if len(products) == 0 {
			return entity.ErrNotFound
		}
		product := products[0]

		// 2. Обновляем имя (с проверкой на дубликат)
		if req.Name != "" {
			existing, err := s.repoProduct.List(ctx, &req.Name, nil, nil, nil)
			if err != nil {
				return err
			}
			for _, e := range existing {
				if e.GUID != guid {
					return entity.ErrAlreadyExists
				}
			}
			product.Name = req.Name
		}

		// 3. Обновляем категорию (с проверкой её существования)
		if !req.CategoryGUID.IsNil() {
			categories, err := s.repoCategory.GetByGUIDs(ctx, []uuid.UUID{req.CategoryGUID})
			if err != nil {
				return err
			}
			if len(categories) == 0 {
				return entity.ErrNotFound
			}
			product.CategoryGUID = req.CategoryGUID
		}

		// 4. Обновляем описание
		if req.Description != nil {
			product.Description = req.Description
		}

		// 5. Обновляем цену
		if req.Price > 0 {
			product.Price = req.Price
		}

		// 6. Обновляем время
		product.UpdatedAt = time.Now()

		// 7. Сохраняем в БД
		if err := s.repoProduct.Update(ctx, product); err != nil {
			return err
		}

		// 8. Сохраняем результат для возврата
		updatedProduct = product
		return nil
	})
	if err != nil {
		return entity.Product{}, err
	}

	return updatedProduct, nil
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
