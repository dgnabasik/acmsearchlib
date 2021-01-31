module github.com/dgnabasik/acmsearchlib/timestampinterval

replace acmsearchlib/nulltime => ../nulltime

go 1.15

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210127040904-c08867ba0f70
	github.com/golang/protobuf v1.4.3
	google.golang.org/protobuf v1.25.0
)
