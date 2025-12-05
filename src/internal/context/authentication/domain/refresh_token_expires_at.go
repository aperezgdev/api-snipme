package domain

import "time"

type RefreshTokenExpiresAt time.Time

func NewRefreshTokenExpiresAt(t time.Time) (RefreshTokenExpiresAt, error) {
	if t.Before(time.Now()) {
		return RefreshTokenExpiresAt{}, nil
	}
	return RefreshTokenExpiresAt(t), nil
}
