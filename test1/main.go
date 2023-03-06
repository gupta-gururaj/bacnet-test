package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	// Get server IP
	var ip string
	ip = "192.168.4.64"
	// addrs, err := net.InterfaceAddrs()
	// if err != nil {
	// 	fmt.Println("Failed to get IP address:", err)
	// 	return
	// }
	// for _, addr := range addrs {
	// 	if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	// 		if ipnet.IP.To4() != nil {
	// 			// fmt.Println("IP address:", ipnet.IP.String())
	// 			ip = ipnet.IP.String()
	// 			fmt.Println(">>")
	// 			return
	// 		}
	// 	}
	// }
	port := "47808"
	address := fmt.Sprintf("%s:%s", ip, port)
	fmt.Println("address", address)
	// Set up UDP connection to BACnet broadcast address
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("err-->", err)
	}

	// Establish a UDP connection to the BACnet network
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("err-->", err)
	}
	defer conn.Close()

	// Create a buffer to store the BACnet message
	var buf bytes.Buffer

	// Add the PDU header to the message
	pduHeader := []byte{0x01, 0x0a, 0x00, 0x0c}
	buf.Write(pduHeader)

	// Add the NPDU header to the message
	npduHeader := []byte{0x01, 0x00, 0x00, 0x0a, 0x00, 0x03, 0xff, 0xff, 0xff}
	buf.Write(npduHeader)

	// Add the APDU header to the message
	apduHeader := []byte{0x00, 0x0d, 0x01, 0x0c}
	buf.Write(apduHeader)

	// Add the object identifier to the message (in this case, object ID 0 analog input 1)
	objectID := []byte{0x0c, 0x01, 0x00}
	buf.Write(objectID)

	// Add the property identifier to the message (in this case, property ID 85 present value)
	propertyID := []byte{0x55}
	buf.Write(propertyID)

	// Add the index to the message (in this case, no index)
	index := []byte{}
	buf.Write(index)

	// Send the BACnet message to the network
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println("err-->", err)
	}

	// Wait for a response from the network
	response := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	n, err := conn.Read(response)
	if err != nil {
		fmt.Println("err-1->", err)
	}
	fmt.Println("n", n)

	// Parse the response from the network
	// responseData := response[22:n]
	var value float32
	binary.Read(bytes.NewReader(response), binary.BigEndian, &value)
	fmt.Printf("Present value: %.2f\n", value)
}
