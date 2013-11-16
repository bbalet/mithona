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
	"code.google.com/p/go.exp/fsnotify"
	"compress/gzip"
	"html/template"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

var templates = make(map[string]*template.Template)

//Standard page
type Page struct {
	Title    string
	Hostname string
	Menu     template.HTML
	Language string
	Message  template.HTML
	Value    []byte
	Version  string
}

// createTemplateMap defines the map of nested templates
func createTemplateMap() {
	templates["home"] = template.Must(template.ParseFiles(XP+"/tmpl/home.html", XP+"/tmpl/base.html", XP+"/tmpl/menu.html", XP+"/tmpl/header.html", XP+"/tmpl/footer.html"))
	templates["connect"] = template.Must(template.ParseFiles(XP+"/tmpl/connect.html", XP+"/tmpl/base.html", XP+"/tmpl/menu.html", XP+"/tmpl/header.html", XP+"/tmpl/footer.html"))
	templates["network"] = template.Must(template.ParseFiles(XP+"/tmpl/network.html", XP+"/tmpl/base.html", XP+"/tmpl/menu.html", XP+"/tmpl/header.html", XP+"/tmpl/footer.html"))
	templates["browser"] = template.Must(template.ParseFiles(XP+"/tmpl/network.html", XP+"/tmpl/base.html", XP+"/tmpl/menu.html", XP+"/tmpl/header.html", XP+"/tmpl/footer.html"))
	templates["sharefile"] = template.Must(template.ParseFiles(XP+"/tmpl/sharefile.html", XP+"/tmpl/base.html", XP+"/tmpl/menu.html", XP+"/tmpl/header.html", XP+"/tmpl/footer.html"))
	templates["browser"] = template.Must(template.ParseFiles(XP+"/tmpl/browser.html", XP+"/tmpl/base.html", XP+"/tmpl/menu.html", XP+"/tmpl/header.html", XP+"/tmpl/footer.html"))
	templates["report"] = template.Must(template.ParseFiles(XP+"/tmpl/report.html", XP+"/tmpl/base.html", XP+"/tmpl/menu.html", XP+"/tmpl/header.html", XP+"/tmpl/footer.html"))
	templates["report-edit"] = template.Must(template.ParseFiles(XP+"/tmpl/report-edit.html", XP+"/tmpl/base.html", XP+"/tmpl/menu.html", XP+"/tmpl/header.html", XP+"/tmpl/footer.html"))
	templates["events"] = template.Must(template.ParseFiles(XP+"/tmpl/events.html", XP+"/tmpl/base.html", XP+"/tmpl/menu.html", XP+"/tmpl/header.html", XP+"/tmpl/footer.html"))
	templates["login"] = template.Must(template.ParseFiles(XP+"/tmpl/login.html", XP+"/tmpl/base.html", XP+"/tmpl/menu.html", XP+"/tmpl/header.html", XP+"/tmpl/footer.html"))
}

// fsTemplatesWatcher watches if the template files have been modified
// and reparse them again dynamically is they have been changed
func fsTemplatesWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logFatal("(fsConfigWatcher) fsnotify.NewWatcher() : ", err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				logInfo("event: %v", ev)
				createTemplateMap()
			case err := <-watcher.Error:
				logInfo("error: %v", err)
			}
		}
	}()

	err = watcher.WatchFlags(XP+"/tmpl/", fsnotify.FSN_MODIFY)
	if err != nil {
		logFatal("(fsConfigWatcher) watcher.WatchFlags(fsnotify.FSN_MODIFY) : ", err)
	}
}

// pageContent fills a standard Page struct
func pageContent(r *http.Request, t string, v []byte) Page {
	hostNameInfo, err := os.Hostname()
	if err != nil {
		hostNameInfo = "ERROR"
		logFatal("get Hostname from Kernel: %s", err)
	}
	var p = Page{Title: t, Hostname: hostNameInfo,
		Language: config.Language, Menu: template.HTML(buildMenu(r)),
		Value: v, Version: VERSION}

	session, _ := store.Get(r, "goServerView")
	if flashes := session.Flashes(); len(flashes) > 0 {
		p.Message = `<div class="alert"><a class="close" data-dismiss="alert">×</a><span>`
		for _, val := range flashes {
			p.Message += template.HTML(val.(string)) + `<br />`
		}
		p.Message += `</span></div>`
	}
	return p
}

// makeHandler serves GZIP content
func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	if "" == w.Header().Get("Content-Type") {
		// If no content type, apply sniffing algorithm to un-gzipped body.
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

// gt returns true if the arguments are of the same type (with int8 and int64 as the same type) and
// the first argument is greater than the second. This is only defined on string, intX, uintX and floatX all
// other types return false.
func gt(a1, a2 interface{}) bool {
	switch a1.(type) {
	case string:
		switch a2.(type) {
		case string:
			return reflect.ValueOf(a1).String() > reflect.ValueOf(a2).String()
		}
	case int, int8, int16, int32, int64:
		switch a2.(type) {
		case int, int8, int16, int32, int64:
			return reflect.ValueOf(a1).Int() > reflect.ValueOf(a2).Int()
		}
	case uint, uint8, uint16, uint32, uint64:
		switch a2.(type) {
		case uint, uint8, uint16, uint32, uint64:
			return reflect.ValueOf(a1).Uint() > reflect.ValueOf(a2).Uint()
		}
	case float32, float64:
		switch a2.(type) {
		case float32, float64:
			return reflect.ValueOf(a1).Float() > reflect.ValueOf(a2).Float()
		}
	}
	return false
}

// eq reports whether the first argument is equal to
// any of the remaining arguments.
func eq(args ...interface{}) bool {
	if len(args) == 0 {
		return false
	}
	x := args[0]
	switch x := x.(type) {
	case string, int, int64, byte, float32, float64:
		for _, y := range args[1:] {
			if x == y {
				return true
			}
		}
		return false
	}

	for _, y := range args[1:] {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}

// parseRange parses a Range header field
func parseRange(data string) int64 {
	stop := (int64)(0)
	part := 0
	for i := 0; i < len(data) && part < 2; i = i + 1 {
		if part == 0 { // part = 0 <=> equal isn't met.
			if data[i] == '=' {
				part = 1
			}

			continue
		}

		if part == 1 { // part = 1 <=> we've met the equal, parse beginning
			if data[i] == ',' || data[i] == '-' {
				part = 2 // part = 2 <=> OK DUDE.
			} else {
				if 48 <= data[i] && data[i] <= 57 { // If it's a digit ...
					// ... convert the char to integer and add it!
					stop = (stop * 10) + (((int64)(data[i])) - 48)
				} else {
					part = 2 // Parsing error! No error needed : 0 = from start.
				}
			}
		}
	}

	return stop
}

// parseCSV parses comma-separated values
func parseCSV(data string) []string {
	splitted := strings.Split(data, ",")

	data_tmp := make([]string, len(splitted))

	for i, val := range splitted {
		data_tmp[i] = strings.TrimSpace(val)
	}

	return data_tmp
}
