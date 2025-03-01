package account

import (
	"app/common"
	"app/db/model"
	"app/lib/log"
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

var logger = log.Logger()

type Store interface {
	Get(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Login(ctx context.Context, email string, password ...string) (string, error)
	Signup(ctx context.Context, data createUser) (user *model.User, err error)
	Update(ctx context.Context, id string, data updateUser) error
	Delete(ctx context.Context, id string) error
}

type userStore struct {
	db *bun.DB
}

func NewStore(db *bun.DB) Store {
	return &userStore{db}
}

func (s *userStore) Get(ctx context.Context, id string) (*model.User, error) {
	user := new(model.User)
	query := s.db.NewSelect().
		Model(user).
		Where("id = ?", id)

	err := query.Scan(ctx, user)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, nil
	}

	return user, nil
}

func (s *userStore) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)
	query := s.db.NewSelect().
		Model(user).
		Where("email = ?", email)

	err := query.Scan(ctx, user)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userStore) Login(ctx context.Context, email string, password ...string) (string, error) {
	user, err := s.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if user.Source == "email" {
		if len(password) == 0 {
			return "", echo.NewHTTPError(http.StatusInternalServerError, "Please provide password for user with email login")
		}

		p := password[0]
		if ok := common.CheckHash(p, *user.Password); !ok {
			return "", echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		return user.Id, nil
	}

	return user.Id, nil
}

func (s *userStore) Signup(ctx context.Context, data createUser) (*model.User, error) {
	// to accomodate oauth signup
	password := new(string)
	if data.Password != nil {
		pass, err := common.Hash(*data.Password)
		if err != nil {
			return nil, err
		}

		password = &pass
	}

	user, err := s.GetByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Email already exist")
	}

	id, _ := uuid.NewV7()
	var verifiedAt *time.Time
	if data.Source == "google" {
		now := time.Now()
		verifiedAt = &now
	} else {
		verifiedAt = nil
	}

	newUser := &model.User{
		Id:         id.String(),
		Email:      data.Email,
		Password:   password,
		Name:       data.Name,
		Source:     data.Source,
		VerifiedAt: verifiedAt,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if _, err := s.db.NewInsert().Model(newUser).Exec(ctx); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userStore) Update(ctx context.Context, id string, data updateUser) error {
	query := s.db.NewUpdate().
		Model((*model.User)(nil)).
		Where("id = ?", id)

	if data.Name != "" {
		query = query.Set("name = ?", data.Name)
	}

	if data.Password != "" {
		password, err := common.Hash(data.Password)
		if err != nil {
			return err
		}

		query = query.Set("password = ?", password)
	}

	if !data.VerifiedAt.IsZero() {
		query = query.Set("verified_at = ?", data.VerifiedAt)
	}

	if _, err := query.Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *userStore) Delete(ctx context.Context, id string) error {
	// TODO: do something before deleting the user

	_, err := s.db.NewDelete().
		Model((*model.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}
