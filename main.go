/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package main

import (
	"context"
	"image"
	"image/color"
	"os"
	"os/signal"
	"q100transmitter/hev10"
	"q100transmitter/logger"
	"q100transmitter/pluto"
	"q100transmitter/spReader"
	"q100transmitter/tuner"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ajstarks/giocanvas"
	"golang.org/x/image/colornames"
)

// application directory
// NOTE: this only works if we are are already in the correct folder
const appFolder = "/home/pi/Q100/q100transmitter-v1/"

// configuration data
var (
	spConfig = spReader.SpConfig{
		Url: "wss://eshail.batc.org.uk/wb/fft/fft_ea7kirsatcontroller:443/",
	}
	heConfig = hev10.HeConfig{
		Audio_codec:   "ACC",
		Audio_bitrate: "64000",
		Video_codec:   "H.265",
		Video_size:    "1280x720",
		Video_bitrate: "330",
		Url:           "udp://192.168.3.10:8282",
	}
	plConfig = pluto.PlConfig{
		Frequency:        "2409.75",
		Mode:             "DBS2",
		Constellation:    "QPSK",
		Symbol_rate:      "333",
		Fec:              "23",
		Gain:             "-10",
		Calibration_mode: "nocalib",   // NOTE: not implemented
		Pcr_pts:          "800",       // NOTE: not implemented
		Pat_period:       "200",       // NOTE: not implemented
		Roll_off:         "0.35",      // NOTE: not implemented
		Pilots:           "off",       // NOTE: not implemented
		Frame:            "LongFrame", // NOTE: not implemented
		H265box:          "undefined", // NOTE: not implemented
		Remux:            "1",         // NOTE: not implemented
		Provider:         "EA7KIR",
		Service:          "Michael",
	}

	tuConfig = tuner.TuConfig{
		Band:                 "Narrow",
		WideFrequency:        "10494.75 / 09",
		WideSymbolrate:       "1000",
		NarrowFrequency:      "10499.25 / 27",
		NarrowSymbolrate:     "333",
		VeryNarrowFrequency:  "10496.00 / 14",
		VeryNarrowSymbolRate: "125",
	}
)

// local data
var (
	spData    spReader.SpData
	spChannel = make(chan spReader.SpData, 5)
)

func main() {
	os.Setenv("DISPLAY", ":0") // required for X11

	spReader.Intitialize(spConfig, spChannel)
	spReader.Start()

	hev10.Initialize(heConfig)

	pluto.Intitialize(plConfig)

	tuner.Intitialize(tuConfig)
	tuner.Start()

	go func() {
		// TODO: add signal to catch interupts
		w := app.NewWindow(app.Fullscreen.Option())
		app.Size(800, 480) // I don't know if this is help in any way
		if err := loop(w); err != nil {
			logger.Fatal.Fatal(err)
			// log.Fatal(err)
		}

		tuner.Stop() // does nothing yet
		// lmReader.Stop()
		spReader.Stop() // does nothing yet

		os.Exit(0)
	}()

	app.Main()
}

func loop(w *app.Window) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ui := UI{
		th: material.NewTheme(gofont.Collection()), // TODO: change text colors, fonts, etc ?
	}
	var ops op.Ops
	// Capture the context done channel in a variable so that we can nil it
	// out after it closes and prevent its select case from firing again.
	done := ctx.Done()

	for {
		select {
		case <-done:
			// When the context cancels, assign the done channel to nil to
			// prevent it from firing over and over.
			done = nil
			// Log something to make it obvious this happened.
			// logger.Info.Printf("context cancelled")
			// Initiate window shutdown.
			tuner.Stop() // TODO: does nothing yet
			// lmReader.Stop() // TODO: does nothing yet
			spReader.Stop() // TODO: does nothing yet - bombs with Control=C
			w.Perform(system.ActionClose)
		// case lmData = <-lmChannel:
		// 	w.Invalidate()
		case spData = <-spChannel:
			w.Invalidate()
		case event := <-w.Events():
			switch event := event.(type) {
			case system.DestroyEvent:
				return event.Err
			case system.FrameEvent:
				if ui.about.Clicked() {
					{
						// TODO: implement an About Box
					}
				}
				if ui.shutdown.Clicked() {
					w.Perform(system.ActionClose)
				}
				if ui.decBand.Clicked() {
					tuner.DecBandSelector(&tuner.Band)
				}
				if ui.incBand.Clicked() {
					tuner.IncBandSelector(&tuner.Band)
				}
				if ui.decSymbolRate.Clicked() {
					tuner.DecSelector(&tuner.SymbolRate)
				}
				if ui.incSymbolRate.Clicked() {
					tuner.IncSelector(&tuner.SymbolRate)
				}
				if ui.decFrequency.Clicked() {
					tuner.DecFrequencySelector(&tuner.Frequency)
				}
				if ui.incFrequency.Clicked() {
					tuner.IncFrequencySelector(&tuner.Frequency)
				}
				if ui.tune.Clicked() {
					tuner.Tune()
				}
				if ui.Ptt.Clicked() {
					tuner.Ptt()
				}

				gtx := layout.NewContext(&ops, event)
				// set background to black
				paint.Fill(gtx.Ops, q100color.scrBgd)
				ui.layoutFlexes(gtx)
				event.Frame(gtx.Ops)
			}
		}
	}
}

