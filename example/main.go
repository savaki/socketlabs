package main

import (
	"os"

	"github.com/savaki/socketlabs"
	"golang.org/x/net/context"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	serverId := os.Getenv("SERVER_ID")
	client := socketlabs.New(apiKey, serverId)
	client.Inject(context.Background(), []socketlabs.Message{
		{
			From: socketlabs.Recipient{
				EmailAddress: "jane.sender@example.com",
			},
			To: []socketlabs.Recipient{
				{
					EmailAddress: "%%DeliveryAddress%%",
					FriendlyName: "%%First%% %%Last%%",
				},
			},
			Subject:  "test message",
			HtmlBody: `<html><h1>hello world</h1><p>argle bargle</p></html>`,
			MergeData: socketlabs.MergeData{
				PerMessage: [][]socketlabs.KV{
					[]socketlabs.KV{
						{
							Field: "DeliveryAddress",
							Value: "joe.public@example.com",
						},
						{
							Field: "First",
							Value: "Joe",
						},
						{
							Field: "Last",
							Value: "Public",
						},
					},
				},
			},
		},
	})
}
