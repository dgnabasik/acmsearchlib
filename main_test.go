package acmsearchlib

// All not-test files init() functions are executed first, then all test files init() functions are executed (hopefully in lexical order).
// go test -v.
import (
	"fmt"
	"strconv"
	"testing"
	"time"

	art "github.com/dgnabasik/acmsearchlib/article"
	dbx "github.com/dgnabasik/acmsearchlib/database"
	fs "github.com/dgnabasik/acmsearchlib/filesystem"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"
	voc "github.com/dgnabasik/acmsearchlib/vocabulary"
)

const (
	SEP = "|"
)

/* nulltime ************************************************************************************/

// Test_nulltime func
func Test_nulltime(t *testing.T) {
	str := nt.GetShortMonthName(13)
	if str != "" {
		t.Error("a1:Expected '', got ", str)
	}
	str = nt.GetShortMonthName(12)
	if str != "dec" {
		t.Error("a2:Expected 'dec', got ", str)
	}

	data := make([]string, 0)
	ndx, found := nt.StringSliceContains(data, "dec")
	if ndx >= 0 || found {
		t.Error("a3:Expected -1, got ", ndx)
	}
	data = append(data, "dec")
	ndx, found = nt.StringSliceContains(data, "dec")
	if ndx < 0 || !found {
		t.Error("a4:Expected 0, got ", ndx)
	}

	ti := nt.New_NullTime("")
	if ti.StandardDate() != nt.NullDate {
		t.Error("a5: Expected " + nt.NullDate)
	}

	ti = nt.New_NullTime("1999-01-01") // Valid=true
	if ti.StandardDate() != nt.NullDate {
		t.Error("a6: Expected " + nt.NullDate)
	}

	ti = nt.New_NullTime("2000-13")
	if ti.StandardDate() != nt.NullDate {
		t.Error("a7: Expected " + nt.NullDate)
	}

	ti = nt.New_NullTime(nt.VeryFirstDate)
	if ti.StandardDate() != nt.VeryFirstDate {
		t.Error("a8: Expected " + nt.VeryFirstDate)
	}

	ti = nt.New_NullTime1("September 23, 2019")
	if ti.StandardDate() != "2019-09-23" {
		t.Error("a9: Expected 2019-09-23")
	}

	/*ti = nt.NullTimeToday()
	if ti.StandardDate() != "2019-11-25" {
		t.Error("a10: Expected 2019-11-25")
	}*/

	ti = nt.New_NullTimeFromFileName("dec-05-2005.html")
	if ti.StandardDate() != "2005-12-05" {
		t.Error("a11: Expected 2005-12-05")
	}

	if ti.FileSystemDate() != "dec-05-2005" {
		t.Error("a12: Expected dec-05-2005")
	}

	if ti.HtmlArchiveDate() != "2005-12-dec" {
		t.Error("a13: Expected 2005-12-dec")
	}

	if ti.NonStandardDate() != "12/05/05" {
		t.Error("a14: Expected 12/05/05")
	}

	if nt.GetStandardDateForm("dec-05-2005") != "2005-12-05" { // Convert mmm-dd-yyyy to yyyy-mm-dd
		t.Error("a15: Expected 2005-12-05")
	}

	ti.AdvanceNextNullTime()
	if ti.StandardDate() != "2005-12-07" {
		t.Error("a16: Expected 2005-12-07")
	}

	nt1, nt2 := nt.GetStartEndOfWeek(ti)
	if nt1.StandardDate() != "2005-12-04" {
		t.Error("a17: Expected 2005-12-04")
	}
	if nt2.StandardDate() != "2005-12-10" {
		t.Error("a18: Expected 2005-12-10")
	}

	year, month, day, hour, min, sec := nt.NullTimeDiff(nt1, nt2)
	if year != 0 || month != 0 || day != 6 || hour != 0 || min != 0 || sec != 0 {
		t.Error("a19: Expected 6 days diff.")
	}

	ntx := nt.NullTimeToday()
	dateSet := make([]nt.NullTime, 0)
	dateSet = append(dateSet, ntx)
	dateSet = append(dateSet, ti)
	dateSet = nt.NullTimeSorter(dateSet)
	if len(dateSet) != 2 {
		t.Error("a20: Sorted dateSet != 2 ")
	}

	ntx = nt.New_NullTime("2019-12-01")
	yes := ntx.IsScheduledDate(nt.TFUnknown)
	if !yes {
		t.Error("a21: Not past 11am.")
	}

	yes = ntx.IsScheduledDate(nt.TFWeek)
	if !yes {
		t.Error("a22: Not start of the week.")
	}

	yes = ntx.IsScheduledDate(nt.TFMonth)
	if !yes {
		t.Error("a23: Not start of the month.")
	}

	ntx = nt.New_NullTime("2020-01-01")
	yes = ntx.IsScheduledDate(nt.TFQuarter)
	if !yes {
		t.Error("a24: Not start of the Quarter.")
	}

	yes = ntx.IsScheduledDate(nt.TFYear)
	if !yes {
		t.Error("a25: Not start of the year.")
	}

	vfd1 := nt.New_NullTime(nt.VeryFirstDate)
	vfd2 := vfd1
	ti2 := nt.New_TimeInterval(nt.TFUnknown, vfd1, vfd2)
	if ti2.StartDate.StandardDate() != ti2.EndDate.StandardDate() {
		t.Error("a26: Expected " + ti2.EndDate.StandardDate())
	}
	// add one time unitL len(list) always 2 except TFUnknown is 1.
	tiList := nt.GetTimeIntervalDatePartitionList(ti2)
	if len(tiList) != 1 && tiList[0].StartDate.StandardDate() != tiList[0].EndDate.StandardDate() {
		t.Error("a27: Expected " + tiList[0].StartDate.StandardDate())
	}

	vfd2 = nt.New_NullTime2(vfd1.DT.AddDate(0, 0, 7)) // y,m,d
	ti3 := nt.New_TimeInterval(nt.TFWeek, vfd1, vfd2)
	tiList = nt.GetTimeIntervalDatePartitionList(ti3)
	if len(tiList) < 1 {
		t.Error("a28: Expected 2 tiList items.")
	}
	year, month, day, hour, min, sec = nt.NullTimeDiff(vfd1, vfd2)
	if year != 0 || month != 0 || day != 7 || hour != 0 || min != 0 || sec != 0 {
		t.Error("a29: Expected 7 days diff.")
	}

	vfd2 = nt.New_NullTime2(vfd1.DT.AddDate(0, 1, 0))
	ti4 := nt.New_TimeInterval(nt.TFMonth, vfd1, vfd2)
	tiList = nt.GetTimeIntervalDatePartitionList(ti4)
	if len(tiList) != 2 {
		t.Error("a30: Expected 2 tiList items.")
	}
	year, month, day, hour, min, sec = nt.NullTimeDiff(vfd1, vfd2)
	if year != 0 || month != 1 || day != 0 || hour != 0 || min != 0 || sec != 0 {
		t.Error("a31: Expected 1 month diff.")
	}

	// skipped TFQuarter & TFSpan

	vfd2 = nt.New_NullTime2(vfd1.DT.AddDate(1, 0, 0))
	ti5 := nt.New_TimeInterval(nt.TFYear, vfd1, vfd2)
	tiList = nt.GetTimeIntervalDatePartitionList(ti5)
	if len(tiList) != 2 {
		t.Error("a32: Expected 2 tiList items.")
	}
	year, month, day, hour, min, sec = nt.NullTimeDiff(vfd1, vfd2)
	if year != 1 || month != 0 || day != 0 || hour != 0 || min != 0 || sec != 0 {
		t.Error("a33: Expected 1 year diff.")
	}

	vfd1 = nt.New_NullTime(nt.VeryFirstDate) // "2000-01-03"
	vfd2 = nt.New_NullTime2(vfd1.DT.AddDate(1, 0, 0))
	ti6 := nt.New_TimeInterval(nt.TFYear, vfd1, vfd2)
	tf := nt.New_TimeFrame(ti6)

	if tf.GivenDate.StandardDate() != "2000-01-03" {
		t.Error("a34: Expected 2000-01-03.")
	}

	if tf.StartOfWeek.StandardDate() != "2000-01-02" {
		t.Error("a35: Expected 2000-01-02.")
	}

	if tf.EndOfWeek.StandardDate() != "2000-01-08" {
		t.Error("a36: Expected 2000-01-08.")
	}

	if tf.StartOfMonth.StandardDate() != "2000-01-01" {
		t.Error("a37: Expected 2000-01-01.")
	}

	if tf.EndOfMonth.StandardDate() != "2000-01-31" {
		t.Error("a38: Expected 2000-01-31.")
	}

	if tf.StartOfQuarter.StandardDate() != "2000-01-01" {
		t.Error("a39: Expected 2000-01-01.")
	}

	if tf.EndOfQuarter.StandardDate() != "2000-03-31" {
		t.Error("a40: Expected 2000-03-31.")
	}

	if tf.StartOfYear.StandardDate() != "2000-01-01" {
		t.Error("a41: Expected 2000-01-01.")
	}

	if tf.EndOfYear.StandardDate() != "2000-12-31" {
		t.Error("a42: Expected 2000-12-31.")
	}

	if tf.StartOfSpan.StandardDate() != nt.VeryFirstDate {
		t.Error("a43: Expected " + nt.VeryFirstDate)
	}

	today := nt.NullTimeToday()
	if tf.EndOfSpan.StandardDate() != today.StandardDate() {
		t.Error("a44: Expected " + today.StandardDate())
	}

	year, month2, day := today.DT.Date()
	tf.Timeframetype = nt.TFSpan
	divisor := float32((year-nt.VeryFirstYear)*148) + float32((int(month2)-1)*12) + float32(day)/float32(3)
	if tf.GetDivisor() != divisor {
		t.Error("a45: Expected " + strconv.FormatFloat(float64(divisor), 'E', -1, 32))
	}

	vfd1, vfd2 = tf.GetTimeFrameDates()
	if vfd1.StandardDate() != nt.VeryFirstDate && vfd2.StandardDate() != today.StandardDate() {
		t.Error("a46: Expected " + today.StandardDate())
	}

	timeInUTC := time.Date(nt.VeryFirstYear, 1, 1, 1, 1, 1, 100, time.UTC)

	uts := nt.GetUnixTimeStampFromTime(timeInUTC)
	if uts != 946688461 {
		t.Error("a47: Expected 946688461.")
	}

	utsstr := nt.FormatUnixTimeStampAsString(uts)
	if utsstr != "946688461" {
		t.Error("a48: Expected 946688461.")
	}

	utsstr = nt.FormatUnixTimeStampAsTime(uts)
	if utsstr != "1999-12-31 18:01:01" {
		t.Error("a49: Expected 1999-12-31 18:01:01.")
	}

}

