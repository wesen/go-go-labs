package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"text/template"
)

func (ei *EmrichenInterpreter) handleFormat(node *yaml.Node) (*yaml.Node, error) {
	formatString := node.Value

	ret, err := ei.renderFormatString(formatString)
	if err != nil {
		return nil, err
	}

	return ValueToNode(ret)
}

func (ei *EmrichenInterpreter) renderFormatString(formatString string) (string, error) {
	tmpl, err := template.New("format").Parse(formatString)
	if err != nil {
		return "", fmt.Errorf("error parsing format string: %v", err)
	}

	var formatted bytes.Buffer
	frame := ei.env.GetCurrentFrame()
	vars := map[string]interface{}{}
	if frame.Variables != nil {
		vars = frame.Variables
	}

	for _, funcMap := range ei.funcmaps {
		tmpl = tmpl.Funcs(funcMap)
	}

	tmpl = tmpl.Funcs(
		map[string]interface{}{
			"lookup": func(path string) interface{} {
				v, err := ei.LookupFirst(path)
				if err != nil {
					return nil
				}
				v_, _ := NodeToInterface(v)
				return v_
			},
			"lookupAll": func(path string) []interface{} {
				v, err := ei.LookupAll(path)
				if err != nil {
					return nil
				}
				v_, _ := NodeToSlice(v)
				return v_
			},
			"exists": func(path string) bool {
				_, err := ei.LookupFirst(path)
				return err == nil
			},
		},
	)
	if err := tmpl.Execute(&formatted, vars); err != nil {
		return "", fmt.Errorf("error executing format template: %v", err)
	}

	return formatted.String(), nil
}
