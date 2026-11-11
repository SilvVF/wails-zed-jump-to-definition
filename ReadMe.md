# Zed Configuration

## 1. Add a task

Open the command palette (`Ctrl+Shift+P`) and select **zed: open tasks**.

Add:

```json
[
  {
    "label": "Go to Wails Go definition",
    "command": "C:\\path\\to\\main.exe",
    "args": [
      "$ZED_FILE",
      "$ZED_ROW",
      "$ZED_COLUMN",
      "$ZED_WORKTREE_ROOT"
    ],
    "use_new_terminal": false,
    "reveal": "never"
  }
]
```

Replace `C:\\path\\to\\main.exe` with the path to the executable.

## 2. Add a keybinding

Open the command palette (`Ctrl+Shift+P`) and select **zed: open keymap**.

Add:

```json
[
  {
    "context": "Editor",
    "bindings": {
      "alt-g": [
        "task::Spawn",
        {
          "task_name": "Go to Wails Go definition",
          "reevaluate_context": true
        }
      ]
    }
  }
]
```

# Usage

Place your cursor on a Wails method call and press **Alt+G** to jump to its Go implementation.
