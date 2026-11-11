package scategory

import (
	"context"
	"time"

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
	var category entity.Category

	// Открываем транзакцию. Вся логика внутри.
	err := s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {
		// 1. Проверка на существование (внутри транзакции!)
		existing, err := s.repoCategory.List(ctx, &req.Name)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return entity.ErrAlreadyExists
		}

		// 2. Создание сущности
		now := time.Now()
		category = entity.Category{
			GUID:      uuid.Must(uuid.NewV4()),
			Name:      req.Name,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// 3. Сохранение в БД
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
	var updatedCategory entity.Category

	err := s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {
		// 1. Получаем текущую категорию
		categories, err := s.repoCategory.GetByGUIDs(ctx, []uuid.UUID{guid})
		if err != nil {
			return err
		}
		if len(categories) == 0 {
			return entity.ErrNotFound
		}

		category := categories[0] // Берем существующую сущность

		// 2. Проверяем уникальность нового имени (если оно меняется)
		if req.Name != "" && req.Name != category.Name {
			existing, err := s.repoCategory.List(ctx, &req.Name)
			if err != nil {
				return err
			}
			for _, e := range existing {
				if e.GUID != guid {
					return entity.ErrAlreadyExists
				}
			}
			category.Name = req.Name // Обновляем имя
		}

		// 3. Обновляем время
		category.UpdatedAt = time.Now()

		// 4. Сохраняем в БД
		if err := s.repoCategory.Update(ctx, category); err != nil {
			return err
		}

		// 5. Сохраняем результат для возврата
		updatedCategory = category
		return nil
	})
	if err != nil {
		return entity.Category{}, err
	}

	return updatedCategory, nil
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
