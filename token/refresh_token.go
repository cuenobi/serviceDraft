package token

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type RefreshClaims struct {
	Username string `json:"username"`
	Place    int    `json:"place"`
	jwt.StandardClaims
}

func RefreshJwtString(username string, place int) (string, error) {
	exp := time.Duration(viper.GetInt(`token.refreshExp`))
	claims := &RefreshClaims{
		Username: username,
		Place:    place,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * exp).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(viper.GetString(`token.refreshSecret`)))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func RefreshParseJWTToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "MISSING_AUTH_HEADER",
			})
		}

		tokenString := strings.Replace(header, "Bearer ", "", 1)
		token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString(`token.refreshSecret`)), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "INVALID_TOKEN_SIGNATURE",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "INVALID_OR_EXPIRED_TOKEN",
			})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "EXPIRED_TOKEN",
			})
		}

		claims, ok := token.Claims.(*RefreshClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "INVALID_TOKEN_CLAIMS",
			})
		}

		c.Locals("user", claims)
		return c.Next()
	}
}
