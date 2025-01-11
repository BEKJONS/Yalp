package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
	"yalp_ulab/config"
	"yalp_ulab/internal/entity"
	"yalp_ulab/pkg/logger"
	"yalp_ulab/pkg/postgres"
)

type BusinessRepo struct {
	pg     *postgres.Postgres
	cfg    *config.Config
	logger *logger.Logger
}

func NewBusinessRepo(pg *postgres.Postgres, cfg *config.Config, logger *logger.Logger) *BusinessRepo {
	return &BusinessRepo{
		pg:     pg,
		cfg:    cfg,
		logger: logger,
	}
}

func (r *BusinessRepo) Create(ctx context.Context, req entity.Business) (entity.Business, error) {
	req.ID = uuid.NewString()
	query, args, err := r.pg.Builder.Insert("businesses").
		Columns(`id, business_name, location, category, description, contact_information, attachments, created_by`).
		Values(req.ID, req.Name, req.Location, req.Category, req.Description, req.ContactInformation, req.Attachments, req.CreatedBy).ToSql()
	if err != nil {
		return entity.Business{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Business{}, err
	}

	return req, nil
}

func (r *BusinessRepo) GetSingle(ctx context.Context, req entity.BusinessSingleRequest) (entity.Business, error) {
	response := entity.Business{}
	var (
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, business_name, location, category, description, contact_information, attachments, created_by, created_at, updated_at`).
		From("businesses")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)
	default:
		return entity.Business{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.Business{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.Name, &response.Location, &response.Category, &response.Description, &response.ContactInformation, &response.Attachments,
			&response.CreatedBy, &createdAt, &updatedAt)
	if err != nil {
		return entity.Business{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *BusinessRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.BusinessList, error) {
	var (
		response             = entity.BusinessList{}
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, business_name, location, category, description, contact_information, attachments, created_by, created_at, updated_at`).
		From("businesses")

	queryBuilder, where := PrepareGetListQuery(queryBuilder, req)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Business
		err = rows.Scan(&item.ID, &item.Name, &item.Location, &item.Category, &item.Description, &item.ContactInformation, &item.Attachments,
			&item.CreatedBy, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("businesses").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *BusinessRepo) Update(ctx context.Context, req entity.Business) (entity.Business, error) {
	mp := map[string]interface{}{
		"business_name":       req.Name,
		"location":            req.Location,
		"category":            req.Category,
		"description":         req.Description,
		"contact_information": req.ContactInformation,
		"attachments":         req.Attachments,
		"updated_at":          time.Now().Format(time.RFC3339),
	}

	query, args, err := r.pg.Builder.Update("businesses").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Business{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Business{}, err
	}

	return req, nil
}

func (r *BusinessRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("businesses").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *BusinessRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	query, args, err := r.pg.Builder.Update("businesses").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
	if err != nil {
		return response, err
	}

	n, err := r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return response, err
	}

	response.RowsEffected = int(n.RowsAffected())

	return response, nil
}
