package wda

import (
	"control/mobile"
	"fmt"
	"github.com/sunls24/gwda"
)

var _ mobile.Driver = (*wdaDriver)(nil)

type wdaDriver struct {
	wd gwda.WebDriver
}

func New() (mobile.Driver, error) {
	wd, err := gwda.NewUSBDriver(nil)
	if err != nil {
		return nil, err
	}
	return &wdaDriver{
		wd: wd,
	}, nil
}

func (w *wdaDriver) CheckHealth() error {
	ok, err := w.wd.IsWdaHealthy()
	if err != nil || !ok {
		return fmt.Errorf("wda not healthy: %v", err)
	}
	return nil
}

func (w *wdaDriver) Scale() int {
	scale, _ := w.wd.Scale()
	if scale <= 0 {
		scale = 1
	}
	return int(scale)
}

func (w *wdaDriver) Screenshot() ([]byte, error) {
	buf, err := w.wd.Screenshot()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (w *wdaDriver) Tap(x, y int) error {
	return w.wd.Tap(x, y)
}

func (w *wdaDriver) Unlock() error {
	lock, err := w.wd.IsLocked()
	if err != nil {
		return err
	}
	if !lock {
		return nil
	}
	err = w.wd.PressButton(gwda.DeviceButtonHome)
	if err != nil {
		return err
	}
	err = w.wd.PressButton(gwda.DeviceButtonHome)
	if err != nil {
		return err
	}
	return w.wd.Homescreen()
}

func (w *wdaDriver) Lock() error {
	return w.wd.Lock()
}

func (w *wdaDriver) AppLaunch(p string) error {
	return w.wd.AppLaunch(p, gwda.NewAppLaunchOption().WithShouldWaitForQuiescence(true))
}

func (w *wdaDriver) AppTerminate(p string) error {
	_, err := w.wd.AppTerminate(p)
	return err
}
