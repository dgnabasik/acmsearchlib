package nulltime

// Provide UTC time functions for ACM time (restricted), Unix time, timestamp, google.protobuf.Timestamp.
// See https://golang.org/pkg/time/ 	Deploy using 'go install'	Not tested on Windows.
// See https://developers.google.com/protocol-buffers/docs/reference/java/com/google/protobuf/Timestamp

import (
	"database/sql/driver" // sb driver Value interface
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	// proto "github.com/golang/protobuf/proto"
	ptypes "github.com/golang/protobuf/ptypes"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

const (
	timeFormat       = "2006-01-02 15:04:05" // RFC3339 format.
	postfixHTML      = ".html"
	Unknown          = "Unknown"
	prefixTimeOffset = "T11:45:26.371Z"
	NullDate         = "0001-01-01" // + " 00:00:00 +0000 UTC"
	VeryFirstDate    = "2000-01-03"
	VeryFirstYear    = 2000
)

// no such thing as a const array.
var shortMonthNames = []string{"", "jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}
var longMonthNames = []string{"", "January ", "February ", "March ", "April ", "May ", "June ", "July ", "August ", "September ", "October ", "November ", "December "}
var leadingZeroNumbers = []string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31"}

// GetShortMonthName func returns short month name.
func GetShortMonthName(m int) string {
	if m >= 1 && m <= 12 {
		return shortMonthNames[m]
	} else {
		return shortMonthNames[0]
	}
}

// indexOf(string) returns default 0th element. ONLY for shortMonthNames & longMonthNames!
func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return 0 //not found.
}

// return index else -1
func stringSliceContains(a []string, x string) (int, bool) {
	for ndx, n := range a {
		if x == n {
			return ndx, true
		}
	}
	return -1, false
}

/********************************************************************************/

// TimeFrameType type returns the enclosing week, month, quarter, year, & span (everything) as pairs of NullTimes given a date.
type TimeFrameType int

// constants	<<< add TFFourYears ???
const (
	TFUnknown TimeFrameType = iota
	TFWeek
	TFMonth
	TFQuarter
	TFYear
	TFSpan
)

// ToString method
func (tft TimeFrameType) ToString() string {
	return [...]string{Unknown, "Week", "Month", "Quarter", "Year", "Span"}[tft]
}

// TimeInterval struct
type TimeInterval struct {
	Timeframetype TimeFrameType
	StartDate     NullTime
	EndDate       NullTime
}

// New_TimeInterval func
func New_TimeInterval(timeframetype TimeFrameType, startDate NullTime, endDate NullTime) TimeInterval {
	p := new(TimeInterval)
	p.Timeframetype = timeframetype
	p.StartDate = startDate
	p.EndDate = endDate
	return *p
}

// GetTimeIntervalDatePartitionList Partition date ranges for concurrent processing for given StartDate, EndDates.
func (ti TimeInterval) GetTimeIntervalDatePartitionList() []TimeInterval {
	var timeIntervalList []TimeInterval

	if ti.Timeframetype == TFUnknown { // same dates
		interval := New_TimeInterval(ti.Timeframetype, ti.StartDate, ti.EndDate)
		timeIntervalList = append(timeIntervalList, interval)
	}

	if ti.Timeframetype == TFWeek {
		dYears, dMonths, dDays, _, _, _ := NullTimeDiff(ti.StartDate, ti.EndDate)
		diffWeeks := dYears*52 + dMonths*4 + dDays/7 + 2
		startOfWeek, endOfWeek := GetStartEndOfWeek(ti.StartDate)
		for weekIndex := 0; weekIndex <= diffWeeks; weekIndex++ {
			newInterval := New_TimeInterval(ti.Timeframetype, startOfWeek, endOfWeek)
			timeIntervalList = append(timeIntervalList, newInterval)
			if endOfWeek.DT.After(ti.EndDate.DT) {
				break
			}
			startOfWeek.DT = startOfWeek.DT.AddDate(0, 0, 7)
			endOfWeek.DT = startOfWeek.DT.AddDate(0, 0, 6)
		}
	}

	if ti.Timeframetype == TFMonth { // get start of month of StartDate
		dYears, dMonths, _, _, _, _ := NullTimeDiff(ti.StartDate, ti.EndDate)
		diffMonths := dYears*12 + dMonths
		startOfMonth := New_NullTime2(time.Date(ti.StartDate.DT.Year(), ti.StartDate.DT.Month(), 1, 0, 0, 0, 0, time.UTC))
		for monthIndex := 0; monthIndex <= diffMonths; monthIndex++ {
			endOfMonth := New_NullTime2(startOfMonth.DT.AddDate(0, 1, -1))
			newInterval := New_TimeInterval(ti.Timeframetype, startOfMonth, endOfMonth)
			timeIntervalList = append(timeIntervalList, newInterval)
			startOfMonth.DT = startOfMonth.DT.AddDate(0, 1, 0)
			if endOfMonth.DT.After(ti.EndDate.DT) {
				break
			}
		}
	}

	if ti.Timeframetype == TFQuarter {
		dYears, dMonths, _, _, _, _ := NullTimeDiff(ti.StartDate, ti.EndDate)
		diffQuarters := dYears*4 + dMonths/3 + 1
		xMonth := ti.StartDate.DT.Month()/3 + 1
		startOfQuarter := New_NullTime2(time.Date(ti.StartDate.DT.Year(), xMonth, 1, 0, 0, 0, 0, time.UTC))
		for quarterIndex := 0; quarterIndex <= diffQuarters; quarterIndex++ {
			endOfQuarter := New_NullTime2(startOfQuarter.DT.AddDate(0, 4, 0)) // add 4 months
			endOfQuarter = New_NullTime2(endOfQuarter.DT.AddDate(0, -1, -1))  // sub 1 day
			newInterval := New_TimeInterval(ti.Timeframetype, startOfQuarter, endOfQuarter)
			timeIntervalList = append(timeIntervalList, newInterval)
			startOfQuarter.DT = startOfQuarter.DT.AddDate(0, 3, 0)
			if endOfQuarter.DT.After(ti.EndDate.DT) {
				break
			}
		}
	}

	if ti.Timeframetype == TFYear {
		for yearIndex := ti.StartDate.DT.Year(); yearIndex <= ti.EndDate.DT.Year(); yearIndex++ {
			startOfYear := New_NullTime(strconv.Itoa(yearIndex) + "-01-01")
			endOfYear := New_NullTime(strconv.Itoa(yearIndex) + "-12-31")
			if yearIndex == ti.EndDate.DT.Year() {
				endOfYear = NullTimeToday()
			}
			newInterval := New_TimeInterval(ti.Timeframetype, startOfYear, endOfYear)
			timeIntervalList = append(timeIntervalList, newInterval)
		}
	}

	if ti.Timeframetype == TFSpan {
		timeIntervalList = append(timeIntervalList, ti)
	}

	return timeIntervalList
}

// ToString method
func (ti TimeInterval) ToString() string {
	return ti.Timeframetype.ToString() + "|" + ti.StartDate.StandardDate() + "|" + ti.EndDate.StandardDate()
}

/********************************************************************************/

// TimeFrame struct
type TimeFrame struct {
	Timeframetype  TimeFrameType
	GivenDate      NullTime
	StartOfWeek    NullTime
	EndOfWeek      NullTime
	StartOfMonth   NullTime
	EndOfMonth     NullTime
	StartOfQuarter NullTime
	EndOfQuarter   NullTime
	StartOfYear    NullTime
	EndOfYear      NullTime
	StartOfSpan    NullTime
	EndOfSpan      NullTime
}

// ToString method
func (tf TimeFrame) ToString() string {
	return tf.Timeframetype.ToString()
}

// GetDivisor method assumes that articles are published consistently 3 times per week for about 154 publish dates
// even though the average is 148. So this could be made more exact by actually counting publish dates.
// Assumes that the requested timeFrame is completely in the past. Only used in database.CalcWordScore().
func (tf TimeFrame) GetDivisor() float32 {
	today := NullTimeToday()
	year, month, day := today.DT.Date()
	switch tf.Timeframetype {
	case TFUnknown:
		return float32(1)
	case TFWeek:
		return float32(3)
	case TFMonth:
		return float32(12)
	case TFQuarter:
		return float32(37)
	case TFYear:
		return float32(148)
	case TFSpan:
		return float32((year-VeryFirstYear)*148) + float32((int(month)-1)*12) + float32(day)/float32(3)
	default:
		return float32(1)
	}
}

// GetTimeFrameDates method
func (tf TimeFrame) GetTimeFrameDates() (NullTime, NullTime) {
	switch tf.Timeframetype {
	case TFWeek:
		return tf.StartOfWeek, tf.EndOfWeek
	case TFMonth:
		return tf.StartOfMonth, tf.EndOfMonth
	case TFQuarter:
		return tf.StartOfQuarter, tf.EndOfQuarter
	case TFYear:
		return tf.StartOfYear, tf.EndOfYear
	case TFSpan:
		return tf.StartOfSpan, tf.EndOfSpan
	default: // TFUnknown
		return tf.GivenDate, tf.GivenDate
	}
}

// Print method
func (tf TimeFrame) Print() {
	fmt.Println("GivenDate     :" + tf.GivenDate.StandardDate())
	fmt.Println("StartOfWeek   :" + tf.StartOfWeek.StandardDate())
	fmt.Println("EndOfWeek     :" + tf.EndOfWeek.StandardDate())
	fmt.Println("StartOfMonth  :" + tf.StartOfMonth.StandardDate())
	fmt.Println("EndOfMonth    :" + tf.EndOfMonth.StandardDate())
	fmt.Println("StartOfQuarter:" + tf.StartOfQuarter.StandardDate())
	fmt.Println("EndOfQuarter  :" + tf.EndOfQuarter.StandardDate())
	fmt.Println("StartOfYear   :" + tf.StartOfYear.StandardDate())
	fmt.Println("EndOfYear     :" + tf.EndOfYear.StandardDate())
	fmt.Println("StartOfSpan   :" + tf.StartOfSpan.StandardDate())
	fmt.Println("EndOfSpan     :" + tf.EndOfSpan.StandardDate())
}

// GetTimeFrameFromUnixTimeStamp func
func GetTimeFrameFromUnixTimeStamp(uts UnixTimeStamp, timeframetype TimeFrameType) TimeFrame {
	tt := time.Unix(int64(uts), 0)
	nt := New_NullTime2(tt)
	timeInterval := New_TimeInterval(timeframetype, nt, nt)
	return New_TimeFrame(timeInterval)
}

/********************************************************************************/

// NullTime struct
type NullTime struct {
	DT      time.Time `json:"dt"` // in UTC
	IsValid bool      `json:"isvalid"`
}

// Scan implements the Scanner interface. Modifies self.
func (nt *NullTime) Scan(value interface{}) error {
	nt.DT, nt.IsValid = value.(time.Time)
	return nil
}

// AdvanceNextNullTime Method to return next Mon-Wed-Fri. Modifies self.
func (nt *NullTime) AdvanceNextNullTime() {
	weekday := nt.DT.Weekday()
	addDays := 0
	if weekday == time.Weekday(1) || weekday == time.Weekday(3) { // Mon || Wed
		addDays = 2
	} else if weekday == time.Weekday(5) { // Fri
		addDays = 3
	}
	nt.DT = nt.DT.AddDate(0, 0, addDays)
}

// FileSystemDate method to return dec-30-2005
func (nt NullTime) FileSystemDate() string {
	year, month, day := nt.DT.Date()
	var m int = int(month)
	tstr := GetShortMonthName(m) + "-" + leadingZeroNumbers[day] + "-" + strconv.Itoa(year)
	return tstr
}

// HtmlArchiveDate method to return 2019-07-jul (month number = month name)
func (nt NullTime) HtmlArchiveDate() string {
	year, month, _ := nt.DT.Date()
	var m int = int(month)
	tstr := strconv.Itoa(year) + "-" + leadingZeroNumbers[m] + "-" + GetShortMonthName(m)
	return tstr
}

// StandardDate method to return yyyy-mm-dd string else default time.
// RFC3339 = "2006-01-02T15:04:05Z07:00"
func (nt NullTime) StandardDate() string {
	var a [20]byte
	var b = a[:0]                           // Using the a[:0] notation converts the fixed-size array to a slice type represented by b that is backed by this array.
	b = nt.DT.AppendFormat(b, time.RFC3339) // AppendFormat() accepts type []byte. The allocated memory a is passed to AppendFormat().
	return string(b[0:10])
}

// NonStandardDate method to return mm/dd/yy string else default time.
func (nt NullTime) NonStandardDate() string {
	year, month, day := nt.DT.Date()
	var m int = int(month)
	var d int = int(day)
	tstr := leadingZeroNumbers[m] + "/" + leadingZeroNumbers[d] + "/" + leadingZeroNumbers[year-VeryFirstYear] //  strconv.Itoa(year)[2:3]
	return tstr
}

// Value implements the driver Value interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.IsValid {
		return nil, nil
	}
	return nt.DT, nil
}

