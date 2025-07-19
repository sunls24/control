package mobile

import (
	"control/helper"
	"encoding/json"
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg := Config{Actions: []Action{{
		TapAction: TapAction{},
		Exit:      Find{},
		Wait:      Find{},
		Popups:    nil,
	}}}
	b, err := json.Marshal(cfg)
	helper.Check(err)
	fmt.Println(string(b))
}
