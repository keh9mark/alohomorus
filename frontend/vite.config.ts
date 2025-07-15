import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    proxy: {
      // '/api': 'http://localhost:34115', // Порт Wails-бэкенда
    },
  },
  plugins: [vue()],
  css: {
    preprocessorOptions: {
      scss: {
        // Дополнительные настройки для SCSS, если требуется
      },
    },
  },
})