/* headers ************************************************************************************/

// Test_headers func
func Test_headers(t *testing.T) {
	start, end := hd.SearchForString("foobar", "bar")
	if start < 0 || end < 0 {
		t.Error("b1: bad SearchForString")
	}

	_, found := hd.SearchForStringIndex("foobar", "BAR")
	if !found {
		t.Error("b2: bad SearchForStringIndex")
	}

	lines1 := []string{"a", "a", "b", "c"}
	lines2 := []string{"A", "c", "d", "e"}

	start, found = hd.StringSliceContains(lines1, "a")
	if !found || start < 0 {
		t.Error("b3: bad StringSliceContains")
	}

	diff := hd.StringSetDifference(lines1, lines2) // in lines1 but not in lines2
	start, found = hd.StringSliceContains(lines1, "b")
	if !found || start < 0 || len(diff) == 0 {
		t.Error("b4: bad StringSetDifference")
	}

	diff = hd.StringSetDifference(lines2, lines1) // in lines2 but not in lines1
	start, found = hd.StringSliceContains(diff, "c")
	if found || start >= 0 || len(diff) == 0 {
		t.Error("b5: bad StringSetDifference")
	}

	str := hd.RandomHex(-1)
	if len(str) != 32 {
		t.Error("b6a:Expected len(32), got ", len(str))
	}
	str = hd.RandomHex(132)
	if len(str) != 256 {
		t.Error("b6b:Expected len(256), got ", len(str))
	}
	str = hd.RandomHex(32)
	if len(str) != 64 {
		t.Error("b6c:Expected len(64), got ", len(str))
	}

	mapIS := hd.GetOrderedMap(lines2)
	if len(mapIS) != len(lines2) {
		t.Error("b7:Expected len(64), got ", len(str))
	}

	len1 := len(lines1)
	lines1 = hd.RemoveDuplicateStrings(lines1)
	len2 := len(lines1)
	if len1 == len2 {
		t.Error("b8:bad RemoveDuplicateStrings")
	}

	lines1 = hd.DeleteStringSliceElement(lines1, "c")
	start, found = hd.StringSliceContains(lines1, "c")
	if found || start >= 0 {
		t.Error("b9: bad StringSetDifference")
	}

	aa := hd.AcmArticle{
		Id:            1,
		ArchiveDate:   nt.New_NullTime(""),
		ArticleNumber: "articlenumber",
		Title:         "title",
		ImageSource:   "imagesource",
		JournalName:   "journalname",
		AuthorName:    "authorname",
		JournalDate:   nt.New_NullTime(""),
		WebReference:  "webreference",
		Summary:       "summary",
	}

	mAA_SS, mAA_IS := aa.GetKeyValuePairs()
	if len(mAA_SS) != len(mAA_IS) {
		t.Error("b10: bad StringSetDifference")
	}

	uniClean := hd.ReplaceSpecialCharacters("'\n<a()-/&ndash;&mdash;&shy;&nbsp;&rsquo;&lsquo;&ldquo;&rdquo;&#151;&rdquo;&ldquo;&ecirc;&egrave;&Eacute;&eacute;&aacute;&oacute;&aring;&szlig;&uuml;&auml;&euml;&ouml;&oslash;&sup1;&hellip;&amp;&pound;&euro;&ntilde;")
	if len(uniClean) > 56 {
		t.Error("b11: bad ReplaceSpecialCharacters")
	}

	uniClean = uniClean + hd.HREF + " protected]"
	uniClean = hd.ReplaceProtected(uniClean)
	if len(uniClean) > 63 {
		t.Error("b12: bad ReplaceProtected")
	}

	words := []string{"3d", "access", "able", "kläui", "att", "beyond", "kommunikationsbüro", "cu", "schrödinger", "tübingen", "either", "four", "iff", "ins", "lin", "björn", "ngn", "éal", "nov", "éciale", "goëry", "göttingen", "loránd", "onto", "sa", "seven", "sf", "lovász", "ramón", "sánchez", "ably", "abroad", "abruptly", "absolutely", "thirteen"}
	vocabList, err := voc.GetVocabularyList(words)
	if err != nil {
		t.Error("b13: bad GetVocabularyList")
	}

	word := "tübingen"
	ndx1 := hd.GetVocabularyItem(word, vocabList)
	ndx2 := hd.GetVocabularyItemIndex(word, vocabList)
	if ndx1 != ndx2 {
		t.Error("b14: bad GetVocabularyItemIndex")
	}

	oaMap := hd.New_OrderedArticleMap()
	key := oaMap.FormatTitle(hd.HREF + "New Title")
	oaMap.Add(key, "New Title")
	title := oaMap.Get(key)
	if title != "New Title" {
		t.Error("b15: bad New_OrderedArticleMap")
	}
}

