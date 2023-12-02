import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vuetify from 'vite-plugin-vuetify'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  base: "/ui/",
  plugins: [
    vue(),
    vuetify(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    rollupOptions: {
      input: {
        login: path.resolve(__dirname, 'login.html'),
        settings: path.resolve(__dirname, 'settings.html'),
        authorize: path.resolve(__dirname, 'authorize.html'),
      },
    },
  }
})
