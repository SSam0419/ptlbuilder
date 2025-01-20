package main

import "github.com/SSam0419/ptlbuilder"

func main() {
	ptlbuilder.Generate(ptlbuilder.ProtocolSpec{
		Package: "protocol",
		Commands: []ptlbuilder.Command{
			{
				Name: "RegisterClient",
				Fields: []ptlbuilder.Field{
					{Name: "ClientAddr", Type: "string"},
				},
			},
			{
				Name: "ListenTopic",
				Fields: []ptlbuilder.Field{
					{Name: "Topic", Type: "string"},
					{Name: "ClientAddr", Type: "string"},
				},
			},
			{
				Name: "SendMessage",
				Fields: []ptlbuilder.Field{
					{Name: "Topic", Type: "string"},
					{Name: "ClientAddr", Type: "string"},
					{Name: "Payload", Type: "[]byte"},
				},
			},
		},
	})
}
