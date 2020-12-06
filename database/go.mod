module github.com/dgnabasik/acmsearchlib/database

go 1.15

replace acmsearchlib/headers => ../headers

replace acmsearchlib/nulltime => ../nulltime

require (
	github.com/deckarep/golang-set v1.7.1
	github.com/dgnabasik/acmsearchlib/headers v0.0.0-20201206193712-f1b276987652
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20201206191427-03bcb92782c7
	github.com/lib/pq v1.9.0
)
