package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"sejastip.id/api"
	"sejastip.id/api/entity"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// StandardHandler is a default HTTP handler for requests
type StandardHandler func(http.ResponseWriter, *http.Request, httprouter.Params) error

// Middleware is our decorator middleware
type Middleware func(StandardHandler) StandardHandler

// HTTP will pass StandardHandler to be compatible with httprouter handlers
func HTTP(handle StandardHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		handle(w, r, p)
	}
}

// AppendMiddlewares appends a group of middlewares to a standard handler
func AppendMiddlewares(handler StandardHandler, middlewares ...Middleware) StandardHandler {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}

// WithLogging encapsulates standard handlers with logging function
func WithLogging(logger *zap.Logger) Middleware {
	return func(handle StandardHandler) StandardHandler {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
			start := time.Now()

			err := handle(w, r, p)

			elapsed := time.Since(start).Seconds() * 1000
			timeElapsedStr := strconv.FormatFloat(elapsed, 'f', -1, 64)
			if err != nil {
				logger.Error(err.Error(),
					zap.String("duration", timeElapsedStr),
					zap.String("request_path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)),
				)
			} else {
				logger.Info("OK",
					zap.String("duration", timeElapsedStr),
					zap.String("request_path", fmt.Sprintf("%s %s", r.Method, r.URL.Path)),
				)
			}

			return err
		}
	}
}

// WithAuthentication encapsulates standard handlers with authentication
func WithAuthentication(privateKey string) Middleware {
	return func(handle StandardHandler) StandardHandler {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
			authHeader := r.Header.Get("Authorization")
			match, _ := regexp.MatchString("^Token ", authHeader)
			if !match {
				log.Println("Unaut coj")
				api.Error(w, api.ErrUnauthorized)
				return api.ErrUnauthorized
			}

			// get claims
			claims := entity.ResourceClaims{}
			tokenString := authHeader[6:]
			token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(privateKey), nil
			})
			if err != nil || !token.Valid {
				log.Println("Unaut yok", err)
				api.Error(w, api.ErrUnauthorized)
				return api.ErrUnauthorized
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, api.ContextKeyName, claims)
			return handle(w, r.WithContext(ctx), p)
		}
	}
}

// DefaultMiddlewares will return default configured middlewares
func DefaultMiddlewares() []Middleware {
	l, _ := zap.NewProduction()
	ms := []Middleware{
		WithLogging(l),
	}
	return ms
}
