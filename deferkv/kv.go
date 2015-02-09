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

		deferclient.PostIt(b, kvUrl)
	}()
}
