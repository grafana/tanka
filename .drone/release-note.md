This is release ${DRONE_TAG} of Tanka (`tk`). Check out the [CHANGELOG](CHANGELOG.md) for detailed release notes.
## Install instructions

#### Binary:
```bash
# download the binary (adapt os and arch as needed)
$ curl -fSL -o "/usr/local/bin/tk" "https://github.com/sh0rez/tanka/releases/download/${DRONE_TAG}/tk-linux-amd64"

# make it executable
$ chmod a+x "/usr/local/bin/tk"

# have fun :)
$ tk --help
```

#### Docker container:
https://hub.docker.com/r/grafana/tanka
```bash
$ docker pull grafana/tanka:${DRONE_TAG#v}
```
