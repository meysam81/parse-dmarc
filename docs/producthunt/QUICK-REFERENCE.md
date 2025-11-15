# Quick Reference Guide

Everything you need at a glance for your Product Hunt launch.

## 🎯 The Essentials

### Tagline (use this exactly)
```
Stop email spoofing with DMARC reports in a beautiful dashboard
```

### Elevator Pitch (30 seconds)
"Parse DMARC transforms cryptic XML reports into actionable insights. It auto-fetches from your inbox, analyzes who's sending email as your domain, and spots spoofing attempts—all in a 14MB Docker image with zero dependencies."

### One-Line Description
"DMARC report analyzer in a 14MB Docker image with auto-fetch and beautiful dashboard."

## 📊 Key Numbers (memorize these)

- **14MB** - Total image size
- **30 seconds** - Setup time
- **143x smaller** - vs traditional stack (~2GB)
- **Zero** - External dependencies
- **5 minutes** - Default fetch interval
- **2 minutes** - From install to dashboard

## 🎨 Brand Colors (hex codes)

- Primary Blue: `#3B82F6`
- Dark Blue: `#1E40AF`
- Success Green: `#10B981`
- Error Red: `#EF4444`
- Text Gray: `#6B7280`

## 💬 Your Talking Points

### Why it exists:
"DMARC reports are critical for email security, but they arrive as compressed XML files that are impossible to read. Existing solutions require Elasticsearch, Kibana, and Python—overkill for most teams."

### What makes it different:
"Single Go binary with embedded Vue.js frontend and SQLite. Everything you need in 14MB. No JVM, no cluster management, no complex setup."

### Who it's for:
- SaaS companies protecting domain reputation
- Security teams monitoring for abuse
- DevOps engineers who hate complexity
- Anyone sending emails who cares about deliverability

### The "wow" factor:
"You can literally deploy email security monitoring in 30 seconds with one Docker command."

## 🚀 Quick Setup Demo Script

Use this in comments/demos:

```bash
# 1. Add DNS record (copy-paste ready)
Type: TXT
Name: _dmarc.yourdomain.com
Value: v=DMARC1; p=none; rua=mailto:dmarc@yourdomain.com

# 2. Run container (that's it!)
docker run -d \
  --name parse-dmarc \
  -p 8080:8080 \
  ghcr.io/meysam81/parse-dmarc:latest

# 3. Open dashboard
open http://localhost:8080
```

## ❓ Common Questions & Quick Answers

**Q: How is this different from ParseDMARC?**
A: ParseDMARC requires Elasticsearch + Kibana + Python. We're a single 14MB binary with embedded frontend and SQLite. Same insights, 1/100th the complexity.

**Q: Is this secure?**
A: Yes. All IMAP connections use TLS. Credentials stay in your config file on your infrastructure. Self-hosted, open source, auditable. We never see your data.

**Q: What about multi-domain support?**
A: Currently one instance per domain, but you can run multiple containers easily. Multi-domain support is on the roadmap!

**Q: Why not just use my email provider's DMARC tools?**
A: Many providers don't offer DMARC reporting dashboards. Those that do often charge premium prices or lock you into their ecosystem. Parse DMARC gives you full control and portability.

**Q: Does this work with [Gmail/Outlook/Custom Server]?**
A: Yes! Any IMAP-compatible inbox works. We include example configs for all major providers.

**Q: Can I contribute?**
A: Absolutely! We're Apache 2.0 licensed. PRs welcome: https://github.com/meysam81/parse-dmarc

## 🎯 Feature Highlights (30-second version)

1. **Auto-Fetch**: Connects to any IMAP inbox every 5 minutes
2. **Parse**: Decompresses and parses XML automatically
3. **Store**: SQLite database (no external DB needed)
4. **Visualize**: Beautiful Vue.js dashboard with real-time stats
5. **Deploy**: Single `docker run` command

## 🔗 Important Links

- **GitHub**: https://github.com/meysam81/parse-dmarc
- **Docker Hub**: ghcr.io/meysam81/parse-dmarc
- **License**: Apache 2.0
- **Docs**: (in GitHub README)

## 📱 Social Media Snippets

