package deferclient

// Trace contains information about this client's trace and its producing package
type Trace struct {
	Out       []byte `json:"Out"`
	Pkg       []byte `json:"Pkg"`
	CRC32     uint32 `json:"CRC32"`
	Size      int64  `json:"Size"`
	CommandId int    `json:"CommandId"`
}

// NewTrace instantitates and returns a new trace
// it is meant to be called once at the after the completing application tracing
func NewTrace(out []byte, pkg []byte, crc32 uint32, size int64, commandid int) *Trace {
	t := &Trace{
		Out:       out,
		Pkg:       pkg,
		CRC32:     crc32,
		Size:      size,
		CommandId: commandid,
	}

	return t
}