/* database ************************************************************************************/

// Test_database func
func Test_database(t *testing.T) {
	str := dbx.GetDatabaseConnectionString() // displays connection.
	if len(str) < 1 {
		t.Error("c1:Expected connection string got ", str)
	}

	str = dbx.Version()
	if len(str) < 1 {
		t.Error("c2:Expected version got ", str)
	}

	db, err := dbx.GetDatabaseReference()
	dbx.CheckErr(err)
	b := dbx.NoRowsReturned(err)
	db.Close()
	if b {
		t.Error("c3:Expected false got ", b)
	}

	words := []string{"3d", "access", "able"}
	str = dbx.CompileInClause(words)
	if len(str) < 1 {
		t.Error("c4:bad CompileInClause ")
	}

	startDate := nt.New_NullTime("2020-01-01")
	endDate := nt.New_NullTime("2020-12-31")
	timeInterval := nt.New_TimeInterval(nt.TFYear, startDate, endDate)

	str = dbx.GetFormattedDatesForProcedure(timeInterval)
	if len(str) < 1 {
		t.Error("c5:bad GetFormattedDatesForProcedure")
	}

	columnName := "columnName"
	str = dbx.GetWhereClause(columnName, words)
	if len(str) < 1 {
		t.Error("c6:bad GetWhereClause")
	}

	str = dbx.GetSingleDateWhereClause(columnName, timeInterval)
	if len(str) < 1 {
		t.Error("c7:bad GetSingleDateWhereClause")
	}

	str = dbx.CompileDateClause(timeInterval, true)
	if len(str) < 1 {
		t.Error("c8a:bad CompileDateClause")
	}

	str = dbx.CompileDateClause(timeInterval, true)
	if len(str) < 1 {
		t.Error("c8a:bad CompileDateClause")
	}

	arr := []int{1, 2, 3}
	intSlice := dbx.FormatArrayForStorage(arr)
	if len(intSlice) < 1 {
		t.Error("c9:bad FormatArrayForStorage")
	}
}

