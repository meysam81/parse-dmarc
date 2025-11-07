#!/usr/bin/env node

/**
 * SVG to PNG Converter
 *
 * This script converts SVG files to PNG format at build time to ensure
 * consistent and deterministic output across different environments.
 *
 * Usage:
 *   node scripts/convert-svg-to-png.js
 *
 * The script will convert:
 *   - assets/favicon.svg → assets/favicon.png (multiple sizes)
 *   - assets/logo.svg → assets/logo.png
 *   - assets/og-image.svg → assets/og-image.png
 *
 * Requirements:
 *   npm install --save-dev sharp
 */

import { mkdirSync, readFileSync, writeFileSync } from "fs";
import { dirname, join } from "path";
import sharp from "sharp";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const rootDir = join(__dirname, "..");
const assetsDir = join(rootDir, "assets");

/**
 * Convert SVG buffer to PNG at specified size
 * @param {Buffer} svgBuffer - SVG file buffer
 * @param {number|null} width - Target width (null to preserve original)
 * @param {number|null} height - Target height (null to preserve original)
 * @returns {Promise<Buffer>} PNG buffer
 */
async function convertSvgToPng(svgBuffer, width = null, height = null) {
  const sharpInstance = sharp(svgBuffer);

  if (width || height) {
    sharpInstance.resize(width, height, {
      fit: "contain",
      background: { r: 0, g: 0, b: 0, alpha: 0 },
    });
  }

  return sharpInstance.png().toBuffer();
}

/**
 * Main conversion logic
 */
async function main() {
  console.log("🎨 Converting SVG assets to PNG...\n");

  const conversions = [
    // Favicon conversions (multiple sizes for different use cases)
    {
      input: join(assetsDir, "favicon.svg"),
      outputs: [
        { path: join(assetsDir, "favicon-16x16.png"), width: 16, height: 16 },
        { path: join(assetsDir, "favicon-32x32.png"), width: 32, height: 32 },
        { path: join(assetsDir, "favicon-48x48.png"), width: 48, height: 48 },
        { path: join(assetsDir, "favicon-64x64.png"), width: 64, height: 64 },
        { path: join(assetsDir, "favicon.png"), width: 256, height: 256 },
      ],
    },
    // Logo conversion
    {
      input: join(assetsDir, "logo.svg"),
      outputs: [
        { path: join(assetsDir, "logo.png"), width: 400, height: 100 },
        { path: join(assetsDir, "logo-2x.png"), width: 800, height: 200 },
      ],
    },
    // OG Image conversion (for social media sharing)
    {
      input: join(assetsDir, "og-image.svg"),
      outputs: [
        { path: join(assetsDir, "og-image.png"), width: 1200, height: 630 },
      ],
    },
  ];

  let successCount = 0;
  let errorCount = 0;

  for (const conversion of conversions) {
    try {
      console.log(`📄 Reading ${conversion.input}...`);
      const svgBuffer = readFileSync(conversion.input);

      for (const output of conversion.outputs) {
        try {
          console.log(
            `  ↳ Converting to ${output.path} (${output.width}x${output.height})...`,
          );

          const pngBuffer = await convertSvgToPng(
            svgBuffer,
            output.width,
            output.height,
          );

          // Ensure directory exists
          writeFileSync(output.path, pngBuffer);

          // Write PNG file
          await sharp(pngBuffer).toFile(output.path);

          console.log(`  ✅ Created ${output.path}`);
          successCount++;
        } catch (error) {
          console.error(
            `  ❌ Failed to convert ${output.path}:`,
            error.message,
          );
          errorCount++;
        }
      }
      console.log("");
    } catch (error) {
      console.error(`❌ Failed to read ${conversion.input}:`, error.message);
      errorCount++;
    }
  }

  console.log("━".repeat(50));
  console.log(`✨ Conversion complete!`);
  console.log(`   Success: ${successCount} files`);
  if (errorCount > 0) {
    console.log(`   Errors: ${errorCount} files`);
    process.exit(1);
  }
}

// Run the conversion
main().catch((error) => {
  console.error("💥 Fatal error:", error);
  process.exit(1);
});
