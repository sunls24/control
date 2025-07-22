package adb

import (
	"bytes"
	"control/mobile"
	"fmt"
	"os/exec"
	"strconv"
)

var _ mobile.Driver = (*adbDriver)(nil)

type adbDriver struct {
}

func New() (mobile.Driver, error) {
	return &adbDriver{}, nil
}

func (a adbDriver) CheckHealth() error {
	//TODO implement me
	panic("implement me")
}

func (a adbDriver) Scale() int {
	return 1
}

func (a adbDriver) Screenshot() ([]byte, error) {
	var buf bytes.Buffer
	cmd := exec.Command("adb", "exec-out", "screencap", "-p")
	cmd.Stdout = &buf
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("screenshot: %v, output: %s", err, stderr.String())
	}
	return buf.Bytes(), nil
}

func (a adbDriver) Tap(x, y int) error {
	return adbShell("input", "tap", strconv.Itoa(x), strconv.Itoa(y))
}

func (a adbDriver) Unlock() error {
	//TODO implement me
	panic("implement me")
}

func (a adbDriver) Lock() error {
	//TODO implement me
	panic("implement me")
}

func (a adbDriver) AppLaunch(p string) error {
	return adbShell("monkey", "-p", p, "-c", "android.intent.category.LAUNCHER", "1")
}

func (a adbDriver) AppTerminate(p string) error {
	return adbShell("am", "force-stop", p)
}

func adbShell(args ...string) error {
	cmd := exec.Command("adb", append([]string{"shell"}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("adb: %s, output: %s", err, string(output))
	}
	return nil
}

func keyInput(key string) error {
	return adbShell("input", "keyevent", key)
}
