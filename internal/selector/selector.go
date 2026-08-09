package selector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// Select provides a raw terminal selector when connected to a TTY and a line-based fallback otherwise.
func Select(in io.Reader, out io.Writer, dirs, initial []string) ([]string, error) {
	input, inputIsFile := in.(*os.File)
	output, outputIsFile := out.(*os.File)

	if inputIsFile && outputIsFile {
		if shouldSelectInteractive(true, true, term.IsTerminal(int(input.Fd())), term.IsTerminal(int(output.Fd()))) {
			return selectInteractive(input, output, dirs, initial)
		}
	}

	return selectLines(in, out, dirs, initial)
}

func shouldSelectInteractive(inputIsFile, outputIsFile, inputIsTerminal, outputIsTerminal bool) bool {
	return inputIsFile && outputIsFile && inputIsTerminal && outputIsTerminal
}

func selectInteractive(in, out *os.File, dirs, initial []string) ([]string, error) {
	state, err := term.MakeRaw(int(in.Fd()))
	if err != nil {
		return nil, fmt.Errorf("enable terminal input: %w", err)
	}
	defer func() { _ = term.Restore(int(in.Fd()), state) }()

	selected := make(map[string]bool)
	setSelected(selected, initial, dirs)

	query := ""
	cursor := 0
	inputBuffer := make([]byte, 1)

	for {
		visible := matchingDirectories(query, dirs)
		if len(visible) == 0 {
			cursor = 0
		} else if cursor >= len(visible) {
			cursor = len(visible) - 1
		}

		render(out, query, visible, selected, cursor)

		if _, err := in.Read(inputBuffer); err != nil {
			return nil, err
		}

		done, err := handleKey(inputBuffer[0], in, selected, &query, &cursor, visible, dirs)
		if err != nil {
			return nil, err
		}

		if done {
			return selectedPaths(dirs, selected), nil
		}
	}
}

//nolint:cyclop // A single key dispatch keeps terminal behavior explicit.
func handleKey(key byte, in *os.File, selected map[string]bool, query *string, cursor *int, visible, dirs []string) (bool, error) {
	switch key {
	case 3:
		return false, fmt.Errorf("selection cancelled")
	case 13, 10:
		return true, nil
	case 32:
		if len(visible) > 0 {
			toggleDirectory(selected, visible[*cursor], dirs)
		}
	case 127, 8:
		if *query != "" {
			*query = (*query)[:len(*query)-1]
			*cursor = 0
		}
	case 27:
		if err := readEscapeKey(in, cursor, len(visible)); err != nil {
			return false, err
		}
	default:
		if key >= 32 && key <= 126 {
			*query += string(key)
			*cursor = 0
		}
	}

	return false, nil
}

func render(out *os.File, query string, visible []string, selected map[string]bool, cursor int) {
	_, _ = fmt.Fprint(out, "\x1b[H\x1b[2J")
	writeTerminalLine(out, "Select production directories (up/down move, Space toggle, Enter confirm)")
	writeTerminalLine(out, "")

	start, end := viewport(len(visible), cursor, terminalHeight(out))
	for i := start; i < end; i++ {
		path := visible[i]

		pointer := " "
		if i == cursor {
			pointer = ">"
		}

		writeTerminalLine(out, fmt.Sprintf("%s [%s] %s", pointer, mark(selected[path]), path))
	}

	if len(visible) > end-start {
		writeTerminalLine(out, fmt.Sprintf("Showing %d-%d of %d directories", start+1, end, len(visible)))
	}

	writeTerminalPrompt(out, fmt.Sprintf("Search: %s", query))
	_, _ = fmt.Fprint(out, "\x1b[J")
}

func terminalHeight(out *os.File) int {
	_, height, err := term.GetSize(int(out.Fd()))
	if err != nil || height < 1 {
		return 24
	}

	return height
}

