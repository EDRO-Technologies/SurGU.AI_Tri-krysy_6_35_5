package ai

type Response struct {
	Answer string `json:"answer"`
	Files  []File `json:"files"`
}

type File struct {
	Name string `json:"name"`
}

type Request struct {
	Question string `json:"question"`
}
