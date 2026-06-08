# 🍪 iCloud Cookie Extractor Extension

Chrome extension to extract ALL cookies (including HttpOnly) from iCloud.com

## Installation

### Chrome / Edge / Brave

1. Open `chrome://extensions/`
2. Enable **Developer mode** (top right toggle)
3. Click **Load unpacked**
4. Select the `cookie-extractor-extension` folder
5. Done!

### Firefox

1. Open `about:debugging`
2. Click **This Firefox** → **Load Temporary Add-on**
3. Select `manifest.json` from the folder

## Usage

1. Open [icloud.com](https://www.icloud.com) and **log in**
2. Click the extension icon 🍪
3. Click **"Extract Cookies from iCloud.com"**
4. Click **"Copy to Clipboard"**
5. Paste into `cookies.txt` in project root
6. Run: `./hidemyemail generate -l test -c 5`

## Why Extension?

JavaScript in console **cannot access HttpOnly cookies** (browser security).  
This extension uses `chrome.cookies` API to get ALL cookies including:

- ✅ `X-APPLE-WEBAUTH-USER` (HttpOnly) - required for DSID
- ✅ `X-APPLE-WEBAUTH-TOKEN` (HttpOnly) - required for clientID
- ✅ All other iCloud cookies

## Security

- ✅ Works **100% locally** - no data sent anywhere
- ✅ Open source - check the code yourself
- ✅ Only accesses `.icloud.com` cookies
- ✅ No tracking, no analytics, no ads

## Files

- `manifest.json` - Extension configuration
- `popup.html` - User interface
- `popup.js` - Cookie extraction logic
- `icon.png` - Extension icon
