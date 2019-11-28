package nullable

import (
	"database/sql"
	"time"
)

type Time struct {
	sql.NullTime
}

func NewTimeFromUnix(unix int64) Time {
	if unix > 0 {
		return Time{
			sql.NullTime{
				Time: time.Unix(unix, 0),
				Valid:  true,
			},
		}
	}
	return Time{
		sql.NullTime{
			Time: time.Time{},
			Valid:  false,
		},
	}
}

func (s *Time) Value() time.Time {
	return s.Time
}

func (s *Time) Scan(value interface{}) error {
	return s.NullTime.Scan(value)
}

func (s *Time) Unix() int64 {
	if s.Valid {
		return s.Time.Unix()
	}
	if s.Time.Unix() == 0 {
		return 0
	}
	return -1
}
