const Dotenv = require('dotenv-webpack');

module.exports = {
    publicPath: "./",
    configureWebpack: {
        plugins: [
            new Dotenv()
        ]
    }
};