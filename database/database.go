package database

/* Do not communicate by sharing memory; instead, share memory by communicating.
   for line:=0; line < len(source); line++ {.}	for k, v := range data {.}; for i := range myconfig{.}

   Channels are a typed conduit through which you can send and receive values with the channel operator, <-.
   By default, sends and receives block until the other side is ready.
   Sends to a buffered channel block only when the buffer is full.
   Closing a channel: v, ok := <-ch  Only the sender should close a channel, never the receiver.
   Closing is only necessary when the receiver must be told there are no more values coming.

   Goroutines run in the same address space, so access to shared memory must be synchronized. A goroutine is context-switched over an OS thread, not a CPU core.
   The Go scheduler (which runs in user space) is cooperative (not preemptive) and uses a work-stealing (not work-sharing) scheduling strategy.
   The select (case) statement lets a goroutine wait on multiple communication operations. A select blocks until one of its cases can run, then it executes that case.

   The compiler uses a technique called escape analysis to decide if a variable is going to be placed on the heap or the stack, but new always allocates on the heap.
   if the compiler cannot prove that the variable is not referenced after the function returns, then the compiler must allocate the variable on the garbage-collected heap to avoid dangling pointer errors. If you need to know where your variables are allocated pass the "-m" gc flag to "go build" or "go run" (e.g., go run -gcflags -m app.go).
   Most memory allocations are served from local thread caches.

   Ddatabase driver: go get -u github.com/lib/pq	(_) include this package even though the package is not explicitly referenced in code.
   pq driver: NullTime implements the sql.Scanner interface so it can be used as a scan destination, similar to sql.NullString.
   s.p. inserts into table Occurrence. The defer statement should come after you check for an error from DB.Query.
*/

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	// comment
	_ "github.com/lib/pq"
)

// mapset https://github.com/deckarep/golang-set/blob/master/README.md & https://godoc.org/github.com/deckarep/golang-set

// DB struct
// dbRef, err := dbase.GetDatabaseReference()
// dbObj := &ArticleDatastore{db: dbRef}
type DB struct {
	*sql.DB
}

// CheckErr database error handler.
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

// GetDatabaseReference opens a database specified by its database driver name and a driver-specific data source name: db,err := GetDatabaseReference()
// defer db.Close() must follow a call to this function in the calling function. sslmode is set to 'required' with lib/pq by default.
func GetDatabaseReference() (*sql.DB, error) {
	const (
		dbHost        = "localhost"
		dbPort        = 5432
		dbUser        = "postgres"
		dbPassword    = "Ski7Vail!"
		dbName        = "postgres"
		dbDriver      = "postgres"
		dbSchema      = "acmsearch"
		dbConnections = 10
	)

	dbConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSchema)

	db, err := sql.Open(dbDriver, dbConn)
	CheckErr(err)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(dbConnections)
	db.SetConnMaxLifetime(0)
	err = db.Ping() // connects
	CheckErr(err)
	return db, err
}

// CallTruncateTables truncates tables with sequences.
func CallTruncateTables() error {
	db, err := GetDatabaseReference()
	defer db.Close()

	_, err = db.Exec("call TruncateTables()")
	CheckErr(err)

	fmt.Println("CallTruncateTables() done.")
	return nil
}

func isHexWord(word string) bool {
	_, err := hex.DecodeString(word)
	return len(word) >= 10 && err == nil
}

// FilteringRules filters output from Postgres ts_stat select. Include 3d prefixes.
// Return 0 for ok, -1 to completely ignore, 1 for modified word.
func FilteringRules(word string) (string, int) {
	if len(strings.TrimSpace(word)) <= 1 || word == "says" {
		return word, -1
	}

	ignore := strings.HasPrefix(word, "0") || strings.HasPrefix(word, "1") || strings.HasPrefix(word, "2") || (strings.HasPrefix(word, "3") && !strings.HasPrefix(word, "3d")) || strings.HasPrefix(word, "4") || strings.HasPrefix(word, "5") || strings.HasPrefix(word, "6") || strings.HasPrefix(word, "7") || strings.HasPrefix(word, "8") || strings.HasPrefix(word, "9") || strings.HasPrefix(word, "-") || strings.HasPrefix(word, "+") || strings.Count(word, "/") > 1 || strings.Count(word, "_") > 1 || strings.HasPrefix(word, "www.") || strings.HasSuffix(word, ".com") || strings.HasSuffix(word, ".org") || isHexWord(word)
	if ignore {
		return word, -1
	}

	newWord := word // Remove leading/trailing . /
	if strings.HasPrefix(newWord, ".") || strings.HasPrefix(newWord, "/") {
		newWord = newWord[1:]
	}
	if strings.HasSuffix(newWord, ".") || strings.HasSuffix(newWord, "/") {
		newWord = newWord[:len(newWord)-1]
	}

	if newWord != word {
		return newWord, 1
	}

	return word, 0
}