func viewport(total, cursor, height int) (int, int) {
	if total == 0 {
		return 0, 0
	}

	// Reserve two header lines, a status line, and the search prompt.
	rows := max(height-4, 1)

	rows = min(rows, total)

	if cursor < 0 {
		cursor = 0
	}

	if cursor >= total {
		cursor = total - 1
	}

	start := 0
	if cursor >= rows {
		start = cursor - rows + 1
	}

	end := start + rows
	if end > total {
		end = total
		start = end - rows
	}

	return start, end
}

func writeTerminalLine(out *os.File, line string) {
	_, _ = fmt.Fprintf(out, "\x1b[2K\r%s\n", line)
}

func writeTerminalPrompt(out *os.File, line string) {
	_, _ = fmt.Fprintf(out, "\x1b[2K\r%s", line)
}

func readEscapeKey(in *os.File, cursor *int, length int) error {
	if length == 0 {
		return nil
	}

	sequence := make([]byte, 2)
	if _, err := in.Read(sequence); err != nil {
		return err
	}

	if sequence[0] != '[' {
		return nil
	}

	switch sequence[1] {
	case 'A':
		*cursor = (*cursor + length - 1) % length
	case 'B':
		*cursor = (*cursor + 1) % length
	}

	return nil
}

func selectLines(in io.Reader, out io.Writer, dirs, initial []string) ([]string, error) {
	selected := map[string]bool{}
	setSelected(selected, initial, dirs)

	scanner := bufio.NewScanner(in)

	fmt.Fprintln(out, "Available directories:")
	printMatches(out, dirs, selected)
	fmt.Fprintln(out)

	for {
		fmt.Fprintln(out, "Production directories (type a search term, blank for all, or 'done' to finish):")
		fmt.Fprint(out, "> ")

		if !scanner.Scan() {
			return selectedPaths(dirs, selected), scanner.Err()
		}

		query := strings.TrimSpace(scanner.Text())
		if strings.EqualFold(query, "done") {
			return selectedPaths(dirs, selected), nil
		}

		visible := matchingDirectories(query, dirs)
		printMatches(out, visible, selected)

		if len(visible) == 0 {
			fmt.Fprintln(out, "No matching directories.")
			continue
		}

		fmt.Fprint(out, "Select indexes (comma-separated; repeat to toggle): ")

		if !scanner.Scan() {
			return nil, scanner.Err()
		}

		toggleSelected(selected, scanner.Text(), visible, dirs)
	}
}

func setSelected(selected map[string]bool, paths, dirs []string) {
	for _, path := range paths {
		for _, dir := range dirs {
			if isPathOrChild(dir, path) {
				selected[dir] = true
			}
		}
	}
}

func toggleDirectory(selected map[string]bool, path string, dirs []string) {
	wantSelected := !selected[path]
	for _, dir := range dirs {
		if isPathOrChild(dir, path) {
			selected[dir] = wantSelected
		}
	}
}

func isPathOrChild(path, parent string) bool {
	return parent == "." || path == parent || strings.HasPrefix(path, parent+"/")
}

func matchingDirectories(query string, dirs []string) []string {
	query = strings.ToLower(query)
	visible := make([]string, 0)

	for _, path := range dirs {
		if strings.Contains(strings.ToLower(path), query) {
			visible = append(visible, path)
		}
	}

	return visible
}

func printMatches(out io.Writer, paths []string, selected map[string]bool) {
	for i, path := range paths {
		fmt.Fprintf(out, "%d: [%s] %s\n", i+1, mark(selected[path]), path)
	}
}

func toggleSelected(selected map[string]bool, input string, visible, dirs []string) {
	for raw := range strings.SplitSeq(input, ",") {
		i, err := strconv.Atoi(strings.TrimSpace(raw))
		if err != nil || i < 1 || i > len(visible) {
			continue
		}

		toggleDirectory(selected, visible[i-1], dirs)
	}
}

func mark(value bool) string {
	if value {
		return "x"
	}

	return " "
}

func selectedPaths(dirs []string, selected map[string]bool) []string {
	result := make([]string, 0)

	for _, path := range dirs {
		if selected[path] {
			result = append(result, path)
		}
	}

	return result
}
