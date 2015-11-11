package deferclient

// Trace contains information about this client's trace and its producing package
type Trace struct {
	Out       []byte `json:"Out,omitempty"`
	Pkg       []byte `json:"Pkg,omitempty"`
	CommandId int    `json:"CommandId"`
	Ignored   bool   `json:"Ignored"`
}

// NewTrace instantitates and returns a new trace
// it is meant to be called once after the completing application tracing
func NewTrace(out []byte, pkg []byte, commandid int, ignored bool) *Trace {
	t := &Trace{
		Out:       out,
		Pkg:       pkg,
		CommandId: commandid,
		Ignored:   ignored,
	}

	return t
}
