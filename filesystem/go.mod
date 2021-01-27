module github.com/dgnabasik/acmsearchlib/filesystem

go 1.15

replace acmsearchlib/headers => ../headers

replace acmsearchlib/nulltime => ../nulltime

require (
	github.com/dgnabasik/acmsearchlib/headers v0.0.0-20210120203449-b95222e5cac7
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210120203449-b95222e5cac7
	google.golang.org/protobuf v1.25.0 // indirect

)
