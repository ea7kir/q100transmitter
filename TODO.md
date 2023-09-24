## TODO:

- implement callsign
- spectrumClient.go needs a timeout. see https://pkg.go.dev/nhooyr.io/websocket
    - which uses: ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
    - also need timeout for paClient.go
- stop using scripts
- add wiring schematic
- test install.sh
- remove unused files and folders
- revisit qLog log levels
- revise what data to monitor
- revise what parameters to use
- tidy the user interface
- rename spectrumClient
- rename encoderClient
- implement calibrater.go
- improve marker widths
- fix service file not working on boot - ok after