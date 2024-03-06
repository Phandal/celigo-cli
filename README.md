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

The API Key will be read from the current working direcotry in the ".celigo" file. This file will
contain the api key on one line, and thats it.

## MVP
- Read api key from users home directory or custom path
- Script resource working for the options listed above
