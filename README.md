**NOTE: This project is still under active development and things will break!**

# celigo-cli
Celigo Cli allows you to interact with the [Celigo Platform](https://integrator.io/), through your
terminal

## API Key
You will need to get your api key from Celigo. You can do this by:
1. Going to [https://integrator.io/accesstokens](https://integrator.io/accesstokens)
2. Click on "Create API token" in the upper right
3. Fill in a Name
4. For Auto purge token select Never
5. For scope select Full.
6. Click Save

After receiving your api key, you can use by setting in your environment variables with the name
`CELIGO_API_KEY`, or using a `.env` file.

### Example .env file
```text
CELIGO_API_KEY=a630749681df4d32a5fcc65a166c3da0
```

## TODO
- tests throughout
- manage .celigo-cli file for file-id mapping instead of using filename
- toml parser for .celigo-cli
