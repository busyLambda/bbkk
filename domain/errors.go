package domain

import "fmt"

type TooLong struct {
	Expected int
	Found    int
}

func (tl *TooLong) Error() string {
	return fmt.Sprintf("ERR TooLong: wanted %d got %d", tl.Expected, tl.Found)
}

type TooShort struct{}

func (ts *TooShort) Error() string {
	return fmt.Sprintf("ERR TooShort")
}
