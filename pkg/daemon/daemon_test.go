package daemon

import (
	"os"
	"testing"

	"github.com/kubic-project/caasp-init/pkg/config"
)

var (
	emptyConfig = &config.KubicInitConfiguration{
		Bootstrap: config.BootstrapConfiguration{
			Registries: []config.Registry{
				{Prefix: "",
					Mirrors: []config.Mirror{
						{URL: ""},
						{URL: ""},
					},
				},
				{Prefix: "",
					Mirrors: []config.Mirror{
						{URL: ""},
					},
				},
			},
		},
	}
	twoRegistries = &config.KubicInitConfiguration{
		Bootstrap: config.BootstrapConfiguration{
			Registries: []config.Registry{
				{Prefix: "mycompany.registry.com",
					Mirrors: []config.Mirror{
						{URL: "https://first.mirror.com"},
						{URL: "second.mirror.com"},
					},
				},
				{Prefix: "somewhere.io",
					Mirrors: []config.Mirror{
						{URL: "https://local.lan.mirror.com"},
					},
				},
			},
		},
	}
)

func TestWriteConfigFile(t *testing.T) {
	type args struct {
		config     *config.KubicInitConfiguration
		daemonFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{nil, "daemon.json"}, true},
		{"2", args{emptyConfig, "daemon.json"}, false},
		{"3", args{twoRegistries, "daemon.json"}, false},
		{"4", args{twoRegistries, "/sys/daemon.json"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			daemonFile = tt.args.daemonFile
			if err := WriteConfigFile(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("WriteConfigFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.RemoveAll("daemon.json")
	os.RemoveAll("/daemon.json")
}
