package main

import (
	"fmt"
	"net"
	"strings"
)

var (
	_                = net.ListenUDP
	ChachedResponses = map[string][]byte{}
)

func main() {
	fmt.Println("🔧 Setting up DNS server...")

	udpAddress, err := net.ResolveUDPAddr("udp", "127.0.0.1:3005")
	if err != nil {
		fmt.Println("Deez Nuts 🥜 there ain't no Freaking one here : ", err)
		return
	}

	UdpConnection, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		fmt.Println("Me Chalu Hou Shakat Nahi wopps sorry he bagh Karan : 😭", err)
		return
	}
	defer UdpConnection.Close()
	fmt.Printf("I'm Ready to serve you (Big Daddy) ��� on: %d\n", udpAddress.Port)
	fmt.Println("Listening to you 💌 On:", udpAddress.String())
	fmt.Println("Ima All ears !!!! ....... Yap Away")
	Buffer := make([]byte, 512)

	for {
		Buffersize, ConnectionSource, err := UdpConnection.ReadFromUDP(Buffer)
		if err != nil {
			fmt.Println("Saman Aaya Nhi arghh Nalleee 😡", err)
			break
		}

		RecivedSaman := string(Buffer[:Buffersize])
		fmt.Printf("Received %d bytes from %s: %s\n", Buffersize, ConnectionSource, RecivedSaman)

		response := CraftTheBigDaddyResponse(Buffer[:Buffersize])

		WrittenBytes, err := UdpConnection.WriteToUDP(response, ConnectionSource)
		if err != nil {
			fmt.Println("Okay Okay I can't write (Ranveer Alahbadiya Hater Vibes).")
			continue
		}
		fmt.Printf("Sent %d bytes 🥜🥜 to %s\n", WrittenBytes, ConnectionSource)
	}
}

func CraftTheBigDaddyResponse(Buffer []byte) []byte {
	Response := make([]byte, 512)

	// Copy request ID
	copy(Response[:2], Buffer[:2])

	// Set flags (Standard query response, Authoritative Answer)
	Response[2] = ﬁlppedIt(Buffer[2])
	Response[3] = 0b10000000

	// Set counts
	Response[4] = 0x00 // QDCOUNT
	Response[5] = 0x01
	Response[6] = 0x00 // ANCOUNT
	Response[7] = 0x01
	Response[8] = 0x00 // NSCOUNT
	Response[9] = 0x00
	Response[10] = 0x00 // ARCOUNT
	Response[11] = 0x00

	// Question section
	QueSection := CraftQuestionSection(Buffer[12:])
	ResponseLen := 12 + copy(Response[12:], QueSection)
	DomainName := "vdcds.cool" // Extract this dynamically from Buffer later

	// Check the cache first
	if cachedResponse, exists := ChachedResponses[DomainName]; exists {
		return cachedResponse // Return the cached response if found
	}

	// Otherwise, craft a new response
	AnswerSection := CraftAnswerSection(DomainName, net.ParseIP("127.0.0.1"))
	ResponseLen += copy(Response[ResponseLen:], AnswerSection)

	// Store the response in the cache
	ChachedResponses[DomainName] = Response[:ResponseLen]

	return Response[:ResponseLen]
}

func ﬁlppedIt(b byte) byte {
	b |= 0b10000000
	return b
}

func CraftQuestionSection(Buffer []byte) []byte {
	QueStionEndPoint := 0
	for QueStionEndPoint < len(Buffer) {
		if Buffer[QueStionEndPoint] == 0 {
			QueStionEndPoint += 5 // Include QTYPE and QCLASS
			break
		}
		QueStionEndPoint++
	}
	return Buffer[:QueStionEndPoint]
}

func CraftAnswerSection(DomainName string, IpAddress net.IP) []byte {
	AnswerSection := make([]byte, 0, 16+len(DomainName)+2+2+4+2+4)

	// NAME: pointer to domain name in question section
	AnswerSection = append(AnswerSection, 0xC0, 0x0C)

	// TYPE: A record
	AnswerSection = append(AnswerSection, 0x00, 0x01)

	// CLASS: IN
	AnswerSection = append(AnswerSection, 0x00, 0x01)

	// TTL: 300 seconds (5 minutes)
	AnswerSection = append(AnswerSection, 0x00, 0x00, 0x01, 0x2C)

	// RDLENGTH: 4 bytes for IPv4
	AnswerSection = append(AnswerSection, 0x00, 0x04)

	// RDATA: IP address
	AnswerSection = append(AnswerSection, IpAddress.To4()...)

	return AnswerSection
}

func FormatTheDomainBuddy(DomainName string) []byte {
	var ForMattedDomainName []byte
	DomainLabels := strings.Split(DomainName, ".")
	for _, SingleLable := range DomainLabels {
		ForMattedDomainName = append(ForMattedDomainName, byte(len(SingleLable)))
		ForMattedDomainName = append(ForMattedDomainName, SingleLable...)
	}
	ForMattedDomainName = append(ForMattedDomainName, 0)
	return ForMattedDomainName
}
