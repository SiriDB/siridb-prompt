# SiriDB Prompt
![alt SiriDB Prompt](/siridb-prompt.png?raw=true)

SiriDB Prompt provides a terminal client for communicating with [SiriDB](https://github.com/SiriDB/siridb-server).

> Note: Since version 2.1.0 SiriDB Prompt is written in Go. For the most recent version in Python you should go to
> this [release tag](https://github.com/SiriDB/siridb-prompt/tree/2.0.6).

---------------------------------------
  * [Installation](#installation)
    * [Pre-compiled](#pre-compiled)
    * [Compile from source](#compile-from-source)
  * [Connect via SSH](#connect-via-ssh)
---------------------------------------

## Installation
SiriDB Prompt can be compiled from source or, for most systems, you can simply download a pre-compiled binary.

### Pre-compiled
Go to https://github.com/SiriDB/siridb-prompt/releases/latest and download the binary for your system.
In this documentation we refer to the binary as `siridb-prompt`. On Linux/OSX it might be required to set the execution flag:
```
$ chmod +x siridb-prompt_X.Y.Z_OS_ARCH.bin
```

You might want to copy the binary to /usr/local/bin and create a symlink like this:
```
$ sudo cp siridb-prompt_X.Y.Z_OS_ARCH.bin /usr/local/bin/
$ sudo ln -s /usr/local/bin/siridb-prompt_X.Y.Z_OS_ARCH.bin /usr/local/bin/siridb-prompt
```
> Note: replace `X.Y.Z_OS_ARCH` with your binary, for example `2.0.0_linux_amd64`

### Compile from source
> Before compiling from source make sure **go** and **git** are installed and your [$GOPATH](https://github.com/golang/go/wiki/GOPATH) is set.

Clone the project using git. (we assume git is installed)
```
git clone https://github.com/SiriDB/siridb-prompt
```

Install the required go packages
```
$ go get -d
```

Build SiriDB Prompt
```
$ go build -o siridb-prompt
```

## Connect via SSH

Build SSH tunnel by running the following command in the terminal window:
```bash
$ ssh login@server -L 9000:127.0.0.1:9000 -N
```

You will get response asking you to enter your SSH password
```bash
login@server password:
```

Now when the tunnel is built stop the previous command and set it into background to be able to initiate other commands.
press `Ctrl+Z` to stop the command and you should get response:
```bash
[1]+ Stopped        ssh login@server -L 9000:127.0.0.1:9000 -N
```
to set it in background, type `bg` and you should get response:
```bash
[1]+ ssh login@server -L 9000:127.0.0.1:9000 -N &
```
Connect to SiriDB as if it is running on your local system, for example by running the following command in the terminal window:
```bash
$ siridb-prompt -d database -u user -s localhost:9000
```
That's it! You are in.

When finished you can close your SSH tunnel:

get it into foreground: press `fg` and you will get the response:
```bash
ssh login@server -L 9000:127.0.0.1:9000 -N
```
press `Ctrl+C` to stop it.
