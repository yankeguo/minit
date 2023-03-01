package merrs

import (
	"strconv"
	"strings"
	"sync"
)

type Errors []error

func (errs Errors) Error() string {
	sb := &strings.Builder{}
	for i, err := range errs {
		if err == nil {
			continue
		}
		if sb.Len() > 0 {
			sb.WriteString("; ")
		}
		sb.WriteRune('#')
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(": ")
		sb.WriteString(err.Error())
	}
	return sb.String()
}

type ErrorGroup interface {
	Add(err error)
	Set(i int, err error)
	Unwrap() error
}

type errorGroup struct {
	errors Errors
	locker *sync.RWMutex
}

func NewErrorGroup() ErrorGroup {
	return &errorGroup{
		locker: &sync.RWMutex{},
	}
}

func (eg *errorGroup) Add(err error) {
	eg.locker.Lock()
	defer eg.locker.Unlock()

	eg.errors = append(eg.errors, err)
}

func (eg *errorGroup) Set(i int, err error) {
	eg.locker.Lock()
	defer eg.locker.Unlock()

	if i >= len(eg.errors) {
		eg.errors = append(eg.errors, make([]error, i+1-len(eg.errors))...)
	}

	eg.errors[i] = err
}

func (eg *errorGroup) Unwrap() error {
	eg.locker.RLock()
	defer eg.locker.RUnlock()

	for _, err := range eg.errors {
		if err != nil {
			return eg.errors
		}
	}

	return nil
}
