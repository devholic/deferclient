package deferkv

import (
	"encoding/json"
	"github.com/deferpanic/deferclient/deferclient"
	"log"
)

// to be DEPRECATED
const (
	// kvUrl is the kv api endpoint
	kvUrl = deferclient.ApiBase + "/kvs/create"
)

// to be DEPRECATED
var (
	// Token is your deferpanic token available in settings
	Token string

	// Environment sets an environment tag to differentiate between separate
	// environments - default is production.
	Environment = "production"

	// AppGroup sets an optional tag to differentiate between your various
	// services - default is default
	AppGroup = "default"
)

// Client is the client used for making k/v request to the
// defer panic api
type Client struct {
	// kvUrl is the kv api endpoint
	kvUrl string

	// Token is your deferpanic token available in settings
	Token string

	// Environment sets an environment tag to differentiate between separate
	// environments - default is production.
	environment string

	// AppGroup sets an optional tag to differentiate between your various
	// services - default is default
	appGroup string

	// BaseClient is the base deferpanic client that all http requests use
	BaseClient *deferclient.DeferPanicClient
}

// DeferKV is a generic k/v struct
type DeferKV struct {
	Key   string `json:"Key"`
	Value int    `json:"Value"`
}

// NewClient instantiates and returns a new client
func NewClient(token string) *Client {

	ds := &Client{
		kvUrl:       deferclient.ApiBase + "/kvs/create",
		Token:       token,
		environment: "production",
		appGroup:    "default",
	}

	ds.BaseClient = deferclient.NewDeferPanicClient(ds.Token)
	ds.BaseClient.Environment = ds.environment
	ds.BaseClient.AppGroup = ds.appGroup

	return ds
}

// Setenvironment sets the environment
// default is 'production'
func (c *Client) Setenvironment(environment string) {
	c.environment = environment
	c.BaseClient.Environment = c.environment
}

// SetappGroup sets the app group
// default is 'default'
func (c *Client) SetappGroup(appGroup string) {
	c.appGroup = appGroup
	c.BaseClient.AppGroup = c.appGroup
}

// Report takes a key and a value and ships it to deferpanic
func (c *Client) Report(key string, value int) {

	kv := DeferKV{
		Key:   key,
		Value: value,
	}

	go func() {
		b, err := json.Marshal(kv)
		if err != nil {
			log.Println(err)
		}

		c.BaseClient.Postit(b, c.kvUrl)
	}()
}

// Report takes a key and a value and ships it to deferpanic
// DEPRECATED
// please use deferkv.NewClient(token)
func Report(key string, value int) {
	log.Println("please consider using deferkv.Client(token)")

	ds := DeferKV{
		Key:   key,
		Value: value,
	}

	go func() {
		b, err := json.Marshal(ds)
		if err != nil {
			log.Println(err)
		}

		// hack
		deferclient.Token = Token
		deferclient.Environment = Environment
		deferclient.AppGroup = AppGroup

		deferclient.PostIt(b, kvUrl)
	}()
}
