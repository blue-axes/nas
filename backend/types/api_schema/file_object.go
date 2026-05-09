package api_schema

type (
	Filename struct {
		Name string `param:"*"`
	}
	FileInfo struct {
		Name     string   `json:"Name"`
		Size     uint64   `json:"Size"`
		FileType string   `json:"FileType"`
		Tags     []string `json:"Tags,omitempty"`
	}

	SearchReq struct {
		Keyword string `json:"Keyword" query:"Keyword"`
		Tag     string `json:"Tag" query:"Tag"`
	}

	UpdateTagsReq struct {
		Tags []string `json:"Tags"`
	}
)
