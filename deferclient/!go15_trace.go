// +build !go1.5

package deferclient

// MakeTrace POST a Trace html to the deferpanic website
func (c *DeferPanicClient) MakeTrace(commandId int, agent *Agent) {
	c.Lock()
	c.RunningCommands[commandId] = true
	c.Unlock()
	defer func() {
		c.Lock()
		delete(c.RunningCommands, commandId)
		c.Unlock()
	}()

	t := NewTrace([]byte{}, []byte{}, commandId)

	b, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		return
	}

	c.Postit(b, traceUrl, false)
}
