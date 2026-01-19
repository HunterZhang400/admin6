package shaid

import (
	"fmt"
	"testing"
)

func TestID(t *testing.T) {
	v := GetSHAID(14626)
	fmt.Println(v)
}
