package main

import (
	"fmt"

	"github.com/mayur-tolexo/distance"
)

func main() {
	if dis, err := distance.GetPinDistanct("201301", []string{"743135"}); err == nil {
		fmt.Println(dis)
	}
}
