package main

import "regexp"

const (
	maxPositionSamples = 30
)

type historicalDataRequest struct {
	ID         string
	NumSamples string
}

type historicalDataReply struct {
	NumSamplesAvailable int
	Position            [maxPositionSamples]position
}

func parseHistRequest(str []byte) (id string, samples string) {

	r, _ := regexp.Compile("REQ_HIST;(\\d*);(\\d*)")
	submatch := r.FindStringSubmatch(string(str))
	id = submatch[1]
	samples = submatch[2]
	// id, _ = strconv.Atoi(submatch[1])
	// samples, _ = strconv.Atoi(submatch[2])
	return
}

func getHistoricalData(id string, samples string) []positionData {

	return []positionData{}
}
