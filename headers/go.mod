module github.com/dgnabasik/acmsearchlib/headers

replace acmsearchlib/nulltime => ../nulltime

go 1.15

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210313144544-5f170d2a702d
	github.com/golang/protobuf v1.5.1 // indirect
	golang.org/x/text v0.3.5
)
