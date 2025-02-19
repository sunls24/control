package task

import (
	"control/adb"
	"control/auto"
	"fmt"
	"log/slog"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Task struct {
	Name    string           `json:"name"`
	Crontab string           `json:"crontab"`
	Package string           `json:"package"`
	Config  auto.ClickConfig `json:"config"`
	Actions []auto.Action    `json:"actions"`
}

func (t *Task) Run() {
	logger := slog.With(slog.String("name", t.Name))
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprint(r))
		}
	}()

	logger.Info("task run")
	check(adb.Unlock())
	check(adb.OpenApp(t.Package))
	check(auto.ClickAuto(t.Config, t.Actions...))
	check(adb.StopApp(t.Package))
	logger.Info("task end")
}
