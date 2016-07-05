package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
)

const (
	maxPositionSamples = 30
)

type historicalDataRequest struct {
	ID         string
	NumSamples string
}

type historicalDataReply struct {
	NumSamplesAvailable string
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

func runHistoricalServerConnection(request []byte) (response string) {

	fmt.Println("Init Historical Addr")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", historicServPort)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Init Historical TCP Connection")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Write to socket, hist request:", string(request))
	_, err = conn.Write(request)
	if err != nil {
		fmt.Println(err.Error())
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Read result from socket", string(result))

	response = string(result)
	return

}
