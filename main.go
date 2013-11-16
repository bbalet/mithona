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
	"bitbucket.org/kardianos/osext"
	"bitbucket.org/kardianos/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var logSrv service.Logger
var name = "goServerView"
var displayName = "Go Server View"
var desc = "Mini webserver allowing to share informations (files, folders, etc...)"
var isService bool = true
var store *sessions.FilesystemStore

// main runs the program as a service or as a command line tool.
// Several verbs allows you to install, start, stop or remove the service.
// "run" verb allows you to run the program as a command line tool.
// e.g. "goServerView install" installs the service
// e.g. "goServerView run" starts the program from the console (blocking)
func main() {
	s, err := service.NewService(name, displayName, desc)
	if err != nil {
		fmt.Printf("%s unable to start: %s", displayName, err)
		return
	}
	logSrv = s

	if len(os.Args) > 1 {
		var err error
		verb := os.Args[1]
		switch verb {
		case "install":
			err = s.Install()
			if err != nil {
				fmt.Printf("Failed to install: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" installed.\n", displayName)
		case "remove":
			err = s.Remove()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", displayName)
		case "run":
			isService = false
			doWork()
		case "start":
			err = s.Start()
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" started.\n", displayName)
		case "stop":
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" stopped.\n", displayName)
		}
		return
	}
	err = s.Run(func() error {
		// start
		go doWork()
		return nil
	}, func() error {
		// stop
		stopWork()
		return nil
	})
	if err != nil {
		s.Error(err.Error())
	}
}

// doWork is the actual main entry of the application whereas main set up
// the context (console program or service)
func doWork() {
	//Load configuration
	logInfo("Load configuration")
	XP, _ = osext.ExecutableFolder()
	CONFIGURATION_FILE = XP + "/conf/config.json"
	PRIVATE_KEY_FILE = XP + "/conf/private.pem"
	PUBLIC_KEY_FILE = XP + "/conf/public.pem"
	CERTIFICATE_FILE = XP + "/conf/cacert.pem"
	DATA_FOLDER = XP + "/data/"
	SESSION_FOLDER = XP + "/sessions/"
	file, err := ioutil.ReadFile(CONFIGURATION_FILE)
	if err != nil {
		logFatal("(main) Configuration file : ", err)
	}
	json.Unmarshal(file, &config)
	fsConfigWatcher()

	store = sessions.NewFilesystemStore(SESSION_FOLDER, []byte(config.Password))
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	}

	//Define the map of nested templates and subscribe to changes
	logInfo("Define the map of nested templates")
	createTemplateMap()
	fsTemplatesWatcher()

	//Start the embedded web server
	logInfo("Start the embedded web server")

	//Define the web application routes
	r := mux.NewRouter()
	r.HandleFunc("/home", makeHandler(homeHandler))
	r.HandleFunc("/share", makeHandler(fileShareHandler))
	r.HandleFunc("/network", makeHandler(networkHandler))
	r.HandleFunc("/lookup", makeHandler(lookupHandler))
	r.HandleFunc("/connect", makeHandler(connectHandler))
	r.HandleFunc("/browser", makeHandler(browserHandler))
	r.HandleFunc("/stat", makeHandler(statHandler))
	r.HandleFunc("/report/{name}", makeHandler(reportHandler))
	r.HandleFunc("/report/edit/{name}", makeHandler(reportEditHandler))
	r.HandleFunc("/report/delete/{name}", makeHandler(reportDeleteHandler))
	r.HandleFunc("/report/save", makeHandler(reportSaveHandler))
	r.HandleFunc("/events", makeHandler(eventsHandler))
	r.HandleFunc("/events/dates", makeHandler(eventsDatesHandler))
	r.HandleFunc("/events/messages", makeHandler(eventsMessagesHandler))
	r.HandleFunc("/logout", makeHandler(logoutHandler))
	r.HandleFunc("/login", makeHandler(loginFormHandler))
	r.HandleFunc("/loginAction", makeHandler(loginActionHandler))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(XP + "/static/")))
	http.Handle("/", r)

	logInfo("Listening on %s\n", config.Port)
	if config.Secured {
		err = http.ListenAndServeTLS(config.Port, CERTIFICATE_FILE, PRIVATE_KEY_FILE, nil)
		checkError(err)
	} else {
		err = http.ListenAndServe(config.Port, nil)
		checkError(err)
	}
}

// stopWork stops the service
func stopWork() {
	logInfo("I'm Stopping!")
}

//------------------------------------------------------------------------------
// Utility functions
//------------------------------------------------------------------------------

// logInfo reports a message in the console or the system log,
// depending on the execution context (console or service)
func logInfo(logMessage string, a ...interface{}) {
	if isService {
		logSrv.Info(logMessage, a...)
	} else {
		log.Printf(logMessage, a...)
	}
}

// logInfo reports an error in the console or the system log,
// depending on the execution context (console or service)
func logFatal(logMessage string, a ...interface{}) {
	if isService {
		logSrv.Error(logMessage, a...)
	} else {
		log.Fatalf(logMessage, a...)
	}
}

// checkError checks and reports any fatal error (errors occuring before the HTTP server is listening)
func checkError(err error) {
	if err != nil {
		logFatal("%v", err)
	}
}

// checkHttpError checks and reports any fatal error. Display an HTTP-500 page
func checkHttpError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logFatal("%v", err)
	}
}
