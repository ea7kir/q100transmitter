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
// const appFolder = "/home/pi/Q100/q100transmitter-v1/"

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
		Band:                    "Narrow",
		WideSymbolrate:          "1000",
		NarrowSymbolrate:        "333",
		VeryNarrowSymbolRate:    "125",
		WideFrequency:           "2405.25 / 09",
		NarrowFrequency:         "2409.75 / 27",
		VeryNarrowFrequency:     "2406.50 / 14",
		WideMode:                "DVB-S2",
		NarrowMode:              "DVB-S2",
		VeryNarrowMode:          "DVB-S2",
		WideCodecs:              "H265 ACC",
		NarrowCdecs:             "H265 ACC",
		VeryNarrowCodecs:        "H265 ACC",
		WideConstellation:       "QPSK",
		NarrowConstellation:     "QPSK",
		VeryNarrorConstellation: "QPSK",
		WideFec:                 "3/4",
		NarrowFec:               "3/4",
		VeryNarrowFec:           "3/4",
		WideVideoBitRate:        "350",
		NarrowVideoBitRate:      "350",
		VeryNarrowVideoBitRate:  "350",
		WideAudioBitRate:        "64000",
		NarrowAudioBitRate:      "64000",
		VeryNarrowAudioBitRate:  "64000",
		WideSpare1:              "sp1-a",
		NarrowSpare1:            "sp1-a",
		VeryNarrowSpare1:        "sp1-a",
		WideSpare2:              "sp2-a",
		NarrowSpare2:            "sp2-a",
		VeryNarrowSpare2:        "sp2-a",
		WideGain:                "-16",
		NarrowGain:              "-16",
		VeryNarrowGain:          "-16",
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
					showAboutBox()
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
					tuner.DecSelector(&tuner.Frequency)
				}
				if ui.incFrequency.Clicked() {
					tuner.IncSelector(&tuner.Frequency)
				}
				if ui.decMode.Clicked() {
					tuner.DecSelector(&tuner.Mode)
				}
				if ui.incMode.Clicked() {
					tuner.IncSelector(&tuner.Mode)
				}
				if ui.decCodecs.Clicked() {
					tuner.DecSelector(&tuner.Codecs)
				}
				if ui.incCodecs.Clicked() {
					tuner.IncSelector(&tuner.Codecs)
				}
				if ui.decConstellation.Clicked() {
					tuner.DecSelector(&tuner.Constellation)
				}
				if ui.incConstellation.Clicked() {
					tuner.IncSelector(&tuner.Constellation)
				}
				if ui.decFec.Clicked() {
					tuner.DecSelector(&tuner.Fec)
				}
				if ui.incFec.Clicked() {
					tuner.IncSelector(&tuner.Fec)
				}
				if ui.decVideoBitRate.Clicked() {
					tuner.DecSelector(&tuner.VideoBitRate)
				}
				if ui.incVideoBitRate.Clicked() {
					tuner.IncSelector(&tuner.VideoBitRate)
				}
				if ui.decAudioBitRate.Clicked() {
					tuner.DecSelector(&tuner.AudioBitRate)
				}
				if ui.incAudioBitRate.Clicked() {
					tuner.IncSelector(&tuner.AudioBitRate)
				}
				if ui.decSpare1.Clicked() {
					tuner.DecSelector(&tuner.Spare1)
				}
				if ui.incSpare1.Clicked() {
					tuner.IncSelector(&tuner.Spare1)
				}
				if ui.decSpare2.Clicked() {
					tuner.DecSelector(&tuner.Spare2)
				}
				if ui.incSpare2.Clicked() {
					tuner.IncSelector(&tuner.Spare2)
				}
				if ui.decGain.Clicked() {
					tuner.DecSelector(&tuner.Gain)
				}
				if ui.incGain.Clicked() {
					tuner.IncSelector(&tuner.Gain)
				}
				if ui.tune.Clicked() {
					tuner.Tune()
				}
				if ui.ptt.Clicked() {
					tuner.Ptt()
				}

				gtx := layout.NewContext(&ops, event)
				// set background to black
				paint.Fill(gtx.Ops, q100color.screenGrey)
				ui.layoutFlexes(gtx)
				event.Frame(gtx.Ops)
			}
		}
	}
}

// custom color scheme
var q100color = struct {
	screenGrey                               color.NRGBA
	labelWhite, labelOrange                  color.NRGBA
	buttonGrey, buttonGreen, buttonRed       color.NRGBA
	gfxBgd, gfxGreen, gfxGraticule, gfxLabel color.NRGBA
	gfxBeacon, gfxMarker                     color.NRGBA
}{
	// see: https://pkg.go.dev/golang.org/x/image/colornames
	// but maybe I should just create my own colors
	screenGrey:   color.NRGBA{R: 16, G: 16, B: 16, A: 255}, // no LightBlack
	labelWhite:   color.NRGBA(colornames.White),
	labelOrange:  color.NRGBA(colornames.Darkorange),       // or Orange or Darkorange or Gold
	buttonGrey:   color.NRGBA{R: 32, G: 32, B: 32, A: 255}, // DarkGrey is too light
	buttonGreen:  color.NRGBA(colornames.Green),
	buttonRed:    color.NRGBA(colornames.Red),
	gfxBgd:       color.NRGBA(colornames.Black),
	gfxGreen:     color.NRGBA(colornames.Green),
	gfxBeacon:    color.NRGBA(colornames.Red),
	gfxMarker:    color.NRGBA{R: 10, G: 10, B: 10, A: 255},
	gfxGraticule: color.NRGBA(colornames.Darkgray),
	gfxLabel:     color.NRGBA{R: 32, G: 32, B: 32, A: 255}, // DarkGrey is too light
}

