**NOTE: This project is still under active development and things will break!**

# celigo-cli
Celigo Cli allows you to interact with the celigo platform, through your terminal

## API Key
Your API Key is read from the following locations in order:
1. `.env` file in the current directory.
2. `CELIGO_API_KEY` environment variable.


## Example .env file
```text
a530749681df4d32a5fcc65a166c3da0
```

## TODO
- tests throughout
- manage .celigo-cli file for file-id mapping instead of using filename
- toml parser for .celigo-cli
