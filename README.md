# Siridb-Prompt


## Connect to SiriDB via SSH Tunnel in Linux


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
[1]+ Stopped        ssh login@server -L 9000:
127.0.0.1:3306 -N
```
to set it in background, type `bg` and you should get response:
```bash
[1]+ ssh login@server -L 3306:127.0.0.1:3306 -N &
```
Connect to SiriDB as if it is running on your local system (localhost:9000) for example by running the following command in the terminal window:
```bash
$ siridb-prompt -p -u user -d database
```
That's it! You are in.

After you are finished close application you've used to access SiriBD (for command line tool type Ctrl+D to leave).
Then return to your SSH tunnel:
get it into foreground: press `fg` and you will get the response:  
```bash
ssh login@server -L 3306:127.0.0.1:3306 -N
```
press `Ctrl+C` to stop it.
