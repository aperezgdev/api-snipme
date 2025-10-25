package domain

type LinkViewsCounter uint

func NewLinkViewsCounter(views uint) LinkViewsCounter {
	return LinkViewsCounter(views)
}

func (lk LinkViewsCounter) Increment(i uint) LinkViewsCounter {
	return LinkViewsCounter(uint(lk) + i)
}
