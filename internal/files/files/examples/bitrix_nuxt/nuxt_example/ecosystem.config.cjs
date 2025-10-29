module.exports = {
  apps: [
    {
      name: 'nuxt-ssr',
      script: '.output/server/index.mjs',
      env: {
        NODE_ENV: 'production',
        PORT: 5174,
        API_BASE_URL: 'https://example.com/api',
      },
    },
  ],
};
