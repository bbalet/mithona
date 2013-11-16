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
	"io/ioutil"
	"net/http"
)

// isFileAvailable checks in config if the file is shared or not
func isFileAvailable(name string) (file FileType, isFound bool) {
	for _, element := range config.Files {
		if element.Name == name {
			return element, true
		}
	}
	return FileType{}, false
}

// fileShareHandler is the HTTP handler for the file sharing feature
func fileShareHandler(w http.ResponseWriter, r *http.Request) {
	//If the requested file is not shared, return a 404 HTTP error
	queryFile := r.FormValue("file")
	file, found := isFileAvailable(queryFile)
	if found == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	//Get the file content
	fileContent, err := ioutil.ReadFile(file.Path)
	checkHttpError(err, w)

	//Render the page
	var p = pageContent(r, "File content", fileContent)
	err = templates["sharefile"].ExecuteTemplate(w, "base", p)
	checkHttpError(err, w)
}
