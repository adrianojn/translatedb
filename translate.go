// Copyright (C) 2014 Adriano Soares
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	dbName   = flag.String("db", "", "database to translate")
	jsonFile = flag.String("using", "wikia.json", "wikia database")
)

var lang string

var data map[string]struct {
	Title     string
	Revisions []struct {
		Text string `json:"*"`
	}
}

func main() {
	flag.StringVar(&lang, "lang", "pt", "output language")
	flag.Parse()

	// load

	f, err := os.Open(*jsonFile)
	check(err)

	json.NewDecoder(f).Decode(&data)
	check(err)

	// parse
	// TODO: English is a very special case

	namePrefix := "|" + lang + "_name = "
	lorePrefix := "|" + lang + "_lore = "

	for _, c := range data {
		if c.Revisions == nil {
			continue
		}
		id := extract(c.Revisions[0].Text, "|number = ")
		name := extract(c.Revisions[0].Text, namePrefix)
		lore := extract(c.Revisions[0].Text, lorePrefix)
		dbUpdate(id, name, lore)
	}
	fmt.Println("name:", namePrefix, "lore:", lorePrefix)
}

func extract(source, prefix string) string {
	for _, s := range strings.Split(source, "\n") {
		if strings.HasPrefix(s, prefix) {
			return strings.TrimPrefix(s, prefix)
		}
	}
	return ""
}

func dbUpdate(id, name, lore string) {
	// stub
	fmt.Println("\tid:", id)
	fmt.Println("\tname:", name)
	fmt.Println("\tlore:", lore)
}

func isEnglish() bool {
	return lang == ""
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
