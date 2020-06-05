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
Copy `sample.env` to `.env` in this directory then update this `.env` according to your needs.
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

Static build of libfap
```
go get
go build --ldflags '-extldflags "/lib/x86_64-linux-gnu/libfap.a"'                                                                                                  
```

Crosscompilation seems not easy due to C lib dependency.
