package impl

import (
	"fmt"
	"testing"
)

func TestNameCreator_CreateName(t *testing.T) {

	namer := NewNamer()
	basePath := "/tmp"
	ext := "png"
	prefix := []string{"1", "2", "3"}

	baseName, paths := namer.CreateName(basePath, ext, prefix)

	pathsCheck := make(map[string]bool)
	pathsCheck[fmt.Sprintf("%s/%s.%s", basePath, baseName, ext)] = false

	for _, p := range prefix {
		pathsCheck [fmt.Sprintf("%s/%s_%s.%s", basePath, p, baseName, ext)] = false
	}

	for _, path := range paths {
		_, ok := pathsCheck[path]
		if !ok {
			t.Fatal("Created wrong paths")
		}
	}
}
