package nullable

import (
	"database/sql"
	"database/sql/driver"
)

type String struct {
	sql.NullString
}

func NewString(str string) String {
	if str == "" {
		return String{
			sql.NullString{
				String: "",
				Valid:  false,
			},
		}
	}
	return String{
		sql.NullString{
			String: str,
			Valid:  true,
		},
	}
}

func (s String) Scan(value interface{}) error {
	return s.NullString.Scan(value)
}

func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.String, nil
}
