package adb

import (
	"control/helper"
	"control/opencv"
	"fmt"
	"github.com/electricbubble/gadb"
	"os"
	"testing"
	"time"

	adb2 "github.com/prife/goadb"
)

func TestADB(t *testing.T) {
	count := 10
	left, err := os.ReadFile("./left.png")
	helper.Check(err)
	right, err := os.ReadFile("./right.png")
	helper.Check(err)

	driver, err := New()
	helper.Check(err)

	for i := 0; i < count; i++ {
		ts := time.Now()
		start := time.Now()
		s, err := driver.Screenshot()
		helper.Check(err)
		fmt.Println("screenshot:", time.Since(start))
		start = time.Now()
		p, err := opencv.Find(s, left, opencv.WithRange(opencv.Bottom))
		helper.Check(err)
		fmt.Println("find left:", time.Since(start))
		start = time.Now()
		err = driver.Tap(p.X, p.Y)
		helper.Check(err)
		fmt.Println("tap left:", time.Since(start))
		s, err = driver.Screenshot()
		helper.Check(err)
		start = time.Now()
		p2, err := opencv.Find(s, right, opencv.WithRange(opencv.Bottom))
		helper.Check(err)
		fmt.Println("find right:", time.Since(start))
		start = time.Now()
		err = driver.Tap(p2.X, p2.Y)
		helper.Check(err)
		fmt.Println("tap right:", time.Since(start))

		fmt.Println("-----", time.Since(ts))
	}
}

func TestTime(t *testing.T) {
	count := 1

	c, err := gadb.NewClient()
	helper.Check(err)
	ds, err := c.DeviceList()
	helper.Check(err)
	d := ds[0]

	start := time.Now()
	for i := 0; i < count; i++ {
		_, err = d.RunShellCommandWithBytes("screencap", "-p")
		helper.Check(err)
	}
	fmt.Println("screencap:", time.Since(start))

	driver, err := New()
	helper.Check(err)
	start = time.Now()
	for i := 0; i < count; i++ {
		_, err = driver.Screenshot()
		helper.Check(err)
	}
	fmt.Println("screenshot:", time.Since(start))

	client, err := adb2.New()
	helper.Check(err)
	dd := client.Device(adb2.AnyDevice())
	start = time.Now()
	_, err = dd.RunCommand("screencap -p")
	helper.Check(err)
	fmt.Println("screencap2:", time.Since(start))
}
