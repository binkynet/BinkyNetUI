module github.com/binkynet/BinkyNetUI

go 1.14

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a

require (
	fyne.io/fyne v1.3.3
	github.com/binkynet/BinkyNet v0.1.1-0.20200628111728-f9596a75610c
	github.com/rs/zerolog v1.18.0
	github.com/spf13/cobra v1.0.0
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	google.golang.org/grpc v1.27.1
)
