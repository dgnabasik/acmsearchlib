package main

import (
	art "github.com/dgnabasik/acmsearchlib/article"
	cond "github.com/dgnabasik/acmsearchlib/conditional"
	dbase "github.com/dgnabasik/acmsearchlib/database"
	fs "github.com/dgnabasik/acmsearchlib/filesystem"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"
	tsi "github.com/dgnabasik/acmsearchlib/timestampinterval"
	ws "github.com/dgnabasik/acmsearchlib/wordscore"
)

func Config() []string {
	return []string{"acmsearchlib config", art.Version(), cond.Version(), dbase.Version(), fs.Version(), hd.Version(), nt.Version(), tsi.Version(), ws.Version()}
}

func main() {
	Config()
}
