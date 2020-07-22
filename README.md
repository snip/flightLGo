# flightLGo / OGN FlightLog Go
Quick and maybe dirty [OGN](https://glidernet.org) automate flightlog written in Golang

This tool is inspired by https://github.com/masone/ogn

This tool is very basic and is **not** doing the following:
- detect runway used
- detect takeoff type (towing/winch/self launch/...)
- detect tawplane towing this particular glider
- ...

For advanced features you may consider to have a look to https://ktrax.kisstech.ch/logbook

## Concept
flightLGo is running on a host as daemon mode (can run on same host as running your OGN receiver).
It analyze in realtime trafic for a small given location (airfield) and send takeoff and landing infromation to: https://flightlog.glidernet.org
It can also send this same data to another website (see web directory content for example of PHP scripts to run this website).

## Usage
```
sudo apt-get install libfap6
wget https://raw.githubusercontent.com/snip/flightLGo/master/sample.env
```
Download binary from https://github.com/snip/flightLGo/releases or build it yourself.
Copy `sample.env` to `.env` in same directory as flightLGo then update this `.env` according to your needs.
Then run `./flightLGo`

## Check activity
flightlog will output to the terminal it is started in.

```
Jul 21 13:42:40 mail flightLGo[22695]: Connected to OGN APRS server to track activity of <ICAO> centered on <LAT> <LOGN> with a radius of <RADIUS> km.
```
This line indicates, that flightLGo has connected to an APRS Server and registered to get motion events within a circle of \<RADIUS> kilometers at the location \<LAT> \<LONG>. When it receives events it will calcluate whether a takeoff or landing is happening. If so, it will output the event to the terminal which looks like
```
Jul 21 14:55:53 mail flightLGo[22695]: 2020-07-21 14:55:46 +0200 CEST> 3EEB98: <D-Callsign> (<CN>) ------------------- Landing
```
This event is also sent to the webserver [https://flightlog.glidernet.org](https://flightlog.glidernet.org), where the event is stored in a database to be shown in [https://flightlog.glidernet.org/?airfield=\<ICAO>](https://flightlog.glidernet.org/?airfield=<ICAO>).

Only takeoff and landing events which are happening while flightLGo is running can be reported to the website.
The website will only show your airport when at lease one takeoff or landing event has been watched and reported to it.

If not registered and started as a [systemd service](contrib/installFlightLGoasService.md) or started via a cron job, flightLGo will terminate when you logoff from the computer.

## Building
Golang installation on Ubuntu/Raspbian:
https://github.com/golang/go/wiki/Ubuntu

Install flightLGo [libfap](http://www.pakettiradio.net/libfap/) dependency:
```
sudo apt-get install libfap-dev libfap6
```

Normal build using libfap dynamicaly linked:
```
go get
go build
```
