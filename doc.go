// Package null provides a generic nullable type that distinguishes between
// unset, null, and valid values.
//
// null solves a common problem in Go: the inability to distinguish between
// "field was not provided" and "field was explicitly set to null" when working
// with JSON APIs, databases, and other data sources.
//
// # The Three-State Problem
//
// Consider a PATCH request where a client wants to:
//   - Leave a field unchanged (don't include it)
//   - Clear a field (set it to null)
//   - Update a field (set it to a value)
//
// With standard Go types, you can't distinguish these cases:
//
//	type Request struct {
//	    Name *string `json:"name"` // nil means... unset? or null?
//	}
//
// With null.Value, you can:
//
//	type Request struct {
//	    Name null.Value[string] `json:"name"`
//	}
//
//	// Client sends: {}
//	// r.Name.IsSet() == false (field was absent)
//
//	// Client sends: {"name": null}
//	// r.Name.IsSet() == true, r.Name.IsNull() == true
//
//	// Client sends: {"name": "Alice"}
//	// r.Name.IsSet() == true, r.Name.IsValid() == true
//
//	// Client sends: {"name": ""}
//	// r.Name.IsSet() == true, r.Name.IsValid() == true (empty string IS a value)
//
// # Quick Start
//
// Creating values:
//
//	// Valid value
//	name := null.New("Alice")
//
//	// Explicit null
//	name := null.NewNull[string]()
//
//	// From pointer (nil becomes null)
//	name := null.NewPtr(namePtr)
//
//	// Unset (zero value)
//	var name null.Value[string]
//
// Checking state:
//
//	if name.IsSet() {
//	    // Field was present (either null or a value)
//	}
//
//	if name.IsNull() {
//	    // Field was explicitly set to null
//	}
//
//	if name.IsValid() {
//	    // Field has a real value
//	}
//
// Extracting values:
//
//	v := name.Get()          // Returns zero value if not valid
//	v := name.GetOr("Bob")   // Returns "Bob" if not valid
//	p := name.Ptr()          // Returns nil if not valid
//
// # JSON Integration
//
// Value[T] implements json.Marshaler and json.Unmarshaler with full three-state
// support:
//
//	type User struct {
//	    Name  null.Value[string] `json:"name"`
//	    Email null.Value[string] `json:"email,omitempty"`
//	}
//
//	// Unmarshal distinguishes all three states
//	json.Unmarshal([]byte(`{"name": "Alice"}`), &u)
//	// u.Name.IsValid() == true, u.Email.IsSet() == false
//
//	json.Unmarshal([]byte(`{"name": null}`), &u)
//	// u.Name.IsNull() == true
//
// Note: When marshaling, both unset and null values produce "null" in JSON
// (JSON has no concept of "unset"). Use omitempty to omit unset fields.
//
// # SQL Integration
//
// Value[T] implements database/sql.Scanner and database/sql/driver.Valuer:
//
//	type User struct {
//	    ID   int64
//	    Name null.Value[string]
//	}
//
//	// Scanning NULL from database
//	row.Scan(&u.ID, &u.Name)
//	// If column is NULL: u.Name.IsNull() == true
//	// If column has value: u.Name.IsValid() == true
//
//	// Inserting NULL into database
//	db.Exec("INSERT INTO users (name) VALUES ($1)", null.NewNull[string]())
//
// Supported SQL types: string, int/int8/int16/int32/int64, uint/uint8/uint16/uint32/uint64,
// float32/float64, bool, time.Time, []byte.
//
// # DynamoDB Integration
//
// For DynamoDB support, use the nullddb subpackage which wraps Value[T] with
// DynamoDB marshaling:
//
//	import "github.com/bjaus/null/nullddb"
//
//	type Item struct {
//	    PK   string                `dynamodbav:"pk"`
//	    Name nullddb.Value[string] `dynamodbav:"name"`
//	}
//
//	item := Item{PK: "123", Name: nullddb.New("Alice")}
//	av, _ := attributevalue.MarshalMap(item)
//
// The nullddb.Value[T] embeds null.Value[T], so all methods are available.
// Convert between them with nullddb.From():
//
//	apiVal := null.New("Alice")
//	ddbVal := nullddb.From(apiVal)
//
// # State Semantics
//
// Value[T] has exactly three states:
//
//	State    | IsSet() | IsNull() | IsValid() | Get()
//	---------|---------|----------|-----------|-------------
//	Unset    | false   | false    | false     | zero value
//	Null     | true    | true     | false     | zero value
//	Valid    | true    | false    | true      | the value
//
// The zero value of Value[T] is Unset, which is the natural state for struct
// fields that weren't explicitly initialized.
//
// # Common Patterns
//
// PATCH request handling:
//
//	func UpdateUser(req UpdateRequest) error {
//	    if req.Name.IsSet() {
//	        if req.Name.IsNull() {
//	            user.Name = nil  // Clear the field
//	        } else {
//	            user.Name = req.Name.Ptr()  // Update the field
//	        }
//	    }
//	    // If !req.Name.IsSet(), leave user.Name unchanged
//	}
//
// Default values:
//
//	config := Config{
//	    Timeout: settings.Timeout.GetOr(30 * time.Second),
//	    Retries: settings.Retries.GetOr(3),
//	}
//
// Conditional queries:
//
//	var conditions []string
//	var args []any
//
//	if filter.Status.IsSet() {
//	    if filter.Status.IsNull() {
//	        conditions = append(conditions, "status IS NULL")
//	    } else {
//	        conditions = append(conditions, "status = $1")
//	        args = append(args, filter.Status.Get())
//	    }
//	}
//
// # Design Decisions
//
// Why not use pointers (*string)?
//   - Pointers can't distinguish between "not set" and "set to nil"
//   - Pointers require heap allocation
//   - Pointers are awkward for literals: you can't write &"hello"
//
// Why not use sql.NullString and friends?
//   - They don't distinguish between unset and null
//   - They marshal to {"String": "value", "Valid": true} in JSON
//   - No generics (separate type for each primitive)
//
// Why a state enum instead of two bools?
//   - Cleaner semantics (impossible to have invalid state combinations)
//   - Same memory footprint after struct padding
//   - Easier to reason about
//
// Why Get() instead of Value()?
//   - Value() conflicts with driver.Valuer interface method
//   - Get() is clear and concise
package null
