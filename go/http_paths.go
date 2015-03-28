package main

import (
	"bufio"
	"encoding/csv"
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

func csvSplit(line string) []string {
	working := strings.Split(line, ",")
	pinpoint := 0
	inquotes := false
	for index, value := range working {
		if inquotes == false {
			if len(value) > 0 && value[0] == '\u0022' { // leading quote
				working[index] = value[1:]
				inquotes = true
				pinpoint = index
			}
		} else {
			if len(value) > 0 && value[len(value)-1] == '\u0022' { // trailing quote
				working[pinpoint] += "," + value[:len(value)-1] // + working[index + 1]
				working = append(working[:index], working[index+1:]...)
				pinpoint = 0
				inquotes = false
			}
			/*   Tweaking the array mid process screws up the index badly TODO  Fix re-indexing */
			if len(value) > 0 && pinpoint > 0 && strings.ContainsRune(value, '\u0022') == false {
				fmt.Printf("working:%v\n", working[index-3:index+4])
				fmt.Printf("Value :%v\n", value)
				fmt.Printf("Index : %v\n", index)
				working[pinpoint] += "," + value
				if (index + 1) > len(working) {
					// working = append(working[:index], working[index+1:]...)
				}
			}
		}

	}
	/*
		}
	*/
	return working
}
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

	r := bufio.NewReader(f)
	// counter := 0 // THERE HAS TO BE A BETTER WAY  TODO: Implement grab first line better

	var records []Record
	var record Record

	rcsv := csv.NewReader(r)
	headers, err := rcsv.Read()
	td.Headers = headers
	values, err := rcsv.Read()
	record.Values = values
	records = append(records, record)
	fmt.Printf("Recsample %v\n", record)
	//td.Records[0].Values = record
	/*
		for s.Scan() {
			if counter == 0 {

				td.Headers = csvSplit(s.Text())
			} else {

				record.Values = csvSplit(s.Text())
				records = append(records, record)
			}
			counter += 1
		}
	*/

	td.Records = records
}

func main() {
	// your http.Handle calls here

	fmt.Print("Starting server\n")
	http.HandleFunc("/data", viewdata)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	log.Fatal(http.ListenAndServe("localhost:4000", nil))
}
