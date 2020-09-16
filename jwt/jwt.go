package jwt

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/xavierror/gowheel/tool"

	"github.com/dgrijalva/jwt-go"
)

// Create 生成 JWT 令牌
func Create(secret string, data map[string]interface{}, exp int64) (token, short string, err error) {
	// Create Token
	Token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := Token.Claims.(jwt.MapClaims)
	for k, v := range data {
		claims[k] = v
	}

	claims["exp"] = time.Now().Unix() + exp

	// Generate encoded Token and send it as response.
	token, err = Token.SignedString([]byte(secret))
	if err != nil {
		return
	}

	// md5
	short = tool.MD5(token)

	return
}

// Parse 解析 JWT 令牌
func Parse(token, secret string) (bytes []byte, err error) {
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (k interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err = errors.New("Unexpected signing method" + token.Header["alg"].(string))
			return
		}
		k = []byte(secret)
		return
	})

	if err != nil {
		return
	}

	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok && tokenObj.Valid {
		bytes, _ = json.Marshal(claims)
	} else {
		err = errors.New("token has expired")
		return
	}
	return
}
