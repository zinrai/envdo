# envdo

A command-line tool that loads environment variables from `.env` or `.envrc` files and executes commands in a subshell.

## Features

- Automatically detects and loads `.env` or `.envrc` files in the current directory
- Allows specifying a custom environment file with the `-f` flag
- Executes commands in a subshell with the loaded environment variables
- Supports quoted values in environment files

## Installation

```bash
$ go build
```

## Usage

Use .env or .envrc from current directory:

```bash
$ envdo command [args...]
```

Specify a custom environment file:

```bash
$ envdo -f .env.local command [args...]
```

### Examples

Run 'echo $HELLO' with variables from .env:

```bash
$ envdo bash -c 'echo $HELLO'
```

Run a script with environment variables from .env.production:

```bash
$ envdo -f .env.production ./deploy.sh
```

## Environment File Format

The environment file uses a simple key-value format:

```
# This is a comment
KEY=value
SECRET="quoted value"
PASSWORD='also quoted'
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
