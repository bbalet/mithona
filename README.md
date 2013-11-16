# What is mithona
*mithona* is a toy project that helps me to explore Go (Google golang) framework and some librairies written in pure go.
It can be run from console or as a service (on *Windows*, *Linux* or *MacOS*).
I decided to share because it appears to be useful for some people and because I want to receive critics so as to improve my understanding of this new language.

*mithona* is currently under development

## Features
* Simply share information from a machine (a file, a folder...).
* Test if another computer can be accessed from the machine where *goServerView* is installed.
* Allow guest users to build their own reports with the available system metrics.
* Share messages and events calendar related to the machine.
* Mapping between system metrics from a machine and Javascript.
* Dynamical configuration.

## Topics covered
* crypto package (hash, partial RSA encoding). (*see login.go*)
* Gorilla Toolkit (Sessions, mux). (*see main.go, login.go*)
* How to run a golang program as a service? (*see main.go*)
* Golang templates and nested templates.  (*see main.go and tmpl folder*)
* Basic network interactions (connect, lookup). (*see network.go*)
* Filesystem notification on MODIFY. (*see config.go*)
* i18n deported on the browser (with i18next, *see tmpl folder*)

## Tasks

* [] Editable reports
* [] Editable calendar (updated by admin)
* [] Javascript file browser

# Getting started
Go to the root folder and run the following command:
    go build .
and then launch the program from the console with:
    goServerView.exe run
Open your browser at the URL http://localhost:7777/home
Use admin/admin to login to the admin interface

# Advanced usage

## Run as a service
Install the service with the following command:
    goServerView.exe install
Start the service from the OS' service control panel or with:
    goServerView.exe start
Stop the service from the OS' service control panel or with:
    goServerView.exe stop
	
## Configuration
Edit *conf/config.json*. Some changes are dynamically applied by the application :
* **port** port number (e.g. "7777") not changed dynamically, require to restart.
* **secured** is the application globally accessed through TLS (e.g. false)  not  changed dynamically, require to restart. Modify the certificate and the public/private keys of the *conf* folder
* **compression** enable GZIP/DEFLATE compression (e.g. true) for file transfert.
* **Language** default language (e.g. "en").
* **password** admin password hashed with SHA512 and encoded in base64 (e.g. "x61Ey612Kl2gpFL56FT9weDnpSo4AV8j8+qx2AuTHdRyY036xxzTTrw10Wq3+4qQyB+XURPWx1ONxp3Y3pB37A==").
	
Menu/feature setup:
* **menu.files** Activate the individual file sharing menu entry (e.g. true)
* **menu.folders** Activate the folder sharing menu entry (e.g. true)
* **menu.queries** Activate the queries list menu entry (e.g. true)
* **menu.network** Activate the network menu entry (e.g. true)
* **menu.connect** Enable/Disable the remote connection feature (e.g. false)
* **menu.reports** Activate the reporting menu entry (e.g. true)
* **menu.events** Activate the calendar menu entry (e.g. true)

Files is an array where you list all the individual files you want to share:
* **files[].name** Menu entry (e.g. "hosts")
* **files[].description** Description of the file (e.g. "Windows Hosts defintion")
* **files[].path** Full path to the file to be shared (e.g. "C:\\Windows\\System32\\drivers\\etc\\hosts")
* **files[].view** Who can view the file content (e.g. "guests"). Not yet implemented.
* **files[].edit** Who can view the file content (e.g. "admin"). Not yet implemented.

Folders is an array where you list all the individual files you want to share:
* **folders[].name** (e.g. "D drive")
* **folders[].description** (e.g. "My D drive")
* **folders[].path** (e.g. "D:\\")

Queries is an array of remote queries (command line executed on the computer where *goServerView* is installed):
* **queries[].name** (e.g. "disk")
* **queries[].type** (e.g. "graph")
* **queries[].description** (e.g. "Disk usage")
* **queries[].lifetime** (e.g. 1000)
* **queries[].cmdLine** (e.g. "C:\\Windows\\System32\\wbem\\wmic.exe path Win32_PerfFormattedData_PerfDisk_PhysicalDisk get Name, PercentDiskTime, AvgDiskQueueLength, DiskReadBytesPerSec, DiskWriteBytesPerSec /format:csv")

Columns is an array of Queries that specify the returned columns
* **queries[].columns[].Name** (e.g. "Node")
* **queries[].columns[].Type** (e.g. "label")
* **queries[].columns[].Ignore** (e.g. true)
