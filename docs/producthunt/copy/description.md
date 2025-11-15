# Product Description

## Short Description (260 chars max)
Transform cryptic DMARC XML reports into actionable insights. Parse DMARC auto-fetches reports from your inbox, analyzes who's sending email as your domain, and spots spoofing attempts—all in a sleek 14MB Docker image with zero dependencies.

## Full Description

**The Problem Every Email Admin Faces:**

You've set up DMARC to protect your domain from spoofing. Great! But now Gmail, Outlook, and Yahoo are flooding your inbox with compressed XML files that look like this:

```xml
<?xml version="1.0"?>
<feedback>
  <report_metadata>
    <org_name>google.com</org_name>
    <email>noreply-dmarc-support@google.com</email>
    <report_id>12345678901234567890</report_id>
```

**Completely unreadable.** Yet these reports contain critical security intelligence about your domain.

**Parse DMARC Changes Everything:**

✨ **Auto-Fetch**: Connects to any IMAP inbox (Gmail, Outlook, you name it) and automatically pulls DMARC reports
📊 **Beautiful Dashboard**: See real-time stats, compliance rates, and authentication results at a glance
🔍 **Spot Threats Fast**: Instantly identify unauthorized senders and spoofing attempts
⚡ **Ridiculously Lightweight**: Single 14MB binary. No Elasticsearch. No Kibana. No Python dependencies.
🚀 **One Command Deploy**: `docker run` and you're done. SQLite backend, embedded frontend, zero config.

**Why Email Teams Love It:**

→ **For Security Teams**: Catch phishing attempts before they damage your reputation
→ **For DevOps**: Deploy in 30 seconds. No database clusters to manage.
→ **For Deliverability Engineers**: Validate SPF/DKIM configs across all your sending sources
→ **For Compliance Officers**: Prove due diligence with comprehensive email authentication monitoring

**Built by email nerds, for email nerds.**

Unlike bloated alternatives that require Elasticsearch clusters and Python environments, Parse DMARC is a single Go binary with an embedded Vue.js frontend. It respects your infrastructure—runs anywhere Docker does, uses minimal resources, and just works.

Perfect for:
- SaaS companies protecting their domain reputation
- Email service providers validating customer configurations
- Security teams monitoring for domain abuse
- DevOps engineers who hate complex deployments

**Start monitoring in 2 minutes:**
1. Add DNS record: `_dmarc.yourdomain.com TXT "v=DMARC1; p=none; rua=mailto:dmarc@yourdomain.com"`
2. Run: `docker run -p 8080:8080 ghcr.io/meysam81/parse-dmarc`
3. Open http://localhost:8080

That's it. Reports start flowing within 24 hours.

---

**Tech Stack:**
- Backend: Go (single binary)
- Frontend: Vue.js (embedded)
- Database: SQLite
- Container: 14MB Alpine-based image
- License: Apache 2.0
