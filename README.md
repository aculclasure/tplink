##  :beginner: About
tplink provides a *minimal* CLI (implemented in Golang) for interacting with a
TP Link Archer C9v1 home wifi router. My goals with this project were to:
- learn how to write a basic API client in Golang
- learn how to write a basic command-line interface in Golang
- cross-compile the tplink application to run on a Raspberry Pi 3

##  :electric_plug: Installation
```
go get -u github.com/aculclasure/tplink
```
This installs the `tplink` binary in your `$GOPATH/bin` location and installs
the entire tplink module in your `$GOPATH/go/pkg/mod` location.

##  :zap: Usage
```
$ tplink --help

tplink is a CLI app that provides a very minimal admin interface to a TP Link
Archer C9 V1 home wifi router. This app allows a user to retrieve information
about connected clients on the router and to also reboot the router.

Usage:
  tplink [command]

Available Commands:
  help        Help about any command
  list        lists information about the router
  reboot      reboots the router
  version     displays the version and exits

Flags:
  -h, --help   help for tplink

Use "tplink [command] --help" for more information about a command.
```

### Example: Listing connected wireless clients
```
$ tplink list wirelessClients --help
wirelessClients queries the wifi router to get the currently connected wireless clients and
prints out the IP address, MAC address, and host name (if known) for each wireless client.

Usage:
  tplink list wirelessClients [flags]

Flags:
  -h, --help   help for wirelessClients

Global Flags:
  -P, --password string   router admin password (required) (default "admin")
      --url string        router URL (required) (default "http://192.168.168.1")
  -U, --user string       router admin user name (required) (default "admin")


$ tplink list wirelessClients --url http://some-tp-link-router --user <ROUTER-USER-NAME> --password <ROUTER-PASSWORD>
INFO: 2020/09/12 14:18:25 connections.go:67: sending request to get wireless connections as (POST http://some-tp-link-rtr/data/map_access_wireless_client_grid.json) ...
IP_ADDRESS         MAC_ADDRESS           HOST_NAME      
10.100.100.100     12-34-56-AA-BB-CC     FakeHost
10.100.100.55      1A-1A-1A-AA-AA-AA     Unknown         
```

##  :computer: Running the tplink client on a Raspberry Pi 3!
One of the cool things I learned about Go when working on this project is that it
supports cross-compilation for a number of different target architectures right out
of the box! To produce a binary that will work on a Raspberry Pi 3, you can do this:
```shell script
$ env GOOS=linux GOARCH=arm GOARM=6 go build
```

##  :page_facing_up: Resources
- [Writing a Go client for your RESTful API](https://medium.com/@marcus.olsson/writing-a-go-client-for-your-restful-api-c193a2f4998c): Blog post
- [Creating an API Client in Go](https://www.scaledrone.com/blog/creating-an-api-client-in-go/): Blog post
- [Godo](https://github.com/digitalocean/godo): Reference source code for a client for 
  accessing the DigitalOcean V2 API
- [Cobra](https://github.com/spf13/cobra): Library for creating command line applications in Go
- [Cross Compiling Golang Applications For Use On A Raspberry Pi](https://www.thepolyglotdeveloper.com/2017/04/cross-compiling-golang-applications-raspberry-pi/): Blog post
