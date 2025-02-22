package tag

type FolderTag struct {
	Tag    string // tag for a folder entry
	Parent *FolderTag
}

var tagMap = make(map[string]*FolderTag)

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
	if tag, exists := tagMap[path]; exists {
		return tag
	} else {
		tag := &FolderTag{Tag: "", Parent: parent}
		tagMap[path] = tag
		return tag
	}
}