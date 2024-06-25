/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"time"

	// _ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"q100transmitter/encoderClient"
	"q100transmitter/paClient"
	"q100transmitter/plutoClient"
	"q100transmitter/pttSwitch"
	"q100transmitter/spectrumClient"
	"q100transmitter/txControl"

	"github.com/ea7kir/qLog"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ajstarks/giocanvas"
	"golang.org/x/image/colornames"
)

// application directory for the configuration data
// const appFolder = "/home/pi/Q100/q100transmitter/"

// configuration data
var (
	spConfig = spectrumClient.SpConfig{
		// Url:    "wss://eshail.batc.org.uk/wb/fft/fft_ea7kirsatcontroller:443/",
		// Origin: "http://eshail.batc.org.uk/wb",
		Origin: "https://eshail.batc.org.uk/",
		Url:    "wss://eshail.batc.org.uk/wb/fft/fft_ea7kirsatcontroller:443/wss",
	}
	svrConfig = paClient.SvrConfig{
		Url:  "paserver.local",
		Port: 9999, //8765,
	}
	heConfig = encoderClient.HeConfig{
		// Codecs:       "H.265 ACC", // H.264 ACC | H.264 G711u | H.265 ACC | H.265 G711u
		// AudioBitRate: "64000",     // 32000 | 64000
		// VideoBitRate: "350",       // 32...16384
		// // alter the following with caution
		StreamIP:   "192.168.3.10",
		StreamPort: "8282",
		ConfigIP:   "192.168.3.1",
	}
	plConfig = plutoClient.PlConfig{
		// configure setting not provided by the GUI
		// Provider: "",
		// Service:  "",
		// alter the following with caution
		// CalibrationMode: "nocalib",
		// Pcr_pts:         "800",
		// Pat_period:      "200",
		// Roll_off:        "0.35",
		// Pilots:          "off",
		// Frame:           "LongFrame",
		// H265box:         "undefined",
		// Remux:           "1",
		Url: "pluto.local", // or maybe "192.168.2.1"
	}
	tuConfig = txControl.TuConfig{
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
		WideCodecs:              "H265 ACC", // H.264 ACC | H.264 G711u | H.265 ACC | H.265 G711u
		NarrowCdecs:             "H265 ACC",
		VeryNarrowCodecs:        "H265 ACC",
		WideConstellation:       "QPSK",
		NarrowConstellation:     "QPSK",
		VeryNarrorConstellation: "QPSK",
		WideFec:                 "3/4",
		NarrowFec:               "3/4",
		VeryNarrowFec:           "3/4",
		WideVideoBitRate:        "440", // 32...16384
		NarrowVideoBitRate:      "340",
		VeryNarrowVideoBitRate:  "310",
		WideAudioBitRate:        "64000", // 32000 | 64000
		NarrowAudioBitRate:      "64000",
		VeryNarrowAudioBitRate:  "32000",
		WideSpare1:              "sp1-a",
		NarrowSpare1:            "sp1-a",
		VeryNarrowSpare1:        "sp1-a",
		WideSpare2:              "sp2-a",
		NarrowSpare2:            "sp2-a",
		VeryNarrowSpare2:        "sp2-a",
		WideGain:                "-15",
		NarrowGain:              "-16",
		VeryNarrowGain:          "-20",
	}
)

// local data
var (
	spData     spectrumClient.SpData
	spChannel  = make(chan spectrumClient.SpData, 3) //, 5)
	svrData    paClient.SvrData
	svrChannel = make(chan paClient.SvrData, 3) //, 5)
)

// profile from the Mac
// go tool pprof http://txtouch.local:6060/debug/pprof/profile
// go tool pprof -http=":" pprof.q100transmitter.samples.cpu.001.pb.gz

