# CLI for Supergloo
## Installation
```
make install-cli    # run from the project root directory
```
## Commands

### Help
Lists the available commands.
#### Usage
```
supergloo help
```

### Get
Displays one or many supergloo resources in table format.
#### Usage
```
supergloo get RESOURCE_TYPE [RESOURCE_NAME] [(-o|--output) wide]
```
#### Flags
- `output`: output format. Currently only the `wide` option, which 
causes additional columns to be displayed, is supported.