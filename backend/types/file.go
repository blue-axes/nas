package types

type (
	File struct {
		ID   uint   `json:"ID"`
		Name string `json:"Name"`
		Ext  string `json:"Ext"`
		Path string `json:"Path"`
		Size uint64 `json:"Size"`
		Md5  string `json:"Md5"`
	}
)
