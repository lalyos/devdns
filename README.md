
this repo is a fork of [robbiev/devdns](https://github.com/robbiev/devdns), wrapped into a docker images


## Usage - osx
One-liner:
```
eval "$(docker run --rm lalyos/devdns --eval)"
```

if you want to see what are the steps run this:
```
docker run --rm lalyos/devdns --eval
```

- configure DNS: creates a `/etc/resolver/dev` file [see details](https://github.com/robbiev/devdns#how)
- starts the container listening on `0.0.0.0:5300`
- test the setup with `ping whatever.dev`

## Usage - linux
On linux `/etc/resolv.conf` is used to configre DNS servers. Unfortunately you can't specify an alternative port number, so you have to expose to port 53.

Start the container:
```
docker run \
  -d \
  --name devdns \
  -p 53:53/udp lalyos/devdns \
    -addr 0.0.0.0:53
```

add a new line to your `/etc/resolv.conf`
```
nameserver 127.0.0.1
```


## Help
To see all options
```
docker run --rm lalyos/devdns --help
```

## Cleanup

Kill the container:
``` 
docker rm -f devdns
```

For osx remove `/etc/resolver` :
```
sudo rm -rf /etc/resolver/
```

For linux remove the line from `/etc/resolv.conf`

```
nameserver 127.0.0.1
```