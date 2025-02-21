package account

import (
	"app/db/models"
	"app/lib/log"
	"context"
	"net/http"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"golang.org/x/crypto/bcrypt"
)

var logger = log.Logger()

type Store interface {
	Get(ctx context.Context, id string) (*models.User, error)
	Login(ctx context.Context, email string, password string) (string, error)
	Signup(ctx context.Context, data createUser) (user *models.User, err error)
}

type userStore struct {
	db *bob.DB
}

func NewStore(db *bob.DB) Store {
	return &userStore{db}
}

func (s *userStore) hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *userStore) checkHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *userStore) Get(ctx context.Context, id string) (*models.User, error) {
	userid, err := uuid.FromString(id)
	if err != nil {
		return nil, err
	}

	return models.Users.Query(
		models.SelectWhere.Users.ID.EQ(userid),
	).One(ctx, s.db)
}

func (s *userStore) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := models.Users.Query(
		sm.Where(
			psql.Quote("email").EQ(psql.Arg(email)),
		),
	).One(ctx, s.db)

	if err != nil {
		return "", err
	}

	if ok := s.checkHash(password, user.Password.GetOr("")); !ok {
		return "", echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return user.ID.String(), nil
}

func (s *userStore) Signup(ctx context.Context, data createUser) (*models.User, error) {
	password, err := s.hash(data.Password)
	if err != nil {
		return nil, err
	}

	v7, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	count, err := models.Users.Query(
		models.SelectWhere.Users.Email.EQ(data.Email),
	).Count(ctx, s.db)

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Email already exist")
	}

	query := models.Users.Insert(
		&models.UserSetter{
			ID:       omit.From(v7),
			Email:    omit.From(data.Email),
			Password: omitnull.From(password),
			Source:   omitnull.From("email"),
		},
		im.Returning("*"),
	)

	return query.One(ctx, s.db)
}
