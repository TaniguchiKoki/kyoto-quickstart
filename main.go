package main

import (
	"encoding/json"
	"log"
	"net/http"

	"git.sr.ht/~kyoto-framework/kyoto"
)

type CUUIDState struct {
	UUID string
}

func CUUID(ctx *kyoto.Context) (state CUUIDState) {
	uuid := func() string {
		resp, _ := http.Get("http://httpbin.org/uuid")
		data := map[string]string{}
		json.NewDecoder(resp.Body).Decode(&data)
		return data["uuid"]
	}
	// Handle action
	handled := kyoto.Action(ctx, "Reload", func(args ...any) {
		// We will just set a new uuid and will print a log
		// It's not makes a lot of sense now, but it's just a demonstration example
		state.UUID = uuid()
		log.Println("New uuid was issued:", state.UUID)
	})
	// Prevent further execution if action handled
	if handled {
		return
	}
	// Default loading behavior
	state.UUID = uuid()
	return state
}

type PIndexState struct {
	UUID1 *kyoto.ComponentF[CUUIDState]
	UUID2 *kyoto.ComponentF[CUUIDState]
}

func PIndex(ctx *kyoto.Context) (state PIndexState) {
	// Define rendering
	kyoto.Template(ctx, "page.index.html")
	// Attach components
	state.UUID1 = kyoto.Use(ctx, CUUID)
	state.UUID2 = kyoto.Use(ctx, CUUID)
	return state
}

func main() {
	kyoto.HandlePage("/", PIndex)
	kyoto.HandleAction(CUUID)
	kyoto.Serve(":8080")
}
