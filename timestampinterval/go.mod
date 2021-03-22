module github.com/dgnabasik/acmsearchlib/timestampinterval

replace acmsearchlib/nulltime => ../nulltime

go 1.16

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210313144544-5f170d2a702d
	github.com/golang/protobuf v1.5.1
	google.golang.org/protobuf v1.26.0
)
