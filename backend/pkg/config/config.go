package config

import (
	"github.com/koding/multiconfig"
	"path"
)

func Load(filePath string, cfg interface{}) error {
	l := getLoader(filePath)
	return l.Load(cfg)
}

func MustLoad(filePath string, cfg interface{}) {
	l := getLoader(filePath)
	err := l.Load(cfg)
	if err != nil {
		panic(err)
	}
}

func getLoader(filePath string) multiconfig.Loader {
	var (
		loaders = make([]multiconfig.Loader, 0)
	)
	ext := path.Ext(filePath)
	switch ext {
	case ".json":
		loaders = append(loaders, &multiconfig.JSONLoader{Path: filePath})
	case ".yaml", ".yml":
		loaders = append(loaders, &multiconfig.YAMLLoader{Path: filePath})
	case ".toml":
		loaders = append(loaders, &multiconfig.TOMLLoader{Path: filePath})
	}

	return multiconfig.DefaultLoader{
		Loader:    multiconfig.MultiLoader(loaders...),
		Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
	}
}
