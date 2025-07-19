package mobile

import (
	"control/helper"
	"control/opencv"
	"errors"
	"image"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type Driver interface {
	CheckHealth() error
	Scale() int

	Screenshot() ([]byte, error)
	Tap(x, y int) error

	Unlock() error
	Lock() error

	AppLaunch(string) error
	AppTerminate(string) error
}

type Control struct {
	cfg    Config
	driver Driver
	cache  map[string][]byte
	logger *slog.Logger
	scale  int
}

func New(driver Driver, cfg Config, logger *slog.Logger) (*Control, error) {
	if err := driver.CheckHealth(); err != nil {
		return nil, err
	}
	return &Control{
		cfg:    cfg,
		driver: driver,
		logger: logger,
		cache:  make(map[string][]byte),
		scale:  driver.Scale(),
	}, nil
}

func (c *Control) Run() error {
	c.logger.Info("start run: " + c.cfg.Dir)
	start := time.Now()
	defer func() {
		c.logger.Info("done: " + time.Since(start).String())
	}()

	if err := c.driver.Unlock(); err != nil {
		return err
	}

	if err := c.driver.AppTerminate(c.cfg.Package); err != nil {
		c.logger.Warn("AppTerminate error")
	}

	c.logger.Info("app launch")
	if err := c.driver.AppLaunch(c.cfg.Package); err != nil {
		return err
	}

	defer func() {
		if err := c.driver.AppTerminate(c.cfg.Package); err != nil {
			c.logger.Warn("AppTerminate error")
		}
		if err := c.driver.Lock(); err != nil {
			c.logger.Warn("Lock error")
		}
	}()

	for _, action := range c.cfg.Actions {
		exit, err := c.waitAndTap(action)
		if err != nil || exit {
			return err
		}
		time.Sleep(time.Millisecond * 500)
	}
	time.Sleep(time.Second)
	return nil
}

func (c *Control) waitAndTap(action Action) (bool, error) {
	c.logger.Info("start waitAndTap => " + action.Path)
	for i := 0; i < 100; i++ {
		screenshot, err := c.driver.Screenshot()
		if err != nil {
			return false, err
		}

		findPopup := false
		for _, popup := range action.Popups {
			if !popup.Find.exist() {
				continue
			}
			_, err = c.find(screenshot, popup.Find)
			if err == nil {
				c.logger.Info("found popup => " + popup.Find.Path)
				err = c.tapAction(screenshot, popup.Tap)
				if err != nil {
					return false, err
				}
				findPopup = true
				break
			}
		}
		if findPopup {
			continue
		}

		if action.Exit.exist() {
			_, err = c.find(screenshot, action.Exit)
			if err == nil {
				c.logger.Info("found exit => " + action.Exit.Path)
				return true, nil
			}
		}
		if action.Wait.exist() {
			_, err = c.find(screenshot, action.Wait)
			if err != nil {
				continue
			}
			c.logger.Info("found wait => " + action.Wait.Path)
		}

		err = c.tapAction(screenshot, action.TapAction)
		if err == nil {
			c.logger.Info("tap => " + action.TapAction.Path)
			return false, nil
		}
	}
	return false, errors.New("wait and tap timeout")
}

func (c *Control) find(source []byte, f Find) (p image.Point, err error) {
	tpl, ok := c.cache[f.Path]
	if !ok {
		tpl, err = os.ReadFile(filepath.Join(c.cfg.Dir, f.Path))
		helper.Check(err)
		c.cache[f.Path] = tpl
	}
	c.logger.Debug("find => " + f.Path)
	return opencv.Find(source, tpl, opencv.WithScale(c.scale), opencv.WithRange(opencv.Range(f.Range)))
}

func (c *Control) findAndTap(source []byte, f Find) error {
	p, err := c.find(source, f)
	if err != nil {
		return err
	}
	return c.driver.Tap(p.X, p.Y)
}

func (c *Control) tapAction(source []byte, a TapAction) error {
	if a.exist() {
		return c.findAndTap(source, a.Find)
	}
	return c.driver.Tap(a.X, a.Y)
}
