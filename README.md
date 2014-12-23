# Go Build in Docker

[gobldock] provides an easy and convenient way to cross-compile Go packages for Docker (Linux x64_86) environments. Internally it does the following: 
 - loads a `golang` docker image with Go and GCC compilers present, so CGO is enabled;
 - mounts the `$GOPATH` from the host;
 - mounts the output directory (`pwd` by default) from the host;
 - runs `go build -i` inside this docker image.

[gobldock]: http://github.com/missionMeteora/gobldock

## Usage

Requires an installed [docker](https://docs.docker.com).

```
./gobldock -h
Usage: ./gobldock <package>
  -o, --output="."      Specify the path for resulting exectuable
  -s, --silent=false    Be silent, no fancy stuff
```

## Example

```
xlab ~/Documents/dev/meteora $ go get github.com/missionMeteora/gobldock
xlab ~/Documents/dev/meteora $ boot2docker shellinit

xlab ~/Documents/dev/meteora $ gobldock github.com/missionMeteora/gobldock
found Docker: /usr/local/bin/docker
Docker version 1.3.2, build 39fa2fa
compiling into: /Users/xlab/Documents/dev/meteora/gobldock

xlab ~/Documents/dev/meteora $ file gobldock
gobldock: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, not stripped
```

## License

[MIT](/LICENSE)
