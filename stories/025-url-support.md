# Story 025: URL Support with Auto-Detection

## Context
Users often have tasks that reference documents, PRs, design docs, etc. at URLs. We want to make it easy to attach URLs to tasks and open them from the TUI.

## Requirements

### Auto-Detection
- Automatically detect `http://` and `https://` URLs in task descriptions
- No special tag required (though `url:https://...` also works)
- Support multiple URLs per task

### Visual Indication
- URLs rendered with special styling (underlined + colored)
- Modern terminals: URLs are clickable (Cmd+Click / Ctrl+Click)
- Visual indicator (🔗) shown for tasks containing URLs

### Keybinding: Press `o` to Open
- When a task is focused/selected, press `o`
- **Single URL**: Opens immediately in default browser
- **Multiple URLs**: Shows selection menu with arrow keys navigation
- **No URLs**: Shows brief message "No URLs in this task"

### Cross-Platform
- macOS: Uses `open` command
- Linux: Uses `xdg-open` command
- Windows: Uses `start` command

## Examples

### Single URL Task
```
(A) Review design doc https://docs.google.com/document/d/abc123 +design
```
Press `o` → Opens URL immediately

### Multiple URLs Task
```
(B) Review PR https://github.com/user/repo/pull/42 and staging https://staging.example.com +review
```
Press `o` → Shows selection menu:
```
┌─────────────────────────────────────────┐
│ Select URL to open:                     │
│                                         │
│ > https://github.com/user/repo/pull/42  │
│   https://staging.example.com           │
│                                         │
│ Press Enter to open • ESC to cancel     │
└─────────────────────────────────────────┘
```

### No URLs
```
(C) Regular task without links +project
```
Press `o` → Shows message briefly, no action

## Acceptance Criteria

1. ✅ Auto-detect HTTP/HTTPS URLs in task descriptions
2. ✅ URLs rendered with underline and color styling
3. ✅ URLs include OSC 8 escape codes for terminal clickability
4. ✅ Tasks with URLs show 🔗 indicator in list view
5. ✅ Press `o` on task with single URL opens it immediately
6. ✅ Press `o` on task with multiple URLs shows selection menu
7. ✅ Selection menu navigable with arrow keys
8. ✅ Enter in selection menu opens selected URL
9. ✅ ESC in selection menu cancels
10. ✅ Press `o` on task without URLs shows "No URLs" message
11. ✅ Cross-platform: Works on macOS, Linux, Windows
12. ✅ Help text updated to show `o` keybinding
13. ✅ Works in all view modes (overview, focus, table)

## Technical Notes

### URL Detection Regex
- Pattern: `https?://[^\s]+` (matches http:// or https:// followed by non-whitespace)
- Should handle common URL chars including query params, anchors, etc.
- Stop at whitespace or end of string

### OSC 8 Escape Sequence
- Format: `\x1b]8;;{URL}\x1b\\{TEXT}\x1b]8;;\x1b\\`
- Supported by: iTerm2, kitty, WezTerm, Windows Terminal, GNOME Terminal 3.50+

### URL Opening
```go
func openURL(url string) error {
    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "darwin":
        cmd = exec.Command("open", url)
    case "linux":
        cmd = exec.Command("xdg-open", url)
    case "windows":
        cmd = exec.Command("cmd", "/c", "start", url)
    default:
        return fmt.Errorf("unsupported platform")
    }
    return cmd.Start()
}
```

### State Management
- Add `urlSelectionMode bool` to Model
- Add `urls []string` to Model (for current task's URLs)
- Add `selectedURL int` to Model (for selection cursor)

## Out of Scope
- URL validation (we trust the user's input)
- URL shortening/truncation in display
- Custom URL schemes (only http/https)
- Opening URLs in specific browsers
