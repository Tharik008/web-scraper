# Google Business Review Scraper (Full-Stack)

A professional tool to extract Google Business reviews using a Chrome Extension and persist them into a local SQLite database via a Go backend.

## 🚀 Features
- **Dynamic DOM Scraping**: Captures text, ratings, and user-uploaded media links.
- **Pure Go Backend**: Uses a CGO-free SQLite driver for easy deployment.
- **Data Integrity**: Uses Review IDs as primary keys to prevent duplicates.

## 🛠️ Installation & Setup

### 1. Backend (Go)
1. Ensure you have Go installed.
2. Run `go mod tidy` to install the Pure Go SQLite driver.
3. Start the server:
   ```bash
   go run main.go