package classpath

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const pathListSeparator = string(os.PathListSeparator)

type Entry interface {
	readClass(className string) ([]byte, Entry, error) // 寻找和加载 class 文件
	String() string                                    // 返回变量的字符串表示，相当于 java 的 toString
}

type DirEntry struct {
	absDir string
}

type CompositeEntry []Entry

type ZipEntry struct {
	absPath string
}

func newEntry(path string) Entry {
	if strings.Contains(path, pathListSeparator) {
		return newCompositeEntry(path)
	}
	if strings.HasSuffix(path, "*") {
		return newWildcardEntry(path)
	}
	if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") ||
		strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".ZIP") {
		return newZipEntry(path)
	}
	return newDirEntry(path)
}

func newDirEntry(path string) *DirEntry {
	var (
		absDir string
		err    error
	)
	absDir, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return &DirEntry{absDir}
}

func (self *DirEntry) readClass(className string) (data []byte, entry Entry, err error) {
	var fileName string
	fileName = filepath.Join(self.absDir, className)
	data, err = ioutil.ReadFile(fileName)
	return data, self, err
}

func (self *DirEntry) String() string {
	return self.absDir
}

func newCompositeEntry(pathList string) (compositeEntry CompositeEntry) {
	var (
		path  string
		entry Entry
	)
	compositeEntry = []Entry{}
	for _, path = range strings.Split(pathList, pathListSeparator) {
		entry = newEntry(path)
		compositeEntry = append(compositeEntry, entry)
	}
	return compositeEntry
}

func (self CompositeEntry) readClass(className string) (data []byte, entry Entry, err error) {
	var from Entry
	for _, entry = range self {
		data, from, err = entry.readClass(className)
		if err == nil {
			return data, from, nil
		}
	}
	return nil, nil, errors.New("class not found: " + className)
}

func (self CompositeEntry) String() string {
	var strs []string
	strs = make([]string, len(self))
	for i, entry := range self {
		strs[i] = entry.String()
	}
	return strings.Join(strs, pathListSeparator)
}

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

func newZipEntry(path string) *ZipEntry {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return &ZipEntry{absPath}
}

func (self *ZipEntry) readClass(className string) (data []byte, entry Entry, err error) {
	var (
		r  *zip.ReadCloser
		f  *zip.File
		rc io.ReadCloser
	)
	r, err = zip.OpenReader(self.absPath)
	if err != nil {
		return nil, nil, err
	}
	defer r.Close()
	for _, f = range r.File {
		if f.Name == className {
			rc, err = f.Open()
			if err != nil {
				return nil, nil, err
			}

			defer rc.Close()
			data, err = ioutil.ReadAll(rc)
			if err != nil {
				return nil, nil, err
			}
			return data, self, nil
		}
	}
	return nil, nil, errors.New("class not found: " + className)
}

func (self *ZipEntry) String() string {
	return self.absPath
}