// GetArticleCount func
func GetArticleCount() int {
	db, err := GetDatabaseReference()
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM AcmData;").Scan(&count)
	CheckErr(err)
	return count
}

// GetLastDateSavedFromDb returns the earliest and latest AcmData.ArchiveDate values else default time.
func GetLastDateSavedFromDb() (nt.NullTime, nt.NullTime, error) {
	articleCount := GetArticleCount()
	if articleCount == 0 {
		return nt.New_NullTime(""), nt.New_NullTime(""), nil // default time.
	}

	db, err := GetDatabaseReference()
	defer db.Close()

	var archiveDate1, archiveDate2 nt.NullTime // NullTime supports Scan() interface.

	err = db.QueryRow("SELECT MIN(ArchiveDate) FROM AcmData;").Scan(&archiveDate1)
	CheckErr(err)

	err = db.QueryRow("SELECT MAX(ArchiveDate) FROM AcmData;").Scan(&archiveDate2)
	CheckErr(err)

	return archiveDate1, archiveDate2, nil
}

// SelectOccurrenceByDate assumes NullTime have zero hours, min, secs, so subtract 1 minute from startDate and add 1 minute to endDate to avoid any time issues. Also assumes occurrenceList is sorted by ArchiveDate
func SelectOccurrenceByDate(occurrenceList []hd.Occurrence, timeinterval nt.TimeInterval) []hd.Occurrence {
	var subList []hd.Occurrence
	//sDate := StartDate.DT.Add(time.Minute * -1)
	//eDate := timeinterval.EndDate.DT.Add(time.Minute * 1)

	for ndx := 0; ndx < len(occurrenceList); ndx++ {
		if occurrenceList[ndx].ArchiveDate.DT.Before(timeinterval.StartDate.DT) {
			continue
		}
		if occurrenceList[ndx].ArchiveDate.DT.Equal(timeinterval.StartDate.DT) {
			subList = append(subList, occurrenceList[ndx])
		}
		if occurrenceList[ndx].ArchiveDate.DT.After(timeinterval.StartDate.DT) && occurrenceList[ndx].ArchiveDate.DT.Before(timeinterval.EndDate.DT) {
			subList = append(subList, occurrenceList[ndx])
		}
		if occurrenceList[ndx].ArchiveDate.DT.Equal(timeinterval.StartDate.DT) {
			subList = append(subList, occurrenceList[ndx])
		}
		if occurrenceList[ndx].ArchiveDate.DT.After(timeinterval.EndDate.DT) {
			break
		}
	}
	return subList
}

// SelectOccurrenceByID assumes occurrenceList is sorted by AcmId.
func SelectOccurrenceByID(occurrenceList []hd.Occurrence, acmId uint32) []hd.Occurrence {
	var subList []hd.Occurrence
	for ndx := 0; ndx < len(occurrenceList); ndx++ {
		if acmId < occurrenceList[ndx].AcmId {
			continue
		}
		if acmId == occurrenceList[ndx].AcmId {
			subList = append(subList, occurrenceList[ndx])
		}
		if acmId > occurrenceList[ndx].AcmId {
			break
		}
	}
	return subList
}

// SelectOccurrenceByWord assumes occurrenceList is sorted by Word.
func SelectOccurrenceByWord(occurrenceList []hd.Occurrence, word string) []hd.Occurrence {
	var subList []hd.Occurrence
	for ndx := 0; ndx < len(occurrenceList); ndx++ {
		order := strings.Compare(occurrenceList[ndx].Word, word)
		if order < 0 {
			continue
		}
		if order == 0 {
			subList = append(subList, occurrenceList[ndx])
		}
		if order > 0 {
			break
		}
	}
	return subList
}

