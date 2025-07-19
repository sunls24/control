package mobile

type Config struct {
	Dir     string   `json:"dir"`
	Package string   `json:"package"`
	Actions []Action `json:"actions"`
}

type Action struct {
	TapAction

	Exit Find `json:"exit"`
	Wait Find `json:"wait"`

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