// Test_filesystem func
func Test_filesystem(t *testing.T) {
	prefix := fs.GetFilePrefixPath()
	if len(prefix) < 1 {
		t.Error("d1:bad GetFilePrefixPath")
	}

	fileInfo, err := fs.ReadDir(prefix + "Documents")
	if err != nil {
		t.Error("d2a: bad ReadDir")
	}
	if len(fileInfo) < 1 {
		t.Error("d2b:bad ReadDir")
	}

	dirPath := prefix + "test"
	err = fs.CreateDirectory(dirPath)
	if err != nil {
		t.Error("d3: bad CreateDirectory")
	}

	err = fs.DeleteDirectory(dirPath)
	if err != nil {
		t.Error("d4: bad DeleteDirectory")
	}

	filePath := prefix + "Documents/The viral universe.txt"
	found, err := fs.FileExists(filePath)
	if err != nil {
		t.Error("d5a: bad FileExists")
	}
	if !found {
		t.Error("d5b: bad FileExists")
	}

	dirPath, err = fs.ReadFileIntoString(filePath)
	if err != nil {
		t.Error("d6a: bad ReadFileIntoString")
	}
	if len(dirPath) < 10 {
		t.Error("d6b: bad ReadFileIntoString")
	}

	lines, err := fs.ReadTextLines(filePath, false)
	if err != nil {
		t.Error("d7a: bad ReadTextLines")
	}
	if len(lines) < 1 {
		t.Error("d7b: bad ReadTextLines")
	}

	err = fs.WriteTextLines(lines, filePath, false)
	if err != nil {
		t.Error("d8a: bad WriteTextLines")
	}
	if len(lines) < 1 {
		t.Error("d8b: bad WriteTextLines")
	}

	dirname := prefix + "acmFiles/" // searches only *.html files (e.g., ../acmFiles)
	since := nt.New_NullTime("2020-01-01")
	lines, err = fs.GetFileList(dirname, since)
	if err != nil {
		t.Error("d9a: bad GetFileList")
	}
	if len(lines) < 1 {
		t.Error("d9b: bad GetFileList")
	}

	since, err = fs.GetMostRecentFileAsNullTime(dirname)
	if err != nil {
		t.Error("d10a: bad GetMostRecentFileAsNullTime")
	}
	if len(lines) < 1 {
		t.Error("d10b: bad GetMostRecentFileAsNullTime")
	}

	fileName := dirname + "apr-14-2021.html"
	i64 := fs.GetFileTime(fileName)
	if i64 < 1 {
		t.Error("d11: bad GetFileTime")
	}

	dirPath = fs.GetSourceDirectory()
	if len(dirPath) < 1 {
		t.Error("d12: bad GetSourceDirectory")
	}

	//func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	//func ZipFiles(pathPrefix string, fileExt string, targetFileName string) error {
	//func ReadOccurrenceListFromCsvFile(filePath string) ([]hd.Occurrence, error) {
	//func (fss *FileService) GetTextFile(ctx *gin.Context) {
}

