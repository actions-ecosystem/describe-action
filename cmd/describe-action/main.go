package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

var (
	actionYAMLPath = flag.String("yaml", "action.yml", "The filepath to action.yml")
	onlyInput      = flag.Bool("input", false, "Whether only print inputs")
	onlyOutput     = flag.Bool("output", false, "Whether only print outputs")
	hasTypes       = flag.Bool("type", false, "Whether table has types")
)

const (
	valueTypeString = "string"
	valueTypeNumber = "number"
	valueTypeBool   = "bool"
)

var valueTypes = []string{valueTypeString, valueTypeNumber, valueTypeBool}

const maxColumnWidth = 256

type Manifest struct {
	Inputs  Inputs  `json:"inputs"`
	Outputs Outputs `json:"outputs"`
}

type Input struct {
	Description string `json:"description"`
	Type        string `json:"type"` // Type is an additional field for Markdown table.
	Required    bool   `json:"required"`
	Default     string `json:"default"`
}

type Inputs map[string]*Input

type Output struct {
	Description string `json:"description"`
	Type        string `json:"type"` // Type is an additional field for Markdown table.
}

type Outputs map[string]*Output

type markdownTableWriter struct {
	tw *tablewriter.Table
}

func main() {
	flag.Parse()

	buf, err := ioutil.ReadFile(*actionYAMLPath)
	if err != nil {
		log.Fatal(err)
	}

	m := &Manifest{}
	if err := yaml.Unmarshal(buf, m); err != nil {
		log.Fatal(err)
	}

	if *hasTypes && !*onlyOutput {
		for name := range m.Inputs {
			prompt := &survey.Select{
				Message: fmt.Sprintf(`Type of "inputs.%s":`, name),
				Options: valueTypes,
			}
			survey.AskOne(prompt, &m.Inputs[name].Type)
		}
	}
	if *hasTypes && !*onlyInput {
		for name := range m.Outputs {
			prompt := &survey.Select{
				Message: fmt.Sprintf(`Type of "outputs.%s":`, name),
				Options: valueTypes,
			}
			survey.AskOne(prompt, &m.Outputs[name].Type)
		}
	}

	mtw := newMarkdownTableWriter(os.Stdout)

	if *onlyInput {
		mtw.writeTableInputs(m.Inputs)
		return
	}

	if *onlyOutput {
		mtw.writeTableOutputs(m.Outputs)
		return
	}

	mtw.writeTableInputs(m.Inputs)
	fmt.Println()
	mtw.writeTableOutputs(m.Outputs)
}

func newMarkdownTableWriter(w io.Writer) *markdownTableWriter {
	tw := tablewriter.NewWriter(w)

	tw.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	tw.SetCenterSeparator("|")
	tw.SetColWidth(maxColumnWidth)

	return &markdownTableWriter{tw: tw}
}

func (w *markdownTableWriter) writeTableInputs(inputs Inputs) {
	var noTypes bool
	for _, input := range inputs {
		if input.Type == "" {
			noTypes = true
			break
		}
	}

	if noTypes {
		w.tw.SetHeader([]string{"Name", "Description", "Required", "Default"})
	} else {
		w.tw.SetHeader([]string{"Name", "Description", "Type", "Required", "Default"})
	}

	data := make([][]string, 0, len(inputs))
	for name, input := range inputs {
		inputDefault := input.Default
		if inputDefault == "" {
			inputDefault = "N/A"
		}

		if noTypes {
			data = append(data, []string{backtickString(name), input.Description, backtickBool(input.Required), backtickString(inputDefault)})
			continue
		}
		data = append(data, []string{backtickString(name), input.Description, backtickString(input.Type), backtickBool(input.Required), backtickString(inputDefault)})
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i][0] < data[j][0] // Sort data with input's name ascending alphabetically
	})

	w.tw.AppendBulk(data)

	w.tw.Render()
}

func (w *markdownTableWriter) writeTableOutputs(outputs Outputs) {
	var noTypes bool
	for _, output := range outputs {
		if output.Type == "" {
			noTypes = true
			break
		}
	}

	if noTypes {
		w.tw.SetHeader([]string{"Name", "Description"})
	} else {
		w.tw.SetHeader([]string{"Name", "Description", "Type"})
	}

	data := make([][]string, 0, len(outputs))
	for name, output := range outputs {
		if noTypes {
			data = append(data, []string{backtickString(name), output.Description})
			continue
		}
		data = append(data, []string{backtickString(name), output.Description, backtickString(output.Type)})
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i][0] < data[j][0] // Sort data with output's name ascending alphabetically
	})

	w.tw.AppendBulk(data)

	w.tw.Render()
}

func backtickString(s string) string {
	return fmt.Sprintf("`%s`", s)
}

func backtickBool(b bool) string {
	return fmt.Sprintf("`%t`", b)
}