// custom color scheme
var q100color = struct {
	scrBgd, scrTxt                             color.NRGBA
	scrTxtData, scrTxtDataSelected             color.NRGBA
	btnTxt, btnBgd, btnBgdSel, btnBgdSelUrgent color.NRGBA
	gfxBgd, gfxGreen, gfxGraticule, gfxLabel   color.NRGBA
	gfxBeacon, gfxMarker                       color.NRGBA
	active                                     bool
}{
	// see: https://pkg.go.dev/golang.org/x/image/colornames
	// but maybe I should just create my own colors
	scrBgd:             color.NRGBA{R: 16, G: 16, B: 16, A: 255}, // no LightBlack
	scrTxt:             color.NRGBA(colornames.White),
	scrTxtData:         color.NRGBA(colornames.Darkorange),  // or Orange or Darkorange or Gold
	scrTxtDataSelected: color.NRGBA(colornames.Greenyellow), // or Green or Greenyellow
	btnTxt:             color.NRGBA(colornames.White),
	btnBgd:             color.NRGBA{R: 32, G: 32, B: 32, A: 255}, // DarkGrey is too light
	btnBgdSel:          color.NRGBA(colornames.Green),
	btnBgdSelUrgent:    color.NRGBA(colornames.Red),
	gfxBgd:             color.NRGBA(colornames.Black),
	gfxGreen:           color.NRGBA(colornames.Green),
	gfxBeacon:          color.NRGBA(colornames.Red),
	gfxMarker:          color.NRGBA{R: 10, G: 10, B: 10, A: 255},
	gfxGraticule:       color.NRGBA(colornames.Darkgray),
	gfxLabel:           color.NRGBA{R: 32, G: 32, B: 32, A: 255}, // DarkGrey is too light
	active:             false,
}

// define all the buttons
type UI struct {
	about, shutdown              widget.Clickable
	decBand, incBand             widget.Clickable
	decSymbolRate, incSymbolRate widget.Clickable
	decFrequency, incFrequency   widget.Clickable
	tune, ptt                    widget.Clickable
	th                           *material.Theme
}

// this make the code more readable2
type (
	C = layout.Context
	D = layout.Dimensions
)

// my customisable button
func (ui *UI) q100_Button(gtx C, button *widget.Clickable, label string, btnActive bool, btnActiveColor color.NRGBA) D {
	inset := layout.Inset{
		Top:    2,
		Bottom: 2,
		Left:   4,
		Right:  4,
	}

	btn := material.Button(ui.th, button, label)
	if btnActive {
		btn.Background = btnActiveColor
	} else {
		btn.Background = q100color.btnBgd
	}
	btn.Color = q100color.btnTxt
	return inset.Layout(gtx, btn.Layout)
}

// my custom label
func (ui *UI) q100_Label(gtx C, label string, txtColor color.NRGBA) D {
	inset := layout.Inset{
		Top:    2,
		Bottom: 2,
		Left:   4,
		Right:  4,
	}

	lbl := material.Body1(ui.th, label)
	lbl.Color = txtColor
	return inset.Layout(gtx, lbl.Layout)
}

