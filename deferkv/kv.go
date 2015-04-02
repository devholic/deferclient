package deferkv

import (
	"encoding/json"
	"github.com/deferpanic/deferclient/deferclient"
	"log"
)

const (
	// kvUrl is the kv api endpoint
	kvUrl = deferclient.ApiBase + "/kvs/create"
)

// Token is your deferpanic token available in settings
var Token string

// Environment sets an environment tag to differentiate between separate
// environments - default is production.
var Environment = "production"

// AppGroup sets an optional tag to differentiate between your various
// services - default is default
var AppGroup = "default"

// DeferKV is a generic k/v struct
type DeferKV struct {
	Key   string `json:"Key"`
	Value int    `json:"Value"`
}

// func Report("pathfinding", pathfindingTime)
func Report(key string, value int) {
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