// define all the buttons
type UI struct {
	about, shutdown                    widget.Clickable
	decBand, incBand                   widget.Clickable
	decSymbolRate, incSymbolRate       widget.Clickable
	decFrequency, incFrequency         widget.Clickable
	decMode, incMode                   widget.Clickable
	decCodecs, incCodecs               widget.Clickable
	decConstellation, incConstellation widget.Clickable
	decFec, incFec                     widget.Clickable
	decVideoBitRate, incVideoBitRate   widget.Clickable
	decAudioBitRate, incAudioBitRate   widget.Clickable
	decSpare1, incSpare1               widget.Clickable
	decSpare2, incSpare2               widget.Clickable
	decGain, incGain                   widget.Clickable
	tune, ptt                          widget.Clickable
	th                                 *material.Theme
}

// this makes the code more readable2
type (
	C = layout.Context
	D = layout.Dimensions
)

func showAboutBox() {
	// TODO: implement an about box
}

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
		btn.Background = q100color.buttonGrey
	}
	btn.Color = q100color.labelWhite
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
				return ui.q100_Button(gtx, &ui.about, "Q-100 Transmitter", false, q100color.buttonGrey)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return ui.q100_Label(gtx, "server date goes here", q100color.labelOrange)
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				return ui.q100_Button(gtx, &ui.shutdown, "Shutdown", false, q100color.buttonGrey)
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
				return ui.q100_Button(gtx, dec, "<", false, q100color.buttonGrey)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(lblWidth)
				return ui.q100_Label(gtx, value, q100color.labelOrange)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				return ui.q100_Button(gtx, inc, ">", false, q100color.buttonGrey)
			})
		}),
	)
}

// returns [    [ button label button ]  [ button label button ]  [ button label button ]   ]
func (ui *UI) q100_TuneRow(gtx C) D {
	const btnWidth = 0

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

// returns a column of 3 rows of [label__  label__]
func (ui *UI) q100_Column3Rows(gtx C, dec, inc [3]widget.Clickable, value [3]string) D {
	const btnWidth = 0
	const lblWidth = 65

	return layout.Flex{
		Axis: layout.Vertical,
		// Spacing: layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			// return ui.q100_LabelValue(gtx, name[0], value[0])
			// return ui.q100_Selector(gtx, &ui.decMode, &ui.incMode, tuner.Mode.Value, btnWidth, lblWidth)
			return ui.q100_Selector(gtx, &dec[0], &inc[0], value[0], btnWidth, lblWidth)
		}),
		layout.Rigid(func(gtx C) D {
			// return ui.q100_LabelValue(gtx, name[1], value[1])
			return ui.q100_Selector(gtx, &dec[1], &inc[1], value[1], btnWidth, lblWidth)
		}),
		layout.Rigid(func(gtx C) D {
			// return ui.q100_LabelValue(gtx, name[2], value[2])
			return ui.q100_Selector(gtx, &dec[2], &inc[2], value[2], btnWidth, lblWidth)
		}),
	)
}

// returns a column with 2 buttons
func (ui *UI) q100_Column2Buttons(gtx C) D {
	const btnWidth = 70
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
				gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				gtx.Constraints.Min.Y = gtx.Dp(btnHeight)
				return ui.q100_Button(gtx, &ui.tune, "TUNE", tuner.IsTuned, q100color.buttonGreen)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				gtx.Constraints.Min.Y = gtx.Dp(btnHeight)
				return ui.q100_Button(gtx, &ui.ptt, "PTT", tuner.IsPtt, q100color.buttonRed)
			})
		}),
	)
}

// returns 3 columns of 3 rows + 1 column with 2 buttons
func (ui *UI) q100_4ColumnsDataWithButtons(gtx C) D {
	dec1 := [3]widget.Clickable{ui.decCodecs, ui.decVideoBitRate, ui.decAudioBitRate}
	inc1 := [3]widget.Clickable{ui.incCodecs, ui.incVideoBitRate, ui.incAudioBitRate}
	val1 := [3]string{tuner.Codecs.Value, tuner.VideoBitRate.Value, tuner.AudioBitRate.Value}

	dec2 := [3]widget.Clickable{ui.decMode, ui.decConstellation, ui.decFec}
	inc2 := [3]widget.Clickable{ui.decMode, ui.decConstellation, ui.decFec}
	val2 := [3]string{tuner.Mode.Value, tuner.Constellation.Value, tuner.Fec.Value}

	dec3 := [3]widget.Clickable{ui.decSpare1, ui.decSpare2, ui.decGain}
	inc3 := [3]widget.Clickable{ui.decSpare1, ui.decSpare2, ui.decGain}
	val3 := [3]string{tuner.Spare1.Value, tuner.Spare2.Value, tuner.Gain.Value}

	return layout.Flex{
		Axis: layout.Horizontal,
		// Spacing: layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return ui.q100_Column3Rows(gtx, dec1, inc1, val1)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Column3Rows(gtx, dec2, inc2, val2)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Column3Rows(gtx, dec3, inc3, val3)
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
