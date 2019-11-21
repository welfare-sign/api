package jwt

import (
	"fmt"
	"testing"
)

func TestCreateToken(t *testing.T) {
	got, _ := CreateToken(3, "NPC", "")
	fmt.Println(got)
}
