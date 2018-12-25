package classpath

import (
	"os"
	"path/filepath"
	"strings"
)

func newWildcardEntry(path string) (compositeEntry CompositeEntry) {
	var (
		baseDir  string
		jarEntry *ZipEntry
	)
	baseDir = path[:len(path)-1] // remove *
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != baseDir {
			return filepath.SkipDir
		}
		if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") {
			jarEntry = newZipEntry(path)
			compositeEntry = append(compositeEntry, jarEntry)
		}
		return nil
	}
	filepath.Walk(baseDir, walkFn)
	return compositeEntry
}
