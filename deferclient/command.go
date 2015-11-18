package deferclient

// CommandType defines command list supported by the clinet
type CommandType byte

const (
	// CommandTypeTrace is a command for generating traces
	CommandTypeTrace CommandType = iota + 1
	// CommandTypeBlockprofile is a command for generating block profile
	CommandTypeBlockprofile
	// CommandTypeCPUProfile is a command for generating cpu profile
	CommandTypeCPUProfile
	// CommandTypeMemprofile is a command for generating memory profile
	CommandTypeMemprofile
)

// Command contains information about this client's command, that has to be executed
type Command struct {
	Id        int         `json:"Id"`
	Type      CommandType `json:"Type"`
	Requested bool        `json:"Requested"`
	Executed  bool        `json:"Executed"`
}

// NewCommand instantitates and returns a new command
// it is meant to be called once before the executing client's command
func NewCommand(id int, commandtype CommandType, requested bool, executed bool) *Command {
	c := &Command{
		Id:        id,
		Type:      commandtype,
		Requested: requested,
		Executed:  executed,
	}

	return c
}
