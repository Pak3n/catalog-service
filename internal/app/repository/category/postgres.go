package pcategory

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"

	"github.com/Pak3n/catalog-service/internal/app/entity"
	"github.com/Pak3n/catalog-service/internal/app/repository"
	rcpostgres "github.com/Pak3n/catalog-service/internal/app/repository/conn/postgres"
)

type (
	repoPg struct {
		*_DB
	}

	_DB = rcpostgres.Client
)

func NewRepoFromPostgres(client *rcpostgres.Client) repository.Category {
	return &repoPg{_DB: client}
}

func (r *repoPg) Create(ctx context.Context, category entity.Category) error {
	_, err := r._DB.GetRawBunDB().NewInsert().Model(&category).Exec(ctx)
	return err
}

func (r *repoPg) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Category, error) {
	var categories []entity.Category
	err := r._DB.GetRawBunDB().
		NewSelect().
		Model(&categories).
		Where("guid IN (?)", bun.List(guids)).
		Scan(ctx)
	return categories, err
}

func (r *repoPg) Update(ctx context.Context, category entity.Category) error {
	result, err := r._DB.GetRawBunDB().
		NewUpdate().
		Model(&category).
		WherePK().
		ExcludeColumn("id", "created_at").
		Exec(ctx)
	if err != nil {
		return err
	}
	return rcpostgres.UpdateErr(result, nil)
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := r._DB.GetRawBunDB().
		NewDelete().
		Model((*entity.Category)(nil)).
		Where("guid = ?", guid).
		Exec(ctx)
	return rcpostgres.DeleteErr(err)
}

func (r *repoPg) List(ctx context.Context, name *string) ([]entity.Category, error) {
	var categories []entity.Category
	query := r._DB.GetRawBunDB().NewSelect().Model(&categories)
	if name != nil {
		query = query.Where("name = ?", *name)
	}
	err := query.Order("name ASC").Scan(ctx)
	return categories, err
}
