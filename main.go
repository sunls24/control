package main

import (
	"control/task"
	"encoding/json"
	"github.com/robfig/cron/v3"
	"log/slog"
	"os"
)

type Config struct {
	SimpleTask []task.Task `json:"simple_task"`
}

func main() {
	f, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	var c Config
	err = json.Unmarshal(f, &c)
	if err != nil {
		panic(err)
	}
	ct := cron.New()
	for _, t := range c.SimpleTask {
		logger := slog.With(slog.String("name", t.Name), slog.String("cron", t.Crontab))
		if t.Crontab != "" {
			_, err = ct.AddFunc(t.Crontab, t.Run)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			logger.Info("任务添加成功")
		}
	}
	ct.Run()
}
