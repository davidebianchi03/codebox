import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  resolve: {
    alias: {
      'bootstrap': path.resolve(__dirname, 'node_modules/bootstrap'),
      '@tabler/core': path.resolve(__dirname, 'node_modules/@tabler/core'),
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        includePaths: [
          path.resolve(__dirname, 'node_modules')
        ],
        api: 'legacy'
      }
    }
  },
  plugins: [react()],
  server: {
    host: "127.0.0.1",
    port: 3000,
  },
  build: {
    outDir: './build',
    emptyOutDir: true,
  }
});