// returns [ [ button ]  [ label_____________________________________________ ]  [ button ] ]
func (ui *UI) q100_TopRow(gtx C) D {
	const btnWidth = 30
	inset := layout.Inset{
		Top:    2,
		Bottom: 2,
		Left:   4,
		Right:  4,
	}

	return layout.Flex{
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				return ui.q100_Button(gtx, &ui.about, "Q-100 Receiver", false, q100color.btnBgd)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return ui.q100_Label(gtx, "server date goes here", q100color.scrTxtData)
			// return inset.Layout(gtx, material.Body1(ui.th, lmData.State).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				return ui.q100_Button(gtx, &ui.shutdown, "Shutdown", false, q100color.btnBgd)
			})
		}),
	)
}

// returns [ button label button ]
func (ui *UI) q100_Selector(gtx C, dec, inc *widget.Clickable, value string, btnWidth, lblWidth unit.Dp) D {
	inset := layout.Inset{
		Top:    2,
		Bottom: 2,
		Left:   4,
		Right:  4,
	}

	return layout.Flex{
		Axis: layout.Horizontal,
		// Spacing: layout.SpaceBetween,
		// Alignment: layout.Middle,
		// WeightSum: 0.3,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				return ui.q100_Button(gtx, dec, "<", false, q100color.btnBgd)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(lblWidth)
				return ui.q100_Label(gtx, value, q100color.scrTxtDataSelected)
				// return inset.Layout(gtx, material.Body1(ui.th, value).Layout)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				return ui.q100_Button(gtx, inc, ">", false, q100color.btnBgd)
			})
		}),
	)
}

// returns [    [ button label button ]  [ button label button ]  [ button label button ]   ]
func (ui *UI) q100_TuneRow(gtx C) D {
	const btnWidth = 50

	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return ui.q100_Selector(gtx, &ui.decBand, &ui.incBand, tuner.Band.Value, btnWidth, 100)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Selector(gtx, &ui.decSymbolRate, &ui.incSymbolRate, tuner.SymbolRate.Value, btnWidth, 50)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Selector(gtx, &ui.decFrequency, &ui.incFrequency, tuner.Frequency.Value, btnWidth, 100)
		}),
	)
}

// returns [ [ ------------------------------- spectrum --------------------------------- ] ]
func (ui *UI) q100_Spectrum(gtx C) D {
	// see: github.com/ajstarks/giocanvas

	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceSides,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				canvas := giocanvas.Canvas{
					Width:   float32(788), //gtx.Constraints.Max.X), //float32(width),  //float32(gtx.Constraints.Max.X),
					Height:  float32(250), //float32(hieght), //float32(500),
					Context: gtx,
				}
				// fmt.Printf("  Canvas: %#v\n", canvas.Context.Constraints)

				canvas.Background(q100color.gfxBgd)
				// tuning marker
				markerCentre, markerWidth := spReader.TuningMarker(tuner.Frequency.Value, tuner.SymbolRate.Value)
				canvas.Rect(markerCentre, 50, markerWidth, 100, q100color.gfxMarker)
				// polygon
				canvas.Polygon(spReader.Xp, spData.Yp, q100color.gfxGreen)
				// graticule
				const fyBase float32 = 3
				const fyInc float32 = 5.88235
				fy := fyBase
				for y := 0; y < 17; y++ {
					switch y {
					case 15:
						canvas.Text(1, fy, 1.5, "15dB", q100color.gfxLabel)
						canvas.HLine(5, fy, 94, 0.01, q100color.gfxGraticule)
					case 10:
						canvas.Text(1, fy, 1.5, "10dB", q100color.gfxLabel)
						canvas.HLine(5, fy, 94, 0.01, q100color.gfxGraticule)
					case 5:
						canvas.Text(1, fy, 1.5, "5dB", q100color.gfxLabel)
						canvas.HLine(5, fy, 94, 0.01, q100color.gfxGraticule)
					default:
						canvas.HLine(5, fy, 94, 0.005, q100color.gfxGraticule)
					}
					fy += fyInc
				}
				// beacon level
				canvas.HLine(5, spData.BeaconLevel, 94, 0.03, q100color.gfxBeacon)

				return layout.Dimensions{
					Size: image.Point{X: int(canvas.Width), Y: int(canvas.Height)},
				}
			},
		),
	)
}

