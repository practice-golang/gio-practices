package main // import "hello-gio"

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	// "gioui.org/font/gofont"
	"giofont"
)

type (
	// D : layout.Dimensions
	D = layout.Dimensions
	// C : layout.Context
	C = layout.Context
)

var (
	editor   = new(widget.Editor)
	topLabel = "Hello, Gio!!"
	longText = `안녕, 지오!!`
	list     = &layout.List{Axis: layout.Vertical}
)

func loop(w *app.Window) error {
	editor.SetText(longText)
	th := material.NewTheme(giofont.Collection())
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			setLayouts(gtx, th)
			e.Frame(gtx.Ops)
		}
	}
}

func setLayouts(gtx layout.Context, th *material.Theme) layout.Dimensions {
	widgets := []layout.Widget{
		material.H3(th, topLabel).Layout,
		func(gtx C) D { return setInputField(gtx, th) },
		func(gtx C) D { return setRectangle(gtx, th) },
	}

	return list.Layout(gtx, len(widgets), func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
	})
}

func setInputField(gtx C, th *material.Theme) D {
	return material.Editor(th, editor, "Hinta").Layout(gtx)
}

func setRectangle(gtx C, th *material.Theme) D {
	gtx.Constraints.Min.Y = gtx.Px(unit.Dp(50))
	gtx.Constraints.Max.Y = gtx.Constraints.Min.Y

	dr := image.Rectangle{Max: gtx.Constraints.Min}
	defer op.Save(gtx.Ops).Load()
	paint.LinearGradientOp{
		Stop1:  layout.FPt(dr.Min),
		Stop2:  layout.FPt(dr.Max),
		Color1: color.NRGBA{R: 0x10, G: 0xff, B: 0x10, A: 0xFF}, // Green
		Color2: color.NRGBA{R: 0x10, G: 0x10, B: 0xff, A: 0xFF}, // Blue
	}.Add(gtx.Ops)
	clip.Rect(dr).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func main() {
	go func() {
		w := app.NewWindow()
		w.Invalidate() // Prevent eat system.FrameEvent
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
