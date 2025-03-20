import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { resolve } from 'path';

export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 9876,
    host: '0.0.0.0',
    proxy: {
      '/api': {
        target: 'http://0.0.0.0:8080',
        changeOrigin: true
      }
    }
  }
}); 