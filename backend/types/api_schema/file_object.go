package api_schema

type (
	Filename struct {
		Name string `param:"*"`
	}
	FileInfo struct {
		Name     string `json:"Name"`
		Size     uint64 `json:"Size"`
		FileType string `json:"FileType"`
	}
)
