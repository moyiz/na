```text
    _______   ______  
   |       \ |      \ 
   | ▓▓▓▓▓▓▓\ \▓▓▓▓▓▓\
   | ▓▓  | ▓▓/      ▓▓
   | ▓▓  | ▓▓  ▓▓▓▓▓▓▓
   | ▓▓  | ▓▓\▓▓    ▓▓
    \▓▓   \▓▓ \▓▓▓▓▓▓▓
```
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/moyiz/na)
[![Go Reference](https://pkg.go.dev/badge/github.com/moyiz/na.svg)](https://pkg.go.dev/github.com/moyiz/na)
![GitHub License](https://img.shields.io/github/license/moyiz/na)
![GitHub release (with filter)](https://img.shields.io/github/v/release/moyiz/na)


**na** (aka: _nested-aliases_, _non-aliases_, _not-[an]-alias_) is a CLI tool to effortlessly manage context aware nested shortcuts for shell commands.

<!-- Demo here -->

## Contents
- [Motivation](#motivation)
- [Features](#features)
- [Installation](#installation)
    - [Binaries](#binaries)
    - [Source](#source)
    - [AUR](#aur)
    - [Build](#build)
- [Usage](#usage)
    - [Adding shortcuts](#adding-shortcuts)
    - [Listing shortcuts](#listing-shortcuts)
    - [Running shortcuts](#running-shortcuts)
    - [Removing shortcuts](#removing-shortcuts)
- [Shell Completions](#shell-completions)
    - [Bash](#bash)
    - [Zsh](#zsh)
    - [Fish](#fish)
    - [Powershell](#powershell)
- [Configuration](#configuration)
    - [Showcase](#showcase)
    - [Example](#example)
- [Known Issues](#known-issues)
- [Future Plans](#future-plans)
- [License](#license)

## Motivation
Shell aliases are fun. They provide an easy and straightforward way to create simple shortcuts without much fuss, but not without few caveats:
- *Aliases must be named with a single word* - It makes grouping related aliases within the same context a bit awkward, e.g:
    ```sh
    alias lab-password-sftpgo="k get secret -n sftpgo sftpgo-admin -ojsonpath='{.data.admin-password}' | base64 -d"
    alias lab-password-gitea="k get secret -n gitea gitea-secret -ojsonpath='{.data.password}' | base64 -d"
    alias lab-create-secret="k create secret --dry-run -oyaml"
    ```
    Using dashes (or underscores) to separate conceptual subcommands interferes with the flow of shell completion and renaming the common prefixes of the aliases names is inconvenient as well. This limitation also prevent aliases from imitating commands.
- Without implementing workarounds, *Aliases (or any shell configuration) are global for the current user*, and thus there is *no context aware toggling of aliases* - Some of the aliases might be relevant only for a single purpose or project.
- *Aliases do not support passing non suffixed arguments* - Aliases are substituted. That is why there is no support for passing arguments that are not suffixed to the alias itself. A possible workaround is to use shell functions or other tools (See [Future Plans](#future-plans)).

Last but not least, this project is a great excuse for me to work with _Go_, as it is a great language for developing CLI tools. Single binaries for CLI tools are awesome.

## Features
- **N**esting **a**liases (ehm, shortcuts).
- Dynamic shell completions.
- Simple and readable configuration.
- Layered configuration (local and home config).

## Installation

### Binaries
Available in [Releases](https://github.com/moyiz/na/releases) page.

### Source
```sh
go install github.com/moyiz/na@latest
```

### AUR
Use your favorite AUR helper:
```sh
yay -S na-bin
```

### Build
```sh
git clone https://github.com/moyiz/na.git
go build .
# And move it to your preferred `bin`
mv na ~/.local/bin
```

## Usage

### Adding shortcuts
Use the `add` subcommand (or its shorter form: `a`) to add new shortcuts:
```sh
na add my shortcut cwd pwd
na a my longcut cwd -- echo A shortcut to show current working directory, which is '$PWD'
na a e echo
```
In the first example (without double dashes), the last argument is the target of the shortcut, i.e `na run my shortcut cwd` will execute `pwd`. This is handy for single word commands.

Multi word commands can be quoted (i.e `na add my shortcut "ls -ltra"`) but using a double dash is also a valid option. In the second example, anything after the double dash is considered target, thus `na run my longcut cwd` will run that long `echo`. Notice that `$PWD` is single quoted. This is to ensure the actual substitution will occur when the shortcut is called, rather than when it is being added.

If the config directory does not exist, `na` will create it.

### Listing shortcuts
Run `na` without arguments, or use the `list` / `ls` subcommand if you prefer.
```sh
na
na ls
na list
```

### Running shortcuts
Use the `run` subcommand (or its shorted form: `r`).
```sh
na run my shortcut cwd
na r e -- Hello World!
```
The first example is self explanatory. Notice the double dash in the second example. In `run` it is mandatory in order to pass arguments to the shortcut itself.

_na_ will try to auto-detect the current calling shell and invoke the command in it. If it fails to detect or current shell is not supported, `sh` will be used as fallback.

### Removing shortcuts
Use the `remove` subcommand (or its shorter form: `rm`).
```sh
na remove my shortcut cwd
na remove my longcut
na rm e
```
Notice that this command accepts partial shortcuts. It will remove the entire subtree of given shortcuts. The first two examples above can be reduced to a single `na rm my` to delete both.

## Shell Completions
Some installation methods already setup completions. To activate them manually, add the following to your shell's configuration.
### Bash
Add this to your `~/.bashrc`:
```sh
source <(na completion bash)
```
### Zsh
Add this to your `~/.zshrc`:
```sh
source <(na completion zsh)
```
In case of `command not found: compdef`, add these too:
```sh
autoload -Uz compinit
compinit
```

### Fish
Add this to `~/.config/fish/config.fish`
```sh
na completion fish | source
```

### Powershell
Exist but untested.

## Configuration
_na_ looks for configuration files in few locations:
- Local directory (`.na.yaml`)
- Current user home config directory (`~/.config/na/na.yaml`)
- XDG config directory (`/etc/xdg/na/na.yaml`)

By default, `na` will merge these configs for `list` and `run`, and use current user's home config directory for `add` and `remove`. 

This behavior will be overridden by passing either:
- `--config FILE` or `-c FILE`: set the only configuration file to use.
- `--local` or `-l`: synonymous to `-c .na.yaml`.

_na_ configuration is a simple dictionary, mapping shortcut names to either commands or other subcommands.

### Showcase
```sh
$ cat .na.yaml
cat: .na.yaml: No such file or directory
$ cat ~/.config/na/na.yaml
cat: /home/moyiz/.config/na/na.yaml: No such file or directory
$ na add my global -- echo Global
$ cat ~/.config/na/na.yaml
my:
    global: echo Global
$ na add -l my local -- echo Local
$ cat .na.yaml 
my:
    local: echo Local
$ na
my global -- echo Global
my local -- echo Local
$ na run my local
Local
$ na run my global
Global
$ na rm -l my
$ cat .na.yaml
{}
$ cat ~/.config/na/na.yaml
my:
    global: echo Global

```

### Example
```yaml
lab:
  secret:
    gitea: k get secret -n gitea gitea-admin-secret -ojsonpath='{.data.password}' | base64 -d
    nextcloud:  get secret -n gitea gitea-admin-secret -ojsonpath='{.data.password}' | base64 -d
    new: k create secret generic --dry-run -oyaml > secret.yaml
    seal: kubeseal --scope cluster-wide -oyaml < secret.yaml > sealed-secret.yaml
    sealete: na r lab secret seal && rm secret.yaml
  renovate: k delete job -n renovate renovate-manual; k create -n renovate job --from cronjob/renovate renovate-manual
  backup:
    prepare: sudo mount -t nfs nas:/pool/backups /media/backups 
    vm: ssh proxmox.home ./backup.sh
    local:
        volumes: na r backup prepare && rsync ...
        home: na r backup prepare && rsync ...
random: echo $RANDOM
```

## Known Issues
- Since `na` spawns a new shell to run commands, some builtin commands will
  have no effect on current shell, e.g: `alias`, `cd`, `declare`, `export`,
  `pushd`, `read`, `set`, `shopt`, `source` and so on.

## Future Plans
- Magic placeholders in commands.
- Read aliases from environment variables.
- Support `.env`
- Rename aliases.
- Dynamically creates actual commands, e.g.
  ```sh
  my alias
  # Instead of:
  na run my alias
  ```

## License
See [LICENSE](./LICENSE).
