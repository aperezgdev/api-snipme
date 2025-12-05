package domain

import (
	"time"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type RefreshToken struct {
	Id        shared_domain.Id
	UserId    shared_domain.Id
	Token     RefreshTokenToken
	ExpiresAt RefreshTokenExpiresAt
	CreatedOn shared_domain.CreatedOn
}

func NewRefreshToken(
	userId string,
	token string,
	expiresAt time.Time,
) (*RefreshToken, error) {
	id, err := shared_domain.NewID()
	if err != nil {
		return nil, err
	}

	userIdVo, err := shared_domain.ParseID(userId)
	if err != nil {
		return nil, err
	}

	tokenVo, err := NewRefreshTokenToken(token)
	if err != nil {
		return nil, err
	}

	expiresAtVO, err := NewRefreshTokenExpiresAt(expiresAt)
	if err != nil {
		return nil, err
	}

	return &RefreshToken{
		Id:        id,
		UserId:    userIdVo,
		Token:     tokenVo,
		ExpiresAt: expiresAtVO,
		CreatedOn: shared_domain.NewCreatedOn(),
	}, nil
}

func RefreshTokenFromPrimitives(
	id string,
	userId string,
	token string,
	expiresAt time.Time,
	createdOn time.Time,
) (*RefreshToken, error) {
	tokenId, err := shared_domain.ParseID(id)
	if err != nil {
		return nil, err
	}

	userIdVO, err := shared_domain.ParseID(userId)
	if err != nil {
		return nil, err
	}
	expiresAtVO, err := NewRefreshTokenExpiresAt(expiresAt)
	if err != nil {
		return nil, err
	}

	tokenVO, err := NewRefreshTokenToken(token)
	if err != nil {
		return nil, err
	}

	return &RefreshToken{
		Id:        tokenId,
		UserId:    userIdVO,
		Token:     tokenVO,
		ExpiresAt: expiresAtVO,
		CreatedOn: shared_domain.CreatedOn(createdOn),
	}, nil
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Time(rt.ExpiresAt).Before(time.Now())
}
