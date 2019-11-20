package mysql

import (
	"database/sql"
	"time"
)

func DecodeNullFloat64(nullFloat64 sql.NullFloat64) float64 {
	if nullFloat64.Valid {
		return nullFloat64.Float64
	}
	return 0
}

func DecodeNullInt64(nullInt64 sql.NullInt64) int64 {
	if nullInt64.Valid {
		return nullInt64.Int64
	}
	return 0
}

func DecodeNullString(nullString sql.NullString) string {
	if nullString.Valid {
		return nullString.String
	}
	return ""
}

func DecodeNullTime(nulltime sql.NullTime) int64 {
	if nulltime.Valid {
		return nulltime.Time.Unix()
	}
	return 0
}

// EncodeNullableFloat32 will encode any given float32 value
// into a sql.NullFloat64.
// If the value passed is 0, the field will be set invalid, resulting in a NULL in MySQL.
func EncodeFloat64Nullable(data float64) sql.NullFloat64 {
	var nullable sql.NullFloat64
	if data != 0 {
		nullable.Float64 = data
		nullable.Valid = true
	} else {
		nullable.Valid = false
	}
	return nullable
}

// EncodeNullableUnsignedFloat32 will encode any given float64 value into a sql.NullFloat64.
// If the value passed is '>= 0', it's considered valid
// If the value passed is '< 0', it's considered invalid (NULL)
func EncodeUnsignedNullFloat64(data float64) sql.NullFloat64 {
	var nullable sql.NullFloat64
	if data >= 0 {
		nullable.Float64 = data
		nullable.Valid = true
	} else {
		nullable.Valid = false
	}
	return nullable
}

// EncodeNullInt64 will encode any int64 value passed into a sql.NullInt64 value.
// If the value passed is '0', the field is set invalid, resulting in a NULL in MySQL
func EncodeNullInt64(data int64) sql.NullInt64 {
	var nullable sql.NullInt64
	if data != 0 {
		nullable.Int64 = data
		nullable.Valid = true
	} else {
		nullable.Valid = false
	}
	return nullable
}

// EncodeNullInt64 will encode any int64 value passed into a sql.NullInt64 value.
// If the value passed is '>= 0', it's considered valid
// If the value passed is '< 0', it's considered invalid (NULL)
func EncodeUnsignedNullInt64(data int64) sql.NullInt64 {
	var nullable sql.NullInt64
	if data >= 0 {
		nullable.Int64 = data
		nullable.Valid = true
	} else {
		nullable.Valid = false
	}
	return nullable
}

// EncodeNullTime encodes a unix (second) based int64 timestamp into a sql.NullTime
// If the value passed is '0', the timestamp is considered invalid (NULL)
func EncodeNullTime(timestamp int64) sql.NullTime {
	var nullable sql.NullTime
	if timestamp != 0 {
		nullable.Time = time.Unix(timestamp, 0)
		nullable.Valid = true
	} else {
		nullable.Valid = false
	}
	return nullable
}

// EncodeUnsignedNullTime encodes a unix (second) based int64 timestamp into a sql.NullTime
// If the value passed is '>= 0', it's considered valid
// If the value passed is '< 0', it's considered invalid (NULL)
func EncodeUnsignedNullTime(timestamp int64) sql.NullTime {
	var nullable sql.NullTime
	if timestamp >= 0 {
		nullable.Time = time.Unix(timestamp, 0)
		nullable.Valid = true
	} else {
		nullable.Valid = false
	}
	return nullable
}

func EncodeNullString(str string) sql.NullString {
	var nullable sql.NullString
	if str != "" {
		nullable.String = str
		nullable.Valid = true
	} else {
		nullable.Valid = false
	}
	return nullable
}
