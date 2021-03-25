module github.com/dgnabasik/acmsearchlib/headers

replace acmsearchlib/nulltime => ../nulltime

go 1.16

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210325002052-704b7df69a8a
	github.com/golang/protobuf v1.5.1 // indirect
	golang.org/x/text v0.3.5
)
