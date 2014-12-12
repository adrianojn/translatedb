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
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbName   = flag.String("db", "", "database to translate")
	jsonFile = flag.String("using", "wikia.json", "wikia database")
	lang     string
)

var data map[string]struct {
	Title     string
	Revisions []struct {
		Text string `json:"*"`
	}
}

var db *sql.DB

func main() {
	flag.StringVar(&lang, "lang", "en", "output language")
	flag.Parse()

	if *dbName == "" {
		fmt.Println("no database specified")
		os.Exit(1)
	}

	isEngligh := lang == "en"

	// load

	f, err := os.Open(*jsonFile)
	catch(err)
	defer f.Close()

	json.NewDecoder(f).Decode(&data)
	catch(err)

	db, err = sql.Open("sqlite3", *dbName)
	catch(err)
	defer db.Close()

	// parse

	namePrefix := "|" + lang + "_name = "
	lorePrefix := "|" + lang + "_lore = "

	for _, card := range data {
		if card.Revisions == nil {
			continue
		}
		text := card.Revisions[0].Text
		id := strings.TrimLeft(extract(text, "|number = "), "0")

		var name, lore string
		if isEngligh {
			name = card.Title
			lore = strip(extract(text, "|lore = "))
		} else {
			name = strip(extract(text, namePrefix))
			lore = strip(extract(text, lorePrefix))
		}

		dbUpdate(id, name, lore)
	}
}

func extract(source, prefix string) string {
	for _, s := range strings.Split(source, "\n") {
		if strings.HasPrefix(s, prefix) {
			return strings.TrimPrefix(s, prefix)
		}
	}
	return ""
}

const updateQuery = `UPDATE texts SET name=?, desc=? WHERE id=?;`

func dbUpdate(id, name, lore string) {
	if id == "" {
		return
	}
	if (name == "") || (lore == "") {
		fmt.Println("incomplete", id, name)
		return
	}
	fmt.Println("updating", name)

	_, err := db.Exec(updateQuery, name, lore, id)
	if err != nil {
		fmt.Println(err)
	}
}

var (
	htmlRegex = regexp.MustCompile(`<.+?>`)
	rubyRegex = regexp.MustCompile(`\{\{.+?\}\}`)
	wikiRegex = regexp.MustCompile(`\[\[.+?\]\]`)
)

func strip(src string) string {
	s := htmlRegex.ReplaceAllString(src, "\n")
	s = wikiRegex.ReplaceAllStringFunc(s, submatch)
	return rubyRegex.ReplaceAllStringFunc(s, submatchRuby)
}

func submatch(s string) string {
	i := strings.Index(s, "|")
	if i < 0 {
		return s[2 : len(s)-2]
	}
	return s[i+1 : len(s)-2]
}

func submatchRuby(s string) string {
	a := strings.Index(s, "|")
	b := strings.LastIndex(s, "|")
	return s[a+1 : b]
}

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
