# minit

![MIT License](https://img.shields.io/github/license/guoyk93/minit)
[![release](https://github.com/guoyk93/minit/actions/workflows/release.yml/badge.svg)](https://github.com/guoyk93/minit/actions/workflows/release.yml)
[![Dockerhub](https://img.shields.io/docker/pulls/guoyk/minit)](https://hub.docker.com/r/guoyk/minit)
[![Patreon Donation](https://img.shields.io/badge/Patreon-Donation-orange)](https://www.patreon.com/guoyk)
[![Buy Me a Coffee Donation](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-Donate-orange)](https://www.buymeacoffee.com/guoyk)

The missing `init` daemon for container

[简体中文](README.zh.md)

## Installation

You can install `minit` to your own container image by a multi-stage `Dockerfile`

```dockerfile
FROM guoyk/minit:VERSION AS minit
# Or using Github Packages
# FROM ghcr.io/guoyk93/minit:VERSION AS minit

# Your own build stage
FROM ubuntu:22.04

# ...

# Copy minit binary
COPY --from=minit /minit /minit

# Set ENTRYPOINT to minit
ENTRYPOINT ["/minit"]

# Add a unit file to /etc/minit.d
ADD my-service.yml /etc/minit.d/my-service.yml
```

## Unit Types

`minit` supports several unit types

### Type: `render`

`render` units execute at the very first stage. It renders template files.

See [pkg/mtmpl/funcs](pkg/mtmpl/funcs) for available functions.

**Example:**

* `/etc/minit.d/render-demo.yaml`

```yaml
kind: render
name: render-demo
files:
  - /opt/*.txt
```

* `/opt/demo.txt`

```text
Hello, {{stringsToUpper .Evn.HOME}}
```

Upon startup, `minit` will render file `/opt/demo.txt`

Since default user for container is `root`, the content of file `/opt/demo.txt` will become:

```text
Hello, ROOT
```

### Type: `once`

`once` units execute after `render` units. It runs command once.

**Example:**

```yaml
kind: once
name: once-demo
command:
  - echo
  - once
```

### Type: `daemon`

`daemon` units execute after `render` and `once`. It runs long-running command.

**Example:**

```yaml
kind: daemon
name: daemon-demo
command:
  - sleep
  - 9999
```

### Type: `cron`

`cron` units execute after `render` and `once`. It runs command at cron basis.

**Example:**

```yaml
kind: cron
name: cron-demo
cron: "* * * * *" # cron expression, support extended syntax by https://github.com/robfig/cron
command:
  - echo
  - cron
```

## Replicas

If `count` field is set, `minit` will replicate this unit with sequence number suffixed

**Example:**

```yaml
kind: once
name: once-demo-replicas
count: 2
command:
  - echo
  - $MINIT_UNIT_SUB_ID
```

Is equal to:

```yaml
kind: once
name: once-demo-replicas-1
command:
  - echo
  - 1
---
kind: once
name: once-demo-replicas-2
command:
  - echo
  - 2
```

## Logging

**Log Files**

`minit` write console logs of every command unit into `/var/log/minit`

This directory can be overridden by environment `MINIT_LOG_DIR`

Set `MINIT_LOG_DIR=none` to disable file logging and optimize performance of `minit`

**Log Encodings**

If `charset` field is set, `minit` will transcode command log from other encodings to `utf8`

**Example:**

```yaml
kind: once
name: once-demo-transcode
charset: gbk # supports gbk, gb18030 only
command:
  - command-that-produces-gbk-logs
```

## Extra Environment Variables

If `env` field is set, `minit` will append extra environment variables while launching command.

**Example:**

```yaml
kind: daemon
name: daemon-demo-env
env:
  AAA: BBB
command:
  - echo
  - $AAA
```

## Render Environment Variables

Any environment with prefix `MINIT_ENV_` will be rendered before passing to command.

**Example:**

```yaml
kind: daemon
name: daemon-demo-render-env
env:
  MINIT_ENV_MY_IP: '{{netResolveIP "google.com"}}'
command:
  - echo
  - $MY_IP
```

## Using `shell`

By default, `command` field will be passed to `exec` syscall with simple environment variable substitution only.

If `shell` field is set, `command` field will act as a simple script file.

**Example:**

```yaml
kind: once
name: once-demo-shell
shell: "/bin/bash -eu"
command: # this is merely a script file
  - if [ -n "${HELLO}" ]; then
  - echo "world"
  - fi
```

## Unit Loading

### From Files

Unit `YAML` files should locate at `/etc/minit.d`

This directory can be overridden by environment variable `MINIT_UNIT_DIR`

Use `---` to separate multiple units in single `YAML` file

### From Environment Variable

```
MINIT_MAIN=redis-server /etc/redis.conf
MINIT_MAIN_DIR=/work
MINIT_MAIN_NAME=main-program
MINIT_MAIN_GROUP=super-main
MINIT_MAIN_KIND=cron
MINIT_MAIN_CRON=* * * * *
MINIT_MAIN_CHARSET=gbk18030
```

### From Arguments

```
ENTRYPOINT ["/minit"]
CMD ["redis-server", "/etc/redis.conf"]
```

## Unit Enabling / Disabling

**Grouping**

Use `group` field to set a group name to units.

Default unit group name is `default`

**Allowlist Mode**

If environment `MINIT_ENABLE` is set, `minit` will run in **Allowlist Mode**, only units with name existed
in `MINIT_ENABLE` will be loaded.

Use format `@group-name` to enable a group of units

Example:

```text
MINIT_ENABLE=once-demo,@demo
```

**Denylist Mode**

If environment `MINIT_DISABLE` is set, `minit` will run in **Denylist Mode**, units with name existed in `MINIT_DISABLE`
will NOT be loaded.

Use format `@group-name` to disable a group of units

Example:

```text
MINIT_DISABLE=once-demo,@demo
```

## Extra Features

### Zombie Processes Cleaning

When running as `PID 1`, `minit` will do zombie process cleaning

This is the responsibility of `PID 1`

### Quick Exit

By default, `minit` will keep running even without `daemon` or `cron` units defined.

If you want to use `minit` in `initContainers` or outside of container, you can set envrionment
variable `MINIT_QUIT_EXIT=true` to let `minit` exit as soon as possible

### Resource limits (ulimit)

**Warning: this feature need container running at Privileged mode**

Use environment variable `MINIT_RLIMIT_XXX` to set resource limits

* `unlimited` means no limitation
* `-` means unchanged

**Supported:**

```text
MINIT_RLIMIT_AS
MINIT_RLIMIT_CORE
MINIT_RLIMIT_CPU
MINIT_RLIMIT_DATA
MINIT_RLIMIT_FSIZE
MINIT_RLIMIT_LOCKS
MINIT_RLIMIT_MEMLOCK
MINIT_RLIMIT_MSGQUEUE
MINIT_RLIMIT_NICE
MINIT_RLIMIT_NOFILE
MINIT_RLIMIT_NPROC
MINIT_RLIMIT_RTPRIO
MINIT_RLIMIT_SIGPENDING
MINIT_RLIMIT_STACK
```

**Example:**

```text
MINIT_RLIMIT_NOFILE=unlimited       # set soft limit and hard limit to 'unlimited'
MINIT_RLIMIT_NOFILE=128:unlimited   # set soft limit to 128，set hard limit to 'unlimited'
MINIT_RLIMIT_NOFILE=128:-           # set soft limit to 128，dont change hard limit
MINIT_RLIMIT_NOFILE=-:unlimited     # don't change soft limit，set hard limit to 'unlimited'
```

### Kernel Parameters (sysctl)

**Warning: this feature need container running at Privileged mode**

Use environment variable `MINIT_SYSCTL` to set kernel parameters

Separate multiple entries with `,`

**Example:**

```
MINIT_SYSCTL=vm.max_map_count=262144,vm.swappiness=60
```

### Transparent Huge Page (THP)

**Warning: this feature need container running at Privileged mode and host `/sys` mounted**

Use environment variable `MINIT_THP` to set THP configuration.

**Example:**

```
# available values: never, madvise, always
MINIT_THP=madvise
```

### Built-in WebDAV server

By setting environment variable `MINIT_WEBDAV_ROOT`, `minit` will start a built-in WebDAV server at port `7486`

Environment Variables:

* `MINIT_WEBDAV_ROOT`, path to serve, `/srv` for example
* `MINIT_WEBDAV_PORT`, port of WebDAV server, default to `7486`
* `MINIT_WEBDAV_USERNAME` and `MINIT_WEBDAV_PASSWORD`, optional basic auth for WebDAV server

### Banner file

By putting a file at `/etc/banner.minit.txt`, `minit` will print it's content at startup

## Donation

View https://guoyk.xyz/donation

## Credits

Guo Y.K., MIT License
