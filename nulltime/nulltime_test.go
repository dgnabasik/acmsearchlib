package nulltime

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

// Test_SupportFunctions func
func Test_SupportFunctions(t *testing.T) {
	fmt.Println("Test_SupportFunctions...")
	str := GetShortMonthName(13)
	if str != "" {
		t.Error("a1:Expected '', got ", str)
	}
	str = GetShortMonthName(12)
	if str != "dec" {
		t.Error("a2:Expected 'dec', got ", str)
	}

	data := make([]string, 0)
	ndx, found := StringSliceContains(data, "dec")
	if ndx >= 0 || found {
		t.Error("a3:Expected -1, got ", ndx)
	}
	data = append(data, "dec")
	ndx, found = StringSliceContains(data, "dec")
	if ndx < 0 || !found {
		t.Error("a4:Expected 0, got ", ndx)
	}

}

// Test_NullTime func
func Test_NullTime(t *testing.T) {
	fmt.Println("Test_NullTime...")

	ti := New_NullTime("")
	if ti.StandardDate() != NullDate {
		t.Error("b1: Expected " + NullDate)
	}

	ti = New_NullTime("1999-01-01") // Valid=true
	if ti.StandardDate() != NullDate {
		t.Error("b2: Expected " + NullDate)
	}

	ti = New_NullTime("2000-13")
	if ti.StandardDate() != NullDate {
		t.Error("b3: Expected " + NullDate)
	}

	ti = New_NullTime("2000-01-08")
	if ti.StandardDate() != NullDate {
		t.Error("b4: Expected " + NullDate)
	}

	ti = New_NullTime(VeryFirstDate)
	if ti.StandardDate() != VeryFirstDate {
		t.Error("b5: Expected " + VeryFirstDate)
	}

	ti = New_NullTime1("September 23, 2019")
	if ti.StandardDate() != "2019-09-23" {
		t.Error("b6: Expected 2019-09-23")
	}

	/*ti = NullTimeToday()
	if ti.StandardDate() != "2019-11-25" {
		t.Error("b7: Expected 2019-11-25")
	}*/

	ti = New_NullTimeFromFileName("dec-05-2005.html")
	if ti.StandardDate() != "2005-12-05" {
		t.Error("b8: Expected 2005-12-05")
	}

	if ti.FileSystemDate() != "dec-05-2005" {
		t.Error("b9: Expected dec-05-2005")
	}

	if ti.HtmlArchiveDate() != "2005-12-dec" {
		t.Error("b10: Expected 2005-12-dec")
	}

	if ti.NonStandardDate() != "12/05/05" {
		t.Error("b11: Expected 12/05/05")
	}

	if GetStandardDateForm("dec-05-2005") != "2005-12-05" { // Convert mmm-dd-yyyy to yyyy-mm-dd
		t.Error("b12: Expected 2005-12-05")
	}

	ti.AdvanceNextNullTime()
	if ti.StandardDate() != "2005-12-07" {
		t.Error("b13: Expected 2005-12-07")
	}

	nt1, nt2 := GetStartEndOfWeek(ti)
	if nt1.StandardDate() != "2005-12-04" {
		t.Error("b14: Expected 2005-12-04")
	}
	if nt2.StandardDate() != "2005-12-10" {
		t.Error("b15: Expected 2005-12-10")
	}

	year, month, day, hour, min, sec := NullTimeDiff(nt1, nt2)
	if year != 0 || month != 0 || day != 6 || hour != 0 || min != 0 || sec != 0 {
		t.Error("b16: Expected 6 days diff.")
	}

	// fmt.Printf("%d %d %d %d %d %d\n", year, month, day, hour, min, sec)

	nt := NullTimeToday()
	dateSet := make([]NullTime, 0)
	dateSet = append(dateSet, nt)
	dateSet = append(dateSet, ti)
	dateSet = NullTimeSorter(dateSet)
	if len(dateSet) != 2 {
		t.Error("b19: Sorted dateSet != 2 ")
	}

	nt = New_NullTime("2019-12-01")
	fmt.Println("Testing with " + nt.StandardDate())
	yes := nt.IsScheduledDate(TFUnknown)
	if !yes {
		t.Error("b20: Not past 11am.")
	}

	yes = nt.IsScheduledDate(TFWeek)
	if !yes {
		t.Error("b21: Not start of the week.")
	}

	yes = nt.IsScheduledDate(TFMonth)
	if !yes {
		t.Error("b22: Not start of the month.")
	}

	nt = New_NullTime("2020-01-01")
	fmt.Println("Testing with " + nt.StandardDate())
	yes = nt.IsScheduledDate(TFQuarter)
	if !yes {
		t.Error("b23: Not start of the Quarter.")
	}

	yes = nt.IsScheduledDate(TFYear)
	if !yes {
		t.Error("b24: Not start of the year.")
	}

}

/*************************************************************************/

