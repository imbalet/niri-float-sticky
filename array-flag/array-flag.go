// Package arrayflag
package arrayflag

import "fmt"

type ArrayFlag []string

func (i *ArrayFlag) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *ArrayFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}
