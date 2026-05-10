package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTService struct {
	secret        []byte
	expiration    time.Duration
	refreshExpiry time.Duration
}

type Claims struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewJWTService(secret string, expiration, refreshExpiry time.Duration) *JWTService {
	return &JWTService{
		secret:        []byte(secret),
		expiration:    expiration,
		refreshExpiry: refreshExpiry,
	}
}

func (s *JWTService) GenerateToken(userID uint64, email, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "disability_system",
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTService) GenerateTokenPair(userID uint64, email, role string) (*TokenPair, error) {
	accessToken, err := s.GenerateToken(userID, email, role)
	if err != nil {
		return nil, err
	}

	refreshClaims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "disability_system",
			Subject:   email,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return s.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

func (s *JWTService) Expiration() time.Duration {
	return s.expiration
}

func (s *JWTService) RefreshToken(tokenString string) (*TokenPair, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	return s.GenerateTokenPair(claims.UserID, claims.Email, claims.Role)
}

var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrTokenExpired         = errors.New("token has expired")
	ErrTokenInvalid        = errors.New("invalid token")
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}