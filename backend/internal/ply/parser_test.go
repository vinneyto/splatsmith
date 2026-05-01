package ply

import (
	"strings"
	"testing"
)

func TestParseASCII(t *testing.T) {
	data := `ply
format ascii 1.0
element vertex 3
property float x
property float y
property float z
end_header
0 0 0
1 2 3
-1 0.5 4
`
	pts, err := ParseASCII(strings.NewReader(data))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(pts) != 3 {
		t.Fatalf("expected 3 points, got %d", len(pts))
	}
}
