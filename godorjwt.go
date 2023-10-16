package godorjwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

const algorithm = "HS256"
const expiry = 60

type Config struct {
	Algorithm string
	Secret    string
	Expiry    int64
}

func New(secret string) *Config {
	return &Config{
		Secret: secret,
	}
}

func Encode(payload map[string]any, config Config) (string, string, int64, error) {
	ex := config.Expiry
	algo := config.Algorithm
	if ex == 0 {
		ex = expiry
	}
	if algo == "" {
		algo = algorithm
	}
	issuedAt := time.Now().UTC().Unix()
	jti, _ := randomHex(32)
	exp := time.Now().UTC().Add(time.Duration(ex) * time.Minute).Unix()
	claims := jwt.MapClaims{
		"iat": issuedAt,
		"jti": jti,
		"exp": exp,
	}
	for key, value := range payload {
		claims[key] = value
	}
	signingMethod := jwt.GetSigningMethod(algo)
	token := jwt.NewWithClaims(signingMethod, claims)
	signedToken, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", "", 0, errors.New("error signing token")
	}
	return signedToken, jti, exp, nil
}

func Decode(tokenString string, config Config) (map[string]any, error) {
	if config.Secret == "" {
		return nil, errors.New("secret is required")
	}
	if tokenString == "" {
		return nil, errors.New("token is required")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (config *Config) Decoder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := ""
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			token = c.Cookies("jwt")
		} else {
			token = strings.Split(authHeader, " ")[1]
		}
		claims, err := Decode(token, *config)

		if token == "" || err != nil || claims["exp"].(int64) < time.Now().UTC().Unix() {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		c.Locals("decodedToken", claims)
		return c.Next()
	}
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