func waitForNetwork() {
	var maxTries = 20
	for {
		client := http.Client{}
		_, err := client.Get("https://google.com")
		if err == nil {
			return
		}
		qLog.Warn("Waiting for network %v", maxTries)
		time.Sleep(time.Second)
		maxTries--
		if maxTries == 0 {
			qLog.Fatal("Unable to conect to network")
			qLog.Close()
			os.Exit(1)
		}
	}
}

func main() {
	logFile, err := os.OpenFile("/home/pi/Q100/transmitter.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("failed to open log file:", err)
		qLog.Close()
		os.Exit(1)
	}

	// log.SetOutput(os.Stderr)
	qLog.SetOutput(logFile)
	defer qLog.Close()

	// read callsign from /home/pi/Q100/callsign
	bytes, err := os.ReadFile("/home/pi/Q100/callsign")
	if err != nil {
		qLog.Fatal("ünable read callsign: %err", err)
	}
	plConfig.Provider = string(bytes)
	// fmt.Println("callsign: ", plConfig.Provider)
	// current Pluto firmware doesn't provide a way to set this
	plConfig.Service = ""

	qLog.Info("----- q100transmitter Opened -----")

	os.Setenv("DISPLAY", ":0") // required for X11

	waitForNetwork()

	spectrumClient.Intitialize(spConfig, spChannel)

	paClient.Initialize(svrConfig, svrChannel)

	encoderClient.Initialize(heConfig)

	plutoClient.Initialize(plConfig)

	pttSwitch.Initialize()

	txControl.Initialize(tuConfig)

	go func() {
		// app.Size(800, 480) // I don't know if this is help in any way
		var w app.Window
		w.Option(app.Fullscreen.Option())

		if err := loop(&w); err != nil {
			qLog.Fatal("failed to start loop: %v", err)
			os.Exit(1)
		}

		// TODO: implement with a d/on channel
		txControl.Stop()
		paClient.Stop()
		spectrumClient.Stop()

		if !true { // change to true for powerdown
			qLog.Info("----- q100transmitter will poweroff -----")
			time.Sleep(1 * time.Second)
			cmd := exec.Command("sudo", "poweroff")
			if err := cmd.Start(); err != nil {
				qLog.Error("failed to poweroff: %v", err)
				os.Exit(1)
			}
			cmd.Wait()
		}

		qLog.Info("----- q100transmitter Closed -----")
		os.Exit(0)
	}()

	app.Main()
}

