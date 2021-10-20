const http = require('http');

const craftPayload = (req, data) => {
  let result = `${req.method} ${req.url} HTTP/${req.httpVersion}\r\n`;
  for (let i = 0; i < req.rawHeaders.length; i += 2) {
    const key = req.rawHeaders[i];
    const value = req.rawHeaders[i+1];
    result += `${key}: ${value}\r\n`;
  }
  result += '\r\n';
  return result + data;
}

const requestListener = async function (req, res) {
  const buffers = [];
  for await (const chunk of req) {
    buffers.push(chunk);
  }
  const data = Buffer.concat(buffers).toString();
  console.log(craftPayload(req, data));
  console.log("============================================");
  res.writeHead(200);
  res.end('Hello, World!');
}

const server = http.createServer(requestListener);
server.listen(5000, '0.0.0.0');
