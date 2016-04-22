package socketlabs

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const (
	InjectUrl = "https://inject.socketlabs.com/api/v1/email"
)

type KV struct {
	Field string `json:",omitempty"`
	Value string `json:",omitempty"`
}

type Recipient struct {
	EmailAddress string `json:",omitempty"`
	FriendlyName string `json:",omitempty"`
}

type MergeData struct {
	PerMessage [][]KV `json:",omitempty"`
	Global     []KV   `json:",omitempty"`
}

type Message struct {
	To          []Recipient `json:",omitempty"`
	From        Recipient   `json:",omitempty"`
	Cc          []Recipient `json:",omitempty"`
	Bcc         []Recipient `json:",omitempty"`
	ReplyTo     *Recipient  `json:",omitempty"`
	Subject     string      `json:",omitempty"`
	TextBody    string      `json:",omitempty"`
	HtmlBody    string      `json:",omitempty"`
	ApiTemplate string      `json:",omitempty"`
	MailingId   string      `json:",omitempty"`
	MessageId   string      `json:",omitempty"`
	Charset     string      `json:",omitempty"`
	MergeData   MergeData   `json:",omitempty"`
}

type AddressResult struct {
	Accepted     bool
	EmailAddress string
	ErrorCode    string
}

type MessageResult struct {
	Index         int
	AddressResult AddressResult
	ErrorCode     string
}

type Response struct {
	ErrorCode          string
	MessageResults     []MessageResult
	TransactionReceipt string
}

type Client struct {
	apiKey   string
	serverId string
	client   *http.Client
}

type Envelope struct {
	ApiKey   string
	ServerId string
	Messages []Message
}

type Option func(*Client)

func New(apiKey, serverId string, options ...Option) *Client {
	client := &Client{
		apiKey:   apiKey,
		serverId: serverId,
		client:   http.DefaultClient,
	}

	for _, opt := range options {
		opt(client)
	}

	return client
}

func (c *Client) Inject(ctx context.Context, messages []Message) (Response, error) {
	envelope := Envelope{
		ApiKey:   c.apiKey,
		ServerId: c.serverId,
		Messages: messages,
	}
	json.NewEncoder(os.Stdout).Encode(envelope)

	data, err := json.Marshal(envelope)
	if err != nil {
		return Response{}, err
	}

	resp, err := ctxhttp.Post(ctx, c.client, InjectUrl, "application/json", bytes.NewReader(data))
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	out := Response{}
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return Response{}, err
	}

	return out, err
}

func HttpClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}