func loop(w *app.Window) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ui := UI{
		//th: material.NewTheme(gofont.Collection()),
		th: material.NewTheme(),
	}
	// Cris says keep using the original font
	ui.th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))

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
			return nil
			// w.Perform(system.ActionClose)
		case svrData = <-svrChannel:
			w.Invalidate()
		case spData = <-spChannel:
			w.Invalidate()
		}

		switch event := w.Event().(type) {
		case app.DestroyEvent:
			return event.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, event)
			if ui.about.Clicked(gtx) {
				showAboutBox()
			}
			if ui.shutdown.Clicked(gtx) {
				return nil
				// w.Perform(system.ActionClose)
			}
			if ui.decBand.Clicked(gtx) {
				txControl.DecBandSelector(&txControl.Band)
			}
			if ui.incBand.Clicked(gtx) {
				txControl.IncBandSelector(&txControl.Band)
			}
			if ui.decSymbolRate.Clicked(gtx) {
				txControl.DecSelector(&txControl.SymbolRate)
			}
			if ui.incSymbolRate.Clicked(gtx) {
				txControl.IncSelector(&txControl.SymbolRate)
			}
			if ui.decFrequency.Clicked(gtx) {
				txControl.DecSelector(&txControl.Frequency)
			}
			if ui.incFrequency.Clicked(gtx) {
				txControl.IncSelector(&txControl.Frequency)
			}
			if ui.decMode.Clicked(gtx) {
				txControl.DecSelector(&txControl.Mode)
			}
			if ui.incMode.Clicked(gtx) {
				txControl.IncSelector(&txControl.Mode)
			}
			if ui.decCodecs.Clicked(gtx) {
				txControl.DecSelector(&txControl.Codecs)
			}
			if ui.incCodecs.Clicked(gtx) {
				txControl.IncSelector(&txControl.Codecs)
			}
			if ui.decConstellation.Clicked(gtx) {
				txControl.DecSelector(&txControl.Constellation)
			}
			if ui.incConstellation.Clicked(gtx) {
				txControl.IncSelector(&txControl.Constellation)
			}
			if ui.decFec.Clicked(gtx) {
				txControl.DecSelector(&txControl.Fec)
			}
			if ui.incFec.Clicked(gtx) {
				txControl.IncSelector(&txControl.Fec)
			}
			if ui.decVideoBitRate.Clicked(gtx) {
				txControl.DecSelector(&txControl.VideoBitRate)
			}
			if ui.incVideoBitRate.Clicked(gtx) {
				txControl.IncSelector(&txControl.VideoBitRate)
			}
			if ui.decAudioBitRate.Clicked(gtx) {
				txControl.DecSelector(&txControl.AudioBitRate)
			}
			if ui.incAudioBitRate.Clicked(gtx) {
				txControl.IncSelector(&txControl.AudioBitRate)
			}
			if ui.decSpare1.Clicked(gtx) {
				txControl.DecSelector(&txControl.Spare1)
			}
			if ui.incSpare1.Clicked(gtx) {
				txControl.IncSelector(&txControl.Spare1)
			}
			if ui.decSpare2.Clicked(gtx) {
				txControl.DecSelector(&txControl.Spare2)
			}
			if ui.incSpare2.Clicked(gtx) {
				txControl.IncSelector(&txControl.Spare2)
			}
			if ui.decGain.Clicked(gtx) {
				txControl.DecSelector(&txControl.Gain)
			}
			if ui.incGain.Clicked(gtx) {
				txControl.IncSelector(&txControl.Gain)
			}
			if ui.tune.Clicked(gtx) {
				txControl.Tune()
			}
			if ui.ptt.Clicked(gtx) {
				txControl.Ptt()
			}

			// gtx := layout.NewContext(&ops, event)
			// set the screen background to to dark grey
			paint.Fill(gtx.Ops, q100color.screenGrey)
			ui.layoutFlexes(gtx)
			event.Frame(gtx.Ops)
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
	gfxMarker:    color.NRGBA{R: 20, G: 20, B: 20, A: 255},
	gfxGraticule: color.NRGBA(colornames.Darkgray),
	gfxLabel:     color.NRGBA{R: 32, G: 32, B: 32, A: 255}, // DarkGrey is too light
}

// define all buttons
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

// makes the code more readable
type (
	C = layout.Context
	D = layout.Dimensions
)

// Returns an About box
func showAboutBox() {
	// TODO: implement an about box
}

// Return a customisable button
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

// Returns a customisable label
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

// Returns 1 row of 2 buttons and a label for About, Status and Shutdown
func (ui *UI) q100_TopStatusRow(gtx C) D {
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
			return ui.q100_Label(gtx, svrData.Status, q100color.labelOrange)
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				return ui.q100_Button(gtx, &ui.shutdown, "Shutdown", false, q100color.buttonGrey)
			})
		}),
	)
}

// returns a single Selector as [ button label button ]
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
		Alignment: layout.Middle, // Chris
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

// Returns 1 row of 3 Selectors for Band SymbolRate and Frequency
func (ui *UI) q100_MainTuningRow(gtx C) D {
	const btnWidth = 0

	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return ui.q100_Selector(gtx, &ui.decBand, &ui.incBand, txControl.Band.Value, btnWidth, 100)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Selector(gtx, &ui.decSymbolRate, &ui.incSymbolRate, txControl.SymbolRate.Value, btnWidth, 50)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Selector(gtx, &ui.decFrequency, &ui.incFrequency, txControl.Frequency.Value, btnWidth, 100)
		}),
	)
}

