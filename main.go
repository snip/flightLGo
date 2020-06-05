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
	"fmt"
	"strconv"
	"net/http"
	"net/url"
	"os"
	"github.com/martinhpedersen/libfap-go"
)

type Beacon struct {
	*fap.Packet
	Aircraft
	Comment
}

type AircraftStatus struct {
	status string
	oldBeacon *Beacon
	previousBeacon *Beacon
}

var aircraftList = map[string]AircraftStatus{}

func main() {
	fmt.Println(softwareName+" "+softwareVersion+" Copyright (C) 2020 Sebastien Chaumontet - https://github.com/snip/flightLGo")
	fmt.Println("This program comes with ABSOLUTELY NO WARRANTY.")
	fmt.Println("This is free software, and you are welcome to redistribute it under GPL v3.0 conditions.\n")
	Load()
	Download()
	Listen(process_message)
}

func takeoff(b *Beacon, flightLogUrl string) {
	http.PostForm(flightLogUrl,
		url.Values{
			"aircraftId": {b.Comment.Id},
			"aircraftReg": {b.Aircraft.Registration},
			"aircraftCN": {b.Aircraft.Callsign},
			"aircraftTracked": {b.Aircraft.Tracked},
			"aircraftIdentified": {b.Aircraft.Identified},
			"takeoffTimestamp": {strconv.FormatInt(b.Packet.Timestamp.Unix(),10)},
			"takeoffAirfield": {os.Getenv("ICAO")},
		})
	if os.Getenv("alternativeFlightLogURL") != "" && os.Getenv("alternativeFlightLogURL") != flightLogUrl {
		takeoff(b,os.Getenv("alternativeFlightLogURL"))
	}
}
func landing(b *Beacon, flightLogUrl string) {
	http.PostForm(flightLogUrl,
		url.Values{
			"aircraftId": {b.Comment.Id},
			"aircraftReg": {b.Aircraft.Registration},
			"aircraftCN": {b.Aircraft.Callsign},
			"aircraftTracked": {b.Aircraft.Tracked},
			"aircraftIdentified": {b.Aircraft.Identified},
			"landingTimestamp": {strconv.FormatInt(b.Packet.Timestamp.Unix(),10)},
			"landingAirfield": {os.Getenv("ICAO")},
		})
	if os.Getenv("alternativeFlightLogURL") != "" && os.Getenv("alternativeFlightLogURL") != flightLogUrl {
		landing(b,os.Getenv("alternativeFlightLogURL"))
	}
}

// packet: https://github.com/martinhpedersen/libfap-go/blob/832c8336185c0a6de6b792ad1531a30eac09398d/packet.go
func process_message(p *fap.Packet) {
	c := ParseComment(p.Comment,p.SrcCallsign)
	a, ok := GetAircraft(c.Id)
	var b Beacon

	if ok {
		b = Beacon{Packet: p, Comment: c, Aircraft: a}
	} else { // Not in ddb
		b = Beacon{Packet: p, Comment: c}
		//fmt.Print("Not found in DDB: ", p.OrigPacket)
	}

	if b.Comment.Id != "" {
		// TODO: ignore if errors > 5
		// TODO: Check to have at least one second (maybe 2 or more) space with previous beacon (case of multiple receivers on same airfield)
		// TODO: delete aircraftList entry which is older than 2mn (Small memory leak)
		/*if b.Aircraft.Callsign == "DRF" || b.Aircraft.Callsign == "RGA" || b.Aircraft.Callsign == "POL" {
			//fmt.Println(b.Packet.Timestamp.Format(time.RFC3339)+"> "+b.Comment.Id +": "+ b.Aircraft.Registration+" ("+b.Aircraft.Callsign+") => Speed:"+fmt.Sprintf("%.0f",b.Speed))
			fmt.Println(strconv.FormatInt(b.Packet.Timestamp.Unix(),10)+"> "+b.Comment.Id +": "+ b.Aircraft.Registration+" ("+b.Aircraft.Callsign+") => Speed:"+fmt.Sprintf("%.0f",b.Speed))
			//fmt.Print("Orig message: ", p.OrigPacket)
		}*/
		if aircraftList[b.Comment.Id].status == "" { // First beacon for this ID
			aircraftList[b.Comment.Id] = AircraftStatus{
				"Unknown",nil,&b,
			}
			//fmt.Println("New")
		} else if aircraftList[b.Comment.Id].oldBeacon == nil { // Second beacon for this ID
			aircraftList[b.Comment.Id] = AircraftStatus{
				"Unknown",aircraftList[b.Comment.Id].previousBeacon,&b,
			}
			//fmt.Println("Second")
		} else {
			var beaconToRoll = true
			if aircraftList[b.Comment.Id].status != "Flying" {
				if b.Speed > 30 && aircraftList[b.Comment.Id].oldBeacon.Speed > 30 && aircraftList[b.Comment.Id].previousBeacon.Speed > 30 {
					if aircraftList[b.Comment.Id].status != "Unknown" {
						fmt.Println(b.Packet.Timestamp.String()+"> "+b.Comment.Id +": "+ b.Aircraft.Registration+" ("+b.Aircraft.Callsign+") ------------------ Takeoff")
						takeoff(aircraftList[b.Comment.Id].oldBeacon,ognFlightLogURL)
					}
					aircraftList[b.Comment.Id] = AircraftStatus{
						"Flying",aircraftList[b.Comment.Id].previousBeacon,&b,
					}
					beaconToRoll = false
				}  // Else we ignore this beacon as 3 beacons are not leading to same satus
			}
			if aircraftList[b.Comment.Id].status != "On ground" {
				if b.Speed < 20 && aircraftList[b.Comment.Id].oldBeacon.Speed < 20 && aircraftList[b.Comment.Id].previousBeacon.Speed < 20 {
					if aircraftList[b.Comment.Id].status != "Unknown" {
						fmt.Println(b.Packet.Timestamp.String()+"> "+b.Comment.Id+": "+ b.Aircraft.Registration+" ("+b.Aircraft.Callsign+") ------------------- Landing")
						landing(aircraftList[b.Comment.Id].oldBeacon,ognFlightLogURL)
					}
					aircraftList[b.Comment.Id] = AircraftStatus{
						"On ground",aircraftList[b.Comment.Id].previousBeacon,&b,
					}
					beaconToRoll = false
				}  // Else we ignore this beacon as 3 beacons are not leading to same satus
			}
			if beaconToRoll == true {
				aircraftList[b.Comment.Id] = AircraftStatus{
					aircraftList[b.Comment.Id].status,aircraftList[b.Comment.Id].previousBeacon,&b,
				}
			}

		}
	}
	//fmt.Printf("%+v", b)
}

func (b Beacon) String() string {
	return fmt.Sprintf("%s (%s/%s) @%f,%f %fm\n", b.Comment.Id, b.Aircraft.Callsign, b.Aircraft.Registration, b.Packet.Latitude, b.Packet.Longitude, b.Altitude)
}
