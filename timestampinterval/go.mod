module github.com/dgnabasik/acmsearchlib/timestampinterval

replace acmsearchlib/nulltime => ../nulltime

go 1.15

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210112195813-1cb5bedbd297
	github.com/golang/protobuf v1.4.3
	google.golang.org/protobuf v1.25.0
)
