package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
)

type t struct{}

func (x t) a() error {
	fmt.Println("this method returns an error")      // MATCH /unchecked error/
	fmt.Println("this method also returns an error") // MATCH /unchecked error/
	return errors.New("")
}

type u struct {
	t t
}

func a() error {
	fmt.Println("this function returns an error") // MATCH /unchecked error/
	return errors.New("")
}

func b() (int, error) {
	fmt.Println("this function returns an int and an error") // MATCH /unchecked error/
	return 0, errors.New("")
}

func c() int {
	fmt.Println("this function returns an int") // MATCH /unchecked error/
	return 7
}

func rec() {
	defer func() {
		recover() // MATCH /unchecked error/
		_ = recover()
	}()
	defer recover() // MATCH /unchecked error/
}

func nilError() error {
	println("")
	return nil
}

type MyError string

func (e MyError) Error() string {
	return string(e)
}

func customError() error {
	println() // not pure
	return MyError("an error occurred")
}

type MyPointerError string

func (e *MyPointerError) Error() string {
	return string(*e)
}

func main() {
	// Single error return
	_ = a()
	a() // MATCH /unchecked error/

	// Return another value and an error
	_, _ = b()
	b() // MATCH /unchecked error/

	// Return a custom error type
	_ = customError()
	customError() // MATCH /unchecked error/

	// Method with a single error return
	x := t{}
	_ = x.a()
	x.a() // MATCH /unchecked error/

	// Method call on a struct member
	y := u{x}
	_ = y.t.a()
	y.t.a() // MATCH /unchecked error/

	m1 := map[string]func() error{"a": a}
	_ = m1["a"]()
	m1["a"]() // MATCH /unchecked error/

	// Additional cases for assigning errors to blank identifier
	z, _ := b()
	_, w := a(), 5

	// Assign non error to blank identifier
	_ = c()

	_ = z + w // Avoid complaints about unused variables

	// Goroutine
	go a()    // MATCH /unchecked error/
	defer a() // MATCH /unchecked error/

	b1 := bytes.Buffer{}
	b2 := &bytes.Buffer{}
	b1.Write(nil)
	b2.Write(nil)
	rand.Read(nil)

	ioutil.ReadFile("main.go") // MATCH /unchecked error/

	nilError()

	err := customError() // MATCH /unchecked error/
	err = customError()
	if err != nil {
		println()
	}
}

// MATCH:118 /never used/
