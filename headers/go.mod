module github.com/dgnabasik/acmsearchlib/headers

replace acmsearchlib/nulltime => ../nulltime

go 1.15

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210131022223-510b54edd542
	golang.org/x/text v0.3.5
	google.golang.org/protobuf v1.25.0 // indirect
)
