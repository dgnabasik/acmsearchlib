module github.com/dgnabasik/acmsearchlib/database

go 1.15

replace acmsearchlib/headers => ../headers

replace acmsearchlib/nulltime => ../nulltime

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210303203715-a3f25ab4f17c
	github.com/lib/pq v1.10.0
	google.golang.org/protobuf v1.25.0 // indirect
)
