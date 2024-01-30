package internal

import (
	"fmt"
	"strings"
)

type Errors []error

func (es Errors) Format() error {
	var ss []string
	for _, e := range es {
		ss = append(ss, e.Error())
	}
	return fmt.Errorf("%s", strings.Join(ss, "\n"))
}

func (es Errors) HasErrors() bool {
	return len(es) > 0
}
