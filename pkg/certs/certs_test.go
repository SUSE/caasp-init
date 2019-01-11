package certs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/kubic-project/caasp-init/pkg/config"
)

var (
	initConfig = &config.KubicInitConfiguration{
		Bootstrap: config.BootstrapConfiguration{},
	}
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
	twoRegistriesWithCerts = &config.KubicInitConfiguration{
		Bootstrap: config.BootstrapConfiguration{
			Registries: []config.Registry{
				{Prefix: "mycompany.registry.com",
					Mirrors: []config.Mirror{
						{URL: "https://first.mirror.com", Certificate: "---- Cert Start ------ ACBDFEBABCDBFDBEBCDBABCDBABCBC"},
						{URL: "http://second.mirror.com", Certificate: "---- Cert Start ------ 9C1DFE7971D7FD7E7CD79B1DB9B1BC"},
					},
				},
				{Prefix: "somewhere.io",
					Mirrors: []config.Mirror{
						{URL: "https://local.lan.mirror.com", Certificate: "---- Cert Start ------ A6B39EBAB63B93BEB63BAB63BAB6B6"},
					},
				},
			},
		},
	}
	testNoHTTP = &config.KubicInitConfiguration{
		Bootstrap: config.BootstrapConfiguration{
			Registries: []config.Registry{
				{Prefix: "mycompany.registry.com",
					Mirrors: []config.Mirror{
						{URL: "second.mirror.com", Certificate: "---- Cert Start ------ 9C1DFE7971D7FD7E7CD79B1DB9B1BC"},
					},
				},
			},
		},
	}
	testURLError = &config.KubicInitConfiguration{
		Bootstrap: config.BootstrapConfiguration{
			Registries: []config.Registry{
				{Prefix: "mycompany.registry.com",
					Mirrors: []config.Mirror{
						{URL: ":", Certificate: "---- Cert Start ------ 9C1DFE7971D7FD7E7CD79B1DB9B1BC"},
					},
				},
			},
		},
	}
)

func TestWriteCertificates(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "caasp-init-certs")
	if err != nil {
		t.Fatalf("creating tmp dir: %s", err)
	}
	defer os.RemoveAll(tmpDir)
	type args struct {
		config     *config.KubicInitConfiguration
		certFolder string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{nil, tmpDir}, true},
		{"2", args{initConfig, tmpDir}, false},
		{"3", args{emptyConfig, tmpDir}, false},
		{"4", args{twoRegistries, tmpDir}, false},
		{"5", args{twoRegistriesWithCerts, tmpDir}, false},
		{"6", args{testNoHTTP, tmpDir}, true},
		{"7", args{twoRegistriesWithCerts, "/sys"}, true},
		{"7", args{testURLError, tmpDir}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			certsFolder = tt.args.certFolder
			if err := WriteCertificates(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("WriteCertificates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
