module.exports = {
  publicPath: './',
  devServer: {
    port: '8000',
    proxy: 'http://localhost:8080',
  },
}
