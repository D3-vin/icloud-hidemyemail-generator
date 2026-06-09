# 📧 iCloud HideMyEmail Generator

CLI tool with interactive menu for generating iCloud Hide My Email addresses with TLS fingerprinting.

[![Version](https://img.shields.io/badge/version-1.0.0-blue)](https://github.com/D3-vin/icloud-hidemyemail-generator/releases)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[![Telegram Channel](https://img.shields.io/badge/Telegram-Channel-blue?logo=telegram)](https://t.me/D3_vin)
[![Telegram Chat](https://img.shields.io/badge/Telegram-Chat-blue?logo=telegram)](https://t.me/D3vin_chat)
[![GitHub](https://img.shields.io/badge/GitHub-Repository-black?logo=github)](https://github.com/D3-vin/icloud-hidemyemail-generator)

[Features](#features) • [Quick Start](#quick-start) • [Usage](#usage) • [Building](#building) • [Contact](#contact)

[English](#english) | [Русский](README_RU.md)

---

## Features

- 🚀 **Fast** - Native Go performance
- 🔒 **Secure** - TLS fingerprinting with Chrome 146 profile
- 📊 **Interactive Menu** - Easy-to-use interface
- 🎯 **Auto-numbering** - Automatically number labels (test1, test2, test3...)
- 📁 **Organized Output** - Generated emails → `generated/`, list results → `results/`
- 🌐 **Cross-platform** - Windows, macOS (Intel & ARM), Linux
- 📦 **Standalone** - Single binary, no dependencies

---

## Quick Start

### 1. Download

Download the latest release for your platform:
- [Windows (64-bit)](https://github.com/D3-vin/icloud-hidemyemail-generator/releases)
- [Linux (64-bit)](https://github.com/D3-vin/icloud-hidemyemail-generator/releases)
- [macOS Intel](https://github.com/D3-vin/icloud-hidemyemail-generator/releases)
- [macOS Apple Silicon](https://github.com/D3-vin/icloud-hidemyemail-generator/releases)

Or build from source:
```bash
git clone https://github.com/D3-vin/icloud-hidemyemail-generator.git
cd icloud-hidemyemail-generator
go build -o hidemyemail .
```

### 2. Extract Cookies

**Using Chrome Extension (Recommended):**

1. Install extension from `cookie-extractor-extension/` folder
2. Open [icloud.com](https://www.icloud.com) and log in
3. Click extension icon 🍪 → Extract → Copy
4. Paste into `cookies.txt`

See [cookie-extractor-extension/README.md](cookie-extractor-extension/README.md) for details.

### 3. Generate Emails

**Interactive Menu:**
```bash
./hidemyemail
```

**CLI Commands:**
```bash
# Generate 5 emails with label "test"
./hidemyemail generate -l test -c 5

# List all emails
./hidemyemail list

# List active emails only
./hidemyemail list --active
```

---

## Usage

### Interactive Menu

Run without arguments to open the menu:

```bash
./hidemyemail
```

```
╔═══════════════════════════════════════╗
║              Menu                     ║
╠═══════════════════════════════════════╣
║  1. Generate emails                   ║
║  2. List all emails                   ║
║  3. Exit                              ║
╚═══════════════════════════════════════╝

Choose option: 1
Enter label: test
Enter count (1-100): 5
Add number to label? (y/n): y

✓ [1/5] email1@privaterelay.appleid.com (label: test1)
✓ [2/5] email2@privaterelay.appleid.com (label: test2)
...
✓ Saved to generated/emails_test_2026-06-09_01-23-45.txt

Choose option: 2
✓ Saved to results/emails_list.txt and results/emails_full.txt
```

### CLI Commands

**Generate:**
```bash
./hidemyemail generate -l <label> -c <count> [options]

Options:
  -l, --label string         Label for emails (required)
  -c, --count int            Number of emails, 1-100 (required)
      --cookie-file string   Cookie file path (default "cookies.txt")
  -o, --output string        Output file (default "emails.txt")
      --no-output-file       Don't save to file
```

**List:**
```bash
./hidemyemail list [options]

Options:
      --label-query string   Regex filter by label
      --active               Show only active emails
      --inactive             Show only inactive emails
      --cookie-file string   Cookie file path (default "cookies.txt")
```

---

## Project Structure

```
icloud-hidemyemail-generator/
├── hidemyemail              # Binary (Linux/macOS)
├── hidemyemail.exe          # Binary (Windows)
├── cookies.txt              # Your iCloud cookies (create this)
├── generated/               # Generated emails with timestamps
├── results/                 # List results (emails_list.txt, emails_full.txt)
├── cookie-extractor-extension/  # Chrome extension for cookies
├── cmd/                     # CLI application
├── internal/                # Core logic
│   ├── api/                 # iCloud API client
│   ├── config/              # Configuration
│   ├── generator/           # Email generation
│   ├── lister/              # Email listing
│   └── output/              # Terminal UI
└── pkg/models/              # Data models
```

---

## Cookie Extraction

### Method 1: Chrome Extension (Recommended)

Install the Chrome extension from `cookie-extractor-extension/` folder to extract HttpOnly cookies.

**Why extension?** JavaScript console scripts cannot access HttpOnly cookies (`X-APPLE-WEBAUTH-USER`, `X-APPLE-WEBAUTH-TOKEN`) due to browser security.

**Installation:**
1. Open Chrome → `chrome://extensions/`
2. Enable "Developer mode" (top right)
3. Click "Load unpacked"
4. Select `cookie-extractor-extension/` folder
5. Pin the extension (puzzle icon → pin 📌)

**Usage:**
1. Go to [icloud.com](https://www.icloud.com) and log in
2. Click extension icon 🍪
3. Click "Extract Cookies"
4. Copy the cookies string
5. Paste into `cookies.txt` in project root

See [cookie-extractor-extension/README.md](cookie-extractor-extension/README.md) for details.

### Method 2: DevTools Network Tab (Manual)

1. Open [icloud.com](https://www.icloud.com) and log in
2. Open DevTools (F12) → Network tab
3. Refresh page (F5)
4. Click any request to `icloud.com`
5. Find **Headers** → **Request Headers** → **Cookie:**
6. Copy the entire cookie string
7. Paste into `cookies.txt`

**Note:** Manual method may miss HttpOnly cookies. Use extension for best results.

---

## Building

### Single Platform

```bash
go build -o hidemyemail .
```

### All Platforms

```bash
# Windows
build.cmd

# Or manually
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/hidemyemail-windows-amd64.exe .
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/hidemyemail-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/hidemyemail-macos-amd64 .
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/hidemyemail-macos-arm64 .
```

---

## Troubleshooting

### "Cookie file not found"
Create `cookies.txt` in project root with your iCloud cookies.

### "Authentication failed"
Cookies expired. Extract fresh cookies from icloud.com.

### "failed to extract DSID from cookies"
Missing HttpOnly cookies. Use the Chrome extension instead of console scripts.

### Rate Limit
iCloud limits generation to **~5 emails per 30 minutes per family member**. Wait and try again.

---

## Contact

- **GitHub**: https://github.com/D3-vin/icloud-hidemyemail-generator
- **Telegram**: [@D3_vin](https://t.me/D3_vin)
- **Author**: [@D3vin_dev](https://t.me/D3vin_dev)

---

## License

MIT License - see [LICENSE](LICENSE) file

---

**⚠️ Disclaimer**: This tool is for educational purposes only. Use at your own risk.
