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
	"encoding/csv"
	"log"
	"net/http"
	"strings"
)

type Aircraft struct {
	Id           string
	Model        string
	Registration string
	Callsign     string
	Tracked      string
	Identified   string
}

type AircraftList map[string]Aircraft

var aircrafts AircraftList

func Download() {
	aircrafts = make(AircraftList)

	response, err := http.Get("http://ddb.glidernet.org/download")
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	parse(response)
	fmt.Println("OGN DDB loaded.")
}

func parse(response *http.Response) {
	csv_reader := csv.NewReader(response.Body)
	csv, err := csv_reader.ReadAll()

	if err != nil {
		panic(err)
	}

	// TODO: iterator
	for _, line := range csv {
		a := Aircraft{
			Id:           strings.Replace(line[1], "'", "", 2),
			Model:        strings.Replace(line[2], "'", "", 2),
			Registration: strings.Replace(line[3], "'", "", 2),
			Callsign:     strings.Replace(line[4], "'", "", 2),
			Tracked:      strings.Replace(line[5], "'", "", 2),
			Identified:   strings.Replace(line[6], "'", "", 2)}

		aircrafts[a.Id] = a
	}

}

func GetAircraft(id string) (Aircraft, bool) {
	a, ok := aircrafts[id]
	return a, ok
}