// CollectWordGrams collects all the words from the same summaries for each word in wordGrams, which is usually the set of words in 1 summary.
// append() does not force new memory allocations each time it is called. This allows users to append inside a loop without thrashing the GC.
// By adding val as a parameter to the closure, val is evaluated at each iteration and placed on the stack for the goroutine, so each slice element is available to the goroutine when it is eventually executed.
// Variables declared within the body of a loop are not shared between iterations, and thus can be used separately in a closure.
// https://medium.com/@cep21/gos-append-is-not-always-thread-safe-a3034db7975
// Do not use shared state as the first variable to append.
// Explicitly make() a new slice with an extra element's worth of capacity, then copy() the old slice to it, then finally append() or add the new value.
func CollectWordGrams(wordGrams []string, timeinterval nt.TimeInterval) ([]hd.Occurrence, mapset.Set) {
	start := time.Now()
	fmt.Print("CollectWordGrams: ")

	var alphaCollection []hd.Occurrence                               // populate in separate goroutine using queue channel.
	occurrenceList, idSet, _ := GetOccurrenceListByDate(timeinterval) // []Occurrence
	sort.Sort(hd.OccurrenceSorterWord(occurrenceList))

	/* serial version:
	for _, word := range wordGrams {
		wordOccurrenceList := SelectOccurrenceByWord(occurrenceList, word)	// []Occurrence
		alphaCollection = append(alphaCollection, wordOccurrenceList...)
	} */

	queue := make(chan hd.Occurrence, 32768) // select count(*) from occurrence where word='says' ==> 26273
	var wg sync.WaitGroup

	for _, word := range wordGrams {
		wg.Add(1)

		go func(word string) {
			defer wg.Done()                                                    // Decrement the counter when the goroutine completes.
			wordOccurrenceList := SelectOccurrenceByWord(occurrenceList, word) // []Occurrence
			for _, wo := range wordOccurrenceList {
				queue <- wo // avoid data race condition.  queue <- Occurrence(i)
			}
		}(word)
	}

	go func() {
		wg.Wait()
		close(queue)
	}()

	for t := range queue {
		if len(strings.TrimSpace(t.Word)) > 0 {
			alphaCollection = append(alphaCollection, t)
		}
	}

	elapsed := time.Since(start)
	fmt.Println(elapsed.String())

	//_ = OutputOccurrenceList(alphaCollection, prefixProcessed + "processed.txt")

	return alphaCollection, idSet
}

// WordGramSubset finds words that OtherWord.IsProperSuperset(alphaWord) : they exist in all the summaries the alphaWord does.
func WordGramSubset(alphaWord string, vocabList []hd.Vocabulary, occurrenceList []hd.Occurrence) []string {
	var wordList []string
	// Assumes ranked vocabList, so start from top (most frequent).
	var alphaVocabList []hd.Vocabulary
	for _, vocab := range vocabList {
		alphaVocabList = append(alphaVocabList, vocab)
		if vocab.Word == alphaWord {
			break // always includes the alphaWord as last in list.
		}
	}

	// build map of all acmIds for each word.
	wordIDMap := make(map[string]mapset.Set) // {word, Set of acmIds}
	for _, vocab := range alphaVocabList {
		idSet := mapset.NewSet() // a word has to exist in every one of these
		for _, item := range occurrenceList {
			if item.Word == vocab.Word {
				idSet.Add(item.AcmId)
			}
		}
		wordIDMap[vocab.Word] = idSet
	}

	// Test with IsProperSuperset().
	alphaIDSet := wordIDMap[alphaWord]
	for key, value := range wordIDMap {
		fmt.Printf("%s : %t\n", key, value.IsProperSuperset(alphaIDSet))
		if value.IsProperSuperset(alphaIDSet) {
			wordList = append(wordList, key)
		}
	}

	if len(wordList) == 0 {
		fmt.Println("There are no wordgram supersets for '" + alphaWord + "'.")
	}

	return wordList
}

