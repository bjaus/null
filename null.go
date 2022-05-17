// null provides "nullable" values with JSON marshaling/unmarshaling abstractions
// on the std database/sql.Null* values.
package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type Nuller interface {
	Null() bool
	Valid() bool
	Set() bool

	driver.Valuer
	sql.Scanner
	json.Marshaler
	json.Unmarshaler
}
