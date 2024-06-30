## TODO:

- display callsign
- spectrumClient.go needs a timeout. see https://pkg.go.dev/nhooyr.io/websocket
    - which uses: ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
    - also need timeout for paClient.go
- stop using the cp2pluto script
- add wiring schematic photos and better info to the doc folder
- test install.sh
- revisit qLog log levels
- revise what data to monitor
- revise what parameters to use
- tidy the user interface
- implement calibrater.go
- improve marker widths
- revist shutdown