// IsScheduledDate method Use time.After() to test if on scheduled date.
func (nt NullTime) IsScheduledDate(when TimeFrameType) bool {
	year, month, _ := nt.DT.Date()
	baseStartTime := NullTimeToday()

	switch when {
	case TFUnknown: // Is it past 11 am today? Articles are published around 10 am.
		baseStartTime.DT = baseStartTime.DT.Add(time.Hour * 11)
		return baseStartTime.DT.After(nt.DT)

	case TFWeek: // Is today the start of the week?  (Sunday)
		startOfWeek, _ := GetStartEndOfWeek(nt)
		return (nt.DT == startOfWeek.DT)

	case TFMonth: // Is today the start of the month?
		startOfMonth := New_NullTime2(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC))
		return (nt.DT == startOfMonth.DT)

	case TFQuarter: // Is today the start of the quarter?
		qStart := []int{0, 1, 1, 1, 4, 4, 4, 7, 7, 7, 10, 10, 10}
		startOfQuarter := New_NullTime2(time.Date(year, time.Month(qStart[month]), 1, 0, 0, 0, 0, time.UTC))
		return (nt.DT == startOfQuarter.DT)

	case TFYear: // Is today the start of the year?
		startOfYear := New_NullTime2(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC))
		return (nt.DT == startOfYear.DT)

	case TFSpan: // Set manually by testing for yesterday.
		return nt.DT.Before(baseStartTime.DT)
	}

	return false
}

