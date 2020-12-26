package conditional

//  manages conditional probabilities and occurrences.

import (
	"encoding/hex"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set"
	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	// comment
	_ "github.com/lib/pq"
)

// mapset https://github.com/deckarep/golang-set/blob/master/README.md & https://godoc.org/github.com/deckarep/golang-set

// Version func
func Version() string {
	return "1.0.10"
}

func isHexWord(word string) bool {
	_, err := hex.DecodeString(word)
	return len(word) >= 10 && err == nil
}

// FormatDate func Put into utils.
func FormatDate(t time.Time) string {
	var a [20]byte
	var b = a[:0]
	b = t.AppendFormat(b, time.RFC3339)
	return string(b[0:10])
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

/*************************************************************************************************/

// SelectOccurrenceByDate assumes NullTime have zero hours and occurrenceList is sorted by ArchiveDate.
func SelectOccurrenceByDate(occurrenceList []hd.Occurrence, timeinterval nt.TimeInterval) []hd.Occurrence {
	var subList []hd.Occurrence

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
func SelectOccurrenceByID(occurrenceList []hd.Occurrence, acmID uint32) []hd.Occurrence {
	var subList []hd.Occurrence
	for ndx := 0; ndx < len(occurrenceList); ndx++ {
		if acmID < occurrenceList[ndx].AcmId {
			continue
		}
		if acmID == occurrenceList[ndx].AcmId {
			subList = append(subList, occurrenceList[ndx])
		}
		if acmID > occurrenceList[ndx].AcmId {
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

// GetOccurrenceListByDate returns result set ordered by ArchiveDate.
// Read []Occurrence values by archiveDate range. This applys FilteringRules(word).
// mapset.Set is the set of distinct AcmId values in the returned list.
func GetOccurrenceListByDate(timeinterval nt.TimeInterval) ([]hd.Occurrence, mapset.Set, error) {
	db, err := dbx.GetDatabaseReference()
	defer db.Close()
	// Invoke stored procedure.
	SELECT := "SELECT * FROM GetOccurrencesByDate('" + timeinterval.StartDate.StandardDate() + "', '" + timeinterval.EndDate.StandardDate() + "')"
	rows, err := db.Query(SELECT)
	dbx.CheckErr(err)
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
		dbx.CheckErr(err)

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
	dbx.CheckErr(err)

	return occurrenceList, idSet, nil
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

	// goroutine version:
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

	return alphaCollection, idSet
}

// GetOccurrencesByAcmid func
func GetOccurrencesByAcmid(xacmid uint32) ([]hd.Occurrence, error) {
	db, err := dbx.GetDatabaseReference()
	defer db.Close()

	SELECT := "SELECT acmId, archiveDate, word, nentry FROM Occurrence WHERE acmId=" + strconv.FormatUint(uint64(xacmid), 10)
	rows, err := db.Query(SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	// fields to read
	var acmID uint32
	var archiveDate nt.NullTime
	var word string
	var nentry int
	var occurrenceList []hd.Occurrence

	for rows.Next() {
		err = rows.Scan(&acmID, &archiveDate, &word, &nentry)
		dbx.CheckErr(err)

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
	dbx.CheckErr(err)

	return occurrenceList, nil
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

/*************************************************************************************************/

// BulkInsertConditionalProbability uses prepared statement.
func BulkInsertConditionalProbability(conditionals []hd.ConditionalProbability) error {
	db, err := dbx.GetDatabaseReference()
	defer db.Close()

	txn, err := db.Begin()
	dbx.CheckErr(err)

	stmt, err := db.Prepare("INSERT INTO Conditional (wordlist, probability, timeframetype, startDate, endDate, firstDate, lastDate) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	dbx.CheckErr(err)

	for _, v := range conditionals {
		_, err = stmt.Exec(v.WordList, v.Probability, v.Timeinterval.Timeframetype, v.Timeinterval.StartDate.DT, v.Timeinterval.EndDate.DT, v.FirstDate.DT, v.LastDate.DT)
		dbx.CheckErr(err)
	}

	err = stmt.Close()
	dbx.CheckErr(err)

	err = txn.Commit()
	dbx.CheckErr(err)

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

/* P(dependent A and B both occurring): Bayes: P(A|B)=P(A∩B)/P(B)=P(B|A)P(A)/P(B)
What is the P of word A given word B (in this interval)? If P(A|B)=P(A) then events A and B are said to be independent.
P(A∩B)=P(A|B)*P(B) is the probability that both events A and B occur; they are present in the same summary.
The imported wordMap has probabilities over timeinterval. startingWordgram allows for restart: must be in wordA|wordB format.
Do for 2 permutations (order matters). Performs FilteringRules(words) Returns len(wordGrams).
Number of permutations for 97022 wordgrams is n!/(n-r)! = 9,413,171,462. */

// CalcConditionalProbability func returns 	wordMap:SELECT Word,Probability FROM vocabulary */
func CalcConditionalProbability(startingWordgram string, wordMap map[string]float32, timeinterval nt.TimeInterval) (int, error) {
	if len(wordMap) < 2 {
		fmt.Println("There must at at least 2 words to compute conditional probabilities.")
		return 0, nil
	}
	permutations := 2
	var cutoffProbability float32 = 0.000001 // 1.0x10^-6
	index := strings.Index(startingWordgram, "|")
	wordAstart := startingWordgram[0:index]
	wordBstart := startingWordgram[index+1:]

	wordGrams := ExtractKeysFromProbabilityMap(wordMap) // []string
	sort.Strings(wordGrams)

	if len(wordGrams) < 10 {
		fmt.Println("Processing: " + strings.Join(wordGrams, " + "))
	} else {
		fmt.Println("Processing: " + strconv.Itoa(len(wordGrams)) + " wordgrams.")
	}

	start := time.Now()
	fmt.Print("CalcConditionalProbability (permutations=" + strconv.Itoa(permutations) + "): ")

	DB1, err := dbx.GetDatabaseReference() // for calling functions
	defer DB1.Close()

	var conditionals []hd.ConditionalProbability
	var condProb1, condProb2 float32 // must match function RETURNS TABLE names.
	var firstDate, lastDate time.Time
	var firstDateValue, lastDateValue nt.NullTime
	var totalInserts int64
	startDateParam := FormatDate(timeinterval.StartDate.DT)
	endDateParam := FormatDate(timeinterval.EndDate.DT)

	if permutations == 2 {
		for wordA := 0; wordA < len(wordGrams)-1; wordA++ {
			if strings.Compare(wordGrams[wordA], wordAstart) < 0 { // not <= !
				continue
			}
			conditionals = nil
			fmt.Print(wordGrams[wordA] + "  ")

			for wordB := wordA + 1; wordB < len(wordGrams); wordB++ {
				if strings.Compare(wordGrams[wordB], wordBstart) <= 0 {
					continue
				}
				err = DB1.QueryRow(`SELECT condProb1, condProb2 FROM GetConditionalProbabilities($1, $2, $3, $4)`, wordGrams[wordA], wordGrams[wordB], startDateParam, endDateParam).Scan(&condProb1, &condProb2)
				dbx.CheckErr(err)
				if condProb1 > cutoffProbability && condProb2 > cutoffProbability {
					err = DB1.QueryRow(`SELECT firstDate, lastDate FROM GetFirstLastArchiveDates($1, $2, $3, $4)`, wordGrams[wordA], wordGrams[wordB], startDateParam, endDateParam).Scan(&firstDate, &lastDate)
					dbx.CheckErr(err)                            // firstDate, lastDate can be null!
					firstDateValue = nt.New_NullTime2(firstDate) // must match function RETURNS TABLE names.
					lastDateValue = nt.New_NullTime2(lastDate)
					wordlist := wordGrams[wordA] + "|" + wordGrams[wordB]
					conditionals = append(conditionals, hd.ConditionalProbability{Id: 0, WordList: wordlist, Probability: condProb1, Timeinterval: timeinterval, FirstDate: firstDateValue, LastDate: lastDateValue})
					wordlist = wordGrams[wordB] + "|" + wordGrams[wordA]
					conditionals = append(conditionals, hd.ConditionalProbability{Id: 0, WordList: wordlist, Probability: condProb2, Timeinterval: timeinterval, FirstDate: firstDateValue, LastDate: lastDateValue})
				}
			}

			if len(conditionals) > 0 {
				_ = BulkInsertConditionalProbability(conditionals)
				totalInserts = totalInserts + int64(len(conditionals))
			}
		}
		fmt.Println(totalInserts)
	}

	elapsed := time.Since(start)
	fmt.Println(elapsed.String())
	return len(wordGrams), nil
}

// GetConditionalByTimeInterval func modifies condProbList pointer which should be declared beforehand. bigramMap does not need to be a pointer.
func GetConditionalByTimeInterval(bigrams []string, timeInterval nt.TimeInterval, condProbList *[]hd.ConditionalProbability, bigramMap map[string]bool) error {
	DB, err := dbx.GetDatabaseReference()
	defer DB.Close()

	inPhrase := dbx.CompileInClause(bigrams)
	query := "SELECT id, wordlist, probability, timeframetype, startDate, endDate, firstDate, lastDate FROM conditional WHERE wordlist IN " + inPhrase + " AND timeframetype=" +
		strconv.Itoa(int(timeInterval.Timeframetype)) + " AND startDate >= '" + timeInterval.StartDate.StandardDate() + "' AND endDate <= '" + timeInterval.EndDate.StandardDate() + "'"

	rows, err := DB.Query(query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetConditionalByTimeInterval(1): %+v\n", err)
		return err
	}
	defer rows.Close()

	var cProb hd.ConditionalProbability
	var timeframetype int
	var startDate time.Time
	var endDate time.Time

	for rows.Next() { // 720,066 total rows per TFTerm.
		err := rows.Scan(&cProb.Id, &cProb.WordList, &cProb.Probability, &timeframetype, &startDate, &endDate, &cProb.FirstDate, &cProb.LastDate)
		if err != nil {
			log.Printf("GetConditionalByTimeInterval(2): %+v\n", err)
			return err
		}
		bigramMap[cProb.WordList] = true
		cProb.Timeinterval = nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startDate), nt.New_NullTime2(endDate))
		*condProbList = append(*condProbList, cProb)
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)
	return err
}

// GetConditionalByProbability func
func GetConditionalByProbability(word string, probabilityCutoff float32, timeInterval nt.TimeInterval, condProbList *[]hd.ConditionalProbability) error {
	DB, err := dbx.GetDatabaseReference()
	defer DB.Close()

	prefix := "'" + word + "|%'"
	postfix := "'%|" + word + "'"
	query := "SELECT id, wordlist, probability, timeframetype, startDate, endDate, firstDate, lastDate FROM conditional WHERE timeframetype=" +
		strconv.Itoa(int(timeInterval.Timeframetype)) + " AND startDate >= '" + timeInterval.StartDate.StandardDate() + "' AND endDate <= '" + timeInterval.EndDate.StandardDate() + "' AND " +
		"(wordlist LIKE " + prefix + " OR wordlist LIKE " + postfix + ")"

	rows, err := DB.Query(query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetConditionalByProbability(1): %+v\n", err)
		return err
	}
	defer rows.Close()

	var cProb hd.ConditionalProbability
	var timeframetype int
	var startDate time.Time
	var endDate time.Time

	for rows.Next() {
		err := rows.Scan(&cProb.Id, &cProb.WordList, &cProb.Probability, &timeframetype, &startDate, &endDate, &cProb.FirstDate, &cProb.LastDate)
		if err != nil {
			log.Printf("GetConditionalByProbability(2): %+v\n", err)
			return err
		}
		if cProb.Probability >= probabilityCutoff && (strings.HasPrefix(cProb.WordList, word+"|") || strings.HasSuffix(cProb.WordList, "|"+word)) { // remove prefix-postfix words.
			cProb.Timeinterval = nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startDate), nt.New_NullTime2(endDate))
			*condProbList = append(*condProbList, cProb)
		}
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)
	return err
}
