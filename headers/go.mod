module github.com/dgnabasik/acmsearchlib/headers

replace acmsearchlib/nulltime => ../nulltime

go 1.16

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210706001359-c82124fe2300
	golang.org/x/text v0.3.7
	google.golang.org/protobuf v1.27.1 // indirect
)
