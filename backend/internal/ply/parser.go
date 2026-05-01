package ply

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/vinneyto/splatmaker/backend/internal/model"
)

// ParseASCII parses a minimal ASCII PLY with vertex x y z as first 3 scalar columns.
func ParseASCII(r io.Reader) ([]model.Vec3, error) {
	s := bufio.NewScanner(r)
	vertexCount := -1
	inHeader := true
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if !inHeader {
			break
		}
		if strings.HasPrefix(line, "format ") && line != "format ascii 1.0" {
			return nil, fmt.Errorf("only ascii ply is supported")
		}
		if strings.HasPrefix(line, "element vertex ") {
			parts := strings.Fields(line)
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid element vertex line")
			}
			n, err := strconv.Atoi(parts[2])
			if err != nil {
				return nil, err
			}
			vertexCount = n
		}
		if line == "end_header" {
			inHeader = false
			break
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if inHeader || vertexCount < 0 {
		return nil, fmt.Errorf("invalid ply header")
	}

	points := make([]model.Vec3, 0, vertexCount)
	for s.Scan() && len(points) < vertexCount {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid vertex line: %q", line)
		}
		x, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, err
		}
		y, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, err
		}
		z, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return nil, err
		}
		points = append(points, model.Vec3{X: x, Y: y, Z: z})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(points) != vertexCount {
		return nil, fmt.Errorf("vertex count mismatch: expected %d, got %d", vertexCount, len(points))
	}
	return points, nil
}
