package timestampinterval

import (
	nt "github.com/dgnabasik/acmsearchlib/nulltime"
	pbx "github.com/dgnabasik/acmsearchlib/timestampinterval/pb" // timestampinterval.proto
)

/* Proto messaging helper functions ************************************************************/

// NewTimeEventRequest func
func NewTimeEventRequest(topic string, pbtft pbx.MTimeStampInterval_MTimeFrameType) *pbx.TimeEventRequest {
	p := new(pbx.TimeEventRequest)
	p.Topic = topic
	theTime := nt.GetTimeStampFromUnixTimeStamp(nt.GetCurrentUnixTimeStamp())
	p.Timestampinterval = new(pbx.MTimeStampInterval)
	p.Timestampinterval.Timeframetype = pbtft
	p.Timestampinterval.StartTime = theTime
	p.Timestampinterval.EndTime = theTime
	return p
}

// NewTimeEventResponse func
func NewTimeEventResponse() *pbx.TimeEventResponse {
	p := new(pbx.TimeEventResponse)
	p.Completed = false
	p.Error = nil
	return p
}

// NewTimeStampInterval func
func NewTimeStampInterval(timeframetype pbx.MTimeStampInterval_MTimeFrameType, startTime nt.UnixTimeStamp, endTime nt.UnixTimeStamp) *pbx.MTimeStampInterval {
	p := new(pbx.MTimeStampInterval)
	p.Timeframetype = timeframetype
	p.StartTime = nt.GetTimeStampFromUnixTimeStamp(startTime)
	p.EndTime = nt.GetTimeStampFromUnixTimeStamp(endTime)
	return p
}

/************************************************************/
/*
func main() {
}*/
