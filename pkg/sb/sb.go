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

func Init() error {
	sbHost := os.Getenv("SUPABASE_URL")
	sbSecret := os.Getenv("SUPABASE_SECRET")
	Client = supabase.CreateClient(sbHost, sbSecret)
	return nil
}

//parse supabase user claims

func decodeSBJWT(tokenString string, claims jwt.Claims) error {
	jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}

type userAccessToken struct {
	Email       string `json:"email"`
	AppMetadata struct {
		Provider string `json:"provider"`
	} `json:"app_metadata"`
	UserMetadata interface{} `json:"user_metadata"`
	Role         string      `json:"role"`
	jwt.RegisteredClaims
}

func GetUserClaims(authToken string) (*userAccessToken, error) {
	userClaims := &userAccessToken{}
	if err := decodeSBJWT(authToken, userClaims); err != nil {
		return nil, err
	}
	return userClaims, nil
}
