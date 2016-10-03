# cloudflare-dyndns

This project is a total work in progress right now.

### Example Usage (in the future)

#### Help:

```bash
$ dyndns --help
dyndns 0.1.0
usage: dyndns <command> [<options>]

Available commands:
   <command>
   cloudflare      Run the Cloudflare DNS updater.
   ip              Check what your public IP is.
   config          Use given config file to run all defined DNS updaters.
   help            Show this message.

See 'dyndns help <command>' for more information on a specific command.
```

```bash
$ dyndns help cloudflare
dyndns 0.1.0
usage: dyndns cloudflare <host> [<ip-address>] [<options>]

Arguments:
   host          Hostname to update on Cloudflare. Multiple hostnames are
                 supperted by comma separating them.
   ip-address    Manually specify the IP address to set DNS entries to. When
                 not given your public IP address is automatically detected.

Available options:
   --email / -e         Cloudflare account email (required).
   --key / -k           Cloudflare API key (required).
   --config / -c        Config file, can be used to specify email and API key.
   --ip-checker / -i    Specify which IP checker to use.
```

#### Update a single DNS record on Cloudflare:

```bash
$ dyndns cloudflare foo.bar.com --email foo@bar.com --key abc
```

#### Update multiple DNS records on Cloudflare:

```bash
$ dyndns cloudflare "foo.bar.com,baz.bar.com" --email foo@bar.com --key abc
```

#### Update Cloudflare DNS record with specific IP address:

```bash
$ dyndns cloudflare foo.bar.com 123.123.123.123 --email foo@bar.com --key abc
```

#### Use a config file (runs all configured DNS updaters with all their hosts):

```bash
echo '
{
  "cloudflare": {
    "email": "foo@bar.com",
    "key": "abc",
    "hosts": [
      "foo.bar.com",
      "baz.bar.com"
    ]
  }
}' > ~/.dyndns.json

$ dyndns config "~/.dyndns.json"
```

#### Specify public IP lookup service:

```bash
# use http://whatismyip.akamai.com/
$ dyndns cloudflare [...] --ip-checker akamai
# use http://icanhazip.com/
$ dyndns cloudflare [...] --ip-checker icanhazip
# use custom checker
$ dyndns cloudflare [...] --ip-checker "http://myip.foo.com/"
```
