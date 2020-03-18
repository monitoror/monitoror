module.exports = {
  devServer: {
    port: '8000',
    proxy: {
      '^/api': {
        target: 'http://localhost:8888',
      },
    },
  },
}
