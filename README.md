# obsidian-bulktag
A small TUI tool written in Go to bulk tag Obsidian notes.

## Main features:
- **Folder based, Hierarchical Tags:** Tags are assigned to folders. Each folder will inherit the parent folder's tags (if any). See [example](#tagging-example).
- **Bulk Tagging:** The *Apply Tags* command will add tags to all notes in the current directory and its subdirectories, all with a single keypress.
- **YAML parsing:** This parses the YAML frontmatter and adds the new tags without removing or modifying existing ones, if present. 
- **TUI:** Simple, keyboard-driven interface. VIM keybindings are supported.

## Tagging Example:
Assuming the following directory structure:
```
MyVault/               (No tag)
├── Projects/          (Tag assigned: `projects`)
│   ├── Active/        (Tag assigned: `active`)
│   │   └── note1.md
│   └── General/       (No tag)
│      └── note2.md
│      └── Unsorted/    (Tag assigned: `unsorted`)
│         └── note3.md
└── note4.md           (No tag)
```
The tags assigned to each note will be:
- `note1.md`: `projects/active`
    (Combines parent folder `Projects` + current folder `Active`)
- `note2.md`: `projects`
    (Inherited from parent folder `Projects`; `General` has no tag to add)
- `note3.md`: `projects/unsorted`
    (Parent: `projects` + current folder `Unsorted`)
- `note4.md`: No tag
    (Root folder has no tag, and the file is not in a tagged folder)

Assuming `note1.md` had the following frontmatter:
```yaml
---
tags: 
    - oldtag
date: 1989-06-04
---
```
It will be updated to:
```yaml
---
tags: 
    - oldtag
    - projects/active
date: 1989-06-04
---
```
# Disclaimer:
This project was quickly put together to solve a personal problem. It has not been extensively tested and should be considered experimental.

Please back up your notes before using the tool. I am not responsible for any data loss or corruption caused by it.

# Usage:
## Installation:
TBD
## Usage
TBD

## Keymaps
You can also press <kbd>?</kbd> or <kbd>i</kbd> while using the program to view the available keymaps.
   Key(s)                    | Description                     |
 |:--------------------------|:--------------------------------|
 | <kbd>g</kbd>              | Go to first                     |
 | <kbd>G</kbd>              | Go to last                      |
 | <kbd>↑</kbd>/<kbd>k</kbd>/<kbd>w</kbd>/<kbd>ctrl+p</kbd> | Move up                       |
 | <kbd>↓</kbd>/<kbd>j</kbd>/<kbd>s</kbd>/<kbd>ctrl+n</kbd> | Move down                     |
 | <kbd>pgup</kbd>/<kbd>K</kbd> | Page up                       |
 | <kbd>pgdown</kbd>/<kbd>J</kbd> | Page down                     |
 | <kbd>→</kbd>/<kbd>l</kbd>/<kbd>d</kbd>/<kbd>enter</kbd> | Open                           |
 | <kbd>←</kbd>/<kbd>h</kbd>/<kbd>a</kbd>/<kbd>backspace</kbd>/<kbd>esc</kbd> | Move back                     |
 | <kbd>q</kbd>/<kbd>ctrl+c</kbd> | Quit                           |
 | <kbd>t</kbd>              | Edit current folder's tag       |
 | <kbd>m</kbd>              | Apply tag to current folder and children |
 | <kbd>?</kbd>              | Toggle help                    |

## Screenshots:
 TBD

## License:
MIT

