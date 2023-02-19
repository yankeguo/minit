package gg

type M = map[string]interface{}

// __BEGIN_OF_GENERATION__

// F00 function with 0 arguments and 0 returns
type F00 func()

func (f F00) Do() {
	f()
}

// D00 interface with a single method 0 with arguments and 0 returns
type D00 interface {
	Do()
}

// F01 function with 0 arguments and 1 returns
type F01[O1 any] func() (o1 O1)

func (f F01[O1]) Do() (o1 O1) {
	return f()
}

// D01 interface with a single method 0 with arguments and 1 returns
type D01[O1 any] interface {
	Do() (o1 O1)
}

// F02 function with 0 arguments and 2 returns
type F02[O1 any, O2 any] func() (o1 O1, o2 O2)

func (f F02[O1, O2]) Do() (o1 O1, o2 O2) {
	return f()
}

// D02 interface with a single method 0 with arguments and 2 returns
type D02[O1 any, O2 any] interface {
	Do() (o1 O1, o2 O2)
}

// F03 function with 0 arguments and 3 returns
type F03[O1 any, O2 any, O3 any] func() (o1 O1, o2 O2, o3 O3)

func (f F03[O1, O2, O3]) Do() (o1 O1, o2 O2, o3 O3) {
	return f()
}

// D03 interface with a single method 0 with arguments and 3 returns
type D03[O1 any, O2 any, O3 any] interface {
	Do() (o1 O1, o2 O2, o3 O3)
}

// F04 function with 0 arguments and 4 returns
type F04[O1 any, O2 any, O3 any, O4 any] func() (o1 O1, o2 O2, o3 O3, o4 O4)

func (f F04[O1, O2, O3, O4]) Do() (o1 O1, o2 O2, o3 O3, o4 O4) {
	return f()
}

// D04 interface with a single method 0 with arguments and 4 returns
type D04[O1 any, O2 any, O3 any, O4 any] interface {
	Do() (o1 O1, o2 O2, o3 O3, o4 O4)
}

// F10 function with 1 arguments and 0 returns
type F10[I1 any] func(i1 I1)

func (f F10[I1]) Do(i1 I1) {
	f(i1)
}

// D10 interface with a single method 1 with arguments and 0 returns
type D10[I1 any] interface {
	Do(i1 I1)
}

// F11 function with 1 arguments and 1 returns
type F11[I1 any, O1 any] func(i1 I1) (o1 O1)

func (f F11[I1, O1]) Do(i1 I1) (o1 O1) {
	return f(i1)
}

// D11 interface with a single method 1 with arguments and 1 returns
type D11[I1 any, O1 any] interface {
	Do(i1 I1) (o1 O1)
}

// F12 function with 1 arguments and 2 returns
type F12[I1 any, O1 any, O2 any] func(i1 I1) (o1 O1, o2 O2)

func (f F12[I1, O1, O2]) Do(i1 I1) (o1 O1, o2 O2) {
	return f(i1)
}

// D12 interface with a single method 1 with arguments and 2 returns
type D12[I1 any, O1 any, O2 any] interface {
	Do(i1 I1) (o1 O1, o2 O2)
}

// F13 function with 1 arguments and 3 returns
type F13[I1 any, O1 any, O2 any, O3 any] func(i1 I1) (o1 O1, o2 O2, o3 O3)

func (f F13[I1, O1, O2, O3]) Do(i1 I1) (o1 O1, o2 O2, o3 O3) {
	return f(i1)
}

// D13 interface with a single method 1 with arguments and 3 returns
type D13[I1 any, O1 any, O2 any, O3 any] interface {
	Do(i1 I1) (o1 O1, o2 O2, o3 O3)
}

// F14 function with 1 arguments and 4 returns
type F14[I1 any, O1 any, O2 any, O3 any, O4 any] func(i1 I1) (o1 O1, o2 O2, o3 O3, o4 O4)

func (f F14[I1, O1, O2, O3, O4]) Do(i1 I1) (o1 O1, o2 O2, o3 O3, o4 O4) {
	return f(i1)
}

// D14 interface with a single method 1 with arguments and 4 returns
type D14[I1 any, O1 any, O2 any, O3 any, O4 any] interface {
	Do(i1 I1) (o1 O1, o2 O2, o3 O3, o4 O4)
}

// F20 function with 2 arguments and 0 returns
type F20[I1 any, I2 any] func(i1 I1, i2 I2)

func (f F20[I1, I2]) Do(i1 I1, i2 I2) {
	f(i1, i2)
}

// D20 interface with a single method 2 with arguments and 0 returns
type D20[I1 any, I2 any] interface {
	Do(i1 I1, i2 I2)
}

