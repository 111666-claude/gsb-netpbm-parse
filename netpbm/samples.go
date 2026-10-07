package netpbm

import (
	"fmt"
	"strings"
)

// Scenario 是一个固定场景。
type Scenario struct {
	Name string
	Run  func() map[string]any
}

// Samples 返回全部场景。
func Samples() []Scenario {
	return []Scenario{
		{Name: "comment", Run: func() map[string]any {
			data := []byte("P2\n# a comment\n2 2\n255\n1 2 3 4\n")
			img, err := Parse(data, &Counter{})
			if err != nil {
				return map[string]any{"error": err.Error()}
			}
			return map[string]any{"width": img.Width}
		}},
		{Name: "pad", Run: func() map[string]any {
			data := append([]byte("P4\n9 2\n"), 0x80, 0xff, 0x00, 0x00)
			img, _ := Parse(data, &Counter{})
			if len(img.Samples) < 10 {
				return map[string]any{"pixel": -1}
			}
			return map[string]any{"pixel": img.Samples[9]}
		}},
		{Name: "maxval", Run: func() map[string]any {
			data := append([]byte("P5\n1 1\n1000\n"), 0x03, 0xe8)
			img, _ := Parse(data, &Counter{})
			if len(img.Samples) == 0 {
				return map[string]any{"sample": -1}
			}
			return map[string]any{"sample": img.Samples[0]}
		}},
		{Name: "rgb", Run: func() map[string]any {
			data := []byte("P3\n1 1\n255\n10 20 30\n")
			img, _ := Parse(data, &Counter{})
			return map[string]any{"samples": len(img.Samples)}
		}},
	}
}

// Work 跑规模线场景。
func Work(n int) map[string]any {
	var b strings.Builder
	b.WriteString("P2\n")
	b.WriteString(fmt.Sprintf("%d 1\n255\n", n))
	for i := 0; i < n; i++ {
		b.WriteString("1 ")
	}
	c := &Counter{}
	img, _ := Parse([]byte(b.String()), c)
	return map[string]any{"samples": len(img.Samples), "scanned": c.Scanned}
}
