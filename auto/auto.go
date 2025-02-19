package auto

import (
	"bytes"
	"control/adb"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	screen  = "./ocr/screen.png"
	screen1 = "./ocr/screen1.png"
	screen2 = "./ocr/screen2.png"
	screen3 = "./ocr/screen3.png"

	ocrRun = `./ocr/OcrLiteOnnx --models ./ocr/models \
--det dbnet.onnx \
--cls angle_net.onnx \
--rec crnn_lite_lstm.onnx \
--keys keys.txt \
--image %s \
--numThread 8 \
--padding 10 \
--maxSideLen 1080 \
--boxScoreThresh 0.6 \
--boxThresh 0.3 \
--unClipRatio 2.0 \
--doAngle 0 \
--mostAngle 0`
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(data []byte) error {
	dur, err := time.ParseDuration(strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

type ClickConfig struct {
	TryCount int      `json:"try_count"`
	TryWait  Duration `json:"try_wait"`
}

func click(action Action) (bool, error) {
	err := adb.Screenshot(screen)
	if err != nil {
		return false, fmt.Errorf("screenshot: %s", err)
	}

	var img = screen
	var offset = 0
	if action.Position != 0 {
		img, offset, err = cropImage(img, action.Position)
		if err != nil {
			return false, fmt.Errorf("crop image: %s", err)
		}
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf(ocrRun, img))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("ocr: %s, output:\n%s", err, string(output))
	}

	items := strings.Split(strings.TrimSpace(string(output)), "---")
	if len(items) == 0 {
		return false, errors.New("not found item")
	}

	if action.Exit != "" || action.Exist != "" {
		for _, v := range items {
			if action.Exit != "" && strings.Contains(v, action.Exit) {
				return true, nil
			}
			if action.Exist != "" && strings.Contains(v, action.Exist) {
				return false, nil
			}
		}
	}

	p, err := calcPoint(action.Text, offset, items)
	if err != nil {
		return false, err
	}

	return false, adb.Click(p.x, p.y)
}

type point struct {
	x, y int
}

func calcPoint(str string, offset int, items []string) (point, error) {
	for _, v := range items {
		if strings.Contains(v, str) {
			sp := strings.Split(strings.TrimSpace(v), "\n")
			points := strings.Split(sp[1], "|")
			var xt, yt int
			for _, p := range points {
				pp := strings.Split(p, ",")
				x, _ := strconv.Atoi(pp[0])
				y, _ := strconv.Atoi(pp[1])
				xt += x
				yt += y
			}

			p := point{
				x: xt / 4,
				y: yt/4 + offset,
			}
			slog.Info(fmt.Sprintf("-> %s (%d,%d) ocr: %sms", str, p.x, p.y, strings.TrimSpace(items[len(items)-1])))
			return p, nil
		}
	}
	return point{}, fmt.Errorf("not found text %s", str)
}

func ClickAuto(cfg ClickConfig, arr ...Action) error {
	time.Sleep(time.Duration(cfg.TryWait))
	for _, action := range arr {
		for i := 0; i < cfg.TryCount; i++ {
			exit, err := click(action)
			if exit {
				return nil
			}
			if err != nil {
				slog.Debug(fmt.Sprintf("click: %s", err), slog.Int("count", i+1))
				if i == cfg.TryCount-1 {
					return fmt.Errorf("-> %s 超过最大重试次数", action.Text)
				}
				time.Sleep(time.Duration(cfg.TryWait))
				continue
			}
			break
		}
		time.Sleep(time.Duration(cfg.TryWait))
	}
	time.Sleep(time.Duration(cfg.TryWait))
	return nil
}

type Action struct {
	Text     string   `json:"text"`     // 目标文本
	Position Position `json:"position"` // 目标文本在屏幕的位置
	Wait     string   `json:"wait"`     // 等待此内容出现后点击
	Exist    string   `json:"exist"`    // 存在此内容，跳过点击
	Exit     string   `json:"exit"`     // 存在此内容，退出流程
}

type Position int

const (
	Top    Position = 1
	Center Position = 2
	Bottom Position = 3
)

func cropImage(path string, p Position) (string, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return "", 0, err
	}

	bounds := img.Bounds()
	h3 := bounds.Dy() / 3

	var cropRect image.Rectangle
	var offset int
	switch p {
	case Top:
		path = screen1
		cropRect = image.Rect(
			bounds.Min.X,
			bounds.Min.Y,
			bounds.Max.X,
			bounds.Min.Y+h3,
		)
	case Center:
		path = screen2
		offset = h3
		cropRect = image.Rect(
			bounds.Min.X,
			bounds.Min.Y+h3,
			bounds.Max.X,
			bounds.Max.Y-h3,
		)
	case Bottom:
		path = screen3
		offset = bounds.Max.Y - h3
		cropRect = image.Rect(
			bounds.Min.X,
			bounds.Max.Y-h3,
			bounds.Max.X,
			bounds.Max.Y,
		)
	}

	newImg := image.NewRGBA(image.Rect(0, 0, cropRect.Dx(), cropRect.Dy()))
	draw.Draw(newImg, newImg.Bounds(), img, cropRect.Min, draw.Src)

	cf, err := os.Create(path)
	if err != nil {
		return "", 0, err
	}
	defer cf.Close()

	if err = png.Encode(cf, newImg); err != nil {
		return "", 0, err
	}

	return path, offset, nil
}
