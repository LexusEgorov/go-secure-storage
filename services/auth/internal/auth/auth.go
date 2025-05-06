package auth

import (
	"auth/internal/models"
	"time"

	"github.com/golang-jwt/jwt"
)

type Storager interface {
	AddUser(string, string) (int, error)
	GetUser(string) (*models.User, error)
	CheckRefresh(string) (bool, error)
	AddToken(userID int, token string) error
}

type AuthProvider struct {
	s          Storager
	tokenTTL   int
	refreshTTL int
	secretKey  []byte
}

func NewAuth(storage Storager) *AuthProvider {
	return &AuthProvider{
		s: storage,
	}
}

func (a AuthProvider) Register(email string, password string) (*models.Credentials, error) {
	user, err := a.s.GetUser(email)

	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, models.ErrConflict
	}

	uID, err := a.s.AddUser(email, password)

	if err != nil {
		return nil, err
	}

	credentials := models.Credentials{
		JWT:     a.createToken(uID, a.tokenTTL),
		Refresh: a.createToken(uID, a.refreshTTL),
	}

	if credentials.JWT == "" || credentials.Refresh == "" {
		return nil, models.ErrInternal
	}

	return &credentials, nil
}

func (a AuthProvider) Auth(email, password string) (*models.Credentials, error) {
	user, err := a.s.GetUser(email)

	if err != nil {
		return nil, err
	}

	credentials := models.Credentials{
		JWT:     a.createToken(user.ID, a.tokenTTL),
		Refresh: a.createToken(user.ID, a.refreshTTL),
	}

	if credentials.JWT == "" || credentials.Refresh == "" {
		return nil, models.ErrInternal
	}

	return &credentials, nil
}

func (a AuthProvider) Refresh(token string) (*models.Credentials, error) {
	ok, err := a.s.CheckRefresh(token)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, models.ErrUnauthorized
	}

	parsedToken := a.parseToken(token)

	if parsedToken == nil {
		return nil, models.ErrUnauthorized
	}

	var uID = -1
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		uID = claims["sub"].(int)
	} else {
		return nil, models.ErrInternal
	}

	credentials := models.Credentials{
		JWT:     a.createToken(uID, a.tokenTTL),
		Refresh: a.createToken(uID, a.refreshTTL),
	}

	if credentials.JWT == "" || credentials.Refresh == "" {
		return nil, models.ErrInternal
	}

	return &credentials, nil
}

func (a AuthProvider) Validate(token string) bool {
	parsed := a.parseToken(token)

	if parsed == nil {
		return false
	}

	return parsed.Valid
}

func (a AuthProvider) parseToken(token string) *jwt.Token {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return a.secretKey, nil
	})

	if err != nil {
		return nil
	}

	return jwtToken
}

func (a AuthProvider) createToken(uID int, duration int) string {
	claims := jwt.MapClaims{
		"sub": uID,
		"exp": time.Now().Add(time.Second * time.Duration(duration)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(a.secretKey)

	if err != nil {
		return ""
	}

	return signed
}
