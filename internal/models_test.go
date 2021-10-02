package internal

import "testing"

func TrimSpaceTest(t *testing.T) {
	testCases := []string {
		"HEllo  ",
		" World  ",
		"   golang  ",
	}
	
	for _, tC := range testCases {
		t.Run(tC, func(t *testing.T) {

		})
	}
}
