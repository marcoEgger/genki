package nullable

import (
	"database/sql"
	"database/sql/driver"
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

func NewTime(t time.Time) Time {
	if t.IsZero() || t.Unix() <= 0 {
		return Time{sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}}
	}
	return Time{sql.NullTime{
		Time:  t,
		Valid: true,
	}}
}

func (s Time) Evaluated() time.Time {
	if s.Valid {
		return s.Time
	}
	return time.Time{}
}


// Value implements the driver Valuer interface.
func (s Time) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.Time, nil
}

func (s Time) Unix() int64 {
	if s.Valid {
		return s.Time.Unix()
	}
	if s.Time.Unix() == 0 {
		return 0
	}
	return -1
}
