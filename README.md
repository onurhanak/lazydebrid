  <h1 align="center">lazydebrid</h1>

lazydebrid is a terminal-based interface for managing Real-Debrid torrents and downloads. It allows you to interact with your Real-Debrid account entirely from the command line, with no browser required.

**It is currently in alpha stage.**

## Features
- Vim-like keybindings
- Add magnet links to Real-Debrid
- Delete torrents from Real-Debrid
- Download unrestricted files directly to your local system
- View and manage active and completed torrents
- Check torrent status
- View detailed file lists and download links

## Installation

### Option 1: Download the binary (recommended)

You can download the latest prebuilt binary for your system from the [Releases](https://github.com/onurhanak/lazydebrid/releases) page.

1. Go to the [Releases](https://github.com/onurhanak/lazydebrid/releases).
2. Download the binary for your platform (`lazydebrid-linux`, `lazydebrid-darwin`, `lazydebrid-windows.exe`, etc.).
3. Make it executable if you are using linux:
```bash
   chmod +x lazydebrid-<your-platform>
```
then move it somewhere in your `PATH`, for example:
```bash
   mv lazydebrid-<your-platform> /usr/local/bin/lazydebrid
```

### Option 2: Build from source

You'll need Go installed (version 1.23 or later).

```bash
git clone https://github.com/onurhanak/lazydebrid.git
cd lazydebrid
go build -o lazydebrid ./cmd/main.go
```

Move the built binary to your preferred location (e.g., `/usr/local/bin/`).

## Configuration
See [configuration](https://github.com/onurhanak/lazydebrid/blob/main/docs/configuration.md) for how to set your API token and customize download path.

## Keybindings
See [keybindings](https://github.com/onurhanak/lazydebrid/blob/main/docs/keybindings.md) for a full list of keyboard shortcuts and navigation controls.

## Acknowledgments

- Thanks to [@jroimartin](https://github.com/jroimartin) for creating the [gocui](https://github.com/jroimartin/gocui) library. The name **lazydebrid** is inspired by [lazygit](https://github.com/jesseduffield/lazygit) and [lazydocker](https://github.com/jesseduffield/lazydocker)


