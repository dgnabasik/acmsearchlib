module github.com/dgnabasik/acmsearchlib/database

go 1.15

replace acmsearchlib/headers => ../headers

replace acmsearchlib/nulltime => ../nulltime

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210215175252-1e36ee979477
	github.com/lib/pq v1.9.0
	google.golang.org/protobuf v1.25.0 // indirect
)
