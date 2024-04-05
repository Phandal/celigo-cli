**NOTE: This project is still under active development and things will break!**

# celigo-cli
Celigo Cli allows you to interact with the celigo platform, through your terminal

## API Key
Your API Key is read from the following locations in order:
1. `.celigo-cli` file in the current directory. This file will list the api key on the first line.
2. `CELIGO_API_KEY` environment variable.

## Usage
Currently, only the script resource type is supported
```bash
celigo script list
celigo script fetch -i <ID> [ -o <Path/To/Output/Directory]
celigo script update -i <Path/To/Source File>
celigo script add -n <ScriptName>
celigo script remove -i <ID>
```