/********************************************************************************/

// GetStartEndOfWeek iterates back to Sunday & iterate forward to Saturday.
func GetStartEndOfWeek(givenDate NullTime) (NullTime, NullTime) {
	p := new(TimeFrame)
	p.StartOfWeek = givenDate

	isoYear, isoWeek := p.StartOfWeek.DT.ISOWeek()
	if isoYear > isoWeek { // dummy
		for p.StartOfWeek.DT.Weekday() != time.Sunday {
			p.StartOfWeek.DT = p.StartOfWeek.DT.AddDate(0, 0, -1)
			isoYear, isoWeek = p.StartOfWeek.DT.ISOWeek()
		}
	}

	p.EndOfWeek = givenDate
	isoYear, isoWeek = p.EndOfWeek.DT.ISOWeek()
	if isoYear > isoWeek { // dummy
		for p.EndOfWeek.DT.Weekday() != time.Saturday {
			p.EndOfWeek.DT = p.EndOfWeek.DT.AddDate(0, 0, 1)
			isoYear, isoWeek = p.EndOfWeek.DT.ISOWeek()
		}
	}

	return p.StartOfWeek, p.EndOfWeek
}

// New_TimeFrame Weeks range 1 to 53. Jan 01 to Jan 03 of year n might belong to week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1 of year n+1.
// time.Date(2019, 10, 1, 0, 0, 0, 0, time.UTC) // Year,Month,Day,  Hour,Minute,Second,Nanosecond,  Location
func New_TimeFrame(timeInterval TimeInterval) TimeFrame {
	qStart := []int{0, 1, 1, 1, 4, 4, 4, 7, 7, 7, 10, 10, 10}
	nt := timeInterval.StartDate // by convention
	p := new(TimeFrame)
	p.Timeframetype = timeInterval.Timeframetype
	p.GivenDate = New_NullTime(nt.StandardDate())
	if nt.DT.Year() < VeryFirstYear {
		return *p
	}

	p.StartOfWeek, p.EndOfWeek = GetStartEndOfWeek(p.GivenDate)
	year, month, _ := p.GivenDate.DT.Date()
	if year < VeryFirstYear || month < 1 {
		fmt.Println("New_TimeFrame(): " + p.GivenDate.StandardDate())
	}

	p.StartOfMonth = New_NullTime2(time.Date(year, month, 1, 0, 0, 0, 0, time.UTC))
	p.EndOfMonth = New_NullTime2(p.StartOfMonth.DT.AddDate(0, 1, -1))
	p.StartOfQuarter = New_NullTime2(time.Date(year, time.Month(qStart[month]), 1, 0, 0, 0, 0, time.UTC))
	monthEnd := New_NullTime2(time.Date(year, time.Month(qStart[month]+2), 1, 0, 0, 0, 0, time.UTC))
	p.EndOfQuarter = New_NullTime2(monthEnd.DT.AddDate(0, 1, -1))
	p.StartOfYear = New_NullTime2(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC))
	p.EndOfYear = New_NullTime2(time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC))

	p.StartOfSpan = New_NullTime(VeryFirstDate)
	p.EndOfSpan = NullTimeToday()

	return *p
}

