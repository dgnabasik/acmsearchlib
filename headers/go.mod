module github.com/dgnabasik/acmsearchlib/headers

replace acmsearchlib/nulltime => ../nulltime

go 1.17

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210830141227-2e506c045c84
	golang.org/x/text v0.3.7
)

require google.golang.org/protobuf v1.27.1 // indirect
