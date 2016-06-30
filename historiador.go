package main

import (
	"regexp"
	"strconv"
)

const (
	maxPositionSamples = 30
)

type historicalDataRequest struct {
	ID         int
	NumSamples int
}

type historicalDataReply struct {
	NumSamplesAvailable int
	Position            [maxPositionSamples]position
}

func parseHistRequest(str []byte) (id int, samples int) {

	r, _ := regexp.Compile("REQ_HIST;(\\d*);(\\d*)")
	submatch := r.FindStringSubmatch(string(str))
	id, _ = strconv.Atoi(submatch[1])
	samples, _ = strconv.Atoi(submatch[2])
	return
}

func getHistoricalData(id int, samples int) []positionData {

	return []positionData{}
}
