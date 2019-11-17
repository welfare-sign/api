package jwt

import (
	"fmt"
	"testing"
)

func TestCreateToken(t *testing.T) {
	got, err := CreateToken(1, "NPC", "")
	fmt.Println(got, err)
}