// T2 tuple with 2 fields
type T2[I1 any, I2 any] struct {
	A I1
	B I2
}

// F21 function with 2 arguments and 1 returns
type F21[I1 any, I2 any, O1 any] func(i1 I1, i2 I2) (o1 O1)

func (f F21[I1, I2, O1]) Do(i1 I1, i2 I2) (o1 O1) {
	return f(i1, i2)
}

// D21 interface with a single method 2 with arguments and 1 returns
type D21[I1 any, I2 any, O1 any] interface {
	Do(i1 I1, i2 I2) (o1 O1)
}

// F22 function with 2 arguments and 2 returns
type F22[I1 any, I2 any, O1 any, O2 any] func(i1 I1, i2 I2) (o1 O1, o2 O2)

func (f F22[I1, I2, O1, O2]) Do(i1 I1, i2 I2) (o1 O1, o2 O2) {
	return f(i1, i2)
}

// D22 interface with a single method 2 with arguments and 2 returns
type D22[I1 any, I2 any, O1 any, O2 any] interface {
	Do(i1 I1, i2 I2) (o1 O1, o2 O2)
}

// F23 function with 2 arguments and 3 returns
type F23[I1 any, I2 any, O1 any, O2 any, O3 any] func(i1 I1, i2 I2) (o1 O1, o2 O2, o3 O3)

func (f F23[I1, I2, O1, O2, O3]) Do(i1 I1, i2 I2) (o1 O1, o2 O2, o3 O3) {
	return f(i1, i2)
}

// D23 interface with a single method 2 with arguments and 3 returns
type D23[I1 any, I2 any, O1 any, O2 any, O3 any] interface {
	Do(i1 I1, i2 I2) (o1 O1, o2 O2, o3 O3)
}

// F24 function with 2 arguments and 4 returns
type F24[I1 any, I2 any, O1 any, O2 any, O3 any, O4 any] func(i1 I1, i2 I2) (o1 O1, o2 O2, o3 O3, o4 O4)

func (f F24[I1, I2, O1, O2, O3, O4]) Do(i1 I1, i2 I2) (o1 O1, o2 O2, o3 O3, o4 O4) {
	return f(i1, i2)
}

// D24 interface with a single method 2 with arguments and 4 returns
type D24[I1 any, I2 any, O1 any, O2 any, O3 any, O4 any] interface {
	Do(i1 I1, i2 I2) (o1 O1, o2 O2, o3 O3, o4 O4)
}

// F30 function with 3 arguments and 0 returns
type F30[I1 any, I2 any, I3 any] func(i1 I1, i2 I2, i3 I3)

func (f F30[I1, I2, I3]) Do(i1 I1, i2 I2, i3 I3) {
	f(i1, i2, i3)
}

// D30 interface with a single method 3 with arguments and 0 returns
type D30[I1 any, I2 any, I3 any] interface {
	Do(i1 I1, i2 I2, i3 I3)
}

// T3 tuple with 3 fields
type T3[I1 any, I2 any, I3 any] struct {
	A I1
	B I2
	C I3
}

// F31 function with 3 arguments and 1 returns
type F31[I1 any, I2 any, I3 any, O1 any] func(i1 I1, i2 I2, i3 I3) (o1 O1)

func (f F31[I1, I2, I3, O1]) Do(i1 I1, i2 I2, i3 I3) (o1 O1) {
	return f(i1, i2, i3)
}

// D31 interface with a single method 3 with arguments and 1 returns
type D31[I1 any, I2 any, I3 any, O1 any] interface {
	Do(i1 I1, i2 I2, i3 I3) (o1 O1)
}

// F32 function with 3 arguments and 2 returns
type F32[I1 any, I2 any, I3 any, O1 any, O2 any] func(i1 I1, i2 I2, i3 I3) (o1 O1, o2 O2)

func (f F32[I1, I2, I3, O1, O2]) Do(i1 I1, i2 I2, i3 I3) (o1 O1, o2 O2) {
	return f(i1, i2, i3)
}

// D32 interface with a single method 3 with arguments and 2 returns
type D32[I1 any, I2 any, I3 any, O1 any, O2 any] interface {
	Do(i1 I1, i2 I2, i3 I3) (o1 O1, o2 O2)
}

// F33 function with 3 arguments and 3 returns
type F33[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any] func(i1 I1, i2 I2, i3 I3) (o1 O1, o2 O2, o3 O3)

func (f F33[I1, I2, I3, O1, O2, O3]) Do(i1 I1, i2 I2, i3 I3) (o1 O1, o2 O2, o3 O3) {
	return f(i1, i2, i3)
}

