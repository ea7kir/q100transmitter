/*
 *  Q-100 q100transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package main

import (
	"context"
	"flag"
	"image"
	"image/color"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"q100transmitter/encoderClient"
	"q100transmitter/paClient"
	"q100transmitter/plutoClient"
	"q100transmitter/pttSwitch"
	"q100transmitter/spClient"
	"q100transmitter/txControl"
	"strings"
	"syscall"
	"time"

	// _ "net/http/pprof"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
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

// local data
var (
	// tuCmd        txControl.TxCmd_t
	txCmdChan  = make(chan txControl.TxCmd_t)
	txData     txControl.TxData_t
	txDataChan = make(chan txControl.TxData_t)
	spData     spClient.SpData_t
	spDataChan = make(chan spClient.SpData_t, 1)
	paData     paClient.SvrData_t
	paDataChan = make(chan paClient.SvrData_t)
)

// profile from the Mac
// go tool pprof http://txtouch.local:6060/debug/pprof/profile
// go tool pprof -http=":" pprof.q100transmitter.samples.cpu.001.pb.gz

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("INFO ----- q100transmitter opened -----")

	var shutdown bool
	flag.BoolVar(&shutdown, "shutdown", false, "close and poweroff")
	flag.Parse()
	// fmt.Println("shudown: ", shutdown)

	// read callsign from /home/pi/Q100/callsign
	bytes, err := os.ReadFile("/home/pi/Q100/callsign")
	if err != nil {
		log.Fatalf("FATAL failed to read callsign: %err", err)
	}
	str := string(bytes)
	provider := strings.TrimSpace(str)
	service := "n/a" // current Pluto firmware doesn't provide a way to set this

	ctx, cancel := context.WithCancel(context.Background())

	go spClient.ReadSpectrumServer(ctx, spDataChan)
	go paClient.ReadPaServer(ctx, paDataChan)

	// TODO: move these into txControl na react to ctx cancel
	encoderClient.Start()                // TODO: implment with ctx
	plutoClient.Start(provider, service) // TODO: implment with ctx
	pttSwitch.Start()                    // TODO: implment with ctx

	go txControl.HandleCommands(ctx, txCmdChan, txDataChan)

	go func() {
		os.Setenv("XDG_RUNTIME_DIR", "/run/user/1000") // TODO: is 1000 corrrect?
		// os.Setenv("DISPLAY", ":0")                     // required for X11. Compile wit: go build --tags nowayland .
		os.Setenv("WAYLAND_DISPLAY", "wayland-1") // required for wayland. Compile with: go build --tags nox11 .
		app.Size(800, 480)                        // I don't know if this is help in any way
		var w app.Window
		w.Option(app.Fullscreen.Option())

		if err := loop(&w); err != nil {
			log.Fatalf("FATAL failed to start loop: %v", err)
		}

		cancel()
		// log.Printf("CANCEL IN MAIN ----- cancel() called")
		// allow time to cancel all functions
		time.Sleep(time.Second * 3)

		// TODO: move these into txControl na react to ctx cancel
		pttSwitch.Stop()     // TODO replace with ctx
		plutoClient.Stop()   // TODO replace with ctx
		encoderClient.Stop() // TODO replace with ctx

		if shutdown {
			log.Printf("INFO ----- q100transmitter will poweroff -----")
			time.Sleep(1 * time.Second)
			cmd := exec.Command("sudo", "poweroff")
			if err := cmd.Start(); err != nil {
				log.Fatalf("FATAL failed to poweroff: %v", err)
			}
			cmd.Wait()
		}

		log.Printf("INFO ----- q100transmitter closed -----")
		// log.Close()
		os.Exit(0)
	}()

	app.Main()
}

func loop(w *app.Window) error {
	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// defer signal.Stop(quit)

	ui := UI{
		//th: material.NewTheme(gofont.Collection()),
		th: material.NewTheme(),
	}
	// Cris says keep using the original font
	ui.th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))

	var ops op.Ops
	// Capture the context done channel in a variable so that we can nil it
	// out after it closes and prevent its select case from firing again.
	// done := ctx.Done()

	for {
		select {
		// case <-ctx.Done():
		case <-interrupt:
			// When the context cancels, assign the done channel to nil to
			// prevent it from firing over and over.
			interrupt = nil
			w.Perform(system.ActionClose)
		case txData = <-txDataChan:
			w.Invalidate()
		case paData = <-paDataChan:
			w.Invalidate()
		case spData = <-spDataChan:
			w.Invalidate()
		}

		switch event := w.Event().(type) {
		case app.DestroyEvent:
			return event.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, event)

			switch {
			case ui.about.Clicked(gtx):
				showAboutBox()
			case ui.shutdown.Clicked(gtx):
				interrupt <- syscall.SIGINT
				// w.Perform(system.ActionClose)
			case ui.decBand.Clicked(gtx):
				txCmdChan <- txControl.CmdDecBand
			case ui.incBand.Clicked(gtx):
				txCmdChan <- txControl.CmdIncBand
			case ui.decSymbolRate.Clicked(gtx):
				txCmdChan <- txControl.CmdDecSymbolRate
			case ui.incSymbolRate.Clicked(gtx):
				txCmdChan <- txControl.CmdIncSymbolRate
			case ui.decFrequency.Clicked(gtx):
				txCmdChan <- txControl.CmdDecFrequency
			case ui.incFrequency.Clicked(gtx):
				txCmdChan <- txControl.CmdIncFrequency
			case ui.decMode.Clicked(gtx):
				txCmdChan <- txControl.CmdDecMode
			case ui.incMode.Clicked(gtx):
				txCmdChan <- txControl.CmdIncMode
			case ui.decCodecs.Clicked(gtx):
				txCmdChan <- txControl.CmdDecCodecs
			case ui.incCodecs.Clicked(gtx):
				txCmdChan <- txControl.CmdIncCodecs
			case ui.decConstellation.Clicked(gtx):
				txCmdChan <- txControl.CmdDecConstellation
			case ui.incConstellation.Clicked(gtx):
				txCmdChan <- txControl.CmdIncConstaellation
			case ui.decFec.Clicked(gtx):
				txCmdChan <- txControl.CmdDecFec
			case ui.incFec.Clicked(gtx):
				txCmdChan <- txControl.CmdIncFec
			case ui.decVideoBitRate.Clicked(gtx):
				txCmdChan <- txControl.CmdDecVideoBitRate
			case ui.incVideoBitRate.Clicked(gtx):
				txCmdChan <- txControl.CmdIncVideoBitRate
			case ui.decAudioBitRate.Clicked(gtx):
				txCmdChan <- txControl.CmdDecAudioBitRate
			case ui.incAudioBitRate.Clicked(gtx):
				txCmdChan <- txControl.CmdIncAudioBitRate
			case ui.decResolution.Clicked(gtx):
				txCmdChan <- txControl.CmdDecResolution
			case ui.incResolution.Clicked(gtx):
				txCmdChan <- txControl.CmdIncResolution
			case ui.decSpare2.Clicked(gtx):
				txCmdChan <- txControl.CmdDecSpare2
			case ui.incSpare2.Clicked(gtx):
				txCmdChan <- txControl.CmdIncSpare2
			case ui.decGain.Clicked(gtx):
				txCmdChan <- txControl.CmdDecGain
			case ui.incGain.Clicked(gtx):
				txCmdChan <- txControl.CmdIncGain
			case ui.tune.Clicked(gtx):
				txCmdChan <- txControl.CmdTune
			case ui.ptt.Clicked(gtx):
				txCmdChan <- txControl.CmdPtt
			}

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
	//TODO: try adding txData to here
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
	decResolution, incResolution       widget.Clickable
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
	const btnWidth = 50
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
				return ui.q100_Button(gtx, &ui.about, "QO-100 Transmitter", false, q100color.buttonGrey)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return ui.q100_Label(gtx, paData.Status, q100color.labelOrange)
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				return ui.q100_Button(gtx, &ui.shutdown, "Shutdown", false, q100color.buttonGrey)
			})
		}),
	)
}

// returns a single Selector_t as [ button label button ]
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
	const btnWidth = 50

	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return ui.q100_Selector(gtx, &ui.decBand, &ui.incBand, txData.CurBand, btnWidth, 100)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Selector(gtx, &ui.decSymbolRate, &ui.incSymbolRate, txData.CurSymbolRate, btnWidth, 50)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Selector(gtx, &ui.decFrequency, &ui.incFrequency, txData.CurFrequency, btnWidth, 100)
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
				// log.Printf("INFO   Canvas: %#v\n", canvas.Context.Constraints)

				canvas.Background(q100color.gfxBgd)
				// tuning marker
				canvas.Rect(txData.MarkerCentre, 50, txData.MarkerWidth, 100, q100color.gfxMarker)
				// polygon
				canvas.Polygon(spClient.Xp, spData.Yp, q100color.gfxGreen)
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
	const btnWidth = 50
	const lblWidth = 85 //65 //123 //65

	return layout.Flex{
		Axis: layout.Vertical,
		// Spacing: layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return ui.q100_Selector(gtx, dec[0], inc[0], value[0], btnWidth, lblWidth)
		}),
		layout.Rigid(func(gtx C) D {
			return ui.q100_Selector(gtx, dec[1], inc[1], value[1], btnWidth, lblWidth)
		}),
		layout.Rigid(func(gtx C) D {
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
				return ui.q100_Button(gtx, &ui.tune, "TUNE", txData.CurIsTuned, q100color.buttonGreen)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Dp(btnWidth)
				gtx.Constraints.Min.Y = gtx.Dp(btnHeight)
				return ui.q100_Button(gtx, &ui.ptt, "PTT", txData.CurIsPtt, q100color.buttonRed)
			})
		}),
	)
}

// Returns a 3x3 matrix of selectors + 1 column with 2 buttons
func (ui *UI) q100_3x3selectorMatrixPlus2buttons(gtx C) D {
	dec1 := [3]*widget.Clickable{&ui.decCodecs, &ui.decVideoBitRate, &ui.decAudioBitRate}
	inc1 := [3]*widget.Clickable{&ui.incCodecs, &ui.incVideoBitRate, &ui.incAudioBitRate}
	val1 := [3]string{txData.CurCodecs, txData.CurVideoBitRate, txData.CurAudioBitRate}

	dec2 := [3]*widget.Clickable{&ui.decMode, &ui.decConstellation, &ui.decFec}
	inc2 := [3]*widget.Clickable{&ui.incMode, &ui.incConstellation, &ui.incFec}
	val2 := [3]string{txData.CurMode, txData.CurConstellation, txData.CurFec}

	dec3 := [3]*widget.Clickable{&ui.decResolution, &ui.decSpare2, &ui.decGain}
	inc3 := [3]*widget.Clickable{&ui.incResolution, &ui.incSpare2, &ui.incGain}
	val3 := [3]string{txData.CurResolution, txData.CurSpare2, txData.CurGain}

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
