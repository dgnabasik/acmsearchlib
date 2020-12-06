module github.com/dgnabasik/acmsearchlib

go 1.15

replace acmsearchlib/headers => ../headers

replace acmsearchlib/nulltime => ../nulltime

require (
	github.com/dgnabasik/acmsearchlib/headers v0.0.0-00010101000000-000000000000
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-00010101000000-000000000000
	github.com/deckarep/golang-set v1.7.1
	github.com/lib/pq v1.9.0
)
