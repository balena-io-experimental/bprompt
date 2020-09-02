# bprompt

A Small but Useful™ utility to switch between different Balena credentials and API endpoints.  Still a WIP.

# Building

Build with go:

```
go build
```

# Usage

- `bprompt -print`: Print a list of accounts that it knows about

- `bprompt -show`: Show the current state of your token and API endpoint

- `bprompt -switch [account]`: Switch to one of those accounts.  This will
  symlink `~/.balena/token` to `~/.balena/token.[account name]`, and
  update the API endpoint in `~/.balenarc`.  **Note:** It will refuse
  to proceed if `~/.balena/token` is not a symlink.

- `bprompt -switchForce [account]`: The same as `-switch`, but don't
  worry about symlinks -- it will delete `~/.balena/token`, and
  re-create the symlink.

- `bprompt -prompt`: Print out a string that can be incorporated into
  your shell prompt.  In Bash, the idea is that you'd add something
  like this to your `.bashrc`:

```
PROMPT_COMMAND='export BALENA_REMINDER=$(path/to/bprompt -prompt)'
PS1="$BALENA_REMINDER \d \t $"
```

and then have a nice visual reminder of the account you're using:

```
prod 🔥⚠😑 Fri Apr 24 14:09:32 $
```

# Future improvements

- Add configuration file for paths and tokens

# License

MIT License; see `LICENSE.md`.
