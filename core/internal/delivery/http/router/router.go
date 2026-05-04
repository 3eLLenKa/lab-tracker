package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"

	api "lab-tracker/internal/delivery/http/gen"
	"lab-tracker/internal/delivery/http/handlers"
	"lab-tracker/internal/service"
)

const (
	ctxKeyUserID   = "user_id"
	ctxKeyUsername = "username"
	ctxKeyRole     = "role"
)

type tokenClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func Init(svc *service.Service, log *zap.Logger, secret string) (*gin.Engine, error) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(zapLogger(log))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	swagger, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("router: get swagger: %w", err)
	}
	swagger.Servers = nil

	r.Use(authMiddleware(secret))

	r.GET("/_info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "lab-tracker", "status": "ok"})
	})

	httpHandlers := handlers.New(svc)
	handler := api.NewStrictHandler(httpHandlers, nil)
	api.RegisterHandlers(r, handler)

	return r, nil
}

func authMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractBearerToken(c.Request)
		if tokenStr == "" {
			c.Next()
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.Next()
			return
		}

		claims, ok := token.Claims.(*tokenClaims)
		if !ok {
			c.Next()
			return
		}

		c.Set(ctxKeyUserID, claims.UserID)
		c.Set(ctxKeyUsername, claims.Username)
		c.Set(ctxKeyRole, claims.Role)

		ctx := context.WithValue(c.Request.Context(), ctxKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ctxKeyUsername, claims.Username)
		ctx = context.WithValue(ctx, ctxKeyRole, claims.Role)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func zapLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
		)
	}
}

func extractBearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(h, "Bearer ")
}