// GetOccurrencesByAcmid func
func GetOccurrencesByAcmid(xacmid uint32) ([]hd.Occurrence, error) {
	db, err := GetDatabaseReference()
	defer db.Close()

	SELECT := "SELECT acmId, archiveDate, word, nentry FROM Occurrence WHERE acmId=" + strconv.FormatUint(uint64(xacmid), 10) + ";"
	rows, err := db.Query(SELECT)
	CheckErr(err)
	defer rows.Close()

	// fields to read
	var acmID uint32
	var archiveDate nt.NullTime
	var word string
	var nentry int
	var occurrenceList []hd.Occurrence

	for rows.Next() {
		err = rows.Scan(&acmID, &archiveDate, &word, &nentry)
		CheckErr(err)

		newWord, rule := FilteringRules(word)
		if rule < 0 {
			continue
		} else if rule > 0 {
			word = newWord
		}

		occurrenceList = append(occurrenceList, hd.Occurrence{AcmId: acmID, ArchiveDate: archiveDate, Word: word, Nentry: nentry})
	}

	// get any iteration errors
	err = rows.Err()
	CheckErr(err)

	return occurrenceList, nil
}

// GetOccurrenceListByDate returns result set ordered by ArchiveDate.
// Read []Occurrence values by archiveDate range. This applys FilteringRules(word).
// mapset.Set is the set of distinct AcmId values in the returned list.
func GetOccurrenceListByDate(timeinterval nt.TimeInterval) ([]hd.Occurrence, mapset.Set, error) {
	db, err := GetDatabaseReference()
	defer db.Close()
	// Invoke stored procedure.
	SELECT := "SELECT * FROM GetOccurrencesByDate('" + timeinterval.StartDate.StandardDate() + "', '" + timeinterval.EndDate.StandardDate() + "');"
	rows, err := db.Query(SELECT)
	CheckErr(err)
	defer rows.Close()

	// fields to read
	var acmID uint32
	var archiveDate nt.NullTime
	var word string
	var nentry int
	var occurrenceList []hd.Occurrence
	idSet := mapset.NewSet()

	for rows.Next() {
		err = rows.Scan(&acmID, &archiveDate, &word, &nentry)
		CheckErr(err)

		newWord, rule := FilteringRules(word)
		if rule < 0 || len(newWord) <= 1 {
			continue
		} else if rule > 0 {
			word = newWord
		}
		idSet.Add(acmID)
		occurrenceList = append(occurrenceList, hd.Occurrence{AcmId: acmID, ArchiveDate: archiveDate, Word: word, Nentry: nentry})
	}

	// get any iteration errors
	err = rows.Err()
	CheckErr(err)

	return occurrenceList, idSet, nil
}

// GetDistinctDates Returns ordered list of distinct dates. Assumes all dates normalized to midnight.
func GetDistinctDates(occurrenceList []hd.Occurrence) []nt.NullTime {
	dateMap := make(map[nt.NullTime]int)
	for _, item := range occurrenceList {
		dateMap[item.ArchiveDate] = 0
	}

	var dateSet []nt.NullTime
	for nt := range dateMap {
		dateSet = append(dateSet, nt)
	}

	dateSet = nt.NullTimeSorter(dateSet) // sort in-place
	return dateSet
}

// GetDistinctWords func
func GetDistinctWords(occurrenceList []hd.Occurrence) []string {
	var wordMap map[string]int
	for _, item := range occurrenceList {
		wordMap[item.Word] = 0
	}

	var wordSet []string
	for w := range wordMap {
		wordSet = append(wordSet, w)
	}

	sort.Strings(wordSet)
	return wordSet
}

// GetIDSetUnion returns union of Ids in mapset.
func GetIDSetUnion(stringMapset map[string]mapset.Set) mapset.Set {
	idSet := mapset.NewSet()
	for _, item := range stringMapset {
		idSet = idSet.Union(item)
	}
	return idSet
}