// New_NullTime1 handles Sep(tember) 23, 2019 format from Journal dates. Needs work.
func New_NullTime1(dt string) NullTime {
	tokens := strings.Split(dt, " ")
	if len(dt) < 12 || len(tokens) != 3 {
		return New_NullTime(NullDate)
	}

	_, longMonth := stringSliceContains(longMonthNames, tokens[0]+" ")
	_, shortMonth := stringSliceContains(shortMonthNames, tokens[0]+" ")
	if longMonth || shortMonth {
		tokens[1] = strings.Trim(tokens[1], ",")
		month, _ := strconv.Atoi(tokens[1])
		str := tokens[2] + "-"
		if longMonth {
			str += leadingZeroNumbers[indexOf(tokens[0]+" ", longMonthNames)] + "-" + leadingZeroNumbers[month]
		} else {
			str += leadingZeroNumbers[indexOf(tokens[0]+" ", shortMonthNames)] + "-" + leadingZeroNumbers[month]
		}
		return New_NullTime(str)
	}

	return New_NullTime(NullDate)
}

// New_NullTime handles time.RFC3339 format. Any date < VeryFirstDate is assigned that date.
// Return new NullTime from datetime yyyy-mm-dd string else declare print error. Test with time.IsZero().
func New_NullTime(dt string) NullTime {
	p := new(NullTime)

	tokens := strings.Split(dt, "-") // use regex: only digits and dash.
	if len(dt) < 10 || len(tokens) != 3 {
		return *p
	}

	str := dt + prefixTimeOffset
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		fmt.Println(err)
	}
	if t.Year() < VeryFirstYear {
		t, _ = time.Parse(time.RFC3339, VeryFirstDate)
	}

	// normalize hours to zero.
	p.DT = t.Truncate(24 * time.Hour)
	p.IsValid = true
	return *p
}

