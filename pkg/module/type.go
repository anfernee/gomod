package module

type moduleType struct {
	Version string     `json:"Version"`
	Time    string     `json:"Time"`
	Origin  originType `json:"Origin"`
}

type originType struct {
	VCS  string `json:"VCS"`
	URL  string `json:"URL"`
	Ref  string `json:"Ref"`
	Hash string `json:"Hash"`
}
