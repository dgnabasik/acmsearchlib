module github.com/dgnabasik/acmsearchlib/timestampinterval

replace acmsearchlib/nulltime => ../nulltime

go 1.15

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210303203715-a3f25ab4f17c
	github.com/golang/protobuf v1.4.3
	google.golang.org/protobuf v1.25.0
)
