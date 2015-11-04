package deferclient

// Trace contains information about this client's trace and its producing package
type Trace struct {
	Out       []byte `json:"Out"`
	Pkg       []byte `json:"Pkg"`
	CommandId int    `json:"CommandId"`
}

// NewTrace instantitates and returns a new trace
// it is meant to be called once after the completing application tracing
func NewTrace(out []byte, pkg []byte, commandid int) *Trace {
	t := &Trace{
		Out:       out,
		Pkg:       pkg,
		CommandId: commandid,
	}

	return t
}
