package main

/*  goServerView allows to simply share information on a server with its users
    Copyright (C) 2013 Benjamin BALET

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.*/

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
)

func networkHandler(w http.ResponseWriter, r *http.Request) {
	hostNameInfo, err := os.Hostname()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logFatal("get Hostname from Kernel: %s", err)
	}
	var p = Page{Title: "Network utils", Hostname: hostNameInfo,
		Language: "fr", Menu: template.HTML(buildMenu(r))}
	err = templates["network"].ExecuteTemplate(w, "base", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func lookupHandler(w http.ResponseWriter, r *http.Request) {
	hostName := r.FormValue("hostName")
	addrs, err := net.LookupHost(hostName)
	fmt.Fprintf(w, "%s", addrs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	if config.Menu.Connect {
		hostName := r.FormValue("hostName")
		portNumber := r.FormValue("portNumber")
		protocol := r.FormValue("protocol")
		sendMode := r.FormValue("sendMode")
		dataToSend := r.FormValue("dataToSend")

		conn, err := net.Dial(protocol, hostName+":"+portNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if sendMode == "SendBase64" {
			data64, err := base64.StdEncoding.DecodeString(dataToSend)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			fmt.Fprintf(conn, "%s\r\n\r\n", data64)
			status, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				fmt.Fprintf(w, "OK. Return : %s", status)
			}
		}
		if sendMode == "SendRaw" || sendMode == "SendBase64" {
			fmt.Fprintf(conn, "%s\r\n\r\n", dataToSend)
			status, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				fmt.Fprintf(w, "OK. Return : %s", status)
			}
		}
		if sendMode == "DoNotSend" {
			fmt.Fprintf(w, "OK: no data to be sent")
		}
		err = conn.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Fprintf(w, "DISABLED: remote connection is disabled")
	}
}
