package category

// category.go manages categories. Derived from ~/websites/dropdownlists/golang/listservice.go

import (
	"strconv"
	"time"

	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"

	"github.com/lib/pq"
)

/*************************************************************************************/

// Version func
func Version() string {
	return "1.16.2"
}

// InsertCategoryWords func. 32k statement limit.
func InsertCategoryWords(categoryID uint64, words []string) error {
	dateupdated := time.Now()

	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return err
	}
	defer db.Close()

	txn, err := db.Begin()
	dbx.CheckErr(err)

	// Must use lowercase column names! First param is table name.
	stmt, err := txn.Prepare(pq.CopyIn("special", "word", "category", "dateupdated"))
	dbx.CheckErr(err)

	for _, word := range words {
		_, err = stmt.Exec(word, categoryID, dateupdated)
		dbx.CheckErr(err)
	}

	_, err = stmt.Exec()
	dbx.CheckErr(err)

	err = stmt.Close()
	dbx.CheckErr(err)

	err = txn.Commit()
	dbx.CheckErr(err)

	return nil
}

// InsertWordCategory func
func InsertWordCategory(description string) (hd.CategoryTable, error) {
	dateupdated := time.Now()
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return hd.CategoryTable{}, err
	}
	defer db.Close()

	var id uint64
	INSERT := "INSERT INTO Wordcategory (description, dateupdated) VALUES ($1, $2) returning id"
	err = db.QueryRow(INSERT, description, dateupdated).Scan(&id)
	dbx.CheckErr(err)

	categoryTable := hd.CategoryTable{Id: id, Description: description, DateUpdated: dateupdated}
	return categoryTable, nil
}

// GetSpecialMap func filters by category
func GetSpecialMap(category int) ([]hd.SpecialTable, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	SELECT := "SELECT id, word, category, dateupdated FROM Special WHERE category=" + strconv.Itoa(category) + " ORDER BY word"
	rows, err := db.Query(SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var special hd.SpecialTable
	specialMap := []hd.SpecialTable{}
	for rows.Next() {
		err = rows.Scan(&special.Id, &special.Word, &special.Category, &special.DateUpdated)
		dbx.CheckErr(err)
		specialMap = append(specialMap, special)
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return specialMap, nil
}

// GetCategoryMap func
func GetCategoryMap() ([]hd.CategoryTable, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	SELECT := "SELECT id, description, dateupdated FROM Wordcategory"
	rows, err := db.Query(SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var catnap hd.CategoryTable
	categoryMap := []hd.CategoryTable{}
	for rows.Next() {
		err = rows.Scan(&catnap.Id, &catnap.Description, &catnap.DateUpdated)
		dbx.CheckErr(err)
		categoryMap = append(categoryMap, catnap)
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return categoryMap, nil
}
