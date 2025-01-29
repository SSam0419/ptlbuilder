package main

import (
	"fmt"

	"github.com/SSam0419/ptlbuilder"
)

func main() {

	spec := ptlbuilder.NewProtocolSpec("protocol", (10))

	cmd := ptlbuilder.NewCommand("RegisterClient")
	cmd.AddField("Address", "string").AddField("Content", "string")

	cmd2 := ptlbuilder.NewCommand("SendMessage")
	cmd2.AddField("Address", "string").AddField("Content", "string")

	spec.AddCommand(cmd).AddCommand(cmd2)

	err := ptlbuilder.NewGenerator(spec).Generate()
	if err != nil {
		fmt.Println("failed to generate template, ", err)
	}
}
