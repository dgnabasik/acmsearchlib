module github.com/dgnabasik/acmsearchlib/filesystem

go 1.15

replace acmsearchlib/headers => ../headers

replace acmsearchlib/nulltime => ../nulltime

require (
	github.com/dgnabasik/acmsearchlib/headers v0.0.0-20210112195813-1cb5bedbd297
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210112195813-1cb5bedbd297
	google.golang.org/protobuf v1.25.0 // indirect

)
