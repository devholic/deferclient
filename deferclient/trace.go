package deferclient

// Trace contains information about this client's trace
type Trace struct {
	HTMLBody string `json:"HtmlBody"`
}

// SetHtmlBody sets a html body for this trace
func (t *Trace) SetHTMLBody() {

	t.HTMLBody = "<html><body>Empty trace</body></html>"
}

// NewTrace instantitates and returns a new trace
// it is meant to be called once at the after the completing application tracing
func NewTrace() *Trace {

	t := &Trace{}

	t.SetHTMLBody()

	return t
}
