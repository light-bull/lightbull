package utils

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/light-bull/lightbull/persistence"
)

const (
	jwtIssuer   = "lightbull"
	jwtValidity = 60 * time.Minute
)

// JWTManager implements the easy creation and validation of JSON Web Tokens
type JWTManager struct {
	key []byte
}

// NewJWTManager initializes a new JWTManager and prepares the key material
func NewJWTManager(persistence *persistence.Persistence) (*JWTManager, error) {
	jwtManager := JWTManager{}

	type format struct {
		Key       []byte    `json:"key"`
		Generated time.Time `json:"generated"`
	}
	data := format{}

	if persistence.HasConfig("jwt") {
		// config is there -> load it or fail
		if err := persistence.LoadConfig("jwt", &data); err != nil {
			return nil, err
		}

		jwtManager.key = data.Key
	} else {
		// generate key and store it
		jwtManager.key = make([]byte, 64)
		if _, err := rand.Read(jwtManager.key); err != nil {
			return nil, errors.New("Failed to generate secret key for JWT: " + err.Error())
		}

		data.Key = jwtManager.key
		data.Generated = time.Now()
		if err := persistence.SaveConfig("jwt", &data, true); err != nil {
			return nil, err
		}
	}

	return &jwtManager, nil
}

// New issues a new JSON Web Token
func (jwtmanager *JWTManager) New() (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    jwtIssuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtValidity)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(jwtmanager.key)
}

// Check validates the given JSON Web Token
func (jwtmanager *JWTManager) Check(tokenString string) bool {
	// parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		// give signing key to parser
		return jwtmanager.key, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	return true
}
