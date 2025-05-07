# Configuration Guide
Lazydebrid uses a simple JSON configuration file to store its settings.

## Config File Location

The config file is stored at:

```bash
$CONFIG/lazydebrid/lazydebrid.json
```

Replace `$CONFIG` with the appropriate configuration directory for your platform:

* Linux/macOS: `$HOME/.config`
* Windows: `%APPDATA%`

## File Structure

```json
{
  "apiToken": "your_real_debrid_api_key",
  "downloadPath": "/your/download/path"
}
```
Default download path is `$HOME/Downloads`.

## Initial Setup
On first run, Lazydebrid will prompt you for the API key and download path. These values will be saved in the config file automatically.

## Updating Config
You can update the config in two ways:

1. Interactively from within the app.
2. Manually by editing lazydebrid.json directly with a text editor.

Changes take effect immediately when the app is restarted.
