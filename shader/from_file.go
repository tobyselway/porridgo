package shader

import (
	"io/fs"
	"path/filepath"
)

func FromFile(filename string, filesystem fs.FS) (*Shader, error) {
	code, err := fs.ReadFile(filesystem, filename)
	if err != nil {
		return nil, err
	}
	return &Shader{
		Name: filepath.Base(filename),
		Code: string(code[:]),
	}, nil
}
