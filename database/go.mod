module github.com/dgnabasik/acmsearchlib/database

go 1.16

replace acmsearchlib/headers => ../headers

replace acmsearchlib/nulltime => ../nulltime

require (
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210706001359-c82124fe2300
	github.com/jackc/pgx/v4 v4.12.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
