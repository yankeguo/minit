package rg

import (
	"fmt"
)

// Guard recover from panic and set err
func Guard(err *error) {
	if r := recover(); r != nil {
		if re, ok := r.(error); ok {
			*err = re
		} else {
			*err = fmt.Errorf("panic: %v", r)
		}
	}
}

// Must0 panic err if not nil
func Must0(err error) {
	if err != nil {
		panic(err)
	}
}

// Must panic err if not nil, else return remaining values
func Must[T any](v T, err error) T {
	if err == nil {
		return v
	} else {
		panic(err)
	}
}

// Must2 panic err if not nil, else return remaining values
func Must2[T1 any, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}

// Must3 panic err if not nil, else return remaining values
func Must3[T1 any, T2 any, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3
}

// Must4 panic err if not nil, else return remaining values
func Must4[T1 any, T2 any, T3 any, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, err error) (T1, T2, T3, T4) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3, v4
}

// Must5 panic err if not nil, else return remaining values
func Must5[T1 any, T2 any, T3 any, T4 any, T5 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, err error) (T1, T2, T3, T4, T5) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3, v4, v5
}

// Must6 panic err if not nil, else return remaining values
func Must6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, err error) (T1, T2, T3, T4, T5, T6) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3, v4, v5, v6
}

// Must7 panic err if not nil, else return remaining values
func Must7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, err error) (T1, T2, T3, T4, T5, T6, T7) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3, v4, v5, v6, v7
}
