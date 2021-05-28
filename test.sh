# Assumes $GOROOT=/usr/local/go
sudo cp *.* $GOROOT/src/acmsearchlib
go test -v acmsearchlib
