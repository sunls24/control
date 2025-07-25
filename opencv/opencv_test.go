package opencv

import (
	"github.com/sunls24/gwda"
	"image"
	"os"
	"testing"
	"time"
)

func TestFind(t *testing.T) {
	driver, err := gwda.NewUSBDriver(nil)
	if err != nil {
		t.Fatal(err)
	}
	s, err := driver.Screenshot()
	if err != nil {
		t.Fatal(err)
	}

	tpl := "./yijian.png"
	count := 1

	sb := s.Bytes()

	var p image.Point
	start := time.Now()
	for i := 0; i < count; i++ {
		p, err = Find(sb, read(tpl))
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("All:", p, time.Since(start))

	start = time.Now()
	for i := 0; i < count; i++ {
		p, err = Find(sb, read(tpl), WithRange(Bottom))
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Bottom:", p, time.Since(start))
}

func read(path string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}

func TestScreenshot(t *testing.T) {
	driver, err := gwda.NewUSBDriver(nil)
	if err != nil {
		t.Fatal(err)
	}
	s, err := driver.Screenshot()
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("./screenshot.png", s.Bytes(), 0o644)
	if err != nil {
		t.Fatal(err)
	}
}
