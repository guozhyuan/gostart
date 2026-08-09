package pkg_test

import (
	"fmt"
	"myserver/internal/pkg"
	"testing"
)

func TestCompare(t *testing.T) {
	s, _ := pkg.Encrypt("guo")
	fmt.Println("Encrypt pwd: ", s)
	b := pkg.Compare("$10$0BX4g0YCr9tANo8i4n2naOa2FyPPnPEiN3GqxuZGTtjUketXyFhoe", "guo")
	if !b {
		t.Errorf("Compare failed")
	} else {
		t.Logf("Compare success")
	}
}