// New_NullTime2 func
func New_NullTime2(dt time.Time) NullTime {
	var a [20]byte
	var b = a[:0]                        // Using the a[:0] notation converts the fixed-size array to a slice type represented by b that is backed by this array.
	b = dt.AppendFormat(b, time.RFC3339) // AppendFormat() accepts type []byte. The allocated memory a is passed to AppendFormat().
	return New_NullTime(string(b[0:10]))
}

// New_NullTimeFromFileName converts dec-30-2005.html to NullTime (default hours...)
func New_NullTimeFromFileName(htmlFile string) NullTime {
	fragments := strings.Split(htmlFile, "/")
	fileName := strings.TrimSuffix(fragments[len(fragments)-1], postfixHTML)
	return New_NullTime(GetStandardDateForm(fileName))
}

// NullTimeDiff func
func NullTimeDiff(startDate NullTime, endDate NullTime) (year, month, day, hour, min, sec int) {
	a := startDate.DT
	b := endDate.DT
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// NullTimeSorter func
func NullTimeSorter(nullTimes []NullTime) []NullTime {
	sort.SliceStable(nullTimes, func(i, j int) bool {
		return nullTimes[i].StandardDate() < nullTimes[j].StandardDate()
	})
	return nullTimes
}

// CurrentTimeString func
func CurrentTimeString() string {
	return time.Now().Format(timeFormat)
}

// NullTimeToday always returns UTC!
func NullTimeToday() NullTime {
	dt := time.Now().UTC()
	// normalize hours to zero.
	dt = dt.Truncate(24 * time.Hour)
	return New_NullTime2(dt)
}

// GetStandardDateForm converts mmm-dd-yyyy to yyyy-mm-dd
func GetStandardDateForm(dt string) string {
	return dt[7:11] + "-" + leadingZeroNumbers[indexOf(dt[0:3], shortMonthNames)] + "-" + dt[4:6]
}

/********************************************************************************/

// UnixTimeStamp type
type UnixTimeStamp int64 // not uint64!

// TimeStampInterval struct
type TimeStampInterval struct { // Enforce exact UTC time.
	Timeframetype TimeFrameType
	StartTime     UnixTimeStamp
	EndTime       UnixTimeStamp
}

// GetCurrentTimeStamp func
func GetCurrentTimeStamp() *timestamp.Timestamp {
	return ptypes.TimestampNow()
}

// GetTimeStampFromUnixTimeStamp func
func GetTimeStampFromUnixTimeStamp(uts UnixTimeStamp) *timestamp.Timestamp {
	tt := GetTimeFromUnixTimeStamp(uts)
	ts, _ := ptypes.TimestampProto(tt)
	return ts
}

// GetNullTimeFromTimeStamp converts
func GetNullTimeFromTimeStamp(tp *timestamp.Timestamp) NullTime {
	tt, _ := ptypes.Timestamp(tp)
	nt := New_NullTime2(tt)
	return nt
}

// GetUnixTimeStampFromTime func.
func GetUnixTimeStampFromTime(t time.Time) UnixTimeStamp {
	uts := t.UTC().Unix() // int64
	return UnixTimeStamp(uts)
}

// GetTimeFromUnixTimeStamp func
func GetTimeFromUnixTimeStamp(uts UnixTimeStamp) time.Time {
	return time.Unix(int64(uts), 0)
}

// GetCurrentUnixTimeStamp func
func GetCurrentUnixTimeStamp() UnixTimeStamp {
	return UnixTimeStamp(time.Now().Unix())
}

// FormatUnixTimeStampAsString func
func FormatUnixTimeStampAsString(uts UnixTimeStamp) string {
	return strconv.FormatInt(int64(uts), 10)
}

// FormatUnixTimeStampAsTime returns 1969-12-31 17:00:00 for zero value.
func FormatUnixTimeStampAsTime(uts UnixTimeStamp) string {
	t := time.Unix(int64(uts), 0)
	return t.Format(timeFormat)
}
