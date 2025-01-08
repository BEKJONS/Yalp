package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"yalp_ulab/config"
	"yalp_ulab/internal/entity"
	"yalp_ulab/pkg/logger"
	"yalp_ulab/pkg/postgres"
)

type UserRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewUserRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UserRepo {
	return &UserRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *UserRepo) Create(ctx context.Context, req entity.User) (entity.User, error) {
	req.ID = uuid.NewString()

	query, args, err := r.pg.Builder.Insert("users").
		Columns(`id, full_name, email, password, user_type, user_role, status`).
		Values(req.ID, req.FullName, req.Email, req.Password, req.UserType, req.UserRole, req.Status).ToSql()
	if err != nil {
		return entity.User{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.User{}, err
	}

	return req, nil
}

func (r *UserRepo) GetSingle(ctx context.Context, req entity.UserSingleRequest) (entity.User, error) {
	response := entity.User{}
	var (
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, full_name, email, password, user_type, user_role, status, created_at, updated_at`).
		From("users")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)
	case req.Email != "":
		queryBuilder = queryBuilder.Where("email = ?", req.Email)
	default:
		return entity.User{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.User{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.FullName, &response.Email, &response.Password,
			&response.UserType, &response.UserRole, &response.Status, &createdAt, &updatedAt)
	if err != nil {
		return entity.User{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

func (r *UserRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error) {
	var (
		response             = entity.UserList{}
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select(`id, full_name, email, password, user_type, user_role, status, created_at, updated_at`).
		From("users")

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
		var item entity.User
		err = rows.Scan(&item.ID, &item.FullName, &item.Email, &item.Password,
			&item.UserType, &item.UserRole, &item.Status, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("users").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *UserRepo) Update(ctx context.Context, req entity.User) (entity.User, error) {
	mp := map[string]interface{}{
		"full_name":  req.FullName,
		"status":     req.Status,
		"email":      req.Email,
		"user_role":  req.UserRole,
		"updated_at": time.Now().Format(time.RFC3339),
	}

	if req.Password != "" {
		mp["password"] = req.Password
	}

	query, args, err := r.pg.Builder.Update("users").SetMap(mp).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.User{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.User{}, err
	}

	return req, nil
}

func (r *UserRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("users").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	query, args, err := r.pg.Builder.Update("users").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
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
