package mobile

type Config struct {
	Dir     string   `json:"dir"`
	Package string   `json:"package"`
	Actions []Action `json:"actions"`
}

func (c *Config) setAfterWait() {
	for i, a := range c.Actions {
		if i == 0 {
			continue
		}
		aw := a.TapAction.Find
		if a.Wait.exist() {
			aw = a.Wait
		}
		c.Actions[i-1].afterWait = aw
		if a.Exit.exist() {
			c.Actions[i-1].afterExit = a.Exit
		}
	}
}

type Action struct {
	TapAction

	Exit Find `json:"exit"`
	Wait Find `json:"wait"`

	afterWait, afterExit Find

	Popups []Popup `json:"popups"`
}

type TapAction struct {
	Find
	X int `json:"x"`
	Y int `json:"y"`
}

type Find struct {
	Path  string `json:"path"`
	Range int    `json:"range"`
}

type Popup struct {
	Tap  TapAction `json:"tap"`
	Find Find      `json:"find"`
}

func (f Find) exist() bool {
	return f.Path != ""
}
