#!/usr/bin/env bun

/**
 * SVG to PNG Converter
 * Converts all SVG brand assets to PNG format for broader compatibility
 *
 * Requirements: npm install sharp
 * Usage: node convert-svgs.js
 */

import { mkdir } from "fs/promises";
import log from "loglevel";
import { join } from "path";
import sharp from "sharp";

var conversions = [
  {
    input: join("./public", "favicon.svg"),
    outputs: [
      { path: join("./public", "favicon-16x16.png"), width: 16, height: 16 },
      { path: join("./public", "favicon-32x32.png"), width: 32, height: 32 },
      { path: join("./public", "favicon-48x48.png"), width: 48, height: 48 },
      { path: join("./public", "favicon-64x64.png"), width: 64, height: 64 },
      { path: join("./public", "favicon.png"), width: 256, height: 256 },
    ],
  },
  {
    input: join("./public", "logo.svg"),
    outputs: [
      { path: join("./public", "logo.png"), width: 400, height: 100 },
      { path: join("./public", "logo-2x.png"), width: 800, height: 200 },
    ],
  },
  // {
  //   input: join("./public", "og-image.svg"),
  //   outputs: [
  //     { path: join("./public", "og-image.png"), width: 1200, height: 630 },
  //   ],
  // },
];

async function convertSVGtoPNG() {
  var publicDir = "./public";
  var coversDir = join(publicDir, "covers");

  try {
    await mkdir(coversDir, { recursive: true });
  } catch (error) {
    log.warn(`⚠️ Could not create directory ${coversDir}:`, error.message);
  }

  log.info("🎨 Converting SVGs to PNGs...\n");

  var successful = 0;
  var failed = 0;

  for (var i = 0; i < conversions.length; i++) {
    var conversion = conversions[i];
    var inputPath = conversion.input;

    for (var j = 0; j < conversion.outputs.length; j++) {
      var output = conversion.outputs[j];

      try {
        await sharp(inputPath)
          .resize(output.width, output.height)
          .png()
          .toFile(output.path);

        log.info(`✅ ${inputPath} → ${output.path}`);
        successful++;
      } catch (error) {
        log.error(`❌ Failed to convert ${inputPath}:`, error.message);
        failed++;
      }
    }
  }

  log.info(`\n📊 Conversion complete:`);
  log.info(`   Successful: ${successful}`);
  log.info(`   Failed: ${failed}`);

  if (successful > 0) {
    log.info(`\n✨ PNG files generated in ${publicDir}/`);
  }
}

// Run the conversion
convertSVGtoPNG().catch(log.error);
