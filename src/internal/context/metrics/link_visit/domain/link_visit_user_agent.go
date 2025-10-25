package domain

type LinkVisitUserAgent string

func NewLinkVisitUserAgent(userAgent string) LinkVisitUserAgent {
	return LinkVisitUserAgent(userAgent)
}
