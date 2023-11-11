package token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GetExpiration(token string) (time.Time, error) {
	p := jwt.NewParser()
	var claims jwt.MapClaims
	_, _, err := p.ParseUnverified(token, &claims)
	if err != nil {
		return time.Time{}, err
	}
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return time.Time{}, err
	}
	return exp.Time, nil
}
