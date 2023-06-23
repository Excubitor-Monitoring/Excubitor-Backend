package plugins

import (
	"errors"
	"strings"
)

var frontendGetters = make(map[string]func(string) []byte)

var ErrFrontendComponentsNotProvided = errors.New("plugin does not provide any frontend components")
var ErrFrontendComponentFileNotFound = errors.New("not found")

func GetExternalFrontendFile(path string) ([]byte, error) {
	if strings.HasPrefix(path, "/static/external/") {
		trimmedPath := strings.TrimPrefix(path, "/static/external/")
		splitPath := strings.SplitN(trimmedPath, "/", 2)

		getter := frontendGetters[splitPath[0]]
		if getter != nil {
			fileContent := getter(splitPath[1])
			if len(fileContent) > 0 {
				return fileContent, nil
			} else {
				return nil, ErrFrontendComponentFileNotFound
			}
		} else {
			return nil, ErrFrontendComponentsNotProvided
		}
	} else {
		return nil, ErrFrontendComponentFileNotFound
	}
}
