package nullable

import "database/sql"

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

func (i *Int64) Value() int64 {
	return i.Int64
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

func (i *Int32) Value() int32 {
	return i.Int32
}
