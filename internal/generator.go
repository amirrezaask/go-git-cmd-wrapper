package main

//go:generate go run generator.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

type JSONModel struct {
	CommandName string         `json:"command_name,omitempty"`
	Enabled     bool           `json:"enabled"`
	Options     []JSONCmdModel `json:"options"`
}

type JSONCmdModel struct {
	MethodName  string `json:"method_name,omitempty"`
	Argument    string `json:"argument"`
	Arguments   string `json:"arguments"`
	Description string `json:"description"`
}

type GenCmdModel struct {
	Name      string
	ImportFMT bool
	Metas     []CmdMeta
}

type CmdMeta struct {
	Type       string
	Method     string
	Argument   string
	Cmd        string
	Comments   []string
	CmdComment string
}

// ByMethodName sort method by name.
type ByMethodName []CmdMeta

func (r ByMethodName) Len() int           { return len(r) }
func (r ByMethodName) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByMethodName) Less(i, j int) bool { return r[i].Method < r[j].Method }

const (
	fileTemplate = `/*
* CODE GENERATED AUTOMATICALLY
* THIS FILE MUST NOT BE EDITED BY HAND
 */
package {{ .Name }}

{{if .ImportFMT }}import (
	"fmt"

	"github.com/ldez/go-git-cmd-wrapper/types"
)
{{else}}import "github.com/ldez/go-git-cmd-wrapper/types"
{{end -}}

{{range .Metas -}}
{{if eq .Type "SIMPLE"}}{{template "templateCmdSimple" .}}
{{else if eq .Type "EQUAL_NO_OPTIONAL" }}{{template "templateCmdEqualNoOptional" .}}
{{else if eq .Type "EQUAL_WITHOUT_NAME" }}{{template "templateCmdEqualNoOptional" .}}
{{else if eq .Type "EQUAL_OPTIONAL_WITHOUT_NAME" }}{{template "templateCmdEqualOptional" .}}
{{else if eq .Type "EQUAL_OPTIONAL_WITH_NAME" }}{{template "templateCmdEqualOptional" .}}
{{else if eq .Type "WITH_PARAMETER" }}{{template "templateCmdWithParameter" .}}
{{else}} BUG
{{end -}}
{{end -}}
`
	templateCmdSimple = `{{- range $index, $element := .Comments}}
// {{if eq $index 0 }}{{ $.Method }} {{end}}{{ $element }} {{end}}
// {{ .CmdComment }}
func {{ .Method }}(g *types.Cmd) {
	g.AddOptions("{{ .Cmd }}")
}`

	templateCmdEqualNoOptional = `{{- range $index, $element := .Comments}}
// {{if eq $index 0 }}{{ $.Method }} {{end}}{{ $element }} {{end}}
// {{ .CmdComment }}
func {{ .Method }}({{ .Argument }} string) func(*types.Cmd) {
	return func(g *types.Cmd) {
		g.AddOptions(fmt.Sprintf("{{ .Cmd }}=%s", {{ .Argument }}))
	}
}`

	templateCmdEqualOptional = `{{- range $index, $element := .Comments}}
// {{if eq $index 0 }}{{ $.Method }} {{end}}{{ $element }} {{end}}
// {{ .CmdComment }}
func {{ .Method }}({{ .Argument }} string) func(*types.Cmd) {
	return func(g *types.Cmd) {
		if len({{ .Argument }}) == 0 {
			g.AddOptions("{{ .Cmd }}")
		} else {
			g.AddOptions(fmt.Sprintf("{{ .Cmd }}=%s", {{ .Argument }}))
		}
	}
}`

	templateCmdWithParameter = `{{- range $index, $element := .Comments}}
// {{if eq $index 0 }}{{ $.Method }} {{end}}{{ $element }} {{end}}
// {{ .CmdComment }}
func {{ .Method }}({{ .Argument }} string) func(*types.Cmd) {
	return func(g *types.Cmd) {
		g.AddOptions("{{ .Cmd }}")
		g.AddOptions({{ .Argument }})
	}
}`
)

