module github.com/dgnabasik/acmsearchlib/headers

replace acmsearchlib/nulltime => ../nulltime

go 1.15

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20201206193712-f1b276987652
	github.com/golang/protobuf v1.4.3
	golang.org/x/text v0.3.4
	google.golang.org/protobuf v1.25.0 // indirect
)
