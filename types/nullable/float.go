package nullable

import (
	"database/sql"
	"database/sql/driver"
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

// Value implements the driver Valuer interface.
func (i Float64) Value() (driver.Value, error) {
	if !i.Valid {
		return float64(-1), nil
	}
	return i.Float64, nil
}

func (i *Float64) Scan(value interface{}) error {
	return i.NullFloat64.Scan(value)
}

