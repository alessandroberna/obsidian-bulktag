package mdProcessor

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"obsidian-tagfmt/internal/tag"
	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

type MdSetSt struct {
	DryRun bool
	Path string
}

var Settings MdSetSt

func (m *MdSetSt) checkEmptyTag () error {
	if tag.CheckEmptyTagMap() {
		return fmt.Errorf("internal Error: empty Tags map")
	}
	tag, exists := tag.LoadTag(m.Path)
	if !exists || tag == nil {
		return fmt.Errorf("internal Error: Tags map does not contain %s", m.Path)
	}
	if tag.FullTagStr() == "" {
		return fmt.Errorf("cannot apply an empty tag")
	}
	Tag = tag
	return nil
} 

var Tag *tag.FolderTag

func Main() error {
	err:= Settings.checkEmptyTag()
	if err != nil {
		return err
	}
	return traverseDir(Settings.Path)
}

func traverseDir(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			newPath := path + "/" + file.Name()
			// Creates proper pointers even if previously unexplored
			// Needed to get eventual child tags in subfolders relative
			// to the path passed to the fucntion
			Tag.NewTagGetter(newPath) 
			traverseDir(newPath)
		} else {
			if strings.HasSuffix(file.Name(), ".md") {
				err := processFile(path + "/" + file.Name(), Tag.FullTagStr())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func processFile(path string, tag string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Prepare a map to hold the front matter.
	var meta map[string]interface{}
	// Parse the front matter from the file.
	body, err := frontmatter.MustParse(bytes.NewReader(data), &meta)
	if err != nil {
		// If there is no front matter, initialize an empty map and use full file as body.
		if err == frontmatter.ErrNotFound {
			meta = make(map[string]any)
			body = data
		} else {
			return err
		}
	}

	// Update the "tags" field in meta.
	
	if existing, ok := meta["tags"]; ok {
		
		switch t := existing.(type) {
		case []interface{}:
			meta["tags"] = append(t, tag)
		case []string:
			meta["tags"] = append(t, tag)
		default:
			//// If "tags" is present but not a slice, overwrite it.
			//meta["tags"] = []string{Settings.Tag}

			// If "tags" is present but not a slice, attempt to convert it to a string and append.
			var existingTag string
			switch v := existing.(type) {
			case string:
				existingTag = v
			default:
				existingTag = fmt.Sprintf("%v", v)
			}
			meta["tags"] = []string{existingTag, tag}
		}
	} else {
		// If no tags exist, create the tags slice with the provided tag.
		meta["tags"] = []string{tag}
	}
	
	// Marshal the metadata into YAML.
	yamlBytes, err := yaml.Marshal(meta)
	if err != nil {
		return err
	}

	// Build the new file content: front matter + body.
	var b strings.Builder
	b.WriteString("---\n")
	b.Write(yamlBytes)
	b.WriteString("---\n\n")
	b.Write(body)

	if Settings.DryRun {
		if _, err := os.Stat("dryrun.md"); os.IsNotExist(err) {
			// File does not exist, create it with initial content.
			if err := os.WriteFile("dryrun.md", []byte(b.String()), 0644); err != nil {
				return err
			}
		} else {
			// File exists, open it in append mode.
			f, err := os.OpenFile("dryrun.md", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close()
		
			if _, err := f.WriteString(b.String()); err != nil {
				return err
			}
		}
	} else {
		// Write the new content back to the file.
		if err := os.WriteFile(path, []byte(b.String()), 0644); err != nil {
			return err
		}
	}
	return nil
}