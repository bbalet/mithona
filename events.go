package main

/*  mithona allows to simply share information on a server with its users
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
	"io"
	"io/ioutil"
	"net/http"
)

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	var p = pageContent(r, "Events on this machine", nil)
	err := templates["events"].ExecuteTemplate(w, "base", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func eventsMessagesHandler(w http.ResponseWriter, r *http.Request) {

	var messagesFile = XP + "/data/messages.json"
	file, err := ioutil.ReadFile(messagesFile)
	if err != nil {
		logFatal("(eventsMessagesHandler) JSON file : ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(file))
}

func eventsDatesHandler(w http.ResponseWriter, r *http.Request) {
	//GET params of fullCalendar are UNIX timestamps
	//startParam
	//endParam
	var datesFile = XP + "/data/events.json"
	file, err := ioutil.ReadFile(datesFile)
	if err != nil {
		logFatal("(eventsDatesHandler) JSON file : ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(file))
}
