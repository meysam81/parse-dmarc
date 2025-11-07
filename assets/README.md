# Parse DMARC Assets

This directory contains the branding assets for the Parse DMARC project.

## Source Files (Version Controlled)

The following SVG files are the **source of truth** and are version controlled:

- **`favicon.svg`** - Shield icon used for favicons (100×100)
- **`logo.svg`** - Full logo with text for documentation and UI (400×100)
- **`og-image.svg`** - Social media share image (1200×630)
- **`demo.png`** - Dashboard screenshot for README

## Generated Files (Build Artifacts)

The following PNG files are **automatically generated** from SVG sources at build time:

- `favicon-16x16.png` - Favicon for browsers
- `favicon-32x32.png` - Favicon for browsers
- `favicon-48x48.png` - Favicon for browsers
- `favicon-64x64.png` - Favicon for browsers
- `favicon.png` - Large favicon (256×256) for Apple touch icon
- `logo.png` - Logo PNG (400×100)
- `logo-2x.png` - High-DPI logo PNG (800×200)
- `og-image.png` - Social media share image PNG (1200×630)

**Note:** Generated PNG files are excluded from version control (see `.gitignore`).

## Build Process

PNG assets are automatically generated during the build process:

```bash
# Install dependencies (sharp for SVG → PNG conversion)
npm install

# Generate PNG files from SVG sources
npm run convert-assets

# Build frontend (automatically runs convert-assets first)
npm run build
```

## Manual Conversion

To manually generate PNG files:

```bash
node scripts/convert-svg-to-png.js
```

## Color Scheme

The project uses a purple-blue gradient theme:

- Primary: `#667eea` (purple-blue)
- Secondary: `#764ba2` (deep purple)
- Accent: White (`#ffffff`)

## Design Guidelines

When modifying assets:

1. **Edit SVG sources only** - Never manually edit generated PNG files
2. **Test conversions** - Run `npm run convert-assets` after editing SVG files
3. **Maintain consistency** - Use the established color scheme and shield icon theme
4. **Keep it simple** - Favicons should be recognizable at small sizes (16×16)

## Asset Usage

### In HTML

```html
<!-- Favicon -->
<link rel="icon" type="image/svg+xml" href="/favicon.svg" />
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png" />

<!-- Social Media -->
<meta property="og:image" content="/og-image.png" />
```

### In Documentation

```markdown
![Parse DMARC Logo](./assets/logo.svg)
```

## Troubleshooting

**PNGs not generating?**

1. Ensure Node.js v18+ is installed
2. Run `npm install` to install sharp
3. Check that SVG files are valid
4. Run the conversion script manually: `node scripts/convert-svg-to-png.js`

**Icons not displaying?**

1. Clear browser cache
2. Rebuild the frontend: `cd frontend && npm run build`
3. Check that publicDir is configured in `vite.config.js`
