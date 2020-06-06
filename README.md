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
