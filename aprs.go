/*
    flightLGo / OGN automate flightlog written in Golang
    Copyright (C) 2020  Sebastien Chaumontet

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

package main

import (
	"bufio"
	"fmt"
	"github.com/martinhpedersen/libfap-go"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func Listen(processor func(packet *fap.Packet)) {
	defer fap.Cleanup()
	var reader *bufio.Reader

	if len(os.Args) > 1 {
		reader = file_reader(os.Args[1])
	} else {
		reader = aprs_reader()
	}

	each_message(reader, processor)
}

func connect() net.Conn {
	connection, err := net.Dial("tcp", ognAPRSserver)
	if err != nil {
		panic(err)
	} else {
		return connection
	}
}

func authenticate(c net.Conn) {
	auth := fmt.Sprintf("user "+softwareName+" pass -1 vers "+softwareName+" "+softwareVersion+" filter r/%s/%s/%s\n",
		os.Getenv("LAT"),
		os.Getenv("LNG"),
		os.Getenv("RADIUS"),
	)
	fmt.Fprintf(c, auth)
}

func keepalive(c net.Conn) {
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		for t := range ticker.C {
			fmt.Fprintf(c, "# "+softwareName+" "+softwareVersion+" "+os.Getenv("ICAO")+" keepalive %s\n", t)
		}
	}()
}

func each_message(r *bufio.Reader, processor func(packet *fap.Packet)) {
	for {
		line, err := r.ReadString('\n')
		//fmt.Println("each_message: "+line)
		if err == io.EOF {
			log.Println(err)
			return
		} else if err != nil {
			log.Println(err)
			return
			//panic(err)
		}

		// APRS,qAS: aircraft beacon
		// APRS,TCPIP*,qAC: ground station beacon
		// ,TCPIP*,qAC,
		if (!strings.Contains(line, ",TCPIP*,qAC,") && strings.Contains(line, ":/")) {
			packet, err := fap.ParseAprs(line, false)
			if (err == nil && packet.Timestamp.Unix() > 0) {
				processor(packet)
			} else {
				log.Print(line)
				log.Println(err)
				//log.Fatal(err)
			}
		} else {
			//log.Print("Ignoring: "+line)
		}
	}
}

func aprs_reader() *bufio.Reader {
	connection := connect()
	authenticate(connection)
	keepalive(connection)
	fmt.Println("Connected to OGN APRS server to track activity of "+os.Getenv("ICAO")+" centered on "+os.Getenv("LAT")+" "+os.Getenv("LNG")+" with a radius of "+os.Getenv("RADIUS")+" km.")

	return bufio.NewReader(connection)
}

func file_reader(fn string) *bufio.Reader {
	fmt.Printf("Reading from %s\n", fn)

	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	return bufio.NewReader(f)
}
