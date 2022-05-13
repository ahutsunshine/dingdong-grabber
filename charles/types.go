package charles

type Session struct {
	Method          string         `json:"method"`
	ProtocolVersion string         `json:"protocolVersion"`
	Scheme          string         `json:"scheme"`
	Path            string         `json:"path"`
	Query           string         `json:"query"`
	RemoteAddress   string         `json:"remoteAddress"`
	ClientAddress   string         `json:"clientAddress"`
	ClientPort      int            `json:"clientPort"`
	Request         SessionRequest `json:"request"`
}

type SessionRequest struct {
	Header SessionRequestHeader `json:"header"`
}

type SessionRequestHeader struct {
	Headers []HeaderEntry `json:"headers"`
}

type HeaderEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
