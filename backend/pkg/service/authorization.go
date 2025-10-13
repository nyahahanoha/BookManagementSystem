package service

import (
	"context"
	"fmt"
	"strings"

	"log/slog"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type AuthInterceptor struct {
	JWKS        jwk.Set
	Logger      *slog.Logger
	addminEmail string
}

func NewAuthInterceptor(jwks jwk.Set, logger *slog.Logger, addminEmail string) connect.Interceptor {
	return &AuthInterceptor{
		JWKS:        jwks,
		Logger:      logger,
		addminEmail: addminEmail,
	}
}

func (i *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		jwtHeader := req.Header().Get("X-Pomerium-Jwt-Assertion")
		if jwtHeader == "" {
			i.Logger.Warn("missing X-Pomerium-Jwt-Assertion header")
			return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing X-Pomerium-Jwt-Assertion header"))
		}

		token, err := jwt.Parse(jwtHeader, func(token *jwt.Token) (interface{}, error) {
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("missing kid in header")
			}
			key, found := i.JWKS.LookupKeyID(kid)
			if !found {
				return nil, fmt.Errorf("key not found: %s", kid)
			}
			var pubkey interface{}
			if err := key.Raw(&pubkey); err != nil {
				return nil, fmt.Errorf("failed to get raw key: %w", err)
			}
			return pubkey, nil
		})

		if err != nil {
			i.Logger.Warn("invalid Pomerium JWT", slog.String("error", err.Error()))
			return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid token: %v", err))
		}

		if !token.Valid {
			i.Logger.Warn("unauthorized access attempt via Pomerium header")
			return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid token"))
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		if email, ok := claims["email"].(string); ok {
			if strings.HasSuffix(req.Spec().Procedure, "GetAllBooks") ||
				strings.HasSuffix(req.Spec().Procedure, "GetBook") ||
				strings.HasSuffix(req.Spec().Procedure, "SearchBook") {
				return next(ctx, req)
			}
			if email != i.addminEmail {
				i.Logger.Warn("forbidden access attempt via Pomerium", slog.String("email", email))
				return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("forbidden"))
			}
			return next(ctx, req)
		} else {
			return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("email claim not found"))
		}
	}
}

func (i *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}
