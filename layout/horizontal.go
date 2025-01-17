package layout

import (
	"math"

	"github.com/leukipp/cortile/common"
	"github.com/leukipp/cortile/store"

	log "github.com/sirupsen/logrus"
)

type HorizontalLayout struct {
	*store.Manager         // Layout store manager
	Proportion     float64 // Master-slave proportion
	WorkspaceNum   uint    // Active workspace index
	Type           string  // Layout name
}

func CreateHorizontalLayout(workspaceNum uint) *HorizontalLayout {
	return &HorizontalLayout{
		Manager:      store.CreateManager(),
		Proportion:   common.Config.Proportion, // TODO: LTR/RTL support
		WorkspaceNum: workspaceNum,
		Type:         "horizontal",
	}
}

func (l *HorizontalLayout) Do() {
	log.Info("Tile ", len(l.Clients()), " windows with ", l.GetType(), " layout [workspace-", l.WorkspaceNum, "]")

	dx, dy, dw, dh := common.DesktopDimensions()
	msize := len(l.Masters)
	ssize := len(l.Slaves)

	my := dy
	mh := int(math.Round(float64(dh) * l.Proportion))
	sy := my + mh
	sh := dh - mh
	gap := common.Config.WindowGap

	asize := len(l.Clients())
	fsize := l.AllowedMasters

	ltr := false // TODO: Load from config

	if ltr && asize > fsize {
		mytmp := my
		mhtmp := mh
		sytmp := sy
		shtmp := sh

		my = sytmp
		mh = shtmp
		sy = mytmp + gap
		sh = mhtmp
	}

	if msize > 0 {
		mw := (dw - (msize+1)*gap) / msize
		if ssize == 0 {
			mh = dh
		}

		for i, c := range l.Masters {
			if !common.Config.WindowDecoration {
				c.UnDecorate()
			}
			c.MoveResize(gap+dx+i*(mw+gap), my+gap, mw, mh-2*gap)
		}
	}

	if ssize > 0 {
		sw := (dw - (ssize+1)*gap) / ssize
		if msize == 0 {
			sy, sh = dy, dh
		}

		for i, c := range l.Slaves {
			if !common.Config.WindowDecoration {
				c.UnDecorate()
			}
			c.MoveResize(gap+dx+i*(sw+gap), sy, sw, sh-gap)
		}
	}

	common.X.Conn().Sync()
}

func (l *HorizontalLayout) Undo() {
	for _, c := range append(l.Masters, l.Slaves...) {
		c.Restore()
	}
}

func (l *HorizontalLayout) IncrementProportion() {
	precision := 1.0 / common.Config.ProportionStep
	proportion := math.Round(l.Proportion*precision)/precision + common.Config.ProportionStep
	l.SetProportion(proportion)
}

func (l *HorizontalLayout) DecrementProportion() {
	precision := 1.0 / common.Config.ProportionStep
	proportion := math.Round(l.Proportion*precision)/precision - common.Config.ProportionStep
	l.SetProportion(proportion)
}

func (l *HorizontalLayout) SetProportion(p float64) {
	l.Proportion = math.Min(math.Max(p, common.Config.ProportionMin), common.Config.ProportionMax)
}

func (l *HorizontalLayout) GetType() string {
	return l.Type
}

func (l *HorizontalLayout) GetManager() *store.Manager {
	return l.Manager
}
