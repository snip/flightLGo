#!/bin/bash
# Author: Peter Simon aka Ventusfahrer
# This program comes with ABSOLUTELY NO WARRANTY.
# This is free software, and you are welcome to redistribute it under GPL v3.0 conditions.
#
# This script 
# - downloads the prebuild executables from https://github.com/snip/flightLGo
# - creates a user fligthLGo with homedirectory /home/fligthLGo where flightLGo is installed to
# - configures flightLGo executable to run as a linux systemd service using this user
# - refer to the last lines of the script to find out which commands to run to 
#   remove it from your computer again
# 
# the script will run for x86_64 and armv7l architectures, maybe we have additional ones (raspis)...


# all commands below have to run using with root priviledges 
# in order to run it, download it to a location on the computer
# where the script should run, enable execution on it and invoke it as root user
# chmod +x flightLGo.Service.sh
# sudo ./flightLGo.Service.sh

#variables
user=flightLGo
group=flightLGo
userpath=/home/$user
flightLGoDirName=flightLGo
flightLGoDirPath=$userpath/$flightLGoDirName


function help {
	echo "usage:   installFlightLGOasServer.sh    :downloads and installs flightLGo as systemd service"
	echo "         installFlightLGOasServer.sh --uninstall :removes flightLGo and all components from the computer"   
        echo "         installFlightLGOasServer.sh --help      :shows this help"
	exit 1
}

#----------------------------------------------------
# uninstall will remove the flightLGo software, 
# service definition, syslog-configuration from the system

function uninstall {
	
	echo "Stopping flightLGo.service"
	systemctl stop flightLGo.service
	echo "Removing flightLGo as service"
	systemctl disable flightLGo.service
	echo "Removing user flightLGo and his home directory"
	userdel -r flightLGo
	echo "Removing the flighLGo service configuration file"
	rm -v /lib/systemd/system/flightLGo.service
	echo "Removing flightLGo syslog configuration file"
	rm -v /etc/rsyslog.d/20-flightLGo.conf
	echo "Restarting rsyslog daemon"
	systemctl restart rsyslog
	echo "Removing flighLGo log files"
	rm -v -Rf /var/log/flightLGo
	exit 0
}

if [[ $1 = "--help" ]];
then
	help
elif [[ $1 = "--uninstall" ]];
then
	uninstall
elif [[ $1 != "" ]];
then
	help
fi

# check supported platforms
platform=`uname -m`

if [[ $platform != "x86_64" ]] && [[ $platform != "armv7l" ]];
then
	echo "your platform \"$platform\" is not (yet) supported by the script"
	echo "please open an issue at https://github.com/snip/flightLGo"
	echo "perhabs we can help you"
	exit 1
fi 

# check whether libfap6 exists in the apt repository
apt-cache pkgnames | grep -q "\<libfap6\>"

if [[ $? != 0 ]];
then
	echo "libfap6 could not be found in your repository"
	echo "\"apt install libfap6\" - will not work"
	echo "installation failed"
	exit 1
fi

# install needed libraries
apt-get install libfap6

# add user flightLGo and to the system
# create his home directories
useradd --groups adm --shell /usr/sbin/nologin flightLGo
mkdir $userpath
mkdir $flightLGoDirPath
chown -R $user:$group $userpath

# determin the latest version of snip/flightLGo for the download
latest=`curl --silent "https://github.com/snip/flightLGo/releases/latest" | sed 's#.*tag/\(.*\)\".*#\1#'`
latestpath="https://github.com/snip/flightLGo/releases/download/$latest"

# download the executable to the destination directory
platform=`uname -m`
if [[ $platform = "x86_64" ]];
then
     echo "downloading X64 version"
     wget $latestpath/flightLGo.X64 --output-document $flightLGoDirPath/flightLGo
else
   if [[ $platform = "armv7l" ]];
   then
     echo "downloading "arm/raspi" version"
     wget $latestpath/flightLGo.arm --output-document $flightLGoDirPath/flightLGo
   else
     echo "your platform $platform is not handled"
     echo ""
     echo "this might happen on raspberry pi which I was not able to test or which are not supported"
	 echo "please modify the script to your needs"
     exit 1
   fi
fi

# fix file ownership and permissions
chown $user:$group $flightLGoDirPath/flightLGo
chmod +x $user:$group $flightLGoDirPath/flightLGo

#---------------------------------------
# create systemd service definition file

echo "
[Unit]
Description=FlightLGo Data Collector
ConditionPathExists=$flightLGoDirPath/flightLGo
After=network.target
StartLimitIntervalSec=60
 
[Service]
Type=simple
User=flightLGo
Group=flightLGo
LimitNOFILE=1024

Restart=always
RestartSec=10

WorkingDirectory=$flightLGoDirPath
ExecStart=$flightLGoDirPath/flightLGo

# make sure log directory exists and owned by root
# still no shure whether we need this directory or not

PermissionsStartOnly=true
ExecStartPre=/bin/mkdir -p /var/log/flightLGo
ExecStartPre=/bin/chown root:adm /var/log/flightLGo
ExecStartPre=/bin/chmod 755 /var/log/flightLGo
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=flightLGo
 
[Install]
WantedBy=multi-user.target
" > /lib/systemd/system/flightLGo.service

#---------------------------------------------------------
# download the configuration file
wget https://raw.githubusercontent.com/snip/flightLGo/master/sample.env --output-document $flightLGoDirPath/.env

#---------------------------------------------------------
# specify your airfields parameters
echo "------------------------------------------------------------------------------------------------

Next this script will invoke an editor to edit the flightLGo configuration file. 
You have to specify the right values for LAT, LONG and ICAO.
IMPORTANT: no blanks behind the \"=\" signs, LAT and LONG with decimal-point not comma!

helpfull hereby: https://www.latlong.net/degrees-minutes-seconds-to-decimal-degrees

If you need to adjust the values after running this script, you edit the file
     /home/flightLGo/flightLGo/.env
After editing, you need to sync the daemon-cache before restarting the service:
  systemctl daemon-reload
  systemctl stop flightLGo.service
  systemctl start fligthLGo.service
------------------------------------------------------------------------------------------------"
read -p "hit return to edit the configuration file" $var
editor $flightLGoDirPath/.env

#---------------------------------------------------------
# create rsyslog-configuration to write output to 
# /var/log/flightLGo/flightLGo.log
echo "if \$syslogtag contains 'flightLGo' then /var/log/flightLGo/flightLGo.log" > /etc/rsyslog.d/20-flightLGo.conf
echo "if \$syslogtag contains 'flightLGo' then stop" >> /etc/rsyslog.d/20-flightLGo.conf

#---------------------------------------------------------
# restart rsyslogd to activate the flightLGo.log configuration

systemctl restart rsyslog

#---------------------------------------------------------
# create the service, that will be started during booting
# systemd will also take control, when the service is crashing

systemctl enable flightLGo.service

#---------------------------------------------
# now you can start the flightLGo service:
systemctl start flightLGo.service

#---------------------------------------------
# and monitor it:
sleep 5
systemctl status flightLGo.service < /dev/null






