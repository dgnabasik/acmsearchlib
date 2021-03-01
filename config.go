package main

import (
	"fmt"

	art "github.com/dgnabasik/acmsearchlib/article"
	cat "github.com/dgnabasik/acmsearchlib/category"
	cond "github.com/dgnabasik/acmsearchlib/conditional"
	dbase "github.com/dgnabasik/acmsearchlib/database"
	fs "github.com/dgnabasik/acmsearchlib/filesystem"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"
	tsi "github.com/dgnabasik/acmsearchlib/timestampinterval"
	voc "github.com/dgnabasik/acmsearchlib/vocabulary"
	ws "github.com/dgnabasik/acmsearchlib/wordscore"
)

// Config func
func Config() []string {
	return []string{"acmsearchlib config",
		"\ndatabase:" + dbase.Version(),
		"\nfilesystem:" + fs.Version(),
		"\nheaders:" + hd.Version(),
		"\nnulltime:" + nt.Version(),
		"\ntimestampinterval:" + tsi.Version(),
		"\nvocabulary:" + voc.Version(),
		"\narticle:" + art.Version(),
		"\nconditional:" + cond.Version(),
		"\nwordscore:" + ws.Version(),
		"\ncategory:" + cat.Version(),
	}
}

func main() {
	fmt.Println(Config())
}
