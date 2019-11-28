package nullable

import (
	"database/sql"
)

type Float64 struct {
	sql.NullFloat64
}

func NewFloat64(i float64) Float64 {
	if i >= 0 {
		return Float64{sql.NullFloat64{
			Float64: i,
			Valid: true,
		}}
	}
	return Float64{sql.NullFloat64{
		Float64: -1,
		Valid: false,
	}}
}

func (i *Float64) Value() float64 {
	if i.Valid {
		return i.NullFloat64.Float64
	} else {
		return -1
	}
}

func (i *Float64) Scan(value interface{}) error {
	return i.NullFloat64.Scan(value)
}

