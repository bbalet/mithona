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
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type PageBrowser struct {
	Title    string
	Hostname string
	Menu     template.HTML
	Language string
	Files    []MyFileInfo
	Folder   string
	Object   string
}

type MyFileInfo struct {
	IsDir   bool
	Name    string
	Size    int64
	ModTime time.Time
	Mode    string
	Path    string
}

var p PageBrowser

func isFolderAvailable(name string) (folder FolderType, isFound bool) {
	for _, element := range config.Folders {
		if element.Name == name {
			return element, true
		}
	}
	return FolderType{}, false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func browserHandler(w http.ResponseWriter, r *http.Request) {

	folderName := r.FormValue("folder")
	objectName := r.FormValue("object")
	var currentPath string

	folder, found := isFolderAvailable(folderName)
	if found == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	hostNameInfo, err := os.Hostname()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logFatal("get Hostname from Kernel: %s", err)
	}
	p = PageBrowser{Title: "Navigateur de fichier", Hostname: hostNameInfo,
		Language: "fr", Menu: template.HTML(buildMenu(r))}

	if objectName != "" {
		if strings.Contains(objectName, folder.Path) {
			currentPath = objectName
			parentPath := currentPath[:strings.LastIndex(currentPath, string(os.PathSeparator))]

			//Add pointer to parrent directory
			if parentPath != "" {
				var aFile = MyFileInfo{IsDir: true,
					Name: "..",
					Path: parentPath + string(os.PathSeparator)}
				p.Files = append(p.Files, aFile)
			}
		} else {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
	} else {
		currentPath = folder.Path
	}

	//download file or browse to the next folder
	fileStat, err := os.Stat(currentPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logFatal("os.Stat(currentPath): %s", err)
	}
	if !fileStat.IsDir() {
		f, err := os.Open(currentPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// Fetching file's mimetype and giving it to the browser
		if mimetype := mime.TypeByExtension(path.Ext(currentPath)); mimetype != "" {
			w.Header().Set("Content-Type", mimetype)
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}
		// Add Content-Length
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileStat.Size()))
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileStat.Name()))
		//TODO : Content-Disposition: attachment; filename=FILENAME.EXT

		// Manage Content-Range (TODO: Manage end byte and multiple Content-Range)
		if r.Header.Get("Range") != "" {
			start_byte := parseRange(r.Header.Get("Range"))

			if start_byte < fileStat.Size() {
				f.Seek(start_byte, 0)
			} else {
				start_byte = 0
			}

			w.Header().Set("Content-Range",
				fmt.Sprintf("bytes %d-%d/%d", start_byte, fileStat.Size()-1, fileStat.Size()))
		}

		// Manage gzip/zlib compression
		output_writer := w.(io.Writer)

		if config.Compression == true && r.Header.Get("Accept-Encoding") != "" {
			encodings := parseCSV(r.Header.Get("Accept-Encoding"))
			for _, val := range encodings {
				if val == "gzip" {
					w.Header().Set("Content-Encoding", "gzip")
					output_writer, _ = gzip.NewWriterLevel(w, gzip.BestSpeed)
					break
				} else if val == "deflate" {
					w.Header().Set("Content-Encoding", "deflate")
					output_writer, _ = zlib.NewWriterLevel(w, zlib.BestSpeed)
					break
				}
			}
		}

		buf := make([]byte, min(MAX_BUFFER_SIZE, fileStat.Size()))
		n := 0
		for err == nil {
			n, err = f.Read(buf)
			output_writer.Write(buf[0:n])
		}
		f.Close()

	} else {
		osFiles, err := ioutil.ReadDir(currentPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logFatal("ioutil.ReadDir(currentPath): %s", err)
		}

		//Copy file information into a custom struct with some transformation
		for _, element := range osFiles {
			var aFile = MyFileInfo{IsDir: element.IsDir(),
				Name:    element.Name(),
				Size:    element.Size(),
				ModTime: element.ModTime(),
				Mode:    element.Mode().String(),
				Path:    filepath.Join(currentPath, element.Name())}
			p.Files = append(p.Files, aFile)
		}

		p.Folder = folderName
		p.Object = objectName

		err = templates["browser"].ExecuteTemplate(w, "base", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
