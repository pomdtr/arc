# Reference

## arc

Arc Companion CLI

### Options

```
  -h, --help   help for arc
```

## arc completion

Generate the autocompletion script for the specified shell

### Synopsis

Generate the autocompletion script for arc for the specified shell.
See each sub-command's help for details on how to use the generated script.


### Options

```
  -h, --help   help for completion
```

## arc completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(arc completion bash)

To load completions for every new session, execute once:

#### Linux:

	arc completion bash > /etc/bash_completion.d/arc

#### macOS:

	arc completion bash > $(brew --prefix)/etc/bash_completion.d/arc

You will need to start a new shell for this setup to take effect.


```
arc completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

## arc completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	arc completion fish | source

To load completions for every new session, execute once:

	arc completion fish > ~/.config/fish/completions/arc.fish

You will need to start a new shell for this setup to take effect.


```
arc completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

## arc completion help

Help about any command

### Synopsis

Help provides help for any command in the application.
Simply type completion help [path to command] for full details.

```
arc completion help [command] [flags]
```

### Options

```
  -h, --help   help for help
```

## arc completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	arc completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
arc completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

## arc completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(arc completion zsh)

To load completions for every new session, execute once:

#### Linux:

	arc completion zsh > "${fpath[1]}/_arc"

#### macOS:

	arc completion zsh > $(brew --prefix)/share/zsh/site-functions/_arc

You will need to start a new shell for this setup to take effect.


```
arc completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

## arc help

Help about any command

### Synopsis

Help provides help for any command in the application.
Simply type arc help [path to command] for full details.

```
arc help [command] [flags]
```

### Options

```
  -h, --help   help for help
```

## arc history

Search history

```
arc history [flags]
```

### Options

```
  -h, --help           help for history
      --json           output as json
  -l, --limit int      limit (default 100)
  -q, --query string   query
```

## arc space

Manage spaces

### Options

```
  -h, --help   help for space
```

## arc space focus

Focus a space

```
arc space focus [flags]
```

### Options

```
  -h, --help   help for focus
```

## arc space help

Help about any command

### Synopsis

Help provides help for any command in the application.
Simply type space help [path to command] for full details.

```
arc space help [command] [flags]
```

### Options

```
  -h, --help   help for help
```

## arc space list

List spaces

```
arc space list [flags]
```

### Options

```
  -h, --help   help for list
      --json   output as json
```

## arc tab

Manage tabs

### Options

```
  -h, --help   help for tab
```

## arc tab close

Close a tab

```
arc tab close [flags]
```

### Options

```
  -h, --help   help for close
```

## arc tab create

Create a new tab.

```
arc tab create <url> [flags]
```

### Options

```
  -h, --help        help for create
      --little      open in little arc
      --space int   space to create tab in
```

## arc tab exec

Execute javascript in the active tab

```
arc tab exec <script> [flags]
```

### Options

```
  -e, --eval string   javascript to evaluate
  -h, --help          help for exec
```

## arc tab focus

Select a tab by id

```
arc tab focus <tab-id> [flags]
```

### Options

```
  -h, --help   help for focus
```

## arc tab get

Get information about the active tab

### Options

```
  -h, --help   help for get
```

## arc tab get help

Help about any command

### Synopsis

Help provides help for any command in the application.
Simply type get help [path to command] for full details.

```
arc tab get help [command] [flags]
```

### Options

```
  -h, --help   help for help
```

## arc tab get title

Get the title of the active tab

```
arc tab get title [flags]
```

### Options

```
  -h, --help   help for title
```

## arc tab get url

Get the url of the active tab

```
arc tab get url [flags]
```

### Options

```
  -h, --help   help for url
```

## arc tab help

Help about any command

### Synopsis

Help provides help for any command in the application.
Simply type tab help [path to command] for full details.

```
arc tab help [command] [flags]
```

### Options

```
  -h, --help   help for help
```

## arc tab list

List tabs

```
arc tab list [flags]
```

### Options

```
      --favorite   only show favorite tabs
  -h, --help       help for list
      --json       output as json
      --pinned     only show pinned tabs
      --unpinned   only show unpinned tabs
```

## arc tab reload

Reload a tab"

```
arc tab reload [flags]
```

### Options

```
  -h, --help   help for reload
```

## arc version

Print the version of Arc

```
arc version [flags]
```

### Options

```
  -h, --help   help for version
```

## arc window

Manage windows

### Options

```
  -h, --help   help for window
```

## arc window close

Close a window

```
arc window close [flags]
```

### Options

```
  -h, --help   help for close
```

## arc window create

Create a new window

```
arc window create [url] [flags]
```

### Options

```
  -h, --help        help for create
      --incognito   open in incognito mode
```

## arc window help

Help about any command

### Synopsis

Help provides help for any command in the application.
Simply type window help [path to command] for full details.

```
arc window help [command] [flags]
```

### Options

```
  -h, --help   help for help
```

## arc window list

List windows

```
arc window list [flags]
```

### Options

```
  -h, --help   help for list
      --json   output as json
```