/* article ************************************************************************************/

// Test_article func
func Test_article(t *testing.T) {
	count, err := art.GetArticleCount()
	if err != nil {
		t.Error("e1a: bad GetArticleCount")
	}
	if count < 1 {
		t.Error("e1b: bad GetArticleCount")
	}

	earliestArchiveDate, latestArchiveDate, err := art.GetLastDateSavedFromDb()
	if err != nil {
		t.Error("e2a: bad GetLastDateSavedFromDb")
	}
	if earliestArchiveDate == latestArchiveDate {
		t.Error("e2b: bad GetLastDateSavedFromDb")
	}

	/*vocabList, err := art.WordFrequencyList()
	if err != nil {
		t.Error("e3a: bad WordFrequencyList")
	}
	if len(vocabList) < 90000 {
		t.Error("e3b: bad WordFrequencyList")
	}*/

	dateList := []string{"2021-04-14"}
	articleList, err := art.GetAcmArticleListByArchiveDates(dateList)
	if err != nil {
		t.Error("e4a: bad GetAcmArticleListByArchiveDates")
	}
	if len(articleList) < 1 {
		t.Error("e4b: bad GetAcmArticleListByArchiveDates")
	}
	count = len(articleList)
	fmt.Println(count)

	testDate := nt.New_NullTime(dateList[0])
	timeinterval := nt.New_TimeInterval(nt.TFYear, testDate, testDate)
	articleList, err = art.GetAcmArticleListByDate(timeinterval)
	if err != nil {
		t.Error("e5a: bad GetAcmArticleListByDate")
	}
	if len(articleList) < 1 {
		t.Error("e5b: bad GetAcmArticleListByDate")
	}
	fmt.Println(len(articleList))

	fmt.Print("    ")
}

