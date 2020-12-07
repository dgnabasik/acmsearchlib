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
	_, mostRecentArchiveDate, _ := dbase.GetLastDateSavedFromDb()
	fmt.Println(mostRecentArchiveDate)
	source, _ := fs.ReadTextLines("config.goo", false)
	fmt.Println(source)
	_, found2 := hd.SearchForStringIndex("one two", "two")
	fmt.Println(found2)
	journalDate := nt.New_NullTime("2020-12-12")
	fmt.Println(journalDate)
}
