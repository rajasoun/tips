// Licensed under the Creative Commons License.

package main

import (
	"log"

	"github.rajasoun/tips/cmd"

	"os"
)

func main() {
	var err error = cmd.Execute(os.Stdout)
	if err != nil {
		log.Fatalln("something went wrong use 'tips --help' ")
	}
}
