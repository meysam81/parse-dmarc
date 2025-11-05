# Quick Start Guide

Get Parse DMARC up and running in 5 minutes!

## Installation

### Option 1: Build from Source (Recommended)

```bash
# Clone the repository
git clone https://github.com/meysam81/parse-dmarc.git
cd parse-dmarc

# Install dependencies
make install-deps

# Build
make build
```

The binary will be at `./bin/parse-dmarc`

### Option 2: Docker

```bash
# Clone the repository
git clone https://github.com/meysam81/parse-dmarc.git
cd parse-dmarc

# Create config directory
mkdir -p data

# Copy example config
cp config.example.json data/config.json

# Edit config with your IMAP credentials
nano data/config.json

# Run with docker-compose
docker-compose up -d
```

## Configuration

### Step 1: Generate Config File

```bash
./bin/parse-dmarc -gen-config
```

This creates `config.json` with default settings.

### Step 2: Configure IMAP

Edit `config.json` with your email server details:

#### Gmail Example

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

**Note for Gmail**:
1. Enable 2-factor authentication
2. Generate an [App Password](https://support.google.com/accounts/answer/185833)
3. Use the app password in the config

#### Office 365 Example

```json
{
  "imap": {
    "host": "outlook.office365.com",
    "port": 993,
    "username": "your-email@company.com",
    "password": "your-password",
    "mailbox": "INBOX",
    "use_tls": true
  }
}
```

## Running

### Mode 1: Continuous Fetch (Recommended)

Fetches reports every 5 minutes and serves the dashboard:

```bash
./bin/parse-dmarc -config=config.json
```

Access dashboard at: http://localhost:8080

### Mode 2: Custom Fetch Interval

Fetch every 10 minutes (600 seconds):

```bash
./bin/parse-dmarc -config=config.json -fetch-interval=600
```

### Mode 3: Fetch Once

Fetch reports once and exit:

```bash
./bin/parse-dmarc -config=config.json -fetch-once
```

### Mode 4: Serve Only

Only serve the dashboard (no fetching):

```bash
./bin/parse-dmarc -config=config.json -serve-only
```

## Using the Dashboard

Once running, open http://localhost:8080 in your browser.

### Dashboard Features

1. **Statistics Overview**
   - Total reports and messages
   - Compliance rate
   - Unique source IPs

2. **Top Sending Sources**
   - View top IPs sending on your behalf
   - See pass/fail rates for each source

3. **Recent Reports**
   - List of all DMARC reports
   - Click any report for details

4. **Report Details**
   - Full authentication results
   - DKIM and SPF outcomes
   - Source IPs and message counts

## Troubleshooting

### "No reports found"

1. Check your IMAP credentials
2. Verify you have DMARC reports in your inbox
3. Check logs for errors

### "Connection failed"

1. Verify IMAP host and port
2. Check if TLS is required
3. Ensure firewall allows outbound IMAP connections

### "Login failed"

1. Verify username and password
2. For Gmail, use App Password (not regular password)
3. Check if 2FA is required

## Production Deployment

### Systemd Service

1. Copy binary:
```bash
sudo cp bin/parse-dmarc /usr/local/bin/
```

2. Create user:
```bash
sudo useradd -r -s /bin/false parse-dmarc
```

3. Create directories:
```bash
sudo mkdir -p /var/lib/parse-dmarc
sudo mkdir -p /etc/parse-dmarc
sudo cp config.json /etc/parse-dmarc/
sudo chown -R parse-dmarc:parse-dmarc /var/lib/parse-dmarc
```

4. Install service:
```bash
sudo cp parse-dmarc.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable parse-dmarc
sudo systemctl start parse-dmarc
```

5. Check status:
```bash
sudo systemctl status parse-dmarc
```

### Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name dmarc.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Next Steps

- Set up DMARC for your domain if you haven't already
- Configure email providers to send reports to your inbox
- Monitor compliance rates over time
- Investigate failed authentication attempts

## Resources

- [RFC 7489 - DMARC](https://datatracker.ietf.org/doc/html/rfc7489)
- [DMARC.org](https://dmarc.org/)
- [Google Postmaster Tools](https://postmaster.google.com/)

## Getting Help

- Open an issue on GitHub
- Check existing documentation
- Review logs for error messages
