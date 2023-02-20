package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tidwall/gjson"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"github.com/ory/x/cmdx"
)

type (
	TableHeader interface {
		Header() []string
	}
	TableRow interface {
		TableHeader
		Columns() []string
		Interface() interface{}
	}
	Table interface {
		TableHeader
		Table() [][]string
		Interface() interface{}
		Len() int
	}
	Nil struct{}

	format string
)

const (
	FormatQuiet      format = "quiet"
	FormatTable      format = "table"
	FormatJSON       format = "json"
	FormatJSONPretty format = "json-pretty"
	FormatYAML       format = "yaml"
	FormatDefault    format = "default"

	FlagFormat = "format"

	None = "<none>"
)

func (Nil) String() string {
	return "null"
}

func (Nil) Interface() interface{} {
	return nil
}

func PrintErrors(cmd *cobra.Command, errs map[string]error) {
	for src, err := range errs {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s: %s\n", src, err.Error())
	}
}

func PrintRow(f format, io io.Writer, cmd *cobra.Command, row TableRow) {
	switch f {
	case FormatQuiet:
		if idAble, ok := row.(interface{ ID() string }); ok {
			_, _ = fmt.Fprintln(io, idAble.ID())
			break
		}
		_, _ = fmt.Fprintln(io, row.Columns()[0])
	case FormatJSON:
		printJSON(io, row.Interface(), false, "")
	case FormatYAML:
		printYAML(io, row.Interface())
	case FormatJSONPretty:
		printJSON(io, row.Interface(), true, "")
	case FormatTable, FormatDefault:
		w := tabwriter.NewWriter(io, 0, 8, 1, '\t', 0)

		fields := row.Columns()
		for i, h := range row.Header() {
			_, _ = fmt.Fprintf(w, "%s\t%s\t\n", h, fields[i])
		}

		_ = w.Flush()
	}
}

func PrintTable(f format, io io.Writer, table Table) {
	switch f {
	case FormatQuiet:
		if table.Len() == 0 {
			fmt.Fprintln(io)
		}

		if idAble, ok := table.(interface{ IDs() []string }); ok {
			for _, row := range idAble.IDs() {
				fmt.Fprintln(io, row)
			}
			break
		}

		for _, row := range table.Table() {
			fmt.Fprintln(io, row[0])
		}
	case FormatJSON:
		printJSON(io, table.Interface(), false, "")
	case FormatJSONPretty:
		printJSON(io, table.Interface(), true, "")
	case FormatYAML:
		printYAML(io, table.Interface())
	default:
		w := tabwriter.NewWriter(io, 0, 8, 1, '\t', 0)

		for _, h := range table.Header() {
			fmt.Fprintf(w, "%s\t", h)
		}
		fmt.Fprintln(w)

		for _, row := range table.Table() {
			fmt.Fprintln(w, strings.Join(row, "\t")+"\t")
		}

		_ = w.Flush()
	}
}

//relationTuples, err := client.NewCollection(resp.RelationTuples)
//client.PrintTable(client.FormatTable, GinkgoWriter, relationTuples)

func PrintTableFromRelationTuples(rels []*rts.RelationTuple, out io.Writer) {
	relationTuples, err := NewCollection(rels)
	PrintTable(FormatTable, out, relationTuples)
	if err != nil {
		panic("Encountered error: " + err.Error())
	}
}

type interfacer interface{ Interface() interface{} }

func PrintJSONAble(f format, io io.Writer, d interface{ String() string }) {
	var path string
	if d == nil {
		d = Nil{}
	}
	switch f {
	default:
		_, _ = fmt.Fprint(io, d.String())
	case FormatJSON:
		var v interface{} = d
		if i, ok := d.(interfacer); ok {
			v = i
		}
		printJSON(io, v, false, "")
	case FormatJSONPretty:
		var v interface{} = d
		if i, ok := d.(interfacer); ok {
			v = i
		}
		printJSON(io, v, true, path)
	case FormatYAML:
		var v interface{} = d
		if i, ok := d.(interfacer); ok {
			v = i
		}
		printYAML(io, v)
	}
}

func getFormat(format string) format {
	f := format

	switch {
	case f == string(FormatTable):
		return FormatTable
	case f == string(FormatJSON):
		return FormatJSON
	case f == string(FormatJSONPretty):
		return FormatJSONPretty
	case f == string(FormatYAML):
		return FormatYAML
	default:
		return FormatDefault
	}
}

func printJSON(w io.Writer, v interface{}, pretty bool, path string) {
	if path != "" {
		temp, err := json.Marshal(v)
		cmdx.Must(err, "Error encoding JSON: %s", err)
		v = gjson.GetBytes(temp, path).Value()
	}

	e := json.NewEncoder(w)
	if pretty {
		e.SetIndent("", "  ")
	}
	err := e.Encode(v)
	// unexpected error
	cmdx.Must(err, "Error encoding JSON: %s", err)
}

func printYAML(w io.Writer, v interface{}) {
	j, err := json.Marshal(v)
	cmdx.Must(err, "Error encoding JSON: %s", err)
	e, err := yaml.JSONToYAML(j)
	cmdx.Must(err, "Error encoding YAML: %s", err)
	_, _ = w.Write(e)
}

func RegisterJSONFormatFlags(flags *pflag.FlagSet) {
	flags.String(FlagFormat, string(FormatDefault), fmt.Sprintf("Set the output format. One of %s, %s, %s, %s and %s.", FormatDefault, FormatJSON, FormatYAML, FormatJSONPretty))
}

func RegisterFormatFlags(flags *pflag.FlagSet) {
	cmdx.RegisterNoiseFlags(flags)
	flags.String(FlagFormat, string(FormatDefault), fmt.Sprintf("Set the output format. One of %s, %s, %s, %s, and %s.", FormatTable, FormatJSON, FormatYAML, FormatJSONPretty))
}

type bodyer interface {
	Body() []byte
}

func PrintOpenAPIError(cmd *cobra.Command, err error) error {
	if err == nil {
		return nil
	}

	var be bodyer
	if !errors.As(err, &be) {
		return err
	}

	var didPrettyPrint bool
	if message := gjson.GetBytes(be.Body(), "error.message"); message.Exists() {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s\n", message.String())
		didPrettyPrint = true
	}
	if reason := gjson.GetBytes(be.Body(), "error.reason"); reason.Exists() {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s\n", reason.String())
		didPrettyPrint = true
	}

	if didPrettyPrint {
		return cmdx.FailSilently(cmd)
	}

	if body, err := json.MarshalIndent(json.RawMessage(be.Body()), "", "  "); err == nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s\nFailed to execute API request, see error above.\n", body)
		return cmdx.FailSilently(cmd)
	}

	return err
}
