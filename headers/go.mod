module github.com/dgnabasik/acmsearchlib/headers

replace acmsearchlib/nulltime => ../nulltime

go 1.15

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20201206191427-03bcb92782c7
	golang.org/x/text v0.3.4
	github.com/golang/protobuf v1.4.3
)
