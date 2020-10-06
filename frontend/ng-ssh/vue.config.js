module.exports = {
  "transpileDependencies": [
    "vuetify"
  ],
  runtimeCompiler: true,
  publicPath: process.env.VUE_APP_STATIC_PREFIX,
  devServer: {
    host: '0.0.0.0',
    port: 8084,
    disableHostCheck: true,
  },
  chainWebpack: config => {
    config
      .plugin('html')
      .tap(args => {
        args[0].title= 'ng-ssh'
        return args
      })
  }
}