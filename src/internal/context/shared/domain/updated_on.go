package domain

import "time"

type UpdatedOn time.Time

func NewUpdatedOn() UpdatedOn {
	return UpdatedOn(time.Now().UTC())
}
