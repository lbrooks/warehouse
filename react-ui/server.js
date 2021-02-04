const { createProxyMiddleware } = require('http-proxy-middleware');
const express = require('express')
const app = express()
const port = 3000

const {initTracer, ZipkinB3TextMapCodec, opentracing} = require("jaeger-client");

function initTrace(serviceName) {
  var config = {
    serviceName: serviceName,
    sampler: {
        type: "const",
        param: 1,
    },
    reporter: {
        collectorEndpoint: process.env.JAEGER_URL || 'http://localhost:14268/api/traces'
    },
  };
  var options = {
    tags: {
      'version': '0.0.1',
    },
    logger: {
      info: function logInfo(msg) {
        console.log("INFO ", msg);
      },
      error: function logError(msg) {
        console.log("ERROR", msg);
      },
    },
  };
  return initTracer(config, options);
}

const tracer = initTrace('warehouse-react');
const codec = new ZipkinB3TextMapCodec({ urlEncoding: true });
tracer.registerInjector(opentracing.FORMAT_HTTP_HEADERS, codec);
tracer.registerExtractor(opentracing.FORMAT_HTTP_HEADERS, codec);

app.use('/api/*', (req, res, next) => {
    const span = tracer.startSpan(req.baseUrl);
    
    const spanCtx = {};
    tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, spanCtx);
    req.spanContext = spanCtx;

    for(let i = 0; i < req.rawHeaders.length; i += 2){
        span.setTag(`request.header.${req.rawHeaders[i]}`, req.rawHeaders[i+1]);
    }
    span.setTag(`request.query`, req.query);

    res.on('finish', () => {
        span.setTag(`http.response_code`, res.statusCode);
        span.finish();
    });

    next();
});

app.use('/api', createProxyMiddleware({
  target: process.env.WAREHOUSE_URL || 'http://localhost:3001',
  changeOrigin: true,
  onProxyReq(proxyReq, req, res) {
    Object.keys(req.spanContext).forEach((k) => {
      proxyReq.setHeader(k, req.spanContext[k]);
    });
  }
}));

app.use(express.static('build'))

app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`)
})
