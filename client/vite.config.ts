import { defineConfig } from 'vite';

export default defineConfig({
  build: {
    // This tells Vite where to put the bundled files.
    outDir: 'dist',
  },
  // Set the base path to ensure the preview server
  base: './',
  // You might also need to explicitly list the public directory
  // if you're having issues with manifest.json.
  publicDir: 'public',
});