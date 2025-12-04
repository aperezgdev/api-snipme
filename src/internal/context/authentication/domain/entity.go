package domain

import (
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
