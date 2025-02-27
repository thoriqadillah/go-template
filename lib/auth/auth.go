package auth

import (
	"app/env"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type jwtClaims struct {
	UserId string `json:"id"`
	jwt.RegisteredClaims
}

func SignToken(userid string) (string, error) {
	claims := &jwtClaims{
		UserId: userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(env.JWT_SECRET))
}

func DecodeToken(tokenStr string) (*jwtClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != echojwt.AlgorithmHS256 {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		return []byte(env.JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	mapClaims := token.Claims.(jwt.MapClaims)
	expiresAt, err := mapClaims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	claims := &jwtClaims{
		UserId: mapClaims["id"].(string),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
		},
	}

	return claims, nil
}

var AuthenticatedMw = echojwt.WithConfig(echojwt.Config{
	SigningKey: []byte(env.JWT_SECRET),
	ErrorHandler: func(c echo.Context, err error) error {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Unauthorized",
		})
	},
	NewClaimsFunc: func(c echo.Context) jwt.Claims {
		return &jwtClaims{}
	},
})

func User(c echo.Context) *jwtClaims {
	user := c.Get("user").(*jwt.Token)
	return user.Claims.(*jwtClaims)
}