/* <<<<
//GetAcmArticlesByID(idMap map[uint32]int, cutoff int) ([]hd.AcmArticle, error) {
CallUpdateOccurrence(timeinterval nt.TimeInterval) error {
CallUpdateTitle(timeinterval nt.TimeInterval) error {
BulkInsertAcmData(articleList []hd.AcmArticle) (int, error) {
*/

/* conditional ************************************************************************************/

/*func Test_GetConditionalByTimeInterval(t *testing.T) { // (bigrams []string, timeInterval nt.TimeInterval, bigramMap map[string]bool, includeTimeframetype bool) error
	words := []string{"3d", "able", "access"}
	permutations := getWordBigramPermutations(words, true) // permute=true
	fmt.Print("    ")
	fmt.Println(permutations)
	gBigramPresenceMap := NewBigramPresenceMap(permutations)
	timeInterval := nt.New_TimeInterval(nt.TimeFrameType(nt.TFSpan), nt.New_NullTime("2008-01-01"), nt.New_NullTime("2011-12-31"))
	fmt.Print("    ")
	fmt.Println(timeInterval)
	condProbList, err := getConditionalByTimeInterval(permutations, timeInterval, gBigramPresenceMap.Presence, false) // includeTimeframetype=false
	n := len(condProbList)
	if err != nil {
		t.Error("Expected > 0, got ", n)
	}
	fmt.Print("    ")
	for _, c := range condProbList {
		fmt.Print(c.WordList + "  ")
	}
	fmt.Println()
} */

// Test_conditional func
func Test_conditional(t *testing.T) {

}

