package certs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/kubic-project/caasp-init/pkg/config"
)

var (
	certsFolder = "/etc/docker/certs.d"
	certName    = "ca.crt"
)

// WriteCertificates this method will
func WriteCertificates(config *config.KubicInitConfiguration) error {
	if config == nil {
		return errors.New("configuration is nil")
	}
	for _, reg := range config.Bootstrap.Registries {
		if reg.Prefix == "" {
			continue
		}
		for _, mirror := range reg.Mirrors {
			if mirror.Certificate == "" {
				continue
			}
			url, err := url.Parse(mirror.URL)
			if err != nil {
				return err
			}

			if url.Scheme == "" {
				return fmt.Errorf("Error in configuration file: malformed Mirror URL \"%s\"", url.String())
			}

			err = os.MkdirAll(path.Join(certsFolder, url.Hostname()), 0744)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(path.Join(certsFolder, url.Hostname(), certName), []byte(mirror.Certificate), os.FileMode(0644))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
