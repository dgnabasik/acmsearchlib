package timestampinterval

//pbx "github.com/dgnabasik/acmsearchlib/timestamp"

/* Proto messaging helper functions ************************************************************/

/* GetTimeStampFromUnixTimeStamp func
func GetTimeStampFromUnixTimeStamp(uts nt.UnixTimeStamp) *timestamp.Timestamp {
	tt := nt.GetTimeFromUnixTimeStamp(uts)
	ts, _ := timestamp.TimestampProto(tt)
	return ts
}

// NewTimeEventRequest func
func NewTimeEventRequest(topic string, pbtft MTimeStampInterval_MTimeFrameType) *TimeEventRequest {
	p := new(TimeEventRequest)
	p.Topic = topic
	var theTime *timestamp.Timestamp //<<< := nt.GetTimeStampFromUnixTimeStamp(nt.GetCurrentUnixTimeStamp())
	p.Timestampinterval = new(MTimeStampInterval)
	p.Timestampinterval.Timeframetype = pbtft
	p.Timestampinterval.StartTime = theTime
	p.Timestampinterval.EndTime = theTime
	return p
}

// NewTimeEventResponse func
func NewTimeEventResponse() *TimeEventResponse {
	p := new(TimeEventResponse)
	p.Completed = false
	p.Error = nil
	return p
}

// NewTimeStampInterval func
func NewTimeStampInterval(timeframetype MTimeStampInterval_MTimeFrameType, startTime nt.UnixTimeStamp, endTime nt.UnixTimeStamp) *MTimeStampInterval {
	p := new(MTimeStampInterval)
	p.Timeframetype = timeframetype
	p.StartTime = nt.GetTimeStampFromUnixTimeStamp(startTime)
	p.EndTime = nt.GetTimeStampFromUnixTimeStamp(endTime)
	return p
} */
