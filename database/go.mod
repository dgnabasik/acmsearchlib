module github.com/dgnabasik/acmsearchlib/database

go 1.16

replace acmsearchlib/headers => ../headers

replace acmsearchlib/nulltime => ../nulltime

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210322142620-1dea691eef44
	github.com/golang/protobuf v1.5.1 // indirect
	github.com/lib/pq v1.10.0
)
