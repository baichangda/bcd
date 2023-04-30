const {createProxyMiddleware} = require("http-proxy-middleware");

module.exports = function (app) {
    app.use(
        createProxyMiddleware('/api', {
            target: 'http://43.138.63.155',
            // target: 'http://127.0.0.1',
            changeOrigin: true,
        }),
    )
}