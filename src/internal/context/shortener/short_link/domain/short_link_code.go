package domain

import "github.com/aperezgdev/api-snipme/src/pkg"

type ShortLinkCode string

func NewCode() (ShortLinkCode, error) {
	code, err := pkg.GenerateShortCode(6)
	if err != nil {
		return "", err
	}
	return ShortLinkCode(code), nil
}
