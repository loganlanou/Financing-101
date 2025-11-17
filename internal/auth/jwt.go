package auth

import (
    "context"
    "errors"
    "time"

    "github.com/go-jose/go-jose/v3/jwt"
)

// VerifySession demonstrates how Clerk JWTs would be validated inside Echo middleware.
func (c *ClerkClient) VerifySession(ctx context.Context, token, signingKey string) (*jwt.Claims, error) {
    if token == "" {
        return nil, errors.New("missing token")
    }

    parsed, err := jwt.ParseSigned(token)
    if err != nil {
        return nil, err
    }

    claims := jwt.Claims{}
    if err := parsed.Claims([]byte(signingKey), &claims); err != nil {
        return nil, err
    }

    if err := claims.Validate(jwt.Expected{Time: time.Now()}); err != nil {
        return nil, err
    }

    return &claims, nil
}