// GetIDSetIntersection NOT USED
func GetIDSetIntersection(stringMapset map[string]mapset.Set) mapset.Set {
	if len(stringMapset) == 0 {
		return mapset.NewSet()
	}
	// populate first set; can't intersect with empty set.
	idSet := mapset.NewSet()
	for _, item := range stringMapset {
		idSet = idSet.Union(item)
		break
	}

	for _, item := range stringMapset {
		idSet = idSet.Intersect(item)
	}
	return idSet
}

// GetMinMaxSetValues For sets of ints
func GetMinMaxSetValues(idSet mapset.Set) (uint32, uint32) {
	min := uint32(4 * 1073741823)
	max := uint32(0)
	it := idSet.Iterator()

	for k := range it.C {
		if k.(uint32) < min {
			min = k.(uint32)
		}
		if k.(uint32) > max {
			max = k.(uint32)
		}
	}
	return min, max
}

// ExtractIDSet func
func ExtractIDSet(word string, stringMapset map[string]mapset.Set) mapset.Set {
	idSet := mapset.NewSet()
	for key, item := range stringMapset {
		if key == word {
			it := item.Iterator()
			for k := range it.C {
				idSet.Add(k.(uint32))
			}
			break
		}
	}

	return idSet
}

// GetIDSetForWordGrams fails with "invalid memory address or nil pointer dereference" if a space is in words.
func GetIDSetForWordGrams(wordGrams []string, occurrenceList []hd.Occurrence) map[string]mapset.Set {
	wordIDMap := make(map[string]mapset.Set, len(wordGrams)) // {word, Set of acmIds}
	for _, word := range wordGrams {
		wordIDMap[word] = mapset.NewSet()
	}

	for _, item := range occurrenceList {
		if _, found := wordIDMap[item.Word]; !found {
			fmt.Println("Key not found for: " + item.Word)
		} else {
			wordIDMap[item.Word].Add(item.AcmId)
		}
	}

	return wordIDMap
}

/*************************************************************************************************/

// BulkInsertConditionalProbability uses prepared statement.
func BulkInsertConditionalProbability(conditionals []hd.ConditionalProbability) error {
	db, err := GetDatabaseReference()
	defer db.Close()

	txn, err := db.Begin()
	CheckErr(err)

	stmt, err := db.Prepare("INSERT INTO Conditional (wordlist, probability, startDate, endDate, firstDate, lastDate) VALUES ($1, $2, $3, $4, $5, $6, $7);")
	CheckErr(err)

	for _, v := range conditionals {
		_, err = stmt.Exec(v.WordList, v.Probability, v.Timeinterval.Timeframetype, v.Timeinterval.StartDate.DT, v.Timeinterval.EndDate.DT, v.FirstDate.DT, v.LastDate.DT)
		CheckErr(err)
	}

	err = stmt.Close()
	CheckErr(err)

	err = txn.Commit()
	CheckErr(err)

	return nil
}

// ExtractKeysFromProbabilityMap func
func ExtractKeysFromProbabilityMap(wordMap map[string]float32) []string {
	words := make([]string, 0)
	for word := range wordMap {
		words = append(words, word)
	}
	return words
}

// GetIDArchiveDateMap func
func GetIDArchiveDateMap(timeinterval nt.TimeInterval) (map[uint32]nt.NullTime, error) {
	var archiveDate nt.NullTime
	var id uint32
	dateMap := make(map[uint32]nt.NullTime)

	db, err := GetDatabaseReference()
	defer db.Close()

	SELECT := "SELECT Id, ArchiveDate FROM AcmData WHERE ArchiveDate >= '" + timeinterval.StartDate.StandardDate() + "' AND ArchiveDate <= '" + timeinterval.EndDate.StandardDate() + "';"
	rows, err := db.Query(SELECT)
	CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &archiveDate)
		CheckErr(err)
		dateMap[id] = archiveDate
	}

	// get any iteration errors
	err = rows.Err()
	CheckErr(err)

	return dateMap, err
}

