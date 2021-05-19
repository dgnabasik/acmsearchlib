package conditional

// Manages conditional probabilities and occurrences.
// NOTE: ConditionalProbability struct does NOT include wordarray.

import (
	"context"
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
	pgx "github.com/jackc/pgx/v4"
)

// comment
const (
	SEP = "|"
)

// mapset https://github.com/deckarep/golang-set/blob/master/README.md & https://godoc.org/github.com/deckarep/golang-set

// Version func
func Version() string {
	return "1.16.2"
}

func isHexWord(word string) bool {
	_, err := hex.DecodeString(word)
	return len(word) >= 10 && err == nil
}

// FilteringRules filters output from Postgres ts_stat select. Include 3d prefixes.
// Return 0 for ok, -1 to completely ignore, 1 for modified word.
func FilteringRules(word string) (string, int) {
	if len(strings.TrimSpace(word)) <= 1 {
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
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	SELECT := "SELECT * FROM GetOccurrencesByDate" + dbx.GetFormattedDatesForProcedure(timeinterval)
	rows, err := db.Query(context.Background(), SELECT)
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
	if err != nil {
		return nil, err
	}
	defer db.Close()

	SELECT := "SELECT acmId, archiveDate, word, nentry FROM Occurrence WHERE acmId=" + strconv.FormatUint(uint64(xacmid), 10)
	rows, err := db.Query(context.Background(), SELECT)
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
	wordMap := make(map[string]int)
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
	if err != nil {
		return err
	}
	defer db.Close()

	txn, err := db.Begin(context.Background())
	dbx.CheckErr(err)

	// Must use lowercase column names!
	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"conditional"}, // tablename
		[]string{"wordlist", "probability", "timeframetype", "startdate", "enddate", "firstdate", "lastdate", "pmi", "dateupdated"},
		pgx.CopyFromSlice(len(conditionals), func(i int) ([]interface{}, error) {
			return []interface{}{conditionals[i].WordList, conditionals[i].Probability, int(conditionals[i].Timeinterval.Timeframetype), conditionals[i].Timeinterval.StartDate.DT,
				conditionals[i].Timeinterval.EndDate.DT, conditionals[i].FirstDate.DT, conditionals[i].LastDate.DT, conditionals[i].Pmi, conditionals[i].DateUpdated}, nil
		}),
	)

	dbx.CheckErr(err)
	if copyCount == 0 {
		log.Printf("BulkInsertConditionalProbability: no rows inserted")
	}
	err = txn.Commit(context.Background())
	dbx.CheckErr(err)

	//// Execute wordarray update stored proc:
	_, err = db.Exec(context.Background(), "call UpdateVocabularyWordarray();")
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
		log.Printf("There must at at least 2 words to compute conditional probabilities.")
		return 0, nil
	}
	permutations := 2
	var cutoffProbability float32 = 0.000001 // 1.0x10^-6
	index := strings.Index(startingWordgram, SEP)
	wordAstart := startingWordgram[0:index]
	wordBstart := startingWordgram[index+1:]

	wordGrams := ExtractKeysFromProbabilityMap(wordMap) // []string
	sort.Strings(wordGrams)

	if len(wordGrams) < 10 {
		fmt.Println("Processing: " + strings.Join(wordGrams, " + "))
	} else {
		fmt.Println("Processing: " + strconv.Itoa(len(wordGrams)) + " wordgrams.")
	}

	DB1, err := dbx.GetDatabaseReference() // for calling functions
	if err != nil {
		return -1, err
	}
	defer DB1.Close()

	start := time.Now()

	var conditionals []hd.ConditionalProbability
	var pAgivenB, pBgivenA, pmi float32 // must match function RETURNS TABLE names.
	var firstDate, lastDate time.Time
	var firstDateValue, lastDateValue nt.NullTime
	var totalInserts int64
	startDateParam := timeinterval.StartDate.StandardDate()
	endDateParam := timeinterval.EndDate.StandardDate()

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
				today := nt.NullTimeToday().DT
				err = DB1.QueryRow(context.Background(), `SELECT pAgivenB, pBgivenA, pmi FROM GetConditionalProbabilities($1, $2, $3, $4)`, wordGrams[wordA], wordGrams[wordB], startDateParam, endDateParam).Scan(&pAgivenB, &pBgivenA, &pmi)
				dbx.CheckErr(err)
				if pAgivenB > cutoffProbability && pBgivenA > cutoffProbability {
					err = DB1.QueryRow(context.Background(), `SELECT firstDate, lastDate FROM GetFirstLastArchiveDates($1, $2, $3, $4)`, wordGrams[wordA], wordGrams[wordB], startDateParam, endDateParam).Scan(&firstDate, &lastDate)
					dbx.CheckErr(err)                            // firstDate, lastDate can be null!
					firstDateValue = nt.New_NullTime2(firstDate) // must match function RETURNS TABLE names.
					lastDateValue = nt.New_NullTime2(lastDate)
					wordlist := wordGrams[wordA] + SEP + wordGrams[wordB]
					conditionals = append(conditionals, hd.ConditionalProbability{Id: 0, WordList: wordlist, Probability: pAgivenB, Timeinterval: timeinterval, FirstDate: firstDateValue, LastDate: lastDateValue, Pmi: pmi, DateUpdated: today})
					wordlist = wordGrams[wordB] + SEP + wordGrams[wordA]
					conditionals = append(conditionals, hd.ConditionalProbability{Id: 0, WordList: wordlist, Probability: pBgivenA, Timeinterval: timeinterval, FirstDate: firstDateValue, LastDate: lastDateValue, Pmi: pmi, DateUpdated: today})
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

// GetConditionalByTimeInterval func modifies condProbList.
func GetConditionalByTimeInterval(bigrams []string, timeInterval nt.TimeInterval, bigramMap map[string]bool, includeTimeframetype bool) ([]hd.ConditionalProbability, error) {
	condProbList := make([]hd.ConditionalProbability, 0)
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return condProbList, err
	}
	defer db.Close()

	inPhrase := dbx.CompileInClause(bigrams)
	query := "SELECT id, wordlist, probability, timeframetype, startDate, endDate, firstDate, lastDate, pmi, dateUpdated FROM conditional WHERE wordlist IN " + inPhrase +
		" AND " + dbx.CompileDateClause(timeInterval, includeTimeframetype)
	rows, err := db.Query(context.Background(), query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetConditionalByTimeInterval(1): %+v\n", err)
		return condProbList, err
	}
	defer rows.Close()

	var cProb hd.ConditionalProbability
	var timeframetype int
	var startDate time.Time
	var endDate time.Time

	for rows.Next() {
		err := rows.Scan(&cProb.Id, &cProb.WordList, &cProb.Probability, &timeframetype, &startDate, &endDate, &cProb.FirstDate, &cProb.LastDate, &cProb.Pmi, &cProb.DateUpdated)
		if err != nil {
			log.Printf("GetConditionalByTimeInterval(2): %+v\n", err)
			return condProbList, err
		}
		bigramMap[cProb.WordList] = true
		cProb.Timeinterval = nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startDate), nt.New_NullTime2(endDate))
		condProbList = append(condProbList, cProb)
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)
	return condProbList, err
}

