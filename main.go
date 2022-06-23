package main

import (
	"github.com/nodauf/web3Toolbox/cmd"
	"log"
)

func init() {
	log.SetFlags(log.Llongfile)
}
func main() {
	cmd.Execute()
}