var (
	// --quiet
	expCmdSimple = regexp.MustCompile(`^(-{1,2}([\w\d\-]+))$`)

	// --strategy=<strategy>
	expCmdEqualNoOptional = regexp.MustCompile(`^(-{1,2}([\w\d\-]+))=<([\w\d\- ]+)>$`)

	// --no-recurse-submodules[=yes|on-demand|no]
	expCmdEqualOptionalWithoutName = regexp.MustCompile(`^(-{1,2}([\w\d\-]+))\[=[\w\d\-()|]+]$`)

	// --recurse-submodules-default=[yes|on-demand]
	// --sign=(true|false|if-asked)
	expCmdEqualWithoutName = regexp.MustCompile(`^(-{1,2}([\w\d\-]+))=[\[(][\w\d\-|()]+[])]$`)

	// --log[=<n>]
	expCmdEqualOptionalWithName = regexp.MustCompile(`^(-{1,2}([\w\d\-]+))\[=<([\w\d\-)]+)>]$`)

	// --foo <bar>
	expCmdWithParameter = regexp.MustCompile(`^(-{1,2}([\w\d\-]+)) ?<([\w\d\-)]+)>$`)
)

func main() {
	filePath := "descriptions.json"
	var jsonModels []JSONModel
	file, err := ioutil.ReadFile(filePath)

	err = json.Unmarshal(file, &jsonModels)
	if err != nil {
		log.Fatal(err)
	}

	for _, jsonModel := range jsonModels {

		if len(jsonModel.CommandName) != 0 && jsonModel.Enabled {

			cmdModel := newGenCmdModel(jsonModel)

			data, err := generateFileContent(cmdModel)
			if err != nil {
				log.Fatal(err)
			}

			genFilePath := fmt.Sprintf("../%[1]s/%[1]s_gen.go", jsonModel.CommandName)

			fmt.Println(genFilePath)
			err = ioutil.WriteFile(genFilePath, []byte(data), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func generateFileContent(model GenCmdModel) (string, error) {

	base := template.New(model.Name)
	base.New("templateCmdSimple").Parse(templateCmdSimple)
	base.New("templateCmdEqualNoOptional").Parse(templateCmdEqualNoOptional)
	base.New("templateCmdEqualOptional").Parse(templateCmdEqualOptional)
	base.New("templateCmdWithParameter").Parse(templateCmdWithParameter)
	tmpl := template.Must(base.Parse(fileTemplate))

	b := &bytes.Buffer{}
	err := tmpl.Execute(b, model)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func newGenCmdModel(jsonModel JSONModel) GenCmdModel {
	return GenCmdModel{
		Name:      jsonModel.CommandName,
		Metas:     jsonCmdModelToCmdMetas(jsonModel.Options),
		ImportFMT: hasImportFMT(jsonModel.Options),
	}
}

func hasImportFMT(jsonCmdModels []JSONCmdModel) bool {
	for _, jsonCmdModel := range jsonCmdModels {
		if expCmdEqualNoOptional.MatchString(jsonCmdModel.Argument) ||
			expCmdEqualOptionalWithoutName.MatchString(jsonCmdModel.Argument) ||
			expCmdEqualWithoutName.MatchString(jsonCmdModel.Argument) ||
			expCmdEqualOptionalWithName.MatchString(jsonCmdModel.Argument) {
			return true
		}
	}
	return false
}

func jsonCmdModelToCmdMetas(jsonCmdModels []JSONCmdModel) []CmdMeta {

	metas := []CmdMeta{}

	for _, jsonCmdModel := range jsonCmdModels {
		if expCmdSimple.MatchString(jsonCmdModel.Argument) {
			metas = append(metas, newMetaCmdSimple(jsonCmdModel))
		} else if expCmdEqualNoOptional.MatchString(jsonCmdModel.Argument) {
			metas = append(metas, newMetaCmdEqualNoOptional(jsonCmdModel))
		} else if expCmdEqualOptionalWithoutName.MatchString(jsonCmdModel.Argument) {
			metas = append(metas, newMetaCmdEqualOptionalWithoutName(jsonCmdModel))
		} else if expCmdEqualWithoutName.MatchString(jsonCmdModel.Argument) {
			metas = append(metas, newMetaCmdEqualWithoutName(jsonCmdModel))
		} else if expCmdEqualOptionalWithName.MatchString(jsonCmdModel.Argument) {
			metas = append(metas, newMetaCmdEqualOptionalWithName(jsonCmdModel))
		} else if expCmdWithParameter.MatchString(jsonCmdModel.Argument) {
			metas = append(metas, newMetaCmdWithParameter(jsonCmdModel))
		} else {
			log.Println("fail", jsonCmdModel)
		}
	}

	sort.Sort(ByMethodName(metas))

	return metas
}

func newMetaCmdSimple(jsonCmdModel JSONCmdModel) CmdMeta {
	return newMetaCmd(expCmdSimple, jsonCmdModel, func(subMatch []string) CmdMeta {
		return newMeta(
			methodName(subMatch[2], jsonCmdModel),
			"",
			subMatch[1],
			"SIMPLE",
			jsonCmdModel.Description,
			jsonCmdModel.Arguments)
	})
}

func newMetaCmdEqualNoOptional(jsonCmdModel JSONCmdModel) CmdMeta {
	return newMetaCmd(expCmdEqualNoOptional, jsonCmdModel, func(subMatch []string) CmdMeta {
		return newMeta(
			methodName(subMatch[2], jsonCmdModel),
			subMatch[3],
			subMatch[1],
			"EQUAL_NO_OPTIONAL",
			jsonCmdModel.Description,
			jsonCmdModel.Arguments)
	})
}

func newMetaCmdEqualOptionalWithoutName(jsonCmdModel JSONCmdModel) CmdMeta {
	return newMetaCmd(expCmdEqualOptionalWithoutName, jsonCmdModel, func(subMatch []string) CmdMeta {
		return newMeta(
			methodName(subMatch[2], jsonCmdModel),
			"value",
			subMatch[1],
			"EQUAL_OPTIONAL_WITHOUT_NAME",
			jsonCmdModel.Description,
			jsonCmdModel.Arguments)
	})
}

func newMetaCmdEqualWithoutName(jsonCmdModel JSONCmdModel) CmdMeta {
	return newMetaCmd(expCmdEqualWithoutName, jsonCmdModel, func(subMatch []string) CmdMeta {
		return newMeta(
			methodName(subMatch[2], jsonCmdModel),
			"value",
			subMatch[1],
			"EQUAL_WITHOUT_NAME",
			jsonCmdModel.Description,
			jsonCmdModel.Arguments)
	})
}

func newMetaCmdEqualOptionalWithName(jsonCmdModel JSONCmdModel) CmdMeta {
	return newMetaCmd(expCmdEqualOptionalWithName, jsonCmdModel, func(subMatch []string) CmdMeta {
		return newMeta(
			methodName(subMatch[2], jsonCmdModel),
			subMatch[3],
			subMatch[1],
			"EQUAL_OPTIONAL_WITH_NAME",
			jsonCmdModel.Description,
			jsonCmdModel.Arguments)
	})
}

func newMetaCmdWithParameter(jsonCmdModel JSONCmdModel) CmdMeta {
	return newMetaCmd(expCmdWithParameter, jsonCmdModel, func(subMatch []string) CmdMeta {
		return newMeta(
			methodName(subMatch[2], jsonCmdModel),
			subMatch[3],
			subMatch[1],
			"WITH_PARAMETER",
			jsonCmdModel.Description,
			jsonCmdModel.Arguments)
	})
}

func methodName(raw string, jsonCmdModel JSONCmdModel) string {
	if len(jsonCmdModel.MethodName) == 0 {
		return raw
	}
	return jsonCmdModel.MethodName
}

type MetaBuilder func(subMatch []string) CmdMeta

func newMetaCmd(regexp *regexp.Regexp, jsonCmdModel JSONCmdModel, builder MetaBuilder) CmdMeta {
	subMatch := regexp.FindStringSubmatch(jsonCmdModel.Argument)
	return builder(subMatch)
}

func newMeta(rawMethodName, rawArg, cmd, cmdType, description, arguments string) CmdMeta {

	method := toGoName(rawMethodName, true)

	var arg string
	if len(rawArg) != 0 {
		arg = toGoName(rawArg, false)
		if strings.ToLower(method) == strings.ToLower(arg) {
			arg = "value"
		}
	}

	return CmdMeta{
		Type:       cmdType,
		Method:     method,
		Argument:   arg,
		Cmd:        cmd,
		Comments:   strings.Split(description, "\n"),
		CmdComment: arguments,
	}
}

func toGoName(kebab string, upperFirst bool) string {
	var camelCase string
	kebabTrim := strings.Replace(strings.TrimSpace(kebab), " ", "-", -1)
	isToUpper := false
	for i, runeValue := range kebabTrim {
		if i == 0 && upperFirst {
			camelCase += strings.ToUpper(string(runeValue))
		} else if isToUpper {
			camelCase += strings.ToUpper(string(runeValue))
			isToUpper = false
		} else {
			if runeValue == '-' {
				isToUpper = true
			} else {
				camelCase += string(runeValue)
			}
		}
	}
	return camelCase
}