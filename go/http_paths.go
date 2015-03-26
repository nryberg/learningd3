package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type String string

type Record struct {
	Values []string
}

type Table_Data struct {
	Headers []string
	Records []Record
}

type Page struct {
	Title     string
	Body      string
	Tabledata Table_Data
}

/*
func (s *Table_Data) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request) {
	t, _ := template.ParseFiles("data.html")

	t.Execute(w)
}
*/

func viewdata(w http.ResponseWriter, r *http.Request) {
	td := new(Table_Data)

	loadData(td)

	p := Page{Title: "Fred", Body: "Body Text", Tabledata: *td}
	t, err := template.ParseFiles("./data.html")
	if err != nil {
		fmt.Println(err)
	} else {
		t.Execute(w, p)
	}

}

func loadData(td *Table_Data) {
	path := "./banks.csv"
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening data file: %v\n", err)
		os.Exit(1)
	}

	s := bufio.NewScanner(f)
	counter := 0 // THERE HAS TO BE A BETTER WAY  TODO: Implement grab first line better

	var records []Record
	var record Record

	for s.Scan() {
		if counter == 0 {
			td.Headers = strings.Split(s.Text(), ",")
		} else {

			record.Values = strings.Split(s.Text(), ",")
			records = append(records, record)
		}
		counter += 1
	}

	td.Records = records
}

func main() {
	// your http.Handle calls here

	fmt.Print("Starting server\n")
	http.HandleFunc("/data", viewdata)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	log.Fatal(http.ListenAndServe("localhost:4000", nil))
}
