package nullable

import "database/sql"

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

func (s *String) Value() string {
	return s.String
}

func (s *String) Scan(value interface{}) error {
	return s.NullString.Scan(value)
}

