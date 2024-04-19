package cli

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	args := []string{"-name", "test", "-age", "20"}

	m := ParseArgs(args)
	if len(m) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(m))
	}
	if m["name"] != "test" {
		t.Errorf("Expected name to be test, got %s", m["name"])
	}
	if m["age"] != "20" {
		t.Errorf("Expected age to be 20, got %s", m["age"])
	}
}

func TestInputFromModelWithArgs(t *testing.T) {
	t.Run("WithArgs", func(t *testing.T) {

		type User struct {
			Name  string `validate:"required"`
			Age   int    `validate:"required"`
			Other string
		}

		user := User{}
		err := InputFromModel(&user, map[string]string{
			"name": "test",
			"age":  "20",
		})
		if err != nil {
			t.Errorf("Error parsing input: %v", err)
		}
		if user.Name == "" {
			t.Errorf("Name should not be empty")
		}
		if user.Age != 20 {
			t.Errorf("Age should not be 0")
		}
		if user.Other != "" {
			t.Errorf("Other should be empty")
		}
	})

	t.Run("Int", func(t *testing.T) {
		type A struct {
			A int  `validate:"required"`
			B *int `validate:"required"`
		}

		a := A{}
		err := InputFromModel(&a, map[string]string{
			"a": "10",
			"b": "20",
		})
		if err != nil {
			t.Errorf("Error parsing input: %v", err)
		}
		if a.A != 10 {
			t.Errorf("A should not be 0")
		}
		if *a.B != 20 {
			t.Errorf("B should not be 0")
		}
	})

	t.Run("Int wrong type", func(t *testing.T) {

		type User struct {
			Name string `validate:"required"`
			Age  int    `validate:"required"`
		}

		user := User{}
		err := InputFromModel(&user, map[string]string{
			"name": "test",
			"age":  "twenty",
		})
		if err == nil {
			t.Errorf("Expected error parsing int")
		}
	})

	t.Run("String", func(t *testing.T) {
		type A struct {
			A string  `validate:"required"`
			B *string `validate:"required"`
		}

		a := A{}
		err := InputFromModel(&a, map[string]string{
			"a": "10",
			"b": "20",
		})
		if err != nil {
			t.Errorf("Error parsing input: %v", err)
		}
		if a.A != "10" {
			t.Errorf("A should not be empty")
		}
		if *a.B != "20" {
			t.Errorf("B should not be empty")
		}
	})
}
