module github.com/dgnabasik/acmsearchlib/headers

replace acmsearchlib/nulltime => ../nulltime

go 1.15

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210102133359-c6c540bba695
	golang.org/x/text v0.3.4
	google.golang.org/protobuf v1.25.0 // indirect
)
