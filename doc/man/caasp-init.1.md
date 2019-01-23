% caasp-init(1) # caasp-init - Initialize docker daemon mirrors configuration
% SUSE LLC
% JANUARY 2019
# NAME
caasp-init - Initialize docker daemon mirrors configuration

# SYNOPSIS
**caasp-init**
[**--help**|**-h**]
[**version**]
[**--config**|**-c**]

# DESCRIPTION
**caasp-init** will create the daemon.json configuration file and the necessary certificates for the mirror you want
to config based on the kubic-init.yaml configuration file.

usage:

`$ caasp-init`

This will use default value for configuration file

`/etc/kubic/kubic-init.yaml`

if you want to indicate a file run

`$ caasp-init -c /path/to/config/file.yaml`

If the configuration file has mirrors declared ot will generate the daemon.json
file with the following structure:

```
{
  "registries": [
    {
    "Mirrors": [
      {
      "URL": "https://airgappedregistry.com"
      }
    ],
    "Prefix": "https://mycompany.registry.com"
    }
  ],
  "iptables":false,
  "log-level": "warn"
}
```

If there is no mirror declared the configuration file will just be the default:

```
{
  "iptables":false,
  "log-level": "warn"
}
```

For help use `caasp-init help`

# GLOBAL OPTIONS

**-h, --help**
  Print usage statement.

**-c, --config**
  kubibc-init.yaml config file (default "/etc/kubic/kubic-init.yaml")

# COMMANDS

**version**
  Print current version of software. See **caasp-init-version**(1) for more detailed
  usage information.

**help**
  Print usage statements. See **caasp-init-help**(1)
  for more detailed usage information.

# SEE ALSO
**caasp-init-help**(1),
**caasp-init-version**(1)

[1]: https://docs.helm.sh
