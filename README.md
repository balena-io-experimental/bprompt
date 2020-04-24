# bprompt

A Small but Useful™ utility to switch between different Balena credentials and API endpoints.  Still a WIP.

# Building

Build with go:

```
go build
```

# Usage

- `bprompt -print`: Print a list of accounts that it knows about

- `bprompt -switch`: Switch to one of those accounts.  This is mean
  to:
  
  - symlink `~/.balena/token` to `~/.balena/token.[account name]`
  
  - update the API endpoint in `~/.balenarc`
  
  **Note:** These features are not done yet.
  
- `bprompt -prompt`: Print out a string that can be incorporated into
  your shell prompt.  In Bash, the idea is that you'd add something
  like this to your `.bashrc`:
  
```
PS1="$BALENA_REMINDER \d \t $"
PROMPT_COMMAND='export BALENA_REMINDER=$(path/to/bprompt)'
```

and then have a nice visual reminder of the account you're using:


```
prod 🔥⚠😑 Fri Apr 24 14:09:32 $
```

# Future improvements

- Actually make the updating work

- Don't hardcode accounts or paths

# License

MIT License; see `LICENSE.md`.