// GetConditionalByProbability func
func GetConditionalByProbability(word string, probabilityCutoff float32, timeInterval nt.TimeInterval, condProbList *[]hd.ConditionalProbability) error {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return err
	}
	defer db.Close()

	//prefix := "'" + word + "|%'"
	//postfix := "'%|" + word + "'"
	query := "SELECT id, wordlist, probability, timeframetype, startDate, endDate, firstDate, lastDate, pmi, dateUpdated FROM conditional WHERE " + dbx.CompileDateClause(timeInterval, false) +
		" AND (wordarray[1]=word OR wordarray[2]=word)"
		//// " AND (wordlist LIKE " + prefix + " OR wordlist LIKE " + postfix + ")"

	rows, err := db.Query(context.Background(), query)
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
		err := rows.Scan(&cProb.Id, &cProb.WordList, &cProb.Probability, &timeframetype, &startDate, &endDate, &cProb.FirstDate, &cProb.LastDate, &cProb.Pmi, &cProb.DateUpdated)
		if err != nil {
			log.Printf("GetConditionalByProbability(2): %+v\n", err)
			return err
		}
		if cProb.Probability >= probabilityCutoff && (strings.HasPrefix(cProb.WordList, word+SEP) || strings.HasSuffix(cProb.WordList, SEP+word)) { // remove prefix-postfix words.
			cProb.Timeinterval = nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startDate), nt.New_NullTime2(endDate))
			*condProbList = append(*condProbList, cProb)
		}
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)
	return err
}

// GetWordBigramPermutations func permutes or combines words.
func GetWordBigramPermutations(words []string, permute bool) []string {
	result := make([]string, 0)
	if len(words) == 1 {
		return words
	}
	for i := 0; i < len(words); i++ {
		for j := i + 1; j < len(words); j++ {
			result = append(result, words[i]+SEP+words[j])
			if permute {
				result = append(result, words[j]+SEP+words[i])
			}
		}
	}
	return result
}

