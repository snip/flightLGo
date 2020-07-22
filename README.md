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
Use the installation script to install flightLGo daemon as systemd service:

```
wget https://raw.githubusercontent.com/snip/flighLGo/master/contrib/InstallFlightLGoasService.sh
sudo ./installFlightLGoService.sh
```

The script downloads the latest flightLGo executable and sample.env from https://github/snip/flightLGo, 
starts an editor to edit flighLGo configuration file (specify the Latitute, Longitude and name of your airfield).
flightLGo will be configured to run as systemd service and started. 

Or perform a manual installation:

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

## Uninstall

If you have used the installFlightLGoasService.sh script to install the daemon you can use:

```
sudo ./installFlightLGoasService.sh --uninstall
```

to remove the daemon as well as all other changes done by the script from your system. Excpetion: libfap6 remains on the system 
