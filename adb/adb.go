package adb

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Screenshot(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	cmd := exec.Command("adb", "exec-out", "screencap", "-p")
	cmd.Stdout = f
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("screenshot: %v, %s", err, stderr.String())
	}
	return nil
}

func Click(x, y int) error {
	return adb("shell", "input", "tap", strconv.Itoa(x), strconv.Itoa(y))
}

func Unlock() error {
	return sh(`set -e
adb shell dumpsys window policy | grep SCREEN_STATE_OFF && adb shell input keyevent 26 && sleep 0.5
adb shell dumpsys window policy | grep showing=true && adb shell input swipe 540 1800 540 800 300
exit 0`)
}

func OpenApp(p string) error {
	_ = StopApp(p)
	return adb("shell", "monkey", "-p", p, "-c", "android.intent.category.LAUNCHER", "1")
}

func StopApp(p string) error {
	return adb("shell", "am", "force-stop", p)
}

const (
	Back  = "4"
	Power = "26"
)

func KeyInput(k string) error {
	return adb("shell", "input", "keyevent", k)
}

func adb(args ...string) error {
	slog.Debug(fmt.Sprintf("adb %s", strings.Join(args, " ")))
	cmd := exec.Command("adb", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("adb: %s, output:\n%s", err, string(output))
	}
	return nil
}

func sh(sh string) error {
	cmd := exec.Command("sh", "-c", sh)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sh: %s, output:\n%s", err, string(output))
	}
	return nil
}
