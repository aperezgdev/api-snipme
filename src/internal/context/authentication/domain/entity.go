package domain

import (
	"time"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type User struct {
	shared_domain.AggregateRoot
	Id            shared_domain.Id
	Email         shared_domain.Email
	OAuthProvider OAuthProvider
	OAuthSubject  OAuthSubject
	CreatedOn     shared_domain.CreatedOn
}

func NewUser(
	email string,
	oauthProvider OAuthProvider,
	oauthSubject string,
) (*User, error) {
	id, err := shared_domain.NewID()
	if err != nil {
		return nil, err
	}

	emailVO, err := shared_domain.NewEmail(email)
	if err != nil {
		return nil, err
	}
	
	oauthSubjectVO, err := NewOAuthSubject(oauthSubject)
	if err != nil {
		return nil, err
	}
	
	user := &User{
		Id:            id,
		Email:         emailVO,
		OAuthProvider: oauthProvider,
		OAuthSubject:  oauthSubjectVO,
		CreatedOn:     shared_domain.NewCreatedOn(),
	}
	
	user.Record(NewUserCreatedEvent(id.String(), email, oauthProvider.String()))
	
	return user, nil
}

func UserFromPrimitives(
	id string,
	email string,
	oauthProvider string,
	oauthSubject string,
	createdOn time.Time,
) (*User, error) {
	userId, err := shared_domain.ParseID(id)
	if err != nil {
		return nil, err
	}

	emailVO, err := shared_domain.NewEmail(email)
	if err != nil {
		return nil, err
	}
	
	provider, err := NewOAuthProvider(oauthProvider)
	if err != nil {
		return nil, err
	}

	oauthSubjectVO, err := NewOAuthSubject(oauthSubject)
	if err != nil {
		return nil, err
	}
	
	return &User{
		Id:            userId,
		Email:         emailVO,
		OAuthProvider: provider,
		OAuthSubject:  oauthSubjectVO,
		CreatedOn:     shared_domain.CreatedOn(createdOn),
	}, nil
}
