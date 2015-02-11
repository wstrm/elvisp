var iptd = require('./lib/iptd.js');

iptdServer = new iptd({
  listen: 4132,
});

iptdServer.listen();
