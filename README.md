# celigo-cli
Celigo Cli allows you to interact with the celigo platform, through your terminal

## How I would like to interact with the cli
I would like the use to be able to enter the celigo command and then choose your resource
type (script, export, import, etc...) then enter the action. For example, the script actions would
look something like the following:

```bash
celigo script list
celigo script fetch -i <ID> [ -o <Path/To/Output File Name>]
celigo script update -i <ID> [ -c <Path/To/Source File Name> -n <Name> -d <Description> ]
celigo script add -n <Name> [ -d <Description> -c <Path/To/Source File Name>]
celigo script remove -i <ID>
```
I will read the api key from the users home directory in the file ".celigo". This file will contain
the api key only in one line. All commands can take an optional -a option that contains the path to
the file that contains the api key in it. This allows the user to switch between celigo accounts
seemlessly, in their "Project Directory"

## MVP
- Read api key from users home directory or custom path
- Script resource working for the options listed above