// CalcConditionalProbability P(dependent A and B both occurring): Bayes: P(A|B)=P(A∩B)/P(B)=P(B|A)P(A)/P(B)
// What is the P of word A given word B (in this interval)?
// The imported wordMap has probabilities over entire TFSpan.
// P(A∩B)=P(A|B)*P(B) is the probability that both events A and B occur; they are present in the same summary.
// If P(A|B)=P(A) then events A and B are said to be independent.
// Do for 2 permutations (order matters). Performs FilteringRules(words) Returns len(wordGrams).
// MICROSERVICE! Number of permutations for 94322 wordgrams is n!/(n-r)! = 8,896,545,362
// Estimated completion time is 1572 hours.
func CalcConditionalProbability(wordMap map[string]float32, timeinterval nt.TimeInterval) int {
	if len(wordMap) < 2 {
		fmt.Println("There must at at least 2 words to compute conditional probabilities.")
		return 0
	}
	permutations := 2
	var cutoffProb float32 = 0.00001 // arbitrary

	wordGrams := ExtractKeysFromProbabilityMap(wordMap) // []string
	sort.Strings(wordGrams)

	if len(wordGrams) < 10 {
		fmt.Println("Processing: " + strings.Join(wordGrams, " + "))
	}

	wordOccurrenceList, totalIDSet := CollectWordGrams(wordGrams, timeinterval) // []Occurrence

	// sort.Sort(OccurrenceSorterId(wordOccurrenceList))  not needed

	start := time.Now()
	fmt.Print("CalcConditionalProbability: ")

	wordIDSets := GetIDSetForWordGrams(wordGrams, wordOccurrenceList) // map[string]mapset.Set
	idDateMap, _ := GetIDArchiveDateMap(timeinterval)                 // map[uint32]NullTime

	var conditionals []hd.ConditionalProbability
	var wordlist string

	// P(A|B)=P(A∩B)/P(B)    P(wordA|wordB) = for those summaries containing wordB, how many contain wordA => intersection
	if permutations == 2 {
		// timeinterval := nt.TimeInterval{Timeinterval startDate, endDate}
		for wordA := 0; wordA < len(wordGrams)-1; wordA++ {
			wordIdSetX := ExtractIDSet(wordGrams[wordA], wordIDSets) // mapset.Set
			conditionals = nil
			// MICROSERVICE: wordMap(94k), wordGrams(94k), wordIDSets(94k), idDateMap, totalIDSet.Cardinality(), wordA, wordB, startDate, endDate => conditionals
			for wordB := wordA + 1; wordB < len(wordGrams); wordB++ {
				wordIdSetY := ExtractIDSet(wordGrams[wordB], wordIDSets) // mapset.Set
				idSetIntersection := wordIdSetX.Intersect(wordIdSetY)
				if idSetIntersection.Cardinality() == 0 {
					//fmt.Println("Empty intersection for " + wordGrams[wordA] + " + " + wordGrams[wordB])
					continue
				}

				minId, maxId := GetMinMaxSetValues(idSetIntersection)
				firstDate := idDateMap[minId]
				lastDate := idDateMap[maxId]

				pAgivenB := (float32(idSetIntersection.Cardinality()) / float32(totalIDSet.Cardinality())) / wordMap[wordGrams[wordB]] // P(wordA ∩ wordB) / P(wordB)
				if pAgivenB > cutoffProb {
					wordlist = wordGrams[wordA] + "|" + wordGrams[wordB]
					conditionals = append(conditionals, hd.ConditionalProbability{Id: 0, WordList: wordlist, Probability: pAgivenB, Timeinterval: timeinterval, FirstDate: firstDate, LastDate: lastDate})
				}

				pBgivenA := (float32(idSetIntersection.Cardinality()) / float32(totalIDSet.Cardinality())) / wordMap[wordGrams[wordA]] // P(wordB ∩ wordA) / P(wordA)
				if pBgivenA > cutoffProb {
					wordlist = wordGrams[wordB] + "|" + wordGrams[wordA]
					conditionals = append(conditionals, hd.ConditionalProbability{Id: 0, WordList: wordlist, Probability: pBgivenA, Timeinterval: timeinterval, FirstDate: firstDate, LastDate: lastDate})
				}
			}

			if len(conditionals) > 0 {
				err := BulkInsertConditionalProbability(conditionals)
				CheckErr(err)
			}
		}
	}

	elapsed := time.Since(start)
	fmt.Println(elapsed.String())

	return len(wordGrams)
}

