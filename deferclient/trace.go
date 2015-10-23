package deferclient

// Trace contains information about this client's trace and its producing package
type Trace struct {
	Out []byte `json:"Out"`
	Pkg []byte `json:"Pkg"`
}

// NewTrace instantitates and returns a new trace
// it is meant to be called once at the after the completing application tracing
func NewTrace() *Trace {

	t := &Trace{}

	return t
}
