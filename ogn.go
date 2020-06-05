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
	"regexp"
	"strconv"
	"strings"
	//"fmt"
)

type Comment struct {
	Id             string
	SignalStrength string
	Frequency      string
	Rot            string
	ClimbRate      float64
	Fpm            string
	Errors         string
}

// example comment: id02DF0A52 -019fpm +0.0rot 55.2dB 0e -9.9kHz gps3x6
func ParseComment(c string,aprsSource string) Comment {

	comment := Comment{}
	items := strings.Split(strings.TrimSpace(c), " ")

	for _, item := range items {
		//fmt.Println("Item: ", item)
		if  strings.HasPrefix(item,"id") {
			comment.Id = extractId(item)
		}
	}
	if comment.Id == "" { // If no idxxx found in comment takes the aprsSource as Id
		comment.Id = aprsSource[3:len(aprsSource)] // Remove the 3 first chars are they are FLR or OGN or ...
	}
	//fmt.Println("Id: ", comment.Id)

	return comment
}

func extractId(s string) string {
	id_matcher := regexp.MustCompile(`id\w{2}(\w+)`)
	return strings.TrimSpace(id_matcher.FindStringSubmatch(s)[1])
}

func extractClimbRate(s string) float64 {
	climb_rate_matcher := regexp.MustCompile(`([+-]\d+\.\d+)rot`)
	climb_rate_str := climb_rate_matcher.FindStringSubmatch(s)[1]
	climb_rate, _ := strconv.ParseFloat(climb_rate_str, 64)
	return climb_rate
}
