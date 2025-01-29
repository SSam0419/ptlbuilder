package main

import "github.com/SSam0419/ptlbuilder"

func main() {

	spec := ptlbuilder.NewProtocolSpec("protocol", (10))

	cmd := ptlbuilder.NewCommand("RegisterClient")
	cmd.AddField("address", "string").AddField("content", "string")

	cmd2 := ptlbuilder.NewCommand("SendMessage")
	cmd2.AddField("address", "string").AddField("content", "string")

	spec.AddCommand(cmd).AddCommand(cmd2)

	ptlbuilder.NewGenerator(spec).Generate()
}
