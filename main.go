package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Return first IP address of the host
func getIP(host string) string {
	ips, err := net.LookupHost(host)
	if err != nil {
		log.Println("Can't find host IP: ", host)
		return "NULL"
	}
	return ips[0]
}

// Process csv file line by line
func processCSV(rc io.Reader) (ch chan []string) {
	// create buffered channel
	ch = make(chan []string, 10)
	go func() {
		r := csv.NewReader(rc)
		// read csv header - column titles
		if _, err := r.Read(); err != nil {
			log.Fatal(err)
		}
		defer close(ch)
		// read the rest of rows and send to the channel
		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			ch <- rec
		}
	}()
	return
}

func main() {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write([]string{"host", "ip"}); err != nil {
		log.Fatal(err)
	}

	// todo: use cli args to specify csv file
	csvFile, _ := os.Open("data.csv")

	for rec := range processCSV(bufio.NewReader(csvFile)) {
		// do IP lookup for each host in the csv row
		ip := getIP(rec[1])
		rec = append(rec, ip)
		if err := w.Write(rec); err != nil {
			log.Fatal(err)
		}
		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(buf.String())
}