func Test_TimeInterval(t *testing.T) {
	fmt.Println("Test_TimeInterval...")

	// GetTimeFrameFromUnixTimeStamp (uts UnixTimeStamp, timeframetype TimeFrameType) TimeFrame {

	vfd1 := New_NullTime(VeryFirstDate)
	vfd2 := vfd1
	ti := New_TimeInterval(TFUnknown, vfd1, vfd2)
	if ti.StartDate.StandardDate() != ti.EndDate.StandardDate() {
		t.Error("c1: Expected " + ti.EndDate.StandardDate())
	}
	// add one time unitL len(list) always 2 except TFUnknown is 1.
	tiList := ti.GetTimeIntervalDatePartitionList()
	if len(tiList) != 1 && tiList[0].StartDate.StandardDate() != tiList[0].EndDate.StandardDate() {
		t.Error("c2: Expected " + tiList[0].StartDate.StandardDate())
	}

	vfd2 = New_NullTime2(vfd1.DT.AddDate(0, 0, 7)) // y,m,d
	ti = New_TimeInterval(TFWeek, vfd1, vfd2)
	tiList = ti.GetTimeIntervalDatePartitionList()
	if len(tiList) != 2 {
		t.Error("c3: Expected 2 tiList items.")
	}
	year, month, day, hour, min, sec := NullTimeDiff(vfd1, vfd2)
	if year != 0 || month != 0 || day != 7 || hour != 0 || min != 0 || sec != 0 {
		t.Error("c4: Expected 7 days diff.")
	}

	vfd2 = New_NullTime2(vfd1.DT.AddDate(0, 1, 0))
	ti = New_TimeInterval(TFMonth, vfd1, vfd2)
	tiList = ti.GetTimeIntervalDatePartitionList()
	if len(tiList) != 2 {
		t.Error("c5: Expected 2 tiList items.")
	}
	year, month, day, hour, min, sec = NullTimeDiff(vfd1, vfd2)
	if year != 0 || month != 1 || day != 0 || hour != 0 || min != 0 || sec != 0 {
		t.Error("c6: Expected 1 month diff.")
	}

	// skip TFQuarter & TFSpan

	vfd2 = New_NullTime2(vfd1.DT.AddDate(1, 0, 0))
	ti = New_TimeInterval(TFYear, vfd1, vfd2)
	tiList = ti.GetTimeIntervalDatePartitionList()
	if len(tiList) != 2 {
		t.Error("c7: Expected 2 tiList items.")
	}
	year, month, day, hour, min, sec = NullTimeDiff(vfd1, vfd2)
	if year != 1 || month != 0 || day != 0 || hour != 0 || min != 0 || sec != 0 {
		t.Error("c8: Expected 1 year diff.")
	}

}

func Test_TimeFrame(t *testing.T) {
	fmt.Println("Test_TimeFrame...")
	vfd1 := New_NullTime(VeryFirstDate) // "2000-01-03"
	vfd2 := New_NullTime2(vfd1.DT.AddDate(1, 0, 0))
	ti := New_TimeInterval(TFYear, vfd1, vfd2)
	tf := New_TimeFrame(ti)

	if tf.GivenDate.StandardDate() != "2000-01-03" {
		t.Error("d1: Expected 2000-01-03.")
	}

	if tf.StartOfWeek.StandardDate() != "2000-01-02" {
		t.Error("d2: Expected 2000-01-02.")
	}

	if tf.EndOfWeek.StandardDate() != "2000-01-08" {
		t.Error("d3: Expected 2000-01-08.")
	}

	if tf.StartOfMonth.StandardDate() != "2000-01-01" {
		t.Error("d4: Expected 2000-01-01.")
	}

	if tf.EndOfMonth.StandardDate() != "2000-01-31" {
		t.Error("d5: Expected 2000-01-31.")
	}

	if tf.StartOfQuarter.StandardDate() != "2000-01-01" {
		t.Error("d6: Expected 2000-01-01.")
	}

	if tf.EndOfQuarter.StandardDate() != "2000-03-31" {
		t.Error("d7: Expected 2000-03-31.")
	}

	if tf.StartOfYear.StandardDate() != "2000-01-01" {
		t.Error("d8: Expected 2000-01-01.")
	}

	if tf.EndOfYear.StandardDate() != "2000-12-31" {
		t.Error("d9: Expected 2000-12-31.")
	}

	if tf.StartOfSpan.StandardDate() != VeryFirstDate {
		t.Error("d10: Expected " + VeryFirstDate)
	}

	today := NullTimeToday()
	if tf.EndOfSpan.StandardDate() != today.StandardDate() {
		t.Error("d11: Expected " + today.StandardDate())
	}

	year, month, day := today.DT.Date()
	tf.Timeframetype = TFSpan
	divisor := float32((year-VeryFirstYear)*148) + float32((int(month)-1)*12) + float32(day)/float32(3)
	if tf.GetDivisor() != divisor {
		t.Error("d12: Expected " + strconv.FormatFloat(float64(divisor), 'E', -1, 32))
	}

	vfd1, vfd2 = tf.GetTimeFrameDates()
	if vfd1.StandardDate() != VeryFirstDate && vfd2.StandardDate() != today.StandardDate() {
		t.Error("d13: Expected " + today.StandardDate())
	}
}

func Test_TimeStamp(t *testing.T) {
	fmt.Println("Test_TimeStamp...")
	timeInUTC := time.Date(VeryFirstYear, 1, 1, 1, 1, 1, 100, time.UTC)

	uts := GetUnixTimeStampFromTime(timeInUTC)
	if uts != 946688461 {
		t.Error("e1: Expected 946688461.")
	}

	utsstr := FormatUnixTimeStampAsString(uts)
	if utsstr != "946688461" {
		t.Error("e2: Expected 946688461.")
	}

	utsstr = FormatUnixTimeStampAsTime(uts)
	if utsstr != "1999-12-31 18:01:01" {
		t.Error("e3: Expected 1999-12-31 18:01:01.")
	}

}
