package paypack

import (
    "context"
    "log"

    paypack "github.com/paypack/go-sdk"
)

type Service struct {
    client *paypack.Client
}


type token struct {
    refresh_token string
    access_token  string
    expires_at    int64
}


func (s *Service) Login(ctx context.Context, client_id, client_secret string) (*token, error) {

    tk, err := s.client.Auth.Login(ctx, client_id, client_secret)
    if err != nil {
        log.Printf("Login Error: %s", err)
        return nil, err
    }

    return &token{
        access_token:  tk.Access,
        refresh_token: tk.Refresh,
        expires_at:    tk.Expires,
    }, nil
}


func (s *Service) Refresh(ctx context.Context, refresh_token string) (*token, error) {

    in := &paypack.Token{Refresh: refresh_token}

    tk, err := s.client.Auth.Refresh(ctx, in)
    if err != nil {
        log.Printf("Refresh Error: %s", err)
        return nil, err
    }

    return &token{
        access_token:  tk.Access,
        refresh_token: tk.Refresh,
        expires_at:    tk.Expires,
    }, nil
}