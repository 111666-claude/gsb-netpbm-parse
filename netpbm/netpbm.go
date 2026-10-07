// Package netpbm 解析 Netpbm 图像。
package netpbm

import (
	"errors"
	"strconv"
	"strings"
)

// Counter 记录比较次数。
type Counter struct {
	Scanned int
}

// Image 是一张图像。
type Image struct {
	Magic   string
	Width   int
	Height  int
	MaxVal  int
	Samples []int
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\v' || b == '\f'
}

func nextToken(data []byte, pos *int, c *Counter) string {
	i := *pos
	for i < len(data) && isSpace(data[i]) {
		i++
	}
	start := i
	for i < len(data) && !isSpace(data[i]) {
		i++
	}
	*pos = i
	c.Scanned++
	return string(data[start:i])
}

// Parse 解析一张图像。
func Parse(data []byte, c *Counter) (*Image, error) {
	pos := 0
	magic := nextToken(data, &pos, c)
	if magic == "" {
		return nil, errors.New("short-header")
	}
	img := &Image{Magic: magic}
	width := nextToken(data, &pos, c)
	height := nextToken(data, &pos, c)
	if width == "" || height == "" {
		return nil, errors.New("short-header")
	}
	img.Width, _ = strconv.Atoi(width)
	img.Height, _ = strconv.Atoi(height)
	img.MaxVal = 1
	if img.Magic != "P1" && img.Magic != "P4" {
		maxval := nextToken(data, &pos, c)
		if maxval == "" {
			return nil, errors.New("short-header")
		}
		img.MaxVal, _ = strconv.Atoi(maxval)
	}
	off := pos
	for off < len(data) && isSpace(data[off]) {
		off++
	}
	count := img.Width * img.Height
	switch img.Magic {
	case "P1", "P2", "P3":
		for i := 0; i < count; i++ {
			parts := strings.Fields(string(data[off:]))
			c.Scanned += len(parts)
			if i >= len(parts) {
				break
			}
			value, _ := strconv.Atoi(parts[i])
			img.Samples = append(img.Samples, value)
		}
	default:
		body := data[off:]
		if img.Magic == "P4" {
			for i := 0; i < count; i++ {
				c.Scanned++
				index := i / 8
				if index >= len(body) {
					break
				}
				bit := (body[index] >> uint(7-i%8)) & 1
				img.Samples = append(img.Samples, int(bit))
			}
			return img, nil
		}
		for i := 0; i < count; i++ {
			c.Scanned++
			if i >= len(body) {
				break
			}
			img.Samples = append(img.Samples, int(body[i]))
		}
	}
	return img, nil
}