// D33 interface with a single method 3 with arguments and 3 returns
type D33[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any] interface {
	Do(i1 I1, i2 I2, i3 I3) (o1 O1, o2 O2, o3 O3)
}

// F34 function with 3 arguments and 4 returns
type F34[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any, O4 any] func(i1 I1, i2 I2, i3 I3) (o1 O1, o2 O2, o3 O3, o4 O4)

func (f F34[I1, I2, I3, O1, O2, O3, O4]) Do(i1 I1, i2 I2, i3 I3) (o1 O1, o2 O2, o3 O3, o4 O4) {
	return f(i1, i2, i3)
}

// D34 interface with a single method 3 with arguments and 4 returns
type D34[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any, O4 any] interface {
	Do(i1 I1, i2 I2, i3 I3) (o1 O1, o2 O2, o3 O3, o4 O4)
}

// F40 function with 4 arguments and 0 returns
type F40[I1 any, I2 any, I3 any, I4 any] func(i1 I1, i2 I2, i3 I3, i4 I4)

func (f F40[I1, I2, I3, I4]) Do(i1 I1, i2 I2, i3 I3, i4 I4) {
	f(i1, i2, i3, i4)
}

// D40 interface with a single method 4 with arguments and 0 returns
type D40[I1 any, I2 any, I3 any, I4 any] interface {
	Do(i1 I1, i2 I2, i3 I3, i4 I4)
}

// T4 tuple with 4 fields
type T4[I1 any, I2 any, I3 any, I4 any] struct {
	A I1
	B I2
	C I3
	D I4
}

// F41 function with 4 arguments and 1 returns
type F41[I1 any, I2 any, I3 any, I4 any, O1 any] func(i1 I1, i2 I2, i3 I3, i4 I4) (o1 O1)

func (f F41[I1, I2, I3, I4, O1]) Do(i1 I1, i2 I2, i3 I3, i4 I4) (o1 O1) {
	return f(i1, i2, i3, i4)
}

// D41 interface with a single method 4 with arguments and 1 returns
type D41[I1 any, I2 any, I3 any, I4 any, O1 any] interface {
	Do(i1 I1, i2 I2, i3 I3, i4 I4) (o1 O1)
}

// F42 function with 4 arguments and 2 returns
type F42[I1 any, I2 any, I3 any, I4 any, O1 any, O2 any] func(i1 I1, i2 I2, i3 I3, i4 I4) (o1 O1, o2 O2)

func (f F42[I1, I2, I3, I4, O1, O2]) Do(i1 I1, i2 I2, i3 I3, i4 I4) (o1 O1, o2 O2) {
	return f(i1, i2, i3, i4)
}

// D42 interface with a single method 4 with arguments and 2 returns
type D42[I1 any, I2 any, I3 any, I4 any, O1 any, O2 any] interface {
	Do(i1 I1, i2 I2, i3 I3, i4 I4) (o1 O1, o2 O2)
}

// F43 function with 4 arguments and 3 returns
type F43[I1 any, I2 any, I3 any, I4 any, O1 any, O2 any, O3 any] func(i1 I1, i2 I2, i3 I3, i4 I4) (o1 O1, o2 O2, o3 O3)

func (f F43[I1, I2, I3, I4, O1, O2, O3]) Do(i1 I1, i2 I2, i3 I3, i4 I4) (o1 O1, o2 O2, o3 O3) {
	return f(i1, i2, i3, i4)
}

// D43 interface with a single method 4 with arguments and 3 returns
type D43[I1 any, I2 any, I3 any, I4 any, O1 any, O2 any, O3 any] interface {
	Do(i1 I1, i2 I2, i3 I3, i4 I4) (o1 O1, o2 O2, o3 O3)
}

// F44 function with 4 arguments and 4 returns
type F44[I1 any, I2 any, I3 any, I4 any, O1 any, O2 any, O3 any, O4 any] func(i1 I1, i2 I2, i3 I3, i4 I4) (o1 O1, o2 O2, o3 O3, o4 O4)

func (f F44[I1, I2, I3, I4, O1, O2, O3, O4]) Do(i1 I1, i2 I2, i3 I3, i4 I4) (o1 O1, o2 O2, o3 O3, o4 O4) {
	return f(i1, i2, i3, i4)
}

// D44 interface with a single method 4 with arguments and 4 returns
type D44[I1 any, I2 any, I3 any, I4 any, O1 any, O2 any, O3 any, O4 any] interface {
	Do(i1 I1, i2 I2, i3 I3, i4 I4) (o1 O1, o2 O2, o3 O3, o4 O4)
}