/*
 isHexWord(word string) bool {
 FilteringRules(word string) (string, int) {
 SelectOccurrenceByDate(occurrenceList []hd.Occurrence, timeinterval nt.TimeInterval) []hd.Occurrence {
 SelectOccurrenceByID(occurrenceList []hd.Occurrence, acmID uint32) []hd.Occurrence {
 SelectOccurrenceByWord(occurrenceList []hd.Occurrence, word string) []hd.Occurrence {
 GetOccurrenceListByDate(timeinterval nt.TimeInterval) ([]hd.Occurrence, mapset.Set, error) {
 CollectWordGrams(wordGrams []string, timeinterval nt.TimeInterval) ([]hd.Occurrence, mapset.Set) {
 GetOccurrencesByAcmid(xacmid uint32) ([]hd.Occurrence, error) {
 WordGramSubset(alphaWord string, vocabList []hd.Vocabulary, occurrenceList []hd.Occurrence) []string {
 GetDistinctDates(occurrenceList []hd.Occurrence) []nt.NullTime {
 GetDistinctWords(occurrenceList []hd.Occurrence) []string {
 BulkInsertConditionalProbability(conditionals []hd.ConditionalProbability) error {
 ExtractKeysFromProbabilityMap(wordMap map[string]float32) []string {
 CalcConditionalProbability(startingWordgram string, wordMap map[string]float32, timeinterval nt.TimeInterval) (int, error) {
 GetConditionalByTimeInterval(bigrams []string, timeInterval nt.TimeInterval, bigramMap map[string]bool, includeTimeframetype bool) ([]hd.ConditionalProbability, error) {
 GetConditionalByProbability(word string, probabilityCutoff float32, timeInterval nt.TimeInterval, condProbList *[]hd.ConditionalProbability) error {
 GetWordBigramPermutations(words []string, permute bool) []string {
 GetConditionalList(words []string, timeInterval nt.TimeInterval, permute bool) ([]hd.ConditionalProbability, error) {
 GetExistingConditionalBigrams(bigrams []string, intervalClause string) ([]string, error) {
 GetProbabilityGraph(words []string, timeInterval nt.TimeInterval) ([]hd.ConditionalProbability, error) {
 GetWordgramConditionalsByInterval(words []string, timeInterval nt.TimeInterval, dimensions int) ([]hd.WordScoreConditionalFlat, error) {
*/
/* simplex ************************************************************************************/
/*
func getTableNames(useTempTable bool) []string {
func GetSimplexByNameUserID(userID int, simplexName, simplexType string, useTempTable bool) ([]hd.SimplexComplex, error) {
func GetSimplexListByUserID(userID int, useTempTable bool) ([]hd.SimplexComplex, error) {
func InsertSimplexComplex(sc hd.SimplexComplex) (hd.SimplexComplex, error) {
func BulkInsertSimplexFacets(facets []hd.SimplexFacet) error {
func PostSimplexComplex(userID int, simplexName, simplexType string, timeInterval nt.TimeInterval) ([]uint64, error) {
func GetSimplexWordDifference(complexIdList []uint64) ([]hd.KeyValueStringPair, error) {
func InsertCategoryWords(categoryID uint64, words []string) error {
func InsertWordCategory(description string) (hd.CategoryTable, error) {
func GetSpecialMap(category int) ([]hd.SpecialTable, error) {
func GetCategoryMap() ([]hd.CategoryTable, error) {
*/

/* profile ************************************************************************************/
/*
	func Encrypt(key, data []byte) ([]byte, error) {
	func Decrypt(key, data []byte) ([]byte, error) {
	func GenerateKey() ([]byte, error) {
	func DeriveKey(password, salt []byte) ([]byte, []byte, error) {
	func EncryptData(password, textdata string) string {
	func DecryptData(password string, ciphertext []byte) (string, error) {
	func GetUserProfile(userName, pwdText string) (hd.UserProfile, error) {
	func InsertUserProfile(userName, userEmail, pwdText string, acmmemberid int) (hd.UserProfile, error) {

	func GetVocabularyByWord(wordX string) (hd.Vocabulary, error) {
	func GetVocabularyList(words []string) ([]hd.Vocabulary, error) {
	func getAcmGraphCount() string {
	func GetWordListMap(prefix string) ([]hd.LookupMap, error) {
	func GetVocabularyListByDate(timeinterval nt.TimeInterval) ([]hd.Vocabulary, error) {
	func GetVocabularyMapProbability(wordGrams []string, timeInterval nt.TimeInterval) (map[string]float32, error) {
	func GetTitleWordsBigramInterval(bigrams []string, timeInterval nt.TimeInterval, useOccurrence bool) ([]hd.Occurrence, error) {
	func UpdateVocabulary(recordList []hd.Vocabulary) (int, error) {
	func GetVocabularyMap(fieldName string) (map[string]int, error) {
	func BulkInsertVocabulary(recordList []hd.Vocabulary) (int, error) {
	func CallUpdateVocabulary() error {
	func GetLookupValues(tableName, columnName string) ([]string, error) {

	wordscore/wordscore.go:func GetWordScores(word string) ([]hd.WordScore, error) {
	wordscore/wordscore.go:func GetWordScoreListByTimeInterval(words []string, timeInterval nt.TimeInterval) ([]hd.WordScore, error) {
	wordscore/wordscore.go:func BulkInsertWordScores(wordScoreList []hd.WordScore) error {
*/
