import { defineConfig } from 'vite';
import basicSsl from '@vitejs/plugin-basic-ssl';

export default defineConfig({
  build: {
    // bundled files location
    outDir: 'dist',
  },
  // base path for preview server
  base: './',
  // You might also need to explicitly list the public directory
  // if you're having issues with manifest.json.
  publicDir: 'public',

  plugins: [
    basicSsl()
  ],
  server: {
    host: true, // This is the magic! It exposes the server to your network.
    https: true // Enable HTTPS
  }
});