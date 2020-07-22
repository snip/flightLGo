# Install flightLGo as Systemd Service
Systemd services are automatically started when your system boots, stopped when the system is shut down. Furthermore the service is monitored and restarted after it crashed.

## Disclaimer

Usage of the below script is on your own risk without ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it under GPL v3.0 conditions

## Concept

The script will 
- install additional libraries (libfap6)
- create a user which account is used to run the service
- download the executable from https://github.com/snip/flightLGo/releases
- download the sample.env configuration file
- invoke an editor to edit the configuration file
- configure the systemd service flightLGo.service
- configure rsyslog to output flightLGo.service's output to /var/log/flightLGo/flightLGo.log
- start the systemd flightLGo.service

## Usage

Open an terminal an execute following commands:
```
wget https://raw.githubusercontent.com/snip/flightLGo/master/contrib/installFlightLGoasService.sh
chmod +x ./installFlightLGoasService.sh
sudo ./installFlightLGoasService.sh
```
The script downloads the latest flightLGo executable and sample.env from https://github/snip/flightLGo, 
starts an editor to edit flighLGo configuration file (specify the Latitute, Longitude and name of your airfield).
flightLGo will be configured to run as systemd service and started. 

## Check flightLGo's Activity

flightLGo's output is written to 
```
/var/log/flightLGo/flightLGo.log
```
when started as a systemd service. Please refere to [Check Activity chapter](../README.md#check-activity) for more details.

## Uninstall

If you have used the installFlightLGoasService.sh you can also use the script to remove flightLGo and all it's components from your server:

```
sudo ./installFlightLGoasService.sh --uninstall
```

The libfap6 package will not been removed, it may be used by other software on your server 
