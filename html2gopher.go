package main

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Conversion struct {
	Website  string
	Location string
	Type     int
}

var Conversions []Conversion

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadFile(dir + "/config.json")
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(body, &Conversions)
	if err != nil {
		log.Println(err)
	}

	for _, conver := range Conversions {
		err := exec.Command("wget", "-O", "temp.txt", conver.Website).Run()
		if err != nil {
			log.Fatal(err)
		}
		switch conver.Type {
		case 2:
			file, err := os.Open("temp.txt")
			if err != nil {
				log.Println(err)
			}
			outfile, err := os.Create(conver.Location)
			if err != nil {
				log.Println(err)
			}

			w := bufio.NewWriter(outfile)
			r := bufio.NewReader(file)
			t := time.Now()
			now := t.String() + "\n" + "\n"
			w.WriteString("Last updated: " + now)
			for {
				s, err := r.ReadString('\n')

				if err == io.EOF {
					break
				}
				if err != nil {
					log.Println(err)
				}
				if !strings.Contains(s, "<s>") {
					if !strings.Contains(s, "*") {
						w.WriteString(s)
					} else {
						work := strings.Replace(s, "* ", "", 2)
						entries := strings.SplitN(work, " ", 2)
						title := ""
						if len(entries) == 1 {
							title = entries[0]
						} else {
							title = entries[1]
						}
						title = strings.Replace(title, "\n", "", -1)
						url := strings.Replace(entries[0], "\n", "", -1)
						output := "h" + title + "\t /URL:" + url + "/\n"
						w.WriteString(output)
						//						fmt.Println(entries[1])
					}
				}
			}
			w.Flush()
			os.Remove("temp.txt")
			log.Println("Gopher update done")
		}
	}
}
