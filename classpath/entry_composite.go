package classpath

import (
	"errors"
	"strings"
)

type CompositeEntry []Entry

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
