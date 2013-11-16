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
	"net/http"
	"os"
)

// homeHandler is an HTTP Handler that displays and runs the logic
// for the homepage.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	hostNameInfo, err := os.Hostname()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logFatal("get Hostname from Kernel: %s", err)
	}
	var p = Page{Title: "Utilitaires pour l'injecteur", Hostname: hostNameInfo,
		Language: "fr", Menu: template.HTML(buildMenu(r))}

	err = templates["home"].ExecuteTemplate(w, "base", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