// GetConditionalList func
func GetConditionalList(words []string, timeInterval nt.TimeInterval, permute bool) ([]hd.ConditionalProbability, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	bigrams := GetWordBigramPermutations(words, permute)
	intervalClause := dbx.CompileDateClause(timeInterval, false)
	compileInClause := dbx.CompileInClause(bigrams)
	SELECT := "SELECT id, wordlist, probability, timeframetype, startDate, endDate, firstDate, lastDate, pmi, dateUpdated FROM Conditional WHERE wordlist IN " +
		compileInClause + " AND " + intervalClause
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var cProb hd.ConditionalProbability
	var timeframetype int
	var startDate time.Time
	var endDate time.Time
	var condProbList []hd.ConditionalProbability

	for rows.Next() {
		err := rows.Scan(&cProb.Id, &cProb.WordList, &cProb.Probability, &timeframetype, &startDate, &endDate, &cProb.FirstDate, &cProb.LastDate, &cProb.Pmi, &cProb.DateUpdated)
		dbx.CheckErr(err)
		cProb.Timeinterval = nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startDate), nt.New_NullTime2(endDate))
		condProbList = append(condProbList, cProb)
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return condProbList, nil
}

// GetExistingConditionalBigrams func tests for existing bigrams in Conditional.WordList
func GetExistingConditionalBigrams(bigrams []string, intervalClause string) ([]string, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	bigramList := make([]string, 0)
	compileInClause := dbx.CompileInClause(bigrams)
	query := "SELECT wordlist FROM Conditional WHERE wordlist IN " + compileInClause + " AND " + intervalClause
	rows, err := db.Query(context.Background(), query)
	dbx.CheckErr(err)
	defer rows.Close()

	var wordlist string
	for rows.Next() {
		err := rows.Scan(&wordlist)
		dbx.CheckErr(err)
		bigramList = append(bigramList, wordlist)
	}
	err = rows.Err()
	dbx.CheckErr(err)
	return bigramList, err
}

// GetProbabilityGraph func returns ordered list of high-probability bigrams for given word.
func GetProbabilityGraph(words []string, timeInterval nt.TimeInterval) ([]hd.ConditionalProbability, error) {
	// build SQL values:
	bigrams := GetWordBigramPermutations(words, true) // permute=true
	intervalClause := dbx.CompileDateClause(timeInterval, false)
	bigrams, _ = GetExistingConditionalBigrams(bigrams, intervalClause) // wordlist[]

	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var SELECT strings.Builder
	for index := 0; index < len(bigrams); index++ {
		bigram := "'" + bigrams[index] + "'"
		ndxSep := strings.Index(bigrams[index], SEP)
		leftWord := "'" + bigrams[index][0:ndxSep] + "'"
		rightWord := "'" + bigrams[index][ndxSep+1:] + "'"
		reverseBigram := "'" + bigrams[index][ndxSep+1:] + SEP + bigrams[index][0:ndxSep] + "'"
		likeWord1 := "'" + bigrams[index][0:ndxSep] + SEP + "%'"
		likeWord2 := "'" + bigrams[index][ndxSep+1:] + SEP + "%'"

		SELECT.WriteString("SELECT id, wordlist, probability, timeframetype, startDate, endDate, firstDate, lastDate, pmi, dateUpdated FROM Conditional ")
		SELECT.WriteString("WHERE wordlist LIKE " + likeWord1 + " AND " + intervalClause)
		SELECT.WriteString("AND pmi >= (SELECT MAX(pmi) FROM Conditional WHERE wordlist=" + bigram + " AND " + intervalClause + ") ")
		SELECT.WriteString("AND probability >= (SELECT MAX(probability) FROM Conditional WHERE wordlist=" + bigram + " AND " + intervalClause + ") ")
		SELECT.WriteString("AND SUBSTRING(wordlist from " + strconv.Itoa(len(leftWord)+2) + " for 32) IN (SELECT word FROM Wordscore WHERE score >= (SELECT MAX(score) FROM Wordscore WHERE word=" + leftWord + ")) ")
		SELECT.WriteString("UNION ")
		SELECT.WriteString("SELECT id, wordlist, probability, timeframetype, startDate, endDate, firstDate, lastDate, pmi, dateUpdated FROM Conditional ")
		SELECT.WriteString("WHERE wordlist LIKE " + likeWord2 + " AND " + intervalClause)
		SELECT.WriteString("AND pmi >= (SELECT MAX(pmi) FROM Conditional WHERE wordlist=" + reverseBigram + " AND " + intervalClause + ") ")
		SELECT.WriteString("AND probability >= (SELECT MAX(probability) FROM Conditional WHERE wordlist=" + reverseBigram + " AND " + intervalClause + ") ")
		SELECT.WriteString("AND SUBSTRING(wordlist FROM " + strconv.Itoa(len(rightWord)+2) + " for 32) IN (SELECT word FROM Wordscore WHERE score >= (SELECT MAX(score) FROM Wordscore WHERE word=" + rightWord + ")) ")

		if index < len(bigrams)-1 {
			SELECT.WriteString("UNION ")
		}
	}
	SELECT.WriteString(";")

	rows, err := db.Query(context.Background(), SELECT.String())
	dbx.CheckErr(err)
	defer rows.Close()

	var cProb hd.ConditionalProbability
	var timeframetype int
	var startDate time.Time
	var endDate time.Time
	var condProbList []hd.ConditionalProbability

	for rows.Next() {
		err := rows.Scan(&cProb.Id, &cProb.WordList, &cProb.Probability, &timeframetype, &startDate, &endDate, &cProb.FirstDate, &cProb.LastDate, &cProb.Pmi, &cProb.DateUpdated)
		dbx.CheckErr(err)
		cProb.Timeinterval = timeInterval
		condProbList = append(condProbList, cProb)
	}

	err = rows.Err()
	dbx.CheckErr(err)

	// fetch reverseWordlist values:
	reverseWordlist := make([]string, 0)
	for _, cp := range condProbList {
		ndxSep := strings.Index(cp.WordList, SEP)
		reverseWordlist = append(reverseWordlist, cp.WordList[ndxSep+1:]+SEP+cp.WordList[0:ndxSep])
	}
	compileInClause := dbx.CompileInClause(reverseWordlist)
	query := "SELECT id, wordlist, probability, timeframetype, startDate, endDate, firstDate, lastDate, pmi, dateUpdated FROM Conditional WHERE wordlist IN " + compileInClause + " AND " + intervalClause
	rows, err = db.Query(context.Background(), query)
	dbx.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&cProb.Id, &cProb.WordList, &cProb.Probability, &timeframetype, &startDate, &endDate, &cProb.FirstDate, &cProb.LastDate, &cProb.Pmi, &cProb.DateUpdated)
		dbx.CheckErr(err)
		cProb.Timeinterval = timeInterval
		condProbList = append(condProbList, cProb)
	}
	err = rows.Err()
	dbx.CheckErr(err)

	return condProbList, nil
}

