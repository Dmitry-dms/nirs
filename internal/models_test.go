package internal

import (
	"fmt"
	"testing"
)

// func TrimSpaceTest(t *testing.T) {
// 	testCases := []string {
// 		"HEllo  ",
// 		" World  ",
// 		"   golang  ",
// 	}
	
// 	for _, tC := range testCases {
// 		t.Run(tC, func(t *testing.T) {

// 		})
// 	}
// }
func Test(t *testing.T) {
	s, _ := loadSettings("../settings.json")
	fmt.Println(s)
}
