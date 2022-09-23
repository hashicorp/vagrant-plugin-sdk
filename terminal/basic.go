package terminal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/bgentry/speakeasy"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// basicUI
type basicUI struct {
	ctx    context.Context
	status *spinnerStatus
}

func BasicUI(ctx context.Context) UI {
	result := &basicUI{
		ctx:    ctx,
		status: nil,
	}

	return result
}

// Input implements UI
func (ui *basicUI) Input(input *Input) (string, error) {
	var buf bytes.Buffer

	// Write the prompt, add a space.
	ui.Output(input.Prompt, WithStyle(input.Style), WithWriter(&buf), WithColor(input.Color))
	fmt.Fprint(color.Output, strings.TrimRight(buf.String(), "\r\n"))
	fmt.Fprint(color.Output, " ")

	// Ask for input in a go-routine so that we can ignore it.
	errCh := make(chan error, 1)
	lineCh := make(chan string, 1)
	go func() {
		var line string
		var err error
		if input.Secret && isatty.IsTerminal(os.Stdin.Fd()) {
			line, err = speakeasy.Ask("")
		} else {
			r := bufio.NewReader(os.Stdin)
			line, err = r.ReadString('\n')
		}
		if err != nil {
			errCh <- err
			return
		}

		lineCh <- strings.TrimRight(line, "\r\n")
	}()

	select {
	case err := <-errCh:
		return "", err
	case line := <-lineCh:
		return line, nil
	case <-ui.ctx.Done():
		// Print newline so that any further output starts properly
		fmt.Fprintln(color.Output)
		return "", ui.ctx.Err()
	}
}

// Interactive implements UI
func (ui *basicUI) Interactive() bool {
	return isatty.IsTerminal(os.Stdin.Fd())
}

// MachineReadable implements UI
func (ui *basicUI) MachineReadable() bool {
	return false
}

// ClearLine implements UI
func (ui *basicUI) ClearLine() {
	_, _, _, w, _ := Interpret("")
	w.Write([]byte("\r\033[K"))
}

// Output implements UI
func (ui *basicUI) Output(msg string, raw ...interface{}) {
	msg, style, disableNewline, w, _ := Interpret(msg, raw...)

	var printer *color.Color
	switch style {
	case HeaderStyle, WarningBoldStyle, ErrorBoldStyle, SuccessBoldStyle, InfoBoldStyle:
		printer = colorInfoBold
	default:
		printer = colorInfo
	}

	switch style {
	case HeaderStyle:
		msg = printer.Sprintf("\n==> " + msg)
	case ErrorStyle, ErrorBoldStyle:
		lines := strings.Split(msg, "\n")
		if len(lines) > 0 {
			printer.Sprintf("! " + lines[0])
			for _, line := range lines[1:] {
				printer.Sprintf("  " + line)
			}
		}
		msg = strings.Join(lines, "\n")
	case WarningStyle, WarningBoldStyle:
		msg = printer.Sprintf("WARNING: " + msg)
	default:
		lines := strings.Split(msg, "\n")
		for i, line := range lines {
			lines[i] = printer.Sprintf("  %s", line)
		}
		msg = strings.Join(lines, "\n")
	}

	st := ui.status
	if st != nil {
		if st.Pause() {
			defer st.Start()
		}
	}

	// Write it
	if disableNewline {
		fmt.Fprint(w, msg)
	} else {
		fmt.Fprintln(w, msg)
	}
}

// NamedValues implements UI
func (ui *basicUI) NamedValues(rows []NamedValue, opts ...Option) {
	cfg := &config{Writer: color.Output}
	for _, opt := range opts {
		opt(cfg)
	}

	var buf bytes.Buffer
	tr := tabwriter.NewWriter(&buf, 1, 8, 0, ' ', tabwriter.AlignRight)
	for _, row := range rows {
		switch v := row.Value.(type) {
		case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
			fmt.Fprintf(tr, "  %s: \t%d\n", row.Name, row.Value)
		case float32, float64:
			fmt.Fprintf(tr, "  %s: \t%f\n", row.Name, row.Value)
		case bool:
			fmt.Fprintf(tr, "  %s: \t%v\n", row.Name, row.Value)
		case string:
			if v == "" {
				continue
			}
			fmt.Fprintf(tr, "  %s: \t%s\n", row.Name, row.Value)
		default:
			fmt.Fprintf(tr, "  %s: \t%s\n", row.Name, row.Value)
		}
	}

	tr.Flush()
	fmt.Fprintln(cfg.Writer, buf.String())
}

// OutputWriters implements UI
func (ui *basicUI) OutputWriters() (io.Writer, io.Writer, error) {
	return os.Stdout, os.Stderr, nil
}

// Status implements UI
func (ui *basicUI) Status() Status {
	if ui.status == nil {
		ui.status = newSpinnerStatus(ui.ctx)
	}

	return ui.status
}

func (ui *basicUI) StepGroup() StepGroup {
	ctx, cancel := context.WithCancel(ui.ctx)
	display := NewDisplay(ctx, color.Output)

	return &fancyStepGroup{
		ctx:     ctx,
		cancel:  cancel,
		display: display,
		done:    make(chan struct{}),
	}
}
