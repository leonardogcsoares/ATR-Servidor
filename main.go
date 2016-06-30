package main

import (
	"fmt"
	"net"
	"time"
)

const (
	// servHistPort = ":7000"
	// histServPort = ":7000"

	gatewayServPort   = ":7000"
	interfaceServPort = ":6000"
)

func main() {

	// Inicializar conexao com fila de mensagens
	// https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/08.1.html
	positionList = make(map[string]position, 1000000)

	fmt.Println("Init Gateway TCP Connection")
	tcpAddrGateway, err := net.ResolveTCPAddr("tcp4", gatewayServPort)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Init Gateway Listener")
	listenerGateway, err := net.ListenTCP("tcp", tcpAddrGateway)
	if err != nil {
		fmt.Println(err.Error())
	}
	// On a separate thread runs gateway connection
	go runGatewayServerConn(listenerGateway)

	fmt.Println("Init Interface TCP Connection")
	tcpAddrInterface, err := net.ResolveTCPAddr("tcp4", interfaceServPort)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Init Interface TCP Listener")
	listenerInterface, err := net.ListenTCP("tcp", tcpAddrInterface)
	if err != nil {
		fmt.Println(err.Error())
	}
	// On a separate thread runs interface connection
	go runInterfaceServerConn(listenerInterface)

	// Infinite loop to represent neverending process
	for {
	}

}

func runInterfaceServerConn(listener *net.TCPListener) {
	for {
		fmt.Println("Accepting Interface Connection")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
		request := make([]byte, 128)                          // set maximum request length to 128KB to prevent flood based attacks
		defer conn.Close()

		for {
			readLen, err := conn.Read(request)

			if err != nil {
				fmt.Println(err.Error())
				break
			}

			if readLen == 0 {
				fmt.Println("Connection Lost for interface")
				break // connection already closed by client

			} else if parseRequestMessage(request) == "REQ_ATIVOS" {
				fmt.Println("REQ_ATIVOS request from interface")
				// Respond with active clients
				// 1. Get active clients from Gateway
				// 2. Mount message with active clients using moutActiveClientsResponse
				// msg := mountActiveClientsResponse([]string{"015225", "115389"})
				msg := mountActiveClientsResponse(retrieveActiveIDs())
				fmt.Println(msg)
				// 3. Send message using conn.Write
				n, err := conn.Write([]byte(msg))
				fmt.Println("Bytes Written: ", n)
				if err != nil {
					fmt.Println("Err on writing to interface\n", err.Error())
				}

			} else if parseRequestMessage(request) == "REQ_HIST" {
				fmt.Println("REQ_HIST request from interface")
				// Respond with N historical data positions for a given ID
				// 1. Get N historical data positions for given ID
				id, samples := parseHistRequest(request)
				// 2. Mount message with active clients using mountHistoricsResponse
				msg := mountHistoricsResponse(getHistoricalData(id, samples), id)
				// 3. Send message using conn.Write
				n, err := conn.Write([]byte(msg))
				fmt.Println("Bytes Written: ", n)
				if err != nil {
					fmt.Println("Err on writing to interface\n", err.Error())
				}

			} else {
				fmt.Println("Connection sent unexpecting request: ", string(request))
				break

			}
		}

	}
}

func runGatewayServerConn(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
		request := make([]byte, 128)                          // set maximum request length to 128KB to prevent flood based attacks
		defer conn.Close()

		for {
			readLen, err := conn.Read(request)

			if err != nil {
				fmt.Println(err.Error())
				break
			}

			if readLen == 0 {
				fmt.Println("Connection Lost for interface")
				break // connection already closed by client

			} else if parseGatewayRequestType(request) == "UPDATE" {
				fmt.Println("Gateway UPDATE")
				pos := parseGatewayPosition(request)
				positionList[pos.ID] = pos

			} else if parseGatewayRequestType(request) == "DELETE" {
				fmt.Println("Gateway DELETE")
				id := parseGatewayID(request)
				delete(positionList, id)

			} else {
				fmt.Println("Connection sent unexpecting request: ", string(request))
				break
			}

		}
	}
}
