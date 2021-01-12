module github.com/dgnabasik/acmsearchlib/filesystem

go 1.15

replace acmsearchlib/headers => ../headers

replace acmsearchlib/nulltime => ../nulltime

require (
	github.com/dgnabasik/acmsearchlib/headers v0.0.0-20210112194524-c533c0e22890
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210104030944-946acfe2561a
	google.golang.org/protobuf v1.25.0 // indirect

)
