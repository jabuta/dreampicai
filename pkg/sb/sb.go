package sb

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nedpals/supabase-go"
)

var (
	Client *supabase.Client
)

func DecodeSBJWT(tokenString string) (jwt.MapClaims, error) {
	jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func Init() error {
	sbHost := os.Getenv("SUPABASE_URL")
	sbSecret := os.Getenv("SUPABASE_SECRET")
	Client = supabase.CreateClient(sbHost, sbSecret)
	return nil
}
