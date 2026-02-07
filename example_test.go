package null_test

import (
	"encoding/json"
	"fmt"

	"github.com/bjaus/null"
)

func Example() {
	type User struct {
		Name  null.Value[string] `json:"name"`
		Email null.Value[string] `json:"email"`
		Age   null.Value[int]    `json:"age"`
	}

	// Simulate a PATCH request with partial data
	jsonData := `{"name": "Alice", "email": null}`

	var user User
	json.Unmarshal([]byte(jsonData), &user)

	fmt.Printf("Name set: %v, value: %q\n", user.Name.IsSet(), user.Name.Get())
	fmt.Printf("Email set: %v, null: %v\n", user.Email.IsSet(), user.Email.IsNull())
	fmt.Printf("Age set: %v\n", user.Age.IsSet())

	// Output:
	// Name set: true, value: "Alice"
	// Email set: true, null: true
	// Age set: false
}

func ExampleNew() {
	v := null.New("hello")
	fmt.Println(v.IsValid(), v.Get())
	// Output: true hello
}

func ExampleNewNull() {
	v := null.NewNull[string]()
	fmt.Println(v.IsSet(), v.IsNull(), v.IsValid())
	// Output: true true false
}

func ExampleNewPtr() {
	s := "hello"
	v1 := null.NewPtr(&s)
	v2 := null.NewPtr[string](nil)

	fmt.Println(v1.IsValid(), v1.Get())
	fmt.Println(v2.IsNull())
	// Output:
	// true hello
	// true
}

func ExampleValue_GetOr() {
	v := null.NewNull[string]()
	fmt.Println(v.GetOr("default"))
	// Output: default
}
