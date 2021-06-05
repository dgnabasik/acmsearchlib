module github.com/dgnabasik/acmsearchlib/database

go 1.16

replace acmsearchlib/headers => ../headers

replace acmsearchlib/nulltime => ../nulltime

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210605011152-f83321b8fda8
	github.com/jackc/pgproto3/v2 v2.0.7 // indirect
	github.com/jackc/pgx/v4 v4.11.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/text v0.3.6 // indirect
)
