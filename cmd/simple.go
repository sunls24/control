package main

import (
	"control/helper"
	"control/mobile"
	"control/mobile/wda"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"log/slog"
	"os"
)

type Config struct {
	Name    string `json:"name"`
	Crontab string `json:"crontab"`
	AtOnce  bool   `json:"atOnce"`
	mobile.Config
}

func init() {
	// TODO TapAction Find string() 实现
	//slog.SetLogLoggerLevel(slog.LevelDebug)
}

func main() {
	f, err := os.ReadFile("./config.json")
	helper.Check(err)
	var tasks = make([]Config, 0)
	helper.Check(json.Unmarshal(f, &tasks))

	ct := cron.New()
	for _, t := range tasks {
		logger := slog.With(slog.String("name", t.Name))
		m := mustNewWDA(t.Config, logger)
		if t.AtOnce {
			helper.Check(m.Run())
			continue
		}
		if t.Crontab == "" {
			continue
		}
		_, err = ct.AddFunc(t.Crontab, func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Error(fmt.Sprint(r))
				}
			}()
			helper.Check(m.Run())
		})
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		logger.Info("crontab added: " + t.Crontab)
	}
	ct.Run()
}

func mustNewWDA(cfg mobile.Config, logger *slog.Logger) *mobile.Control {
	driver, err := wda.New()
	helper.Check(err)
	m, err := mobile.New(driver, cfg, logger)
	helper.Check(err)
	return m
}
