package packet

type Player struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Players struct {
	Max    int `json:"max"`
	Online int `json:"online"`
	//Sample []Player `json:"sample"`
}

type Description struct {
	Text string `json:"text"`
}

type StatusRequestResponse struct {
	Version            `json:"version"`
	Players            `json:"players"`
	Description        `json:"description"`
	Favicon            string `json:"favicon"`
	EnforcesSecureChat bool   `json:"enforcesSecureChat"`
}
