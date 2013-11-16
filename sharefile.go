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
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

func isFileAvailable(name string) (file FileType, isFound bool) {
	for _, element := range config.Files {
		if element.Name == name {
			return element, true
		}
	}
	return FileType{}, false
}

func fileShareHandler(w http.ResponseWriter, r *http.Request) {
	//Check if the file is shared or not
	queryFile := r.FormValue("file")
	file, found := isFileAvailable(queryFile)
	if found == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	//Get the content
	hostsData, err := ioutil.ReadFile(file.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logFatal("read hosts file: %s", err)
	}
	//Render the page
	hostNameInfo, err := os.Hostname()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logFatal("get Hostname from Kernel: %s", err)
	}
	var p = Page{Title: "File content", Hostname: hostNameInfo,
		Language: "fr", Menu: template.HTML(buildMenu(r)),
		Value: hostsData}
	err = templates["sharefile"].ExecuteTemplate(w, "base", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
