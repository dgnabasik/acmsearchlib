module github.com/dgnabasik/acmsearchlib/filesystem

go 1.16

replace acmsearchlib/headers => ../headers

replace acmsearchlib/nulltime => ../nulltime

require (
	github.com/dgnabasik/acmsearchlib/headers v0.0.0-20210825170524-53e3e64f7179
	github.com/dgnabasik/acmsearchlib/nulltime v0.0.0-20210825170524-53e3e64f7179
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/ugorji/go v1.2.6 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect

)