### Twitter (280 chars)
```
🚀 Just launched Parse DMARC on @ProductHunt!

Stop email spoofing with a 14MB Docker image that turns DMARC XML reports into actionable insights.

✨ Auto-fetch from inbox
📊 Real-time dashboard
🔒 Self-hosted
⚡ Zero dependencies

[PH Link]
```

### Twitter Thread Opener
```
Email security shouldn't require a PhD in XML parsing.

Today I'm launching Parse DMARC – a tool that makes DMARC reports actually readable.

Here's why it exists: 🧵
```

### LinkedIn Opening
```
📧 Email security shouldn't be this hard.

If you've ever set up DMARC, you know the pain: Gmail and Outlook send you XML reports that are impossible to read manually.

Today I'm launching Parse DMARC on Product Hunt...
```

## 🎬 Video Script (if you make one)

**0:00-0:10** - Problem: "DMARC reports look like this [show XML]"
**0:10-0:20** - Solution: "Parse DMARC turns them into this [show dashboard]"
**0:20-0:30** - How: "Auto-fetches, parses, visualizes. All in 14MB."
**0:30-0:40** - Setup: "One Docker command. 30 seconds."
**0:40-0:50** - Features: Quick dashboard tour
**0:50-0:60** - CTA: "Try it now. Link in comments."

## 🏆 Competitive Advantages (vs alternatives)

| Feature | Parse DMARC | ParseDMARC | Postmark | Dmarcian |
|---------|-------------|------------|----------|----------|
| Size | 14MB | ~2GB | N/A | N/A |
| Setup | 30s | Hours | N/A | N/A |
| Self-hosted | ✅ | ✅ | ❌ | ❌ |
| Cost | Free | Free | $$ | $$$ |
| Dependencies | 0 | Many | N/A | N/A |

## 🎯 Target Personas

### 1. DevOps Engineer Danny
- **Pain**: Tired of complex setups
- **Hook**: "30-second deployment"
- **Benefit**: Minimal resource usage

### 2. Security Engineer Sarah
- **Pain**: Can't read XML reports
- **Hook**: "Spot threats instantly"
- **Benefit**: Real-time monitoring

### 3. Email Admin Emma
- **Pain**: Expensive SaaS tools
- **Hook**: "Self-hosted, free forever"
- **Benefit**: Full control of data

### 4. Startup CTO Chris
- **Pain**: Limited budget & time
- **Hook**: "14MB, zero dependencies"
- **Benefit**: Just works™

## 📈 Success Indicators

Hour by hour goals:

- **Hour 1**: 20 upvotes, 5 comments
- **Hour 3**: 50 upvotes, 15 comments, trending
- **Hour 6**: 100 upvotes, 30 comments, top 10
- **Hour 12**: 150 upvotes, 40 comments, top 5
- **End of day**: 200+ upvotes, 50+ comments, top 3

## ⚡ Emergency Responses

### If servers go down:
"We're experiencing high traffic (great problem to have!). Working on scaling now. In the meantime, you can self-host: [link]"

### If there's a bug:
"Great catch! Issue created: [GitHub link]. Will push fix within [timeframe]. Thanks for helping make this better!"

### If someone's being negative:
"Thanks for the feedback! Could you elaborate on [specific point]? Always looking to improve."

## 🎨 Visual Assets Quick Pick

**Need a quick graphic?** Use these:

- **Twitter**: `social/twitter-card.svg`
- **LinkedIn**: `social/linkedin-card.svg`
- **Blog post**: `images/hero-banner.svg`
- **Documentation**: `images/architecture-diagram.svg`
- **Comparison**: `images/comparison-chart.svg`

## 📞 Pre-Launch Mental Checklist

- [ ] Have I tested the product thoroughly?
- [ ] Can I explain it in 30 seconds?
- [ ] Do I know my key numbers?
- [ ] Can I respond to common objections?
- [ ] Am I ready to be online all day?
- [ ] Do I have my supporter team ready?
- [ ] Have I cleared my calendar?
- [ ] Am I excited and authentic?

## 💪 Confidence Boosters

Remember:
- You've solved a real problem
- Your solution is genuinely simpler
- The technical execution is solid
- You're helping people, not just promoting
- Open source benefits everyone
- You've prepared thoroughly

**You've got this! 🚀**

---

## 📝 Print This Page

Consider printing this page on launch day for quick reference when responding to comments. Having key numbers and responses at your fingertips helps you respond faster and more consistently.

Good luck! 🎉
