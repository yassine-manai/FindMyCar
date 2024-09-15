package config

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// CustomDate for handling date in DD-MM-YYYY format
type CustomDate struct {
	time.Time
}

const customDateFormat = "02/01/2006"

// UnmarshalJSON implements the custom unmarshaling for JSON
func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	dateString := string(data)
	parsedTime, err := time.Parse(`"`+customDateFormat+`"`, dateString)
	if err != nil {
		return fmt.Errorf("invalid date format, expected DD/MM/YYYY: %v", err)
	}
	cd.Time = parsedTime
	return nil
}

// MarshalJSON implements custom marshaling for JSON
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + cd.Time.Format(customDateFormat) + `"`), nil
}

// Scan implements the Scanner interface for reading from the DB
func (cd *CustomDate) Scan(value interface{}) error {
	dateString, ok := value.(string)
	if !ok {
		return fmt.Errorf("could not scan type %T into CustomDate", value)
	}
	parsedTime, err := time.Parse(customDateFormat, dateString)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}
	cd.Time = parsedTime
	return nil
}

// Value implements the Valuer interface for writing to the DB
func (cd CustomDate) Value() (driver.Value, error) {
	return cd.Time.Format(customDateFormat), nil
}
