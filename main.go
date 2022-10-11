package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/ivch/twoktest/eval"
)

func main() {
	flag.Parse()
	files := flag.Args()

	for i := range files {
		data, err := ioutil.ReadFile(files[i])
		if err != nil {
			log.Fatal(err)
		}

		e := eval.New(os.Stdout)
		if err := e.Run(data); err != nil {
			log.Println(err)
		}
	}
}
