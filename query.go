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
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

//TODO : buffer query with a max life time (to be parametrized)

func csvToJson(content string, query QueryType) ResultType {

	reader := csv.NewReader(strings.NewReader(content))
	//First line contains headers
	header, err := reader.Read()
	if err != nil {
		fmt.Println("Error:", err)
		return ResultType{}
	}
	result := ResultType{}

	for i, _ := range header {
		if !query.Columns[i].Ignore {
			var value ColumnType
			value.Name = query.Columns[i].Name
			value.Type = query.Columns[i].Type
			result.Columns = append(result.Columns, value)
		}
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return ResultType{}
		}
		var line RecordType
		for i, _ := range record {
			if !query.Columns[i].Ignore {
				line.Columns = append(line.Columns, strings.TrimSpace(record[i]))
			}
		}
		result.Records = append(result.Records, line)
	}
	return result
}

func isQueryAvailable(name string) (query QueryType, isFound bool) {
	for _, element := range config.Queries {
		if element.Name == name {
			return element, true
		}
	}
	return QueryType{}, false
}

//---------------------------------------------------------------------------------------------------
//Stat endpoint : check if query exists, return HTTP-404 if it doesn't exist
//---------------------------------------------------------------------------------------------------
func statHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var out []byte
	var content string

	queryName := r.FormValue("query")

	query, found := isQueryAvailable(queryName)
	if found == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	//Run the query through the command line
	switch runtime.GOOS {
	case "windows":
		out, err = exec.Command("cmd", "/c", query.CmdLine).Output()
		content = string(bytes.Trim(out, " \n\r"))
	default:
		out, err = exec.Command(query.CmdLine).Output()
		content = string(out)
	}
	if err != nil {
		logFatal("%T", err)
	}

	switch strings.ToLower(query.Type) {
	case "graph":
		w.Header().Set("Content-Type", "application/json")
		resultJson := csvToJson(content, query)
		b, err := json.Marshal(resultJson)
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Fprintf(w, "%s", b)
	case "datatable":
		w.Header().Set("Content-Type", "application/json")
		resultJson := csvToJson(content, query)
		fmt.Fprint(w, resultJson)
	case "values":
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, content)
	default:
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, content)
	}
}
