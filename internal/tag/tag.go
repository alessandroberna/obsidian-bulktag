package tag

import "sync"

type FolderTag struct {
	Tag    string // tag for a folder entry
	Parent *FolderTag
}

//var TagMap = make(map[string]*FolderTag)
var TagMap sync.Map

func StoreTag(path string, tag *FolderTag) {
	TagMap.Store(path, tag)
}

func LoadTag(key string) (*FolderTag, bool) {
    value, ok := TagMap.Load(key)
    if !ok {
        return nil, false
    }
    folderTag, ok := value.(*FolderTag)
    if !ok {
        return nil, false
    }
    return folderTag, true
}

func ConditionalSlashJoin(string1 string, string2 string) string {
	if string1 == "" {
		return string2
	}
	if string2 == "" {
		return string1
	}
	return string1 + "/" + string2
}

func (f *FolderTag) ParentTagsStr() string {
	if f.Parent != nil {
		return ConditionalSlashJoin(f.Parent.ParentTagsStr(), f.Parent.Tag) // ../../. + "/" + ../.
	} else {
		return ""
	}
}

func (f *FolderTag) FullTagStr() string {
	if f.Tag != "" {
		return ConditionalSlashJoin(f.ParentTagsStr(), f.Tag)
	} else {
		return f.ParentTagsStr()
	}
}

func NewTagGetter(path string, parent *FolderTag) *FolderTag {
	if tag, exists := LoadTag(path); exists {
		return tag
	} else {
		tag := &FolderTag{Tag: "", Parent: parent}
		StoreTag(path, tag)
		return tag
	}
}

func (f * FolderTag) NewTagGetter (path string) *FolderTag {
	if tag, exists := LoadTag(path); exists {
		return tag
	} else {
		tag := &FolderTag{Tag: "", Parent: f}
		StoreTag(path, tag)
		return tag
	}
}

// returns true if the map is empty
func CheckEmptyTagMap() bool {
	empty := true
	TagMap.Range(func(key, value any) bool {
		empty = false
		return false
	})
	return empty
}
