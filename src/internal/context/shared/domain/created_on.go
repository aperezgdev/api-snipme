package domain

import "time"

type CreatedOn time.Time

func NewCreatedOn() CreatedOn {
	return CreatedOn(time.Now().UTC())
}
