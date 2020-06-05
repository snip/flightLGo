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

const softwareName = "flightLGo"
const softwareVersion = "0.0.0b1"
const ognAPRSserver = "aprs.glidernet.org:14580" // Filtered
//const ognAPRSserver = "aprs.glidernet.org:10152" // Full feed
const ognFlightLogURL = "https://flightlog.glidernet.org/push"
