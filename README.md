# Parse DMARC

A minimal dependency Go application that parses DMARC aggregate reports and presents them in a delightful Vue.js dashboard. Built as a single binary, similar to [listmonk](https://listmonk.app/).

## Features

- 📧 **IMAP Integration** - Automatically fetches DMARC reports from your email inbox
- 🔍 **RFC 7489 Compliant** - Fully compliant with DMARC aggregate report standards
- 📊 **Beautiful Dashboard** - Modern Vue.js 3 SPA with real-time statistics
- 🗄️ **SQLite Storage** - Lightweight embedded database, no external dependencies
- 📦 **Single Binary** - Everything embedded in one executable
- 🚀 **Minimal Dependencies** - Pure Go + Vue.js, no external services required
- 🔒 **Secure** - TLS support for IMAP connections

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm (for building frontend)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/meysam81/parse-dmarc.git
cd parse-dmarc
```

2. Install dependencies and build:
```bash
make install-deps
make build
```

The compiled binary will be at `./bin/parse-dmarc`.

### Configuration

1. Generate a sample configuration file:
```bash
./bin/parse-dmarc -gen-config
```

2. Edit `config.json` with your IMAP credentials:
```json
{
  "imap": {
    "host": "imap.gmail.com",
    "port": 993,
    "username": "your-email@gmail.com",
    "password": "your-app-password",
    "mailbox": "INBOX",
    "use_tls": true
  },
  "database": {
    "path": "./dmarc.db"
  },
  "server": {
    "port": 8080,
    "host": "0.0.0.0"
  }
}
```

### Running

**Fetch reports and start dashboard:**
```bash
./bin/parse-dmarc -config=config.json
```

**Fetch once and exit:**
```bash
./bin/parse-dmarc -config=config.json -fetch-once
```

**Only serve dashboard (no fetching):**
```bash
./bin/parse-dmarc -config=config.json -serve-only
```

**Custom fetch interval:**
```bash
./bin/parse-dmarc -config=config.json -fetch-interval=600
```

Access the dashboard at `http://localhost:8080`

## Architecture

### Backend (Go)

- **IMAP Client** (`internal/imap`) - Connects to email server and fetches DMARC reports
- **Parser** (`internal/parser`) - Parses DMARC XML reports (handles gzip/zip compression)
- **Storage** (`internal/storage`) - SQLite database layer for persisting reports
- **API** (`internal/api`) - REST API endpoints for dashboard
- **Config** (`internal/config`) - Configuration management

### Frontend (Vue.js 3)

- **Single Page Application** - Responsive, modern UI
- **Real-time Statistics** - Compliance rates, message counts, source IPs
- **Report Viewer** - Detailed view of individual DMARC reports
- **Auto-refresh** - Dashboard updates every 5 minutes

### DMARC Report Structure (RFC 7489)

The parser handles the complete DMARC aggregate report schema:

- **Report Metadata** - Organization, date range, report ID
- **Policy Published** - DMARC policy (p, sp, pct, alignment modes)
- **Records** - Individual authentication results:
  - Source IP and message count
  - DKIM authentication results
  - SPF authentication results
  - Policy evaluation (disposition)

## API Endpoints

- `GET /api/statistics` - Dashboard statistics
- `GET /api/reports` - List of reports (paginated)
- `GET /api/reports/:id` - Detailed report view
- `GET /api/top-sources` - Top sending source IPs

## Development

### Frontend Development

Run frontend dev server with hot reload:
```bash
make frontend-dev
```

### Backend Development

Run backend in development mode:
```bash
make dev
```

### Running Tests

```bash
make test
```

## Building

### Build Everything
```bash
make build
```

### Build Frontend Only
```bash
make frontend
```

### Build Backend Only
```bash
make backend
```

## Deployment

The application compiles to a single binary that can be deployed anywhere:

```bash
# Build
make build

# Copy binary to server
scp bin/parse-dmarc user@server:/usr/local/bin/

# Run on server
parse-dmarc -config=/etc/parse-dmarc/config.json
```

### Systemd Service

Create `/etc/systemd/system/parse-dmarc.service`:

```ini
[Unit]
Description=DMARC Report Parser
After=network.target

[Service]
Type=simple
User=parse-dmarc
WorkingDirectory=/var/lib/parse-dmarc
ExecStart=/usr/local/bin/parse-dmarc -config=/etc/parse-dmarc/config.json
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable parse-dmarc
sudo systemctl start parse-dmarc
```

## IMAP Configuration Tips

### Gmail
- Use [App Passwords](https://support.google.com/accounts/answer/185833)
- Host: `imap.gmail.com`, Port: `993`

### Outlook/Office 365
- Host: `outlook.office365.com`, Port: `993`
- May require app password or OAuth (OAuth not yet supported)

### Generic IMAP
- Most providers use port `993` for TLS
- Check your email provider's IMAP settings

## Compatibility

Compatible with DMARC reports from major email service providers:

- Google Workspace
- Microsoft Office 365
- Yahoo
- SendGrid
- Mailchimp
- Amazon SES
- And any RFC 7489 compliant reporter

## Comparison with ParseDMARC

This implementation is inspired by [ParseDMARC](https://github.com/domainaware/parsedmarc) but differs in:

- **Language**: Go instead of Python
- **Deployment**: Single binary vs Python package
- **Storage**: Embedded SQLite vs External Elasticsearch/OpenSearch
- **UI**: Built-in Vue.js dashboard vs Kibana/Grafana
- **Dependencies**: Minimal vs Heavy (no Elasticsearch, Kibana, etc.)

## License

MIT License - see LICENSE file for details

## Contributing

Contributions welcome! Please open an issue or submit a pull request.

## Roadmap

- [ ] Forensic report support (RUF)
- [ ] OAuth2 support for IMAP
- [ ] Export reports (CSV, JSON)
- [ ] Email alerts for compliance issues
- [ ] Multi-domain support
- [ ] DMARC policy recommendations
- [ ] Historical trend analysis

## References

- [RFC 7489 - DMARC](https://datatracker.ietf.org/doc/html/rfc7489)
- [DMARC.org](https://dmarc.org/)
- [ParseDMARC](https://github.com/domainaware/parsedmarc)
