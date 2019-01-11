package config

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/golang/glog"
)

const (
	// DefaultEnvVarSeeder The environment variable used for passing the seeder
	DefaultEnvVarSeeder = "SEEDER"

	// DefaultEnvVarToken The environment variable used for passing the token
	DefaultEnvVarToken = "TOKEN"

	// DefaultEnvVarManager The environment variable used for passing the kubic-manager image
	DefaultEnvVarManager = "MANAGER_IMAGE"

	// DefaultAPIServerPort Default API server port
	DefaultAPIServerPort = 6443
)

// CniConfiguration The CNI configuration
// Subnets details are specified in the kubeadm configuration file
type CniConfiguration struct {
	BinDir  string `yaml:"binDir,omitempty"`
	ConfDir string `yaml:"confDir,omitempty"`
	Driver  string `yaml:"driver,omitempty"`
	Image   string `yaml:"image,omitempty"`
}

// ClusterFormationConfiguration struct
type ClusterFormationConfiguration struct {
	Seeder      string `yaml:"seeder,omitempty"`
	Token       string `yaml:"token,omitempty"`
	AutoApprove bool   `yaml:"autoApprove,omitempty"`
}

// OIDCConfiguration struct
type OIDCConfiguration struct {
	Issuer   string `yaml:"issuer,omitempty"`
	ClientID string `yaml:"clientID,omitempty"`
	CA       string `yaml:"ca,omitempty"`
	Username string `yaml:"username,omitempty"`
	Groups   string `yaml:"groups,omitempty"`
}

// AuthConfiguration struct
type AuthConfiguration struct {
	OIDC OIDCConfiguration `yaml:"OIDC,omitempty"`
}

// CertsConfiguration struct
type CertsConfiguration struct {
	// TODO
	Directory string `yaml:"directory,omitempty"`
	CaHash    string `yaml:"caCrtHash,omitempty"`
}

// DNSConfiguration struct
type DNSConfiguration struct {
	Domain       string `yaml:"domain,omitempty"`
	ExternalFqdn string `yaml:"externalFqdn,omitempty"`
}

// ProxyConfiguration struct
type ProxyConfiguration struct {
	HTTP       string `yaml:"http,omitempty"`
	HTTPS      string `yaml:"https,omitempty"`
	NoProxy    string `yaml:"noProxy,omitempty"`
	SystemWide bool   `yaml:"systemWide,omitempty"`
}

// BindConfiguration struct
type BindConfiguration struct {
	Address   string `yaml:"address,omitempty"`
	Interface string `yaml:"interface,omitempty"`
}

// PathsConfigration struct
type PathsConfigration struct {
	Kubeadm string `yaml:"kubeadm,omitempty"`
}

// LocalEtcdConfiguration struct
type LocalEtcdConfiguration struct {
	ServerCertSANs []string `yaml:"serverCertSANs,omitempty"`
	PeerCertSANs   []string `yaml:"peerCertSANs,omitempty"`
}

// EtcdConfiguration struct
type EtcdConfiguration struct {
	LocalEtcd *LocalEtcdConfiguration `yaml:"local,omitempty"`
}

// NetworkConfiguration struct
type NetworkConfiguration struct {
	Bind          BindConfiguration  `yaml:"bind,omitempty"`
	Cni           CniConfiguration   `yaml:"cni,omitempty"`
	DNS           DNSConfiguration   `yaml:"dns,omitempty"`
	Proxy         ProxyConfiguration `yaml:"proxy,omitempty"`
	PodSubnet     string             `yaml:"podSubnet,omitempty"`
	ServiceSubnet string             `yaml:"serviceSubnet,omitempty"`
}

// RuntimeConfiguration struct
type RuntimeConfiguration struct {
	Engine string `yaml:"engine,omitempty"`
}

// FeaturesConfiguration struct
type FeaturesConfiguration struct {
	PSP bool `yaml:"PSP,omitempty"`
}

// ServicesConfiguration struct
type ServicesConfiguration struct {
}

// BootstrapConfiguration se the required configuration
// for bootstraping kubic-init
type BootstrapConfiguration struct {
	Registries []Registry `yaml:"registries,omitempty"`
}

// Registry struct
// Defines a registry mirror
// Prefix: string of the registry that will be replaced
// Mirrors: array with the values to replace the `Prefix`
type Registry struct {
	Prefix  string   `yaml:"prefix"`
	Mirrors []Mirror `yaml:"mirrors"`
}

// Mirror struct
// Defines the Mirrors to be used
// URL: url of the mirror registry.
// Certificate: certificate content for the registry.
// Fingerprint: fingerprint of the certificate to check validity.
// HashAlgorithm: hash algorithm used.
type Mirror struct {
	URL           string `yaml:"url"`
	Certificate   string `yaml:"certificate,omitempty"`
	Fingerprint   string `yaml:"fingerprint,omitempty"`
	HashAlgorithm string `yaml:"hashalgorithm,omitempty"`
}

// KubicInitConfiguration The kubic-init configuration
type KubicInitConfiguration struct {
	Network          NetworkConfiguration          `yaml:"network,omitempty"`
	Paths            PathsConfigration             `yaml:"paths,omitempty"`
	ClusterFormation ClusterFormationConfiguration `yaml:"clusterFormation,omitempty"`
	Certificates     CertsConfiguration            `yaml:"certificates,omitempty"`
	Etcd             EtcdConfiguration             `yaml:"etcd,omitempty"`
	Runtime          RuntimeConfiguration          `yaml:"runtime,omitempty"`
	Features         FeaturesConfiguration         `yaml:"features,omitempty"`
	Services         ServicesConfiguration         `yaml:"services,omitempty"`
	Auth             AuthConfiguration             `yaml:"auth,omitempty"`
	Bootstrap        BootstrapConfiguration        `yaml:"bootstrap,omitempty"`
}

// FileAndDefaultsToKubicInitConfig Load a Kubic configuration file, setting some default values
func FileAndDefaultsToKubicInitConfig(cfgPath string) (*KubicInitConfiguration, error) {
	var err error

	internalcfg := &KubicInitConfiguration{}

	if len(cfgPath) > 0 {
		glog.V(1).Infof("[caasp-init] loading kubic-init configuration from '%s'", cfgPath)
		if _, err = os.Stat(cfgPath); err != nil {
			return nil, fmt.Errorf("%q does not exist: %v", cfgPath, err)
		}

		b, err := ioutil.ReadFile(cfgPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read config from %q [%v]", cfgPath, err)
		}

		if err = yaml.Unmarshal(b, &internalcfg); err != nil {
			return nil, fmt.Errorf("unable to decode config from bytes: %v", err)
		}
	}

	return internalcfg, nil
}
