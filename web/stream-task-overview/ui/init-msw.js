// MSW initialization script
import { writeFileSync, mkdirSync, existsSync } from 'fs';
import { resolve } from 'path';
import { fileURLToPath } from 'url';

// Get the directory name in ESM
const __filename = fileURLToPath(import.meta.url);
const __dirname = new URL('.', import.meta.url).pathname;

async function init() {
  // Create public directory if it doesn't exist
  const publicDir = resolve(__dirname, 'public');
  if (!existsSync(publicDir)) {
    mkdirSync(publicDir, { recursive: true });
  }

  // Download the worker script from MSW CDN
  const workerUrl = 'https://unpkg.com/msw@latest/mockServiceWorker.js';
  console.log(`[MSW] Downloading worker from ${workerUrl}`);
  
  const response = await fetch(workerUrl);
  if (!response.ok) {
    throw new Error(`Failed to download MSW worker: ${response.statusText}`);
  }
  
  const workerScript = await response.text();
  const outFile = resolve(publicDir, 'mockServiceWorker.js');
  
  // Write worker script to public directory
  writeFileSync(outFile, workerScript);
  console.log(`[MSW] Worker file successfully created at ${outFile}`);
}

init().catch(console.error); 