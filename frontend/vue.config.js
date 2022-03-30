module.exports = {
  "publicPath": './',
  "transpileDependencies": [
    "vuetify"
  ],
  devServer: {
    proxy: {
      '/api': {
          target: 'http://localhost:3001'
      }
    }
  }
}