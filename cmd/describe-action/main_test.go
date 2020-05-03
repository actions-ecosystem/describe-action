package main

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_markdownTableWriter_writeTableInputs(t *testing.T) {
	type args struct {
		inputs Inputs
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "without types",
			args: args{
				inputs: map[string]*Input{
					"github_token": {
						Description: "A GitHub token.",
						Required:    true,
					},
					"repo": {
						Description: "The owner and repository name. e.g.) Codertocat/Hello-World.",
						Required:    false,
						Default:     "${{ github.repository }}",
					},
				},
			},
			want: `|     NAME     |                         DESCRIPTION                          | REQUIRED |          DEFAULT           |
|--------------|--------------------------------------------------------------|----------|----------------------------|
| github_token | A GitHub token.                                              | ` + "`true`" + `   | ` + "`N/A`" + `                      |
| repo         | The owner and repository name. e.g.) Codertocat/Hello-World. | ` + "`false`" + `  | ` + "`${{ github.repository }}`" + ` |
`,
		},
		{
			name: "with types",
			args: args{
				inputs: map[string]*Input{
					"github_token": {
						Description: "A GitHub token.",
						Type:        valueTypeString,
						Required:    true,
					},
					"repo": {
						Description: "The owner and repository name. e.g.) Codertocat/Hello-World.",
						Type:        valueTypeString,
						Required:    false,
						Default:     "${{ github.repository }}",
					},
				},
			},
			want: `|     NAME     |                         DESCRIPTION                          |   TYPE   | REQUIRED |          DEFAULT           |
|--------------|--------------------------------------------------------------|----------|----------|----------------------------|
| github_token | A GitHub token.                                              | ` + "`string`" + ` | ` + "`true`" + `   | ` + "`N/A`" + `                      |
| repo         | The owner and repository name. e.g.) Codertocat/Hello-World. | ` + "`string`" + ` | ` + "`false`" + `  | ` + "`${{ github.repository }}`" + ` |
`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := &bytes.Buffer{}
			w := newMarkdownTableWriter(got)
			w.writeTableInputs(tt.args.inputs)
			if diff := cmp.Diff(tt.want, got.String()); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}

func Test_markdownTableWriter_writeTableOutputs(t *testing.T) {
	type args struct {
		outputs Outputs
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "without types",
			args: args{
				outputs: map[string]*Output{
					"result": {
						Description: "The result of the action.",
					},
					"note": {
						Description: "The note about the action.",
					},
				},
			},
			want: `|  NAME  |        DESCRIPTION         |
|--------|----------------------------|
| note   | The note about the action. |
| result | The result of the action.  |
`,
		},
		{
			name: "with types",
			args: args{
				outputs: map[string]*Output{
					"result": {
						Description: "The result of the action.",
						Type:        valueTypeString,
					},
					"note": {
						Description: "The note about the action.",
						Type:        valueTypeString,
					},
				},
			},
			want: `|  NAME  |        DESCRIPTION         |   TYPE   |
|--------|----------------------------|----------|
| note   | The note about the action. | ` + "`string`" + ` |
| result | The result of the action.  | ` + "`string`" + ` |
`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := &bytes.Buffer{}
			w := newMarkdownTableWriter(got)
			w.writeTableOutputs(tt.args.outputs)
			if diff := cmp.Diff(tt.want, got.String()); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
