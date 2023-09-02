## TODO:

- spectrumClient.go needs a timeout. see https://pkg.go.dev/nhooyr.io/websocket
    - which uses: ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
    - also need timeout for paClient.go
- stop using scripts
- add wiring schematic
- automate the install process
- remove unused files and folders
- revisit qLog log levelse
- revise what data to monitor
- revise what parameters to use
- tidy the user interface
- rename spectrumClient
- rename encoderClient
- implement calibrater.go and marker widths
- add 3 more buttons and split [H265 ACC] into separate buttons
- add install UDEV_RULES to install script