// returns [ label__  label__ ]
func (ui *UI) q100_LabelValue(gtx C, label, value string) D {
	const lblWidth = 105
	const valWidth = 110
	inset := layout.Inset{
		Top:    2,
		Bottom: 2,
		Left:   4,
		Right:  4,
	}

	return layout.Flex{
		Axis: layout.Horizontal,
		// Spacing: layout.SpaceEnd,
		// Alignment: layout.Middle,
		// WeightSum: 0.3,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(lblWidth)
				gtx.Constraints.Max.X = gtx.Dp(lblWidth)
				return ui.q100_Label(gtx, label, q100color.scrTxt)
				// return inset.Layout(gtx, material.Body1(ui.th, label).Layout)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(valWidth)
				gtx.Constraints.Max.X = gtx.Dp(valWidth)
				return ui.q100_Label(gtx, value, q100color.scrTxtData)
				// return inset.Layout(gtx, material.Body1(ui.th, value).Layout)
			})
		}),
	)
}

// returns a column of 4 rows of [label__  label__]
func (ui *UI) q100_Column4Rows(gtx C, name, value [4]string) D {
	// const btnWidth = 50

	return layout.Flex{
		Axis: layout.Vertical,
		// Spacing: layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return ui.q100_LabelValue(gtx, name[0], value[0])
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_LabelValue(gtx, name[1], value[1])
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_LabelValue(gtx, name[2], value[2])
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_LabelValue(gtx, name[3], value[3])
		}),
	)
}

// returns a column with 2 buttons
func (ui *UI) q100_Column2Buttons(gtx C) D {
	// const btnWidth = 80
	const btnHeight = 50
	inset := layout.Inset{
		Top:    2,
		Bottom: 2,
		Left:   4,
		Right:  4,
	}
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				// gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				gtx.Constraints.Min.Y = gtx.Dp(btnHeight)
				return ui.q100_Button(gtx, &ui.tune, " TUNE ", tuner.IsTuned, q100color.btnBgdSelUrgent)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				// gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				gtx.Constraints.Min.Y = gtx.Dp(btnHeight)
				return ui.q100_Button(gtx, &ui.ptt, " PTT ", tuner.IsPtt, q100color.btnBgdSel)
			})
		}),
	)
}

// returns 3 columns of 4 rows + 1 column with 2 buttons
func (ui *UI) q100_4ColumnsDataWithButtons(gtx C) D {
	names1 := [4]string{"Frequency", "Symbol Rate", "Mode", "Constellation"}
	values1 := [4]string{lmData.Frequency, lmData.SymbolRate, lmData.Mode, lmData.Constellation}
	names2 := [4]string{"FEC", "Codecs", "dB MER", "dB Margin"}
	values2 := [4]string{lmData.Fec, lmData.VideoCodec + " " + lmData.AudioCodec, lmData.DbMer, lmData.DbMargin}
	names3 := [4]string{"dBm Power", "Null Ratio", "Provider", "Service"}
	values3 := [4]string{lmData.DbmPower, lmData.NullRatio, lmData.Provider, lmData.Service}

	return layout.Flex{
		Axis: layout.Horizontal,
		// Spacing: layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return ui.q100_Column4Rows(gtx, names1, values1)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Column4Rows(gtx, names2, values2)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Column4Rows(gtx, names3, values3)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Column2Buttons(gtx)
		}),
	)
}

// returns the entire display
func (ui *UI) layoutFlexes(gtx C) D {
	// inset := layout.Inset{
	// 	Top:    2,
	// 	Bottom: 2,
	// 	Left:   4,
	// 	Right:  4,
	// }

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{
				Axis: layout.Vertical,
				// Spacing:   layout.SpaceEnd,
				// Alignment: layout.Alignment(layout.N),
			}.Layout(gtx,
				// top row displays [ [ button ]  [ label_____________________________________________ ]  [ button ] ]
				layout.Rigid(ui.q100_TopRow),
				// spectrum displays [ [ ------------------------------- spectrum --------------------------------- ] ]
				layout.Rigid(ui.q100_Spectrum),
				// tuning row displays [    [ button label button ]  [ button label button ]  [ button label button ]   ]
				layout.Rigid(ui.q100_TuneRow),
				// data + buttons [ [ label__  label__ ]   [ label__  label__ ]   [ label__  label__ ]  [ button ] ]
				layout.Rigid(ui.q100_4ColumnsDataWithButtons),
			)
		}),
	)
}
