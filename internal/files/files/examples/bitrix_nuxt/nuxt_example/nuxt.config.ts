import tailwindcss from "@tailwindcss/vite";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: {
    enabled: true
  },
  css: ['./app/assets/css/main.css'],
  devServer: {
    host: '0.0.0.0',
    port: 5173,
  },
  vite: {
    plugins: [
      tailwindcss(),
    ],
    server: {
      proxy: {
        '/upload': {
          target: 'http://nginx', // docker service name
          changeOrigin: true,
          secure: false,
        },
      },
    },
  },
})
