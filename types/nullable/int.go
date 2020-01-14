package nullable

import (
	"database/sql"
	"database/sql/driver"
)

type Int64 struct {
	sql.NullInt64
}

func NewInt64(i int64) Int64 {
	if i >= 0 {
		return Int64{sql.NullInt64{
			Int64: i,
			Valid: true,
		}}
	}
	return Int64{sql.NullInt64{
		Int64: -1,
		Valid: false,
	}}
}

// Value implements the driver Valuer interface.
func (i Int64) Value() (driver.Value, error) {
	if !i.Valid {
		return -1, nil
	}
	return i.Int64, nil
}

type Int32 struct {
	sql.NullInt32
}

func NewInt32(i int32) Int32 {
	if i >= 0 {
		return Int32{sql.NullInt32{
			Int32: i,
			Valid: true,
		}}
	}
	return Int32{sql.NullInt32{
		Int32: -1,
		Valid: false,
	}}
}

// Value implements the driver Valuer interface.
// The integer is converted to an int64 as Sql does not support int32 types
func (i Int32) Value() (driver.Value, error) {
	if !i.Valid {
		return int64(-1), nil
	}
	return int64(i.Int32), nil // convert to int64 as mysql does not support int/int32 types
}