// Returns the Spectrum display
//
// see: github.com/ajstarks/giocanvas for docs
func (ui *UI) q100_SpectrumDisplay(gtx C) D {
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
					Theme:   ui.th,
				}
				// qLog.Info("  Canvas: %#v\n", canvas.Context.Constraints)

				canvas.Background(q100color.gfxBgd)
				// tuning marker
				canvas.Rect(spData.MarkerCentre, 50, spData.MarkerWidth, 100, q100color.gfxMarker)
				// polygon
				canvas.Polygon(spectrumClient.Xp, spData.Yp, q100color.gfxGreen)
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
func (ui *UI) q100_Column3Rows(gtx C, dec, inc [3]*widget.Clickable, value [3]string) D {
	const btnWidth = 0
	const lblWidth = 65 //123 //65

	return layout.Flex{
		Axis: layout.Vertical,
		// Spacing: layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			// return ui.q100_LabelValue(gtx, name[0], value[0])
			// return ui.q100_Selector(gtx, &ui.decMode, &ui.incMode, txController.Mode.Value, btnWidth, lblWidth)
			return ui.q100_Selector(gtx, dec[0], inc[0], value[0], btnWidth, lblWidth)
		}),
		layout.Rigid(func(gtx C) D {
			// return ui.q100_LabelValue(gtx, name[1], value[1])
			return ui.q100_Selector(gtx, dec[1], inc[1], value[1], btnWidth, lblWidth)
		}),
		layout.Rigid(func(gtx C) D {
			// return ui.q100_LabelValue(gtx, name[2], value[2])
			return ui.q100_Selector(gtx, dec[2], inc[2], value[2], btnWidth, lblWidth)
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
				return ui.q100_Button(gtx, &ui.tune, "TUNE", txControl.IsTuned, q100color.buttonGreen)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				gtx.Constraints.Min.Y = gtx.Dp(btnHeight)
				return ui.q100_Button(gtx, &ui.ptt, "PTT", txControl.IsPtt, q100color.buttonRed)
			})
		}),
	)
}

// Returns a 3x3 matrix of selectors + 1 column with 2 buttons
func (ui *UI) q100_3x3selectorMatrixPlus2buttons(gtx C) D {
	dec1 := [3]*widget.Clickable{&ui.decCodecs, &ui.decVideoBitRate, &ui.decAudioBitRate}
	inc1 := [3]*widget.Clickable{&ui.incCodecs, &ui.incVideoBitRate, &ui.incAudioBitRate}
	val1 := [3]string{txControl.Codecs.Value, txControl.VideoBitRate.Value, txControl.AudioBitRate.Value}

	dec2 := [3]*widget.Clickable{&ui.decMode, &ui.decConstellation, &ui.decFec}
	inc2 := [3]*widget.Clickable{&ui.incMode, &ui.incConstellation, &ui.incFec}
	val2 := [3]string{txControl.Mode.Value, txControl.Constellation.Value, txControl.Fec.Value}

	dec3 := [3]*widget.Clickable{&ui.decSpare1, &ui.decSpare2, &ui.decGain}
	inc3 := [3]*widget.Clickable{&ui.incSpare1, &ui.incSpare2, &ui.incGain}
	val3 := [3]string{txControl.Spare1.Value, txControl.Spare2.Value, txControl.Gain.Value}

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
		// layout.Rigid(func(gtx C) D {
		// 	return ui.q100_Column3Rows(gtx, dec4, inc4, val4) // TODO: add 3 more buttons here
		// }),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Column2Buttons(gtx)
		}),
	)
}

// layoutFlexes returns the entire display
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
				layout.Rigid(ui.q100_TopStatusRow),
				layout.Rigid(ui.q100_SpectrumDisplay),
				layout.Rigid(ui.q100_MainTuningRow),
				layout.Rigid(ui.q100_3x3selectorMatrixPlus2buttons),
			)
		}),
	)
}
