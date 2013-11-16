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
	"net/url"
)

// buildMenu builds the menu (using Bootstrap class names) of the application
func buildMenu(r *http.Request) string {
	var MenuEntry string
	MenuEntry += `	<ul class="nav navbar-nav">`

	if config.Menu.Network {
		MenuEntry += `<li><a href="/network" data-i18n="Network">Network</a></li>`
	}
	if config.Menu.Events {
		MenuEntry += `<li><a href="/events" data-i18n="Events">Events</a></li>`
	}
	if config.Menu.Reports {
		MenuEntry += buildMenuReports()
	}
	if config.Menu.Files {
		MenuEntry += buildMenuFiles()
	}
	if config.Menu.Folders {
		MenuEntry += buildMenuFolders()
	}
	if config.Menu.Queries {
		MenuEntry += buildMenuQueries()
	}
	MenuEntry += `				</ul>`

	session, _ := store.Get(r, "goServerView")
	if session.Values["isConnected"] == true {
		MenuEntry += `
			<ul class="nav navbar-nav pull-right">
				<li><a href="/logout">`
		MenuEntry += session.Values["user"].(string)
		MenuEntry += `</a></li>
			</ul>`
	} else {
		MenuEntry += `
			<ul class="nav navbar-nav pull-right">
				<li><a href="/login" data-i18n="Login">Login</a></li>
			</ul>`
	}
	return MenuEntry
}

// buildMenuFiles builds the "share files" menu entry
func buildMenuFiles() string {
	if len(config.Files) > 0 {
		var MenuEntry = `
				    <li class="dropdown">
	                  <a href="#" class="dropdown-toggle" data-toggle="dropdown" data-i18n="Files">Files<b class="caret"></b></a>
	                  <ul class="dropdown-menu">`
		//List all shared files
		for _, element := range config.Files {
			MenuEntry += `<li><a href="/share?file=` + element.Name +
				`" title="` + element.Description + `">` +
				element.Name + `</a></li>`
		}
		MenuEntry += `
						</ul>
	                </li>`
		return MenuEntry
	} else {
		return ""
	}
}

// buildMenuFiles builds the "share folders" menu entry
func buildMenuFolders() string {
	if len(config.Folders) > 0 {
		var MenuEntry = `
				    <li class="dropdown">
	                  <a href="#" class="dropdown-toggle" data-toggle="dropdown" data-i18n="Folders">Folders<b class="caret"></b></a>
	                  <ul class="dropdown-menu">`
		//List all shared files
		for _, element := range config.Folders {
			MenuEntry += `<li><a href="/browser?folder=` + element.Name +
				`" title="` + element.Description + `">` +
				element.Name + `</a></li>`
		}
		MenuEntry += `
						</ul>
	                </li>`
		return MenuEntry
	} else {
		return ""
	}
}

// buildMenuFiles builds the "Queries" menu entry
func buildMenuQueries() string {
	if len(config.Queries) > 0 {
		var MenuEntry = `
				    <li class="dropdown">
	                  <a href="#" class="dropdown-toggle" data-toggle="dropdown" data-i18n="Queries">Queries<b class="caret"></b></a>
	                  <ul class="dropdown-menu">`
		//List all shared files
		for _, element := range config.Queries {
			MenuEntry += `<li><a href="/stat?query=` + element.Name +
				`" title="` + element.Description + `">` +
				element.Name + `</a></li>`
		}
		MenuEntry += `
						</ul>
	                </li>`
		return MenuEntry
	} else {
		return ""
	}
}

// buildMenuReports builds the "Reports" menu entry
func buildMenuReports() string {
	if len(config.Queries) > 0 {
		var MenuEntry = `
				    <li class="dropdown">
	                  <a href="#" class="dropdown-toggle" data-toggle="dropdown" data-i18n="Reports">Reports<b class="caret"></b></a>
	                  <ul class="dropdown-menu">`
		//List all available reports
		files, _ := ioutil.ReadDir(XP + "/reports/")
		for _, f := range files {
			MenuEntry += `<li><a href="/report/` +
				url.QueryEscape(f.Name()) + `">` + f.Name() + `</a></li>`
		}
		MenuEntry += `
						</ul>
	                </li>`
		return MenuEntry
	} else {
		return ""
	}
}
