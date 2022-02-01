package main

import (
	"bytes"
	"fmt"
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Contact struct {
	FirstName string `csv:"First Name"`
	HomePhone string `csv:"Home Phone"`
	MobilePhone string `csv:"Mobile Phone"`
	DisplayName string `csv:"Display Name"`
}
func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Join(wd, "contacts.csv")
	fi, err := os.OpenFile(dir, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	values := []*Contact{}
	if err := gocsv.Unmarshal(fi, &values); err != nil {
		log.Fatal(err)
	}
	log.Printf("Total => %d", len(values))
	out := &bytes.Buffer{}
	for _, next := range values {
		log.Println(next)
		line := fmt.Sprintf("%s,%s\n", next.DisplayName, trimPhone(next.MobilePhone))
		out.Write([]byte(line))
	}
	ioutil.WriteFile("contacts_data.txt", out.Bytes(), 0644)
}

func trimPhone(p string) string {
	s := strings.ReplaceAll(p, " ", "")
	s = strings.ReplaceAll(s, "+", "")
	newValue := s
	if strings.HasPrefix(s, "0") {
		newValue = "234" + s[1:]
	}
	return strings.TrimSpace(newValue)
}
