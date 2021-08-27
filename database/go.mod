module github.com/dgnabasik/acmsearchlib/database

go 1.16

replace acmsearchlib/headers => ../headers

replace acmsearchlib/nulltime => ../nulltime

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210825170524-53e3e64f7179
	github.com/jackc/pgx/v4 v4.13.0
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/text v0.3.7 // indirect
)
