package netpbm

import "testing"

func TestHeader(t *testing.T) {
	data := []byte("P2\n2 2\n255\n1 2 3 4\n")
	img, err := Parse(data, &Counter{})
	if err != nil {
		t.Fatalf("解析失败：%v", err)
	}
	if img.Width != 2 || img.Height != 2 || img.MaxVal != 255 {
		t.Fatalf("头不对：%+v", img)
	}
}

func TestAsciiSamples(t *testing.T) {
	data := []byte("P2\n2 1\n255\n3 4\n")
	img, _ := Parse(data, &Counter{})
	if len(img.Samples) != 2 || img.Samples[0] != 3 {
		t.Fatalf("样本不对：%+v", img.Samples)
	}
}

func TestShortHeader(t *testing.T) {
	if _, err := Parse([]byte("P2\n"), &Counter{}); err == nil {
		t.Fatal("头不完整应报错")
	}
}

func TestBinarySingleByte(t *testing.T) {
	data := append([]byte("P5\n2 1\n255\n"), 7, 8)
	img, _ := Parse(data, &Counter{})
	if len(img.Samples) != 2 || img.Samples[0] != 7 {
		t.Fatalf("样本不对：%+v", img.Samples)
	}
}

func TestMagicStored(t *testing.T) {
	img, _ := Parse([]byte("P2\n1 1\n255\n1\n"), &Counter{})
	if img.Magic != "P2" {
		t.Fatalf("魔数应是 P2：%q", img.Magic)
	}
}

func TestCounterCounts(t *testing.T) {
	c := &Counter{}
	Parse([]byte("P2\n1 1\n255\n1\n"), c)
	if c.Scanned == 0 {
		t.Fatal("应统计比较次数")
	}
}
