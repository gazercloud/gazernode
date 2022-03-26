const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/api/request',
    createProxyMiddleware({
      target: 'http://localhost:8084',
        changeOrigin: true,
		secure: false,
        cookieDomainRewrite: "localhost",
      onProxyReq: (proxyReq) => {
        if (proxyReq.getHeader('origin')) {
          proxyReq.setHeader('origin', 'http://localhost:8084')
        }
      }
    })
  );
};