// GetWordgramConditionalsByInterval func assigns consecutive id values.  Common:bool column not in database.
// Id values start at 10000 to avoid js Select Id conflicts.
// ISSUE: the index on WordArray[1] will return a row where word1 is a prefix of word2: (e.g., '3d|3ds')
func GetWordgramConditionalsByInterval(words []string, timeInterval nt.TimeInterval, dimensions int) ([]hd.WordScoreConditionalFlat, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	inPhrase := dbx.CompileInClause(words) // Can't use dbx.CompileDateClause() because of w alias.
	SELECT := `SELECT w.word, c.wordlist, w.score, c.probability, c.pmi, c.timeframetype, c.startDate, c.endDate, c.firstdate, c.lastdate FROM Wordscore AS w 
		INNER JOIN Conditional AS c ON w.word=c.wordarray[1] WHERE w.startdate=c.startDate AND w.endDate=c.endDate`
	SELECT += " AND w.word IN " + inPhrase + " AND w.startDate='" + timeInterval.StartDate.StandardDate() + "' AND w.endDate='" + timeInterval.EndDate.StandardDate() + "' ORDER BY c.wordlist"

	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var score, pmi, probability float32
	var timeframetype int
	var startDate, endDate, firstDate, lastDate time.Time
	var word, wordlist string
	var id int = 10000
	wordScoreConditionalList := make([]hd.WordScoreConditionalFlat, 0)

	for rows.Next() {
		err = rows.Scan(&word, &wordlist, &score, &probability, &pmi, &timeframetype, &startDate, &endDate, &firstDate, &lastDate)
		dbx.CheckErr(err)
		wordArray := strings.Split(wordlist, SEP)
		id++
		wordScoreConditionalList = append(wordScoreConditionalList, hd.WordScoreConditionalFlat{ID: id, WordArray: wordArray,
			Wordlist: wordlist, Score: score, Probability: probability, Pmi: pmi, Timeframetype: timeframetype, StartDate: startDate, EndDate: endDate, FirstDate: firstDate, LastDate: lastDate})
	}
	err = rows.Err()
	dbx.CheckErr(err)

	return wordScoreConditionalList, err
}
