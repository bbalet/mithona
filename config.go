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
	"code.google.com/p/go.exp/fsnotify"
	"encoding/json"
	"io/ioutil"
	"time"
)

type Configuration struct {
	Port        string
	Secured     bool
	Compression bool
	Language    string
	Password    string
	Menu        MenuType
	Files       []FileType
	Folders     []FolderType
	Queries     []QueryType
}

type MenuType struct {
	Files   bool
	Folders bool
	Queries bool
	Network bool
	Connect bool
	Events  bool
	Reports bool
}

type QueryType struct {
	Name        string
	Type        string
	Description string
	LifeTime    uint
	LastUpdate  time.Time
	CmdLine     string
	Columns     []ColumnType
}

type ColumnType struct {
	Name   string
	Type   string
	Ignore bool
}

type ResultType struct {
	Columns []ColumnType
	Records []RecordType
}

type RecordType struct {
	Columns []string
}

type FileType struct {
	Name        string
	Description string
	Path        string
}

type FolderType struct {
	Name        string
	Description string
	Path        string
}

var config Configuration
var CONFIGURATION_FILE string
var PRIVATE_KEY_FILE string
var PUBLIC_KEY_FILE string
var CERTIFICATE_FILE string
var DATA_FOLDER string
var SESSION_FOLDER string
var XP string

const VERSION = "0.1"
const MAX_BUFFER_SIZE = 2048

// fsConfigWatcher watches if the configuration file is modified and reload the
// configuration dynamically is the file has been changed
func fsConfigWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logFatal("(fsConfigWatcher) fsnotify.NewWatcher() : ", err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				logInfo("event: %v", ev)
				file, err := ioutil.ReadFile(CONFIGURATION_FILE)
				if err != nil {
					logFatal("(main) Configuration file : ", err)
				}
				json.Unmarshal(file, &config)
			case err := <-watcher.Error:
				logInfo("error: %v", err)
			}
		}
	}()

	err = watcher.WatchFlags(CONFIGURATION_FILE, fsnotify.FSN_MODIFY)
	if err != nil {
		logFatal("(fsConfigWatcher) watcher.WatchFlags(fsnotify.FSN_MODIFY) : ", err)
	}
}
