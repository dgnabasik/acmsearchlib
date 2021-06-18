package basicgraph

// graph.go GraphStructure handles graph queries.
// Every hypergraph may be represented by a bipartite graph.
import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	"github.com/jackc/pgx/v4" // The vertices of all graphs in this package are numbered 0..n-1.
)

/*************************************************************************************/

func getTableNames(useTempTable bool) []string {
	tableIndex := 0
	if useTempTable {
		tableIndex = 1
	}
	return []string{[]string{"Simplex", "temp_Simplex"}[tableIndex], []string{"Facet", "temp_Facet"}[tableIndex]}
}

// GetSimplexByNameUserID func : simplexName, simplextype are case-sensitive.
func GetSimplexByNameUserID(userID int, simplexName, simplexType string, useTempTable bool) ([]hd.SimplexComplex, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tableNames := getTableNames(useTempTable)
	query := `SELECT s.ID, s.UserID, s.SimplexName, s.SimplexType, s.EulerCharacteristic, s.Dimension, s.FiltrationValue, s.NumSimplices, s.NumVertices, s.BettiNumbers, s.Timeframetype, 
		s.StartDate, s.EndDate, s.Enabled, s.DateCreated, s.DateUpdated, f.ComplexID, f.SourceVertexID, f.TargetVertexID, f.SourceWord, f.TargetWord, f.Weight FROM `
	query += tableNames[0] + " s RIGHT OUTER JOIN " + tableNames[1] + " f ON f.ComplexID=s.ID WHERE s.UserID=" + strconv.Itoa(userID) + " AND LOWER(s.SimplexName)='" +
		strings.ToLower(simplexName) + "' AND s.SimplexType='" + simplexType + "' ORDER BY s.StartDate, f.ComplexID"
	// query returns as many rows per simplex as there are facets.
	rows, err := db.Query(context.Background(), query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetSimplexByNameUserID(1): %+v\n", err)
		return nil, err
	}
	defer rows.Close()

	var sc hd.SimplexComplex
	var f hd.SimplexFacet
	var timeframetype int
	var oldID uint64
	var startDate, endDate time.Time
	scList := make([]hd.SimplexComplex, 0)
	facets := make([]hd.SimplexFacet, 0)

	// CreateSimplexComplex(simplexName, simplexType, facets, timeInterval, userProfile.ID, eulerCharacteristic, nDimension, nNumSimplices, nNumVertices, float32(filtrationValue))
	for rows.Next() {
		err := rows.Scan(&sc.ID, &sc.UserID, &sc.SimplexName, &sc.SimplexType, &sc.EulerCharacteristic, &sc.Dimension, &sc.FiltrationValue, &sc.NumSimplices, &sc.NumVertices, &sc.BettiNumbers,
			&timeframetype, &startDate, &endDate, &sc.Enabled, &sc.DateCreated, &sc.DateUpdated, &f.ComplexID, &f.SourceVertexID, &f.TargetVertexID, &f.SourceWord, &f.TargetWord, &f.Weight)
		if err != nil {
			log.Printf("GetSimplexByNameUserID(2): %+v\n", err)
			return nil, err
		}
		sc.Timeinterval = nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startDate), nt.New_NullTime2(endDate))
		if sc.ID != oldID {
			scList = append(scList, sc)
			oldID = sc.ID
		}
		facets = append(facets, f)
	}
	err = rows.Err()
	dbx.CheckErr(err)

	// insert facets into each simplex:
	for ndx := range scList {
		scList[ndx].FacetVector = make([]hd.SimplexFacet, 0)
		for _, f := range facets {
			if scList[ndx].ID == f.ComplexID {
				scList[ndx].FacetVector = append(scList[ndx].FacetVector, f)
			}
		}
	}

	return scList, err
}

// GetSimplexListByUserID func fetches all of a user's simplices but not the facets.
func GetSimplexListByUserID(userID int, useTempTable bool) ([]hd.SimplexComplex, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return []hd.SimplexComplex{}, err
	}
	defer db.Close()

	tableNames := getTableNames(useTempTable)
	query := `SELECT s.ID, s.UserID, s.SimplexName, s.SimplexType, s.EulerCharacteristic, s.Dimension, s.FiltrationValue, s.NumSimplices, s.NumVertices, s.BettiNumbers, s.Timeframetype, 
		s.StartDate, s.EndDate, s.Enabled, s.DateCreated, s.DateUpdated FROM `
	query += tableNames[0] + " s WHERE s.UserID=" + strconv.Itoa(userID) + " ORDER BY s.ID"

	rows, err := db.Query(context.Background(), query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetSimplexListByUserID(1): %+v\n", err)
		return []hd.SimplexComplex{}, err
	}
	defer rows.Close()

	complexes := make([]hd.SimplexComplex, 0)
	var s hd.SimplexComplex
	var timeframetype int
	var startDate, endDate time.Time

	for rows.Next() {
		err := rows.Scan(&s.ID, &s.UserID, &s.SimplexName, &s.SimplexType, &s.EulerCharacteristic, &s.Dimension, &s.FiltrationValue, &s.NumSimplices, &s.NumVertices, &s.BettiNumbers,
			&timeframetype, &startDate, &endDate, &s.Enabled, &s.DateCreated, &s.DateUpdated)
		if err != nil {
			log.Printf("GetSimplexListByUserID(2): %+v\n", err)
			return complexes, err
		}
		s.Timeinterval = nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startDate), nt.New_NullTime2(endDate))
		s.FacetVector = make([]hd.SimplexFacet, 0)
		complexes = append(complexes, s)
	}
	err = rows.Err()
	dbx.CheckErr(err)

	return complexes, err
}

