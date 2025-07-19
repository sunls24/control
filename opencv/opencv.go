package opencv

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
)

const (
	defMethod    = gocv.TmCcoeffNormed
	defThreshold = 0.95
)

func Find(source, tpl []byte, opts ...FindOption) (p image.Point, err error) {
	cfg := defFindConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	var sourceMat, tplMat gocv.Mat
	if sourceMat, err = gocv.IMDecode(source, gocv.IMReadGrayScale); err != nil {
		return
	}
	defer sourceMat.Close()
	if tplMat, err = gocv.IMDecode(tpl, gocv.IMReadGrayScale); err != nil {
		return
	}
	defer tplMat.Close()

	offset := 0
	if cfg.sourceRange != All {
		start, end := cfg.sourceRange.get(sourceMat.Rows())
		sourceMat = sourceMat.RowRange(start, end)
		offset = start
	}

	result, mask := gocv.NewMat(), gocv.NewMat()
	defer func() {
		_ = mask.Close()
		_ = result.Close()
	}()
	err = gocv.MatchTemplate(sourceMat, tplMat, &result, defMethod, mask)
	if err != nil {
		return
	}

	_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)
	if maxVal < defThreshold {
		err = fmt.Errorf("threshold < %f", defThreshold)
		return
	}

	return image.Point{
		X: (maxLoc.X + tplMat.Cols()/2) / cfg.scale,
		Y: ((maxLoc.Y + tplMat.Rows()/2) + offset) / cfg.scale,
	}, nil
}
