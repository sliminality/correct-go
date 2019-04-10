// Package assert provides helpers for testing.
package assert

import "testing"

type Any = interface{}

func Nil(t *testing.T) func(p *string, v ...Any) {
	return func(p *string, v ...Any) {
		if p != nil {
			yikes(t, "Expected nil pointer", v)
		}
	}
}

func NotNil(t *testing.T) func(p *string, v ...Any) {
	return func(p *string, v ...Any) {
		if p == nil {
			yikes(t, "Expected non-nil pointer", v)
		}
	}
}

func True(t *testing.T) func(cond bool, v ...Any) {
	return func(cond bool, v ...Any) {
		if !cond {
			yikes(t, "Expected true, got false", v)
		}
	}
}

func False(t *testing.T) func(cond bool, v ...Any) {
	return func(cond bool, v ...Any) {
		if cond {
			yikes(t, "Expected false, got true", v)
		}
	}
}

func Truthy(t *testing.T) func(cond Any, v ...Any) {
	return func(cond Any, v ...Any) {
		if !cond.(bool) {
			yikes(t, "Expected value to be truthy", v)
		}
	}
}

func Equal(t *testing.T) func(a Any, b Any, v ...Any) {
	return func(a Any, b Any, v ...Any) {
		if a != b {
			yikes(t, "Expected values to be equal", v)
		}
	}
}

func yikes(t *testing.T, why string, v ...Any) {
	a := []Any{"Assertion failed:"}
	if len(v) == 0 {
		a = append(a, why)
	} else {
		a = append(a, v...)
	}
	t.Error(a...)
}
