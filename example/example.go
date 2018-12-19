package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mayur-tolexo/distance"
)

func main() {
	source := flag.String("s", "201301", "source")
	dest := flag.String("d", "110042", "destination , seperated")
	flag.Parse()
	if dis, err := distance.GetPinDistanct(*source, strings.Split(*dest, ",")); err == nil {
		fmt.Println(dis)
	}
}
