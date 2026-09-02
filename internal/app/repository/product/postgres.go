package pproduct

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

func NewRepoFromPostgres(client *rcpostgres.Client) repository.Product {
	return &repoPg{_DB: client}
}

func (r *repoPg) Create(ctx context.Context, product entity.Product) error {
	_, err := r._DB.GetRawBunDB().NewInsert().Model(&product).Exec(ctx)
	return err
}

func (r *repoPg) GetByGUIDs(ctx context.Context, guids []uuid.UUID) ([]entity.Product, error) {
	var products []entity.Product
	err := r._DB.GetRawBunDB().
		NewSelect().
		Model(&products).
		Where("guid IN (?)", bun.List(guids)).
		Scan(ctx)
	return products, err
}

func (r *repoPg) Update(ctx context.Context, product entity.Product) error {
	result, err := r._DB.GetRawBunDB().
		NewUpdate().
		Model(&product).
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
		Model((*entity.Product)(nil)).
		Where("guid = ?", guid).
		Exec(ctx)
	return rcpostgres.DeleteErr(err)
}

func (r *repoPg) List(ctx context.Context, name *string, categoryGUID *uuid.UUID) ([]entity.Product, error) {
	var products []entity.Product
	query := r._DB.GetRawBunDB().NewSelect().Model(&products)

	if name != nil {
		query = query.Where("name = ?", *name)
	}
	if categoryGUID != nil {
		query = query.Where("category_guid = ?", *categoryGUID)
	}

	err := query.Order("name ASC").Scan(ctx)
	return products, err
}
