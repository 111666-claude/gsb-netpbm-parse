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

// skipHeader 跳过头中的空白与注释（# 到行尾），返回下一个非空白位置。
func skipHeader(data []byte, pos int) int {
	for pos < len(data) {
		if isSpace(data[pos]) {
			pos++
			continue
		}
		if data[pos] == '#' {
			for pos < len(data) && data[pos] != '\n' && data[pos] != '\r' {
				pos++
			}
			continue
		}
		break
	}
	return pos
}

// nextToken 读取头中的下一个记号，自动跳过其前的空白与注释。
func nextToken(data []byte, pos *int) (string, bool) {
	*pos = skipHeader(data, *pos)
	start := *pos
	for *pos < len(data) && !isSpace(data[*pos]) {
		*pos++
	}
	return string(data[start:*pos]), start < *pos
}

// nextAscii 从 ASCII 栅格中读取一个十进制样本，只向前扫描一次。
func nextAscii(data []byte, pos *int, c *Counter) (int, bool) {
	*pos = skipHeader(data, *pos)
	start := *pos
	for *pos < len(data) && !isSpace(data[*pos]) {
		*pos++
	}
	if start == *pos {
		return 0, false
	}
	c.Scanned++
	value, _ := strconv.Atoi(string(data[start:*pos]))
	return value, true
}

// Parse 解析一张图像。
func Parse(data []byte, c *Counter) (*Image, error) {
	if c == nil {
		c = &Counter{}
	}
	pos := 0
	magic, ok := nextToken(data, &pos)
	if !ok {
		return nil, errors.New("short-header")
	}
	img := &Image{Magic: magic}

	spp := 1
	switch magic {
	case "P1", "P2", "P3", "P5", "P6":
	case "P4":
	default:
		return nil, errors.New("bad-magic")
	}
	if magic == "P3" || magic == "P6" {
		spp = 3
	}

	width, wok := nextToken(data, &pos)
	height, hok := nextToken(data, &pos)
	if !wok || !hok {
		return nil, errors.New("short-header")
	}
	img.Width, _ = strconv.Atoi(width)
	img.Height, _ = strconv.Atoi(height)

	img.MaxVal = 1
	if magic != "P1" && magic != "P4" {
		maxval, mok := nextToken(data, &pos)
		if !mok {
			return nil, errors.New("short-header")
		}
		img.MaxVal, _ = strconv.Atoi(maxval)
	}

	// 二进制格式的头与栅格之间恰好隔一个空白字符（栅格首字节本身可能也是空白）。
	if magic == "P4" || magic == "P5" || magic == "P6" {
		if pos >= len(data) || !isSpace(data[pos]) {
			return nil, errors.New("short-data")
		}
		pos++
	}

	total := img.Width * img.Height * spp
	switch magic {
	case "P1", "P2", "P3":
		img.Samples = make([]int, 0, total)
		for i := 0; i < total; i++ {
			value, ok := nextAscii(data, &pos, c)
			if !ok {
				return nil, errors.New("short-data")
			}
			img.Samples = append(img.Samples, value)
		}
	case "P4":
		stride := (img.Width + 7) / 8
		img.Samples = make([]int, 0, total)
		for row := 0; row < img.Height; row++ {
			for col := 0; col < img.Width; col++ {
				index := row*stride + col/8
				if index >= len(data)-pos {
					return nil, errors.New("short-data")
				}
				c.Scanned++
				bit := (data[pos+index] >> uint(7-col%8)) & 1
				img.Samples = append(img.Samples, int(bit))
			}
		}
	case "P5", "P6":
		widthBytes := 1
		if img.MaxVal >= 256 {
			widthBytes = 2
		}
		img.Samples = make([]int, 0, total)
		for i := 0; i < total; i++ {
			base := pos + i*widthBytes
			if base+widthBytes > len(data) {
				return nil, errors.New("short-data")
			}
			c.Scanned++
			value := int(data[base])
			if widthBytes == 2 {
				value = value<<8 | int(data[base+1])
			}
			img.Samples = append(img.Samples, value)
		}
	}
	return img, nil
}
