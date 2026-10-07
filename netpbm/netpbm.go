// Package netpbm 解析 Netpbm 图像。
package netpbm

import (
	"errors"
	"strconv"
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

// nextToken 从头位置跳过空白与注释（# 到行尾），取下一个记号。
func nextToken(data []byte, pos *int, c *Counter) (string, bool) {
	i := *pos
	for i < len(data) {
		ch := data[i]
		c.Scanned++
		if isSpace(ch) {
			i++
			continue
		}
		if ch == '#' {
			for i < len(data) && data[i] != '\n' {
				c.Scanned++
				i++
			}
			continue
		}
		break
	}
	start := i
	for i < len(data) {
		c.Scanned++
		if isSpace(data[i]) {
			break
		}
		i++
	}
	*pos = i
	return string(data[start:i]), i > start
}

// Parse 解析一张图像。
func Parse(data []byte, c *Counter) (*Image, error) {
	if c == nil {
		c = &Counter{}
	}
	pos := 0
	magic, ok := nextToken(data, &pos, c)
	if !ok {
		return nil, errors.New("short-header")
	}
	img := &Image{Magic: magic, MaxVal: 1}

	samplesPerPixel := 0
	switch magic {
	case "P1", "P2", "P3", "P4", "P5", "P6":
		if magic == "P3" || magic == "P6" {
			samplesPerPixel = 3
		} else {
			samplesPerPixel = 1
		}
	default:
		return nil, errors.New("bad-magic")
	}

	widthTok, okW := nextToken(data, &pos, c)
	heightTok, okH := nextToken(data, &pos, c)
	if !okW || !okH {
		return nil, errors.New("short-header")
	}
	img.Width, _ = strconv.Atoi(widthTok)
	img.Height, _ = strconv.Atoi(heightTok)

	if magic != "P1" && magic != "P4" {
		maxvalTok, ok := nextToken(data, &pos, c)
		if !ok {
			return nil, errors.New("short-header")
		}
		img.MaxVal, _ = strconv.Atoi(maxvalTok)
	}

	count := img.Width * img.Height * samplesPerPixel

	switch magic {
	case "P1", "P2", "P3":
		for i := 0; i < count; i++ {
			tok, ok := nextToken(data, &pos, c)
			if !ok {
				return img, nil
			}
			value, _ := strconv.Atoi(tok)
			img.Samples = append(img.Samples, value)
		}
	case "P4":
		if pos < len(data) && isSpace(data[pos]) {
			pos++
		}
		body := data[pos:]
		rowBytes := (img.Width + 7) / 8
		for row := 0; row < img.Height; row++ {
			for col := 0; col < img.Width; col++ {
				index := row*rowBytes + col/8
				if index >= len(body) {
					return img, nil
				}
				bit := (body[index] >> uint(7-col%8)) & 1
				img.Samples = append(img.Samples, int(bit))
			}
		}
		c.Scanned += rowBytes * img.Height
	case "P5", "P6":
		if pos < len(data) && isSpace(data[pos]) {
			pos++
		}
		body := data[pos:]
		bytesPerSample := 1
		if img.MaxVal >= 256 {
			bytesPerSample = 2
		}
		for i := 0; i < count; i++ {
			offset := i * bytesPerSample
			if offset+bytesPerSample > len(body) {
				return img, nil
			}
			value := int(body[offset])
			if bytesPerSample == 2 {
				value = value<<8 | int(body[offset+1])
			}
			img.Samples = append(img.Samples, value)
		}
		c.Scanned += count * bytesPerSample
	}
	return img, nil
}
