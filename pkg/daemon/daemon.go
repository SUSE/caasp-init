package daemon

import (
	"html/template"
	"os"

	"github.com/kubic-project/caasp-init/pkg/config"
)

const (
	daemonTemplate = `{ {{if .Bootstrap.Registries }} {{ $first := index .Bootstrap.Registries 0 }}
  {{ if ne $first.Prefix ""}}"registries": [  {{range $r, $registry := .Bootstrap.Registries}}   {{ if $r }},{{end}}
    {
    "Mirrors": [ {{range $m, $mirror := $registry.Mirrors}} {{if $m}},{{end}}
      {
        "URL": "{{$mirror.URL}}"
      }
      {{end}}],
      "Prefix": "{{$registry.Prefix}}"
    }
  {{end}}],
  {{end}}{{end}}"iptables":false,
  "log-level": "warn"
}
`
)

var (
	daemonFile = "/etc/docker/daemon.json"
)

// WriteConfigFile writes the daemon config file
// will be generated from a basic template
// and will include any mirror specified in the configuration
func WriteConfigFile(config *config.KubicInitConfiguration) error {
	tmpl, err := template.New("kubic").Parse(daemonTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(daemonFile)
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, config)
	if err != nil {
		os.RemoveAll(daemonFile)
		return err
	}
	file.Sync()
	file.Close()
	return nil
}