// InsertSimplexComplex func. Insert row into temp_Simplex before inserting []hd.SimplexFacet rows with InsertSimplexFacets().
// Assigns hd.SimplexComplex.ID = each hd.SimplexFacet.ComplexID. temp_Simplex rows are NOT in StartDate order!
func InsertSimplexComplex(sc hd.SimplexComplex) (hd.SimplexComplex, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return sc, err
	}
	defer db.Close()

	var id uint64 // DateCreated & DateUpdated use default server time.
	INSERT := `INSERT INTO temp_Simplex (UserID, SimplexName, SimplexType, EulerCharacteristic, Dimension, FiltrationValue, NumSimplices, NumVertices, 
		BettiNumbers, Timeframetype, StartDate, EndDate, Enabled) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) returning id`

	err = db.QueryRow(context.Background(), INSERT, sc.UserID, sc.SimplexName, sc.SimplexType, sc.EulerCharacteristic, sc.Dimension, sc.FiltrationValue,
		sc.NumSimplices, sc.NumVertices, sc.BettiNumbers, sc.Timeinterval.Timeframetype, sc.Timeinterval.StartDate.DT, sc.Timeinterval.EndDate.DT, sc.Enabled).Scan(&id)
	dbx.CheckErr(err)

	sc.ID = id
	for ndx := 0; ndx < len(sc.FacetVector); ndx++ {
		sc.FacetVector[ndx].ComplexID = sc.ID
	}
	err = BulkInsertSimplexFacets(sc.FacetVector)

	return sc, err
}

// BulkInsertSimplexFacets func inserts into temp_Facet
func BulkInsertSimplexFacets(facets []hd.SimplexFacet) error {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return err
	}
	defer db.Close()

	txn, err := db.Begin(context.Background())
	dbx.CheckErr(err)

	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"temp_facet"}, // tablename
		[]string{"complexid", "sourcevertexid", "targetvertexid", "sourceword", "targetword", "weight"}, // Must use lowercase column names!
		pgx.CopyFromSlice(len(facets), func(i int) ([]interface{}, error) {
			return []interface{}{facets[i].ComplexID, facets[i].SourceVertexID, facets[i].TargetVertexID, facets[i].SourceWord, facets[i].TargetWord, facets[i].Weight}, nil
		}),
	)

	dbx.CheckErr(err)
	if copyCount == 0 {
		log.Printf("BulkInsertSimplexFacets: no rows inserted")
	}
	err = txn.Commit(context.Background())
	dbx.CheckErr(err)

	return nil
}

// GetSimplexWordDifference func returns words that are the same, gained, and lost between two SimplexComplex-Facet sets. KeyValueStringPair format: word|type={S,G,L}
func GetSimplexWordDifference(userID int, simplexName, simplexType string, useTempTable bool) (map[nt.TimeInterval][]hd.KeyValueStringPair, []hd.SimplexComplex, error) {
	simplexList, err := GetSimplexByNameUserID(userID, simplexName, simplexType, useTempTable) // result set ordered by StartDate
	if dbx.NoRowsReturned(err) || err != nil {
		return nil, nil, err
	}

	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	// PostgreSQL functions invoked with SELECT; stored procs invoked with CALL.
	var sourceword, wordtype string
	sMap := make(map[nt.TimeInterval][]hd.KeyValueStringPair)

	for ndx := 0; ndx < len(simplexList)-1; ndx++ {
		SELECT := "SELECT sourceword, wordtype FROM WordDifference(" + strconv.FormatUint(simplexList[ndx].ID, 10) + "," + strconv.FormatUint(simplexList[ndx+1].ID, 10) + ")"
		rows, err := db.Query(context.Background(), SELECT)
		dbx.CheckErr(err)
		defer rows.Close()

		kvsp := make([]hd.KeyValueStringPair, 0)
		for rows.Next() {
			err = rows.Scan(&sourceword, &wordtype)
			dbx.CheckErr(err)
			kvsp = append(kvsp, hd.KeyValueStringPair{Key: sourceword, Value: wordtype})
		}
		err = rows.Err()
		dbx.CheckErr(err)
		sMap[simplexList[ndx].Timeinterval] = kvsp
	}

	return sMap, simplexList, nil
}

// PostSimplexComplex func moves [temp_Simplex] & [temp_Facet] data into [Simplex] & [Facet] tables.  Returns number of simplexes moved.
func PostSimplexComplex(userID int, simplexName, simplexType string) (int, error) {
	simplexList, err := GetSimplexByNameUserID(userID, simplexName, simplexType, true) // useTempTable
	dbx.CheckErr(err)

	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var newsimplexid uint64
	count := 0
	for _, simplex := range simplexList {
		err = db.QueryRow(context.Background(), "SELECT newsimplexid FROM PostSimplexComplex("+strconv.FormatUint(simplex.ID, 10)+")").Scan(&newsimplexid)
		dbx.CheckErr(err)
		if err != nil {
			count++
		}
	}

	return count, err
}
