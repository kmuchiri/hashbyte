# HashByte

An Activation Byte Recovery Tool implemented in Go. 

## Installation

### Method 1: Using the Install Script (Linux & macOS)
You can quickly install the latest pre-compiled binary via install script:
```bash
curl -sSL https://raw.githubusercontent.com/kmuchiri/hashbyte/main/scripts/install.sh | sudo bash
```

### Method 2: Using the Install Script (Windows)
For Windows users, open PowerShell and run:
```powershell
Invoke-Expression (Invoke-RestMethod -Uri "https://raw.githubusercontent.com/kmuchiri/hashbyte/main/scripts/install.ps1")
```
*Or using the shorthand:*
```powershell
irm https://raw.githubusercontent.com/kmuchiri/hashbyte/main/scripts/install.ps1 | iex
```
This script downloads the binary, places it in your local AppData folder (`%LOCALAPPDATA%\hashbyte`), and automatically adds it to your system PATH.

### Method 3: Build from Source
Ensure you have [Go](https://go.dev/) installed.

```bash

git clone https://github.com/kmuchiri/hashbyte.git
cd hashbyte
go build -o hashbyte main.go
sudo mv hashbyte /usr/local/bin/

```
*(Note: The rainbow tables are embedded into the binary using `//go:embed`, meaning the resulting executable is fully standalone and can be moved anywhere on your system.)*

## Usage

### 1. Brute-Force Mode
Search the full keyspace without relying on rainbow tables. This will utilize all available CPU core, making it computing intensive.

> [!IMPORTANT]
> Reference: time completed on Ryzen 7 5850U computer was 81 seconds.
> If running multiple instances or on a low power processor use the rainbow table lookup.

```bash
hashbyte brute-force <sha1_hash>
```
*Example:* `hashbyte brute-force 05768689284d85b1abe78176134439220f87c209`

### 2. Rainbow Table Lookup
Look up a hash instantly using the pre-computed rainbow tables embedded directly within the executable.

```bash
hashbyte rainbow lookup <sha1_hash>
```

#### Verbose Output
By default, `lookup` only prints the resulting activation bytes in lowercase hex. If you want detailed metrics (like table used, step, and elapsed time), use the `-v` or `--verbose` flag:
```bash
hashbyte rainbow lookup -v <sha1_hash>
```

#### Using External Tables
If you generated a custom rainbow table file and want to use it instead of the embedded ones, simply provide the path before the hash:
```bash
hashbyte rainbow lookup [table_file] <sha1_hash>
```
*Example:* `hashbyte rainbow lookup -v my_tables.bin 05768689284d85b1abe78176134439220f87c209`

### 3. Generate New Tables
You can build your own rainbow tables and save them to a binary file:
```bash
hashbyte rainbow generate <output_file>
```

## Uninstall

Because `hashbyte` is a single standalone executable, uninstalling it is as simple as deleting the file.

**Linux & macOS:**
```bash
sudo rm /usr/local/bin/hashbyte
```

**Windows:**
Open PowerShell and remove the installation directory:
```powershell
Remove-Item -Recurse -Force $env:LOCALAPPDATA\hashbyte
```
*(You can also safely remove the folder from your system PATH via the Windows Environment Variables menu).*

