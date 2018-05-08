# gnotail

  A light-weight inotify based tailer and pattern matcher for tailing and pattern matching multiple files in parallel.
  Handles logrotation and file truncation. 
  
##### Requirement:

Linux kernel 2.6.13 or higher

##### How to install:

`go get -u github.com/nohupped/gnotail`

##### How to run:

```
Usage of /tmp/go-build911163385/command-line-arguments/_obj/exe/main:
  -conf string
    	path to configuration of json format. Example format would be 
{
   "/var/log/syslog": [
       "cron",
       "su.*"
   ],
   "/var/log/auth.log": [
       "sudo",
       "pam_unix"
   ]
}
  -loglevel int
    	Log level when printing to STDOUT. ErrorLevel: 0, WarnLevel: 1, InfoLevel: 2, DebugLevel: 3. Defaults to InfoLevel. (default 2)
  -port int
    	Portnumber to which the output be sent (default 9999)
  -udp_conn_addr string
    	address to start udp server (default "127.0.0.1")


```

The above conf file will make gnotail tail the logs `/var/log/syslog` and `/var/log/auth.log` in 2 separate goroutines and does a regex pattern matching for each of the patterns and writes to a `UDP` port as json data.

######An example message read with `netcat -ul 9999`:
```
{"filename":"/var/log/auth.log","rule":"sudo","message":"May  9 01:28:33 localghost sudo:  someuser : TTY=pts/2 ; PWD=/home/someuser ; USER=root ; COMMAND=/usr/bin/whoami"}{"filename":"/var/log/auth.log","rule":"sudo","message":"May  9 01:28:33 localghost sudo: pam_unix(sudo:session): session opened for user root by (uid=0)"}{"filename":"/var/log/auth.log","rule":"pam_unix","message":"May  9 01:28:33 localghost sudo: pam_unix(sudo:session): session opened for user root by (uid=0)"}{"filename":"/var/log/auth.log","rule":"sudo","message":"May  9 01:28:33 localghost sudo: pam_unix(sudo:session): session closed for user root"}{"filename":"/var/log/auth.log","rule":"pam_unix","message":"May  9 01:28:33 localghost sudo: pam_unix(sudo:session): session closed for user root"}


```
