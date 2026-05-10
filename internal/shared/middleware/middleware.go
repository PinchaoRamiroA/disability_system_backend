package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	RequestIDHeader = "X-Request-ID"
	AuthorizationHeader = "Authorization"
	BearerPrefix = "Bearer "
	RequestIDContextKey = "request_id"
	UserClaimsContextKey = "user_claims"
)

type JWTMiddleware struct {
	secret []byte
	expiration time.Duration
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTMiddleware(secret string, expiration time.Duration) *JWTMiddleware {
	return &JWTMiddleware{
		secret: []byte(secret),
		expiration: expiration,
	}
}

func (m *JWTMiddleware) GenerateToken(userID uint, email, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "disability_system",
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTMiddleware) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

var (
	ErrInvalidSigningMethod = jwt.ErrSignatureInvalid
	ErrInvalidToken = jwt.ErrTokenSignatureInvalid
)

func JWTAuth(secret string) gin.HandlerFunc {
	middleware := NewJWTMiddleware(secret, 24*time.Hour)
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Authorization header required",
				"error": gin.H{"code": "MISSING_AUTH_HEADER"},
			})
			return
		}

		if len(authHeader) <= len(BearerPrefix) {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Invalid authorization header format",
				"error": gin.H{"code": "INVALID_AUTH_FORMAT"},
			})
			return
		}

		tokenString := authHeader[len(BearerPrefix):]
		claims, err := middleware.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Token inválido o expirado",
				"error": gin.H{"code": "INVALID_TOKEN"},
			})
			return
		}

		c.Set(UserClaimsContextKey, claims)
		c.Next()
	}
}

func GetUserClaims(c *gin.Context) *Claims {
	claims, exists := c.Get(UserClaimsContextKey)
	if !exists {
		return nil
	}
	return claims.(*Claims)
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set(RequestIDContextKey, requestID)
		c.Header(RequestIDHeader, requestID)
		c.Next()
	}
}

func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDContextKey); exists {
		return requestID.(string)
	}
	return ""
}

func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if isOriginAllowed(origin, allowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
			c.Header("Access-Control-Expose-Headers", "X-Request-ID")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func isOriginAllowed(origin string, allowed []string) bool {
	for _, o := range allowed {
		if o == "*" || o == origin {
			return true
		}
	}
	return false
}

func Recovery(logger any) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"message": "Error interno del servidor",
					"error": gin.H{
						"code": "INTERNAL_ERROR",
						"request_id": GetRequestID(c),
					},
				})
			}
		}()
		c.Next()
	}
}

func Logger(logFunc func(msg string, args ...any)) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logFunc("HTTP Request",
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"status", status,
			"latency", latency.String(),
			"ip", c.ClientIP(),
			"request_id", GetRequestID(c),
		)
	}
}