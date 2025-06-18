module.exports = {
  pluginOptions: {
    vuetify: {
      // https://github.com/vuetifyjs/vuetify-loader/tree/next/packages/vuetify-loader
    }
  },
  publicPath: '/',
  devServer: {
    proxy: {
      '/api/v1': {'target':'http://localhost:3001'},
    }
  }
}
