package mysql

import (
	"database/sql"
	"time"
)

func DecodeNullableFloat64(nullFloat64 sql.NullFloat64) float64 {
	if nullFloat64.Valid {
		return nullFloat64.Float64
	}
	return 0
}

func DeccodeNullableInt64(nullInt64 sql.NullInt64) int64 {
	if nullInt64.Valid {
		return nullInt64.Int64
	}
	return 0
}

func DecodeNullableString(nullString sql.NullString) string {
	if nullString.Valid {
		return nullString.String
	}
	return ""
}

func DecodeNullableTimestampInt64(nulltime sql.NullTime) int64 {
	if nulltime.Valid {
		return nulltime.Time.Unix()
	}
	return 0
}

func EncodeNullableFloat32(data float32) sql.NullFloat64 {
	var nullable sql.NullFloat64
	if data != 0 {
		nullable.Float64 = float64(data)
		nullable.Valid = true
	} else {
		nullable.Valid = false
	}
	return nullable
}

func EncodeNullableInt64(data int64) sql.NullInt64 {
	var nullable sql.NullInt64
	if data != 0 {
		nullable.Int64 = data
		nullable.Valid = true
	} else {
		nullable.Valid = false
	}
	return nullable
}

func EncodeNullableTimestamp(timestamp int64) sql.NullTime {
	var nullable sql.NullTime
	if timestamp != 0 {
		nullable.Time = time.Unix(timestamp, 0)
		nullable.Valid = true
	} else {
		nullable.Valid = false
	}
	return nullable
}

func EncodeNullableString(str string) sql.NullString {
	var nullable sql.NullString
	if str != "" {
		nullable.String = str
		nullable.Valid = true
	} else {
		nullable.Valid = false
	}
	return nullable
}
