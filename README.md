# waitfordircontents
Simplified/specialized inotifywait

```
NAME:
   waitfordircontents - Wait until given directories are not empty

USAGE:
   waitfordircontents [global options] command [command options] [arguments...]

AUTHOR:
   Odd Eivind Ebbesen <oddebb@gmail.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --env-var variable, -e variable  Environment variable to get paths from. Values should be colon separated.
   --directories value, -d value    Directories to watch. Separate by commas, or specify multiple times.  (accepts multiple inputs)
   --timeout value, -t value        How long to wait before giving up. 0 means wait forever. (default: 0s)
   --exit-on-watch-failure, -x      Exit with error if a watch failure happens (default: false)
   --help, -h                       show help (default: false)

COPYRIGHT:
   (C) 2022 Odd Eivind Ebbesen

```
