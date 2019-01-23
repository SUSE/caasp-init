/*
 * Copyright 2019 SUSE LINUX GmbH, Nuernberg, Germany..
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

const (
	configContent = `---
apiVersion: kubic.suse.com/v1alpha2
kind: KubicInitConfiguration
features:
  PSP: true
runtime:
  engine: crio
paths:
  kubeadm: /usr/bin/kubeadm
auth:
  oidc:
   issuer: https://some.name.com:32000
   clientID: kubernetes
   ca: /etc/kubernetes/pki/ca.crt
   username: email
   groups: groups
certificates:
  directory: /etc/kubernetes/pki
  caCrt:
  caCrtHash:
etcd:
  local:
    serverCertSANs: []
    peerCertSANs: []
manager:
  image: "kubic-init:latest"
clusterFormation:
  seeder: some-node.com
  token: 94dcda.c271f4ff502789ca
  autoApprove: false
network:
  bind:
    address: 0.0.0.0
    interface: eth0
  podSubnet: "172.16.0.0/13"
  serviceSubnet: "172.24.0.0/16"
  proxy:
    http: my-proxy.com:8080
    https: my-proxy.com:8080
    noProxy: localdomain.com
    systemwide: false
  dns:
    domain: someDomain.local
    externalFqdn: some.name.com
  cni:
    driver: flannel
    image: registry.opensuse.org/devel/caasp/kubic-container/container/kubic/flannel:0.9.1
bootstrap:
  registries:
    - prefix: https://mycompany.registry.com
      mirrors:
        - url: https://mycompany.airgapped.com
        - url: https://mycompany2.airgapped.com
          certificate: "-----BEGIN CERTIFICATE-----
MIeGJzCCBA+gAweBAgeBATANBgkqhkeG9w0BAQUFADCBsjELMAkGA1UEBhMCRlex
DzANBgNVBAgMBkFsc2FjZTETMBEGA1UEBwwKU3RyYXNeb3VyZzEYMBYGA1UECgwP
hnx8SB3sVJZHeer8f/UQQwqbAO+Kdy70NmbSaqaVtp8jOxLeedWkwSyRTsuU6D8e
DeH5uEqBXExjrj0FslxcVKdVj5glVcSmkLwZKbEU1OKwleT/eXFhvooWhQ==
-----END CERTIFICATE-----"
          fingerprint: "E8:73:0C:C5:84:B1:EB:17:2D:71:54:4D:89:13:EE:47:36:43:8D:BF:5D:3C:0F:5B:FC:75:7E:72:28:A9:7F:73"
          hashalgorithm: "SHA256"
    - prefix: https://registry.io
      mirrors:
        - url: https://mycompany.airgapped.com
        - url: https://mycompany2.airgapped.com
          certificate: "-----BEGIN CERTIFICATE-----
MIItJzCCBA+tAwIBAtIBATANBtkqhkit9w0BAQUFADCBsjELMAktA1UEBhMCRlIx
DzANBtNVBAtMBkFsc2FjZTETMBEtA1UEBwwKU3RyYXNib3VyZzEYMBYtA1UECtwP
hnx8SB3sVJZHeer8f/UQQwqbAO+Kdy70NmbSaqaVtp8jOxLiidWkwSyRTsuU6D8i
DiH5uEqBXExjrj0FslxcVKdVj5tlVcSmkLwZKbEU1OKwleT/iXFhvooWhQ==
-----END CERTIFICATE-----"
          fingerprint: "E8:73:0C:C5:84:B1:EB:17:2D:71:54:4D:89:13:EE:47:36:43:8D:BF:5D:3C:0F:5B:FC:75:7E:72:28:A9:7F:73"
          hashalgorithm: "SHA256"
`
	configContentWithError = `---
bootstrap:
  registries:
    - 
    prefix: https://mycompany.registry.com
    mirrors:
        - url: https://mycompany.airgapped.com
        - url: https://mycompany2.airgapped.com
          certificate: "-----BEGIN CERTIFICATE-----
MIIGpzCCBA+gAwIBAgIBATANBgkqhkiG9w0BAQUFADCBspELMAkGA1UEBhMCRlIx
DzANBgNVBAgMBkFsc2FpZTETMBEGA1UEBwwKU3RyYXNib3VyZzEYMBYGA1UECgwP
hnx8SB3sVpZHeer8f/UQQwqbAO+Kdy70NmbSaqaVtp8pOxLiidWkwSyRTsuU6D8i
DiH5uEqBXExprp0FslxcVKdVp5glVcSmkLwZKbEU1OKwleT/iXFhvooWhQ==
-----END CERTIFICATE-----"
          fingerprint: "E8:73:0C:C5:84:B1:EB:17:2D:71:54:4D:89:13:EE:47:36:43:8D:BF:5D:3C:0F:5B:FC:75:7E:72:28:A9:7F:73"
          hashalgorithm: "SHA256"
    - prefix: https://registry.io
      mirrors:
        - url: https://mycompany.airgapped.com
        - url: https://mycompany2.airgapped.com
          certificate: "-----BEGIN CERTIFICATE-----
LIIGJzCCBA+gAwIBAgIBATANBgkqhkiG9w0BAQUFADCBsjELlAkGA1UEBhlCRlIx
DzANBgNVBAglBkFsc2FjZTETlBEGA1UEBwwKU3RyYXNib3VyZzEYlBYGA1UECgwP
hnx8SB3sVJZHeer8f/UQQwqbAO+Kdy70NlbSaqaVtp8jOxLiidWkwSyRTsuU6D8i
DiH5uEqBXExjrj0FslxcVKdVj5glVcSlkLwZKbEU1OKwleT/iXFhvooWhQ==
-----END CERTIFICATE-----"
          fingerprint: "E8:73:0C:C5:84:B1:EB:17:2D:71:54:4D:89:13:EE:47:36:43:8D:BF:5D:3C:0F:5B:FC:75:7E:72:28:A9:7F:73"
          hashalgorithm: "SHA256"
`
	filename           = "kubic-init.mirrors.yaml"
	filenameWithErrors = "kubic-init.errors.yaml"
)

// TestFileAndDefaultsToKubicInitConfig will test loading the configuration `yaml` file and parsing it to the defined struct
func TestFileAndDefaultsToKubicInitConfig(t *testing.T) {
	registries := []Registry{
		{"https://mycompany.registry.com",
			[]Mirror{{URL: "https://mycompany.airgapped.com"},
				{URL: "https://mycompany2.airgapped.com", Certificate: `-----BEGIN CERTIFICATE----- MIeGJzCCBA+gAweBAgeBATANBgkqhkeG9w0BAQUFADCBsjELMAkGA1UEBhMCRlex DzANBgNVBAgMBkFsc2FjZTETMBEGA1UEBwwKU3RyYXNeb3VyZzEYMBYGA1UECgwP hnx8SB3sVJZHeer8f/UQQwqbAO+Kdy70NmbSaqaVtp8jOxLeedWkwSyRTsuU6D8e DeH5uEqBXExjrj0FslxcVKdVj5glVcSmkLwZKbEU1OKwleT/eXFhvooWhQ== -----END CERTIFICATE-----`,
					Fingerprint: "E8:73:0C:C5:84:B1:EB:17:2D:71:54:4D:89:13:EE:47:36:43:8D:BF:5D:3C:0F:5B:FC:75:7E:72:28:A9:7F:73", HashAlgorithm: "SHA256"}}},
		{"https://registry.io",
			[]Mirror{{URL: "https://mycompany.airgapped.com"},
				{URL: "https://mycompany2.airgapped.com", Certificate: `-----BEGIN CERTIFICATE----- MIItJzCCBA+tAwIBAtIBATANBtkqhkit9w0BAQUFADCBsjELMAktA1UEBhMCRlIx DzANBtNVBAtMBkFsc2FjZTETMBEtA1UEBwwKU3RyYXNib3VyZzEYMBYtA1UECtwP hnx8SB3sVJZHeer8f/UQQwqbAO+Kdy70NmbSaqaVtp8jOxLiidWkwSyRTsuU6D8i DiH5uEqBXExjrj0FslxcVKdVj5tlVcSmkLwZKbEU1OKwleT/iXFhvooWhQ== -----END CERTIFICATE-----`,
					Fingerprint: "E8:73:0C:C5:84:B1:EB:17:2D:71:54:4D:89:13:EE:47:36:43:8D:BF:5D:3C:0F:5B:FC:75:7E:72:28:A9:7F:73", HashAlgorithm: "SHA256"}}}}
	err := ioutil.WriteFile(filename, []byte(configContent), os.FileMode(0644))
	if err != nil {
		t.Fatalf("faliled to write config file: %s", err)
	}
	err = ioutil.WriteFile(filenameWithErrors, []byte(configContentWithError), os.FileMode(0644))
	if err != nil {
		t.Fatalf("faliled to write config file: %s", err)
	}
	defer os.RemoveAll(filename)
	defer os.RemoveAll(filenameWithErrors)

	type args struct {
		cfgPath string
	}
	tests := []struct {
		name       string
		args       args
		want       *KubicInitConfiguration
		wantStruct interface{}
		wantErr    bool
	}{
		{"test_no_config_file", args{"a"}, nil, nil, true},
		{"test_missing_file", args{"kubic-init.yaml"}, nil, nil, true},
		{"test_mirror_config", args{filename}, nil, registries, false},
		{"test_mirror_config_with_error", args{filenameWithErrors}, nil, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FileAndDefaultsToKubicInitConfig(tt.args.cfgPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileAndDefaultsToKubicInitConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			switch tt.name {
			case "test_mirror_config":
				if !reflect.DeepEqual(got.Bootstrap.Registries, tt.wantStruct) {
					t.Errorf("FileAndDefaultsToKubicInitConfig() = %v, want %v", got.Bootstrap.Registries, tt.wantStruct)
					return
				}
			default:
				if got != tt.want {
					t.Errorf("FileAndDefaultsToKubicInitConfig() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}