// GetWhereClause Don't know PostgreSQL limit of IN values.
func GetWhereClause(columnName string, wordGrams []string) string {
	var sb strings.Builder
	sb.WriteString(" WHERE " + columnName + " IN (")
	for ndx := 0; ndx < len(wordGrams); ndx++ {
		sb.WriteString("'" + wordGrams[ndx] + "'")
		if ndx < len(wordGrams)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString(");")
	return sb.String()
}

// GetVocabularyMapProbability Read all Vocabulary.Word,Probability values if wordGrams is empty. Applys filtering.
func GetVocabularyMapProbability(wordGrams []string) (map[string]float32, error) {
	db, err := GetDatabaseReference()
	defer db.Close()

	wordIDMap := make(map[string]float32)
	var word string
	var floatField float32

	SELECT := "SELECT Word,Probability FROM vocabulary"
	if len(wordGrams) > 0 {
		SELECT = SELECT + GetWhereClause("Word", wordGrams)
	} else {
		SELECT = SELECT + ";"
	}

	rows, err := db.Query(SELECT)
	CheckErr(err)

	for rows.Next() {
		err = rows.Scan(&word, &floatField)
		CheckErr(err)

		newWord, rule := FilteringRules(word)
		if rule < 0 {
			continue
		} else if rule > 0 {
			word = newWord
		}

		wordIDMap[word] = floatField
	}

	// get any iteration errors
	err = rows.Err()
	CheckErr(err)

	return wordIDMap, err
}

// GetVocabularyByWord func
func GetVocabularyByWord(wordX string) hd.Vocabulary {
	db, err := GetDatabaseReference()
	defer db.Close()

	var word, speechPart string
	var id uint32
	var rowCount, frequency, wordRank int
	var probability float32
	SELECT := "SELECT id, word, rowcount, frequency, wordrank, probability, speechpart FROM Vocabulary WHERE Word='" + wordX + "';"
	err = db.QueryRow(SELECT).Scan(&id, &word, &rowCount, &frequency, &wordRank, &probability, &speechPart)
	CheckErr(err)

	return hd.Vocabulary{Id: id, Word: word, RowCount: rowCount, Frequency: frequency, WordRank: wordRank, Probability: probability, SpeechPart: speechPart}
}

// GetVocabularyListByDate reads Vocabulary table filtered by articleList.
func GetVocabularyListByDate(timeinterval nt.TimeInterval) ([]hd.Vocabulary, error) {
	start := time.Now()
	fmt.Print("GetVocabularyListByDate() ")

	db, err := GetDatabaseReference()
	defer db.Close()

	SELECT := "SELECT * FROM GetVocabularyByDate('" + timeinterval.StartDate.StandardDate() + "', '" + timeinterval.EndDate.StandardDate() + "');"
	// SELECT * FROM vocabulary WHERE word IN (SELECT word FROM GetOccurrencesByDate(startDate, endDate));
	rows, err := db.Query(SELECT)
	CheckErr(err)
	defer rows.Close()

	// fields to read
	var word, speechPart string
	var id uint32
	var rowCount, frequency, wordRank int
	var probability float32
	var vocabList []hd.Vocabulary

	for rows.Next() { // this order follows the \d Vocabulary description:
		err = rows.Scan(&id, &word, &rowCount, &frequency, &wordRank, &probability, &speechPart)
		CheckErr(err)

		newWord, rule := FilteringRules(word)
		if rule < 0 {
			continue
		} else if rule > 0 {
			word = newWord
		}

		vocabList = append(vocabList, hd.Vocabulary{Id: id, Word: word, RowCount: rowCount, Frequency: frequency, WordRank: wordRank, Probability: probability, SpeechPart: speechPart})
	}

	// get any iteration errors
	err = rows.Err()
	CheckErr(err)

	elapsed := time.Since(start)
	fmt.Println(elapsed.String())

	return vocabList, err
}

// BulkUpdateVocabularySpeechpart concatentates parts into output. Unknown list is returned.
// Change to root word. See https://www.datamuse.com/api/ & https://www.wordsapi.com/
// curl "https://wordsapiv1.p.mashape.com/words/soliloquy" -H "X-Mashape-Key: <APIkey>"
/*
func BulkUpdateVocabularySpeechpart() []string {
	fmt.Print("BulkInsertVocabulary_Speechpart: ")
	var wordNetSpeechParts = hd.New_WordNetSpeechParts()
	var word string
	var wordSet []string

	start := time.Now()
	db, err := GetDatabaseReference()
	defer db.Close()

	// part 1: WHERE SpeechPart not assigned.
	SELECT := "SELECT word FROM Vocabulary WHERE SpeechPart='';"
	rows, err := db.Query(SELECT)
	CheckErr(err)
	for rows.Next() {
		err = rows.Scan(&word)
		CheckErr(err)
		wordSet = append(wordSet, word)
	}

	// part 2:
	txn, err := db.Begin()
	CheckErr(err)

	stmt, err := db.Prepare("UPDATE vocabulary SET SpeechPart= $1 WHERE Word= $2;")
	CheckErr(err)

	for _, w := range wordSet {
		speechPart := wordNetSpeechParts.GetSpeechpart(w)
		if speechPart == "" {
			wordNetSpeechParts.Unknown = append(wordNetSpeechParts.Unknown, w)
		}
		_, err = stmt.Exec(speechPart, w)
		CheckErr(err)
	}

	err = stmt.Close()
	CheckErr(err)

	err = txn.Commit()
	CheckErr(err)

	elapsed := time.Since(start)
	fmt.Println(elapsed.String())

	return wordNetSpeechParts.Unknown
}
*/
/*************************************************************************************************/

// GetAcmArticleListByDate PostgreSql allows for defining a generic get-all-rows stored proc and appending the WHERE clause to the select instead of defining the WHERE clause inside the stored proc, but it is slower.
func GetAcmArticleListByDate(timeinterval nt.TimeInterval) ([]hd.AcmArticle, error) {
	db, err := GetDatabaseReference()
	defer db.Close()

	SELECT := "SELECT * FROM GetAcmArticles() WHERE ArchiveDate >= '" + timeinterval.StartDate.StandardDate() + "' AND ArchiveDate <= '" + timeinterval.EndDate.StandardDate() + "';"
	rows, err := db.Query(SELECT)
	CheckErr(err)
	defer rows.Close()

	// fields to read
	var id uint32
	var archiveDate, journalDate nt.NullTime
	var articleNumber, title, imageSource, journalName, authorName, webReference, summary string

	var articleList []hd.AcmArticle

	for rows.Next() { // this order follows the \d AcmData description:
		err = rows.Scan(&id, &archiveDate.DT, &articleNumber, &title, &imageSource, &journalName, &authorName, &journalDate.DT, &webReference, &summary)
		CheckErr(err)
		articleList = append(articleList, hd.AcmArticle{Id: id, ArchiveDate: archiveDate, ArticleNumber: articleNumber, Title: title, ImageSource: imageSource, JournalName: journalName, AuthorName: authorName, JournalDate: journalDate, WebReference: webReference, Summary: summary})
	}

	// get any iteration errors
	err = rows.Err()
	CheckErr(err)

	return articleList, err
}

// ProcessAcm Re-calculate all conditional probabilities and wordscores. <<<<
func ProcessAcm() error {
	/*
		wordMap, err := db.GetVocabularyMapProbability( []string{} )
		n = db.CalcConditionalProbability(wordMap, timeInterval)
		err = db.CalcAllWordScoresConcurrently(timeInterval)
		if err != nil {
			fmt.Println(err)
			return err
		}

		articles, err := db.GetAcmArticleListByDate(timeInterval) // []AcmArticle

		uploadScript, err := rg.ProduceArticleInsertRDF(articles) // defined in rdfgraph.go
		if err != nil {
			fmt.Println(err)
			return err
		} else {
			fmt.Println("Upload script 1: " + uploadScript)
		}

		// The time it takes to run this query (20s) allows the O/S to properly delete the target folder.
		vocabList, err := db.GetVocabularyListByDate(timeInterval) // []Vocabulary

		uploadScript, err = rg.ProduceVocabularyInsertRDF(vocabList)
		if err != nil {
			fmt.Println(err)
			return err
		} else {
			fmt.Println("Upload script 2: " + uploadScript)
		} */

	return nil
}
