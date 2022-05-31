package terminal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// Type is an enum of all the available machine readable formats
type MachineReadableFormat int64

const (
	TableFormat MachineReadableFormat = iota
)

type machineReadableUI struct {
	mu     sync.Mutex
	format MachineReadableFormat
}

func MachineReadableUI(ctx context.Context, format MachineReadableFormat) UI {
	result := &machineReadableUI{
		format: format,
	}
	return result
}

// Input implements UI
func (ui *machineReadableUI) Input(input *Input) (string, error) {
	return "", ErrNonInteractive
}

// Interactive implements UI
func (ui *machineReadableUI) Interactive() bool {
	return false
}

// Output implements UI
func (ui *machineReadableUI) Output(msg string, raw ...interface{}) {
	if ui.format == TableFormat {
		msg = strings.Replace(msg, "\n", "\\n", -1)
		msg = strings.Replace(msg, "\r", "\\r", -1)
		msg = strings.Replace(msg, ",", "%!(VAGRANT_COMMA)", -1)
		tbl := NewTable()
		trow := []TableEntry{
			{Value: strconv.FormatInt(time.Now().Unix(), 10)},
			{Value: "ui"},
			{Value: "output"},
			{Value: msg},
		}
		tbl.Rows = append(tbl.Rows, trow)
		ui.Table(tbl)
	}
}

func (ui *machineReadableUI) ClearLine() {
	// NO-OP
}

// NamedValues implements UI
func (ui *machineReadableUI) NamedValues(rows []NamedValue, opts ...Option) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

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
func (ui *machineReadableUI) OutputWriters() (io.Writer, io.Writer, error) {
	return os.Stdout, os.Stderr, nil
}

// Status implements UI
func (ui *machineReadableUI) Status() Status {
	return &nonInteractiveStatus{mu: &ui.mu}
}

func (ui *machineReadableUI) StepGroup() StepGroup {
	return &nonInteractiveStepGroup{mu: &ui.mu}
}

// Table implements UI
func (ui *machineReadableUI) Table(tbl *Table, opts ...Option) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	// Build our config and set our options
	cfg := &config{Writer: color.Output}
	for _, opt := range opts {
		opt(cfg)
	}

	table := tablewriter.NewWriter(cfg.Writer)
	table.SetHeader(tbl.Headers)
	table.SetBorder(false)
	table.SetAutoWrapText(false)

	for _, row := range tbl.Rows {
		colors := make([]tablewriter.Colors, len(row))
		entries := make([]string, len(row))

		for i, ent := range row {
			entries[i] = ent.Value

			color, ok := colorMapping[ent.Color]
			if ok {
				colors[i] = tablewriter.Colors{color}
			}
		}

		table.Rich(entries, colors)
	}

	table.Render()
}
