// var debug = process.env.NODE_ENV !== "production";
var webpack = require('webpack');
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')

module.exports = (env,argv) => ({
  context: __dirname,
  // devtool: debug ? "inline-sourcemap" : null,
  entry: __dirname+"/scripts/*",
  output: {
    path: __dirname+"/web/public/js/",
    filename: "mkgame-go.min.js"
  },
  plugins: [],
  optimization: {
    // minimize: true,
    minimizer: argv.mode === 'production' ? [
      new UglifyJsPlugin({
        uglifyOptions: {
          compress: false,
          mangle: true,
        },
      }),
    ] : [],
  }
});