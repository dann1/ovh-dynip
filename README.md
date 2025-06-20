# ovh-dynip

A simple CLI tool to update an A record in your OVH DNS zone to match your current public IP. Useful for home where the ISP allows DMZ access but doesn't provide a static IP. This IP changes over time and it requires to be kept updated on your zone records. This program helps automating that.

## Requirements

- An Application Key and Application Secret configured at `~/.ovh.conf`. See https://github.com/ovh/go-ovh?tab=readme-ov-file#application-keyapplication-secret
- An ovh domain
- golang if you wish to build the binary.
  - For this simply issue `make build` and the executable `ovh-dynip` will appear in `./bin`
  - Run `make install` to have it on `/usr/sbin/`

## Usage

First generate a consumer key for your given key pair, you only need to do this once. For example

```
ovh-dynip -generate-key                                                                                                                         2s
Visit this URL to authorize: https://ca.ovh.com/auth/sso/api?credentialToken=2ab50e5c3d0cc5f06433ee6e91d96a29f454611e17d67f579e6ba1149da0b8cf

Then set this consumer key in your OVH config:
consumer_key=d9da48c33e1e6714cc9d3532f917358f⏎
```

Then set the `consumer_key` on `~/.ovh.conf`. For example, if using OVH Canada

```
[default]
endpoint=ovh-ca

[ovh-ca]
application_key=<your app key>
application_secret=<your app secret>
consumer_key=d9da48c33e1e6714cc9d3532f917358f
```

Then update the record by passing an fqdn to the `--update` argument. For example, with the fqdn `home.example.ovh`

```
ovh-dynip --update home.example.ovh
2025/06/20 17:18:42 Updating home.example.ovh -> 189.162.134.3
2025/06/20 17:18:43 No update needed
```

If will only update if the A record public IP mismatches the public IP of the host where `ovh-dynip` is executed.

### Automate

Example Crontab Entry

To check and update your DNS record every 10 minutes, issue the command `crontab -e` and add

```
*/10 * * * * ovh-dynip -update home.example.ovh >> /var/log/ovh-dynip.log 2>&1
```
