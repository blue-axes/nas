package api_schema

type (
	Filename struct {
		Name string `param:"Name"`
	}
	FileInfo struct {
		Name     string `json:"Name"`
		Size     uint64 `json:"Size"`
		FileType string `json:"FileType"`
	}
)
