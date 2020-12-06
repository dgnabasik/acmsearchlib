package main

import (
	//proto "github.com/dgnabasik/acmsearchlib/TimeEventService/pb"
	//pb "github.com/dgnabasik/acmsearchlib/WebpageService/pb"

	"fmt"

	dbase "github.com/dgnabasik/acmsearchlib/database"
	fs "github.com/dgnabasik/acmsearchlib/filesystem"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"
)

func main() {
	_, mostRecentArchiveDate, err := dbase.GetLastDateSavedFromDb()
	source, _ := fs.ReadTextLines("config.go", false)
	_, found2 = hd.SearchForStringIndex(source, "main")
	journalDate = nt.New_NullTime("2020-12-12")
	fmt.Println(journalDate)
}
