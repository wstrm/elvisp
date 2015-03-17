"use strict";

var iptd = require('./lib/iptd.js'),
    fs   = require('fs');

var config = {
  password: 'testPassword',
  listen: 4132,
  iptunnel: {
    range: [
      ['2a03', 'b0c0', '0002', '00d0', '0000', '0000', '01c0', 'f000'],
      ['2a03', 'b0c0', '0002', '00d0', '0000', '0000', '01c0', 'f00f']
    ],
    prefix: 0
  },
  db: __dirname + '/db'
};

var cjdnsadmin = JSON.parse(fs.readFileSync(process.env['HOME'] + '/.cjdnsadmin'));
var cjdroute = fs.readFileSync(cjdnsadmin.config);

try {
  cjdroute = JSON.parse(cjdroute);
} catch (err) {
  console.log('Failed to parse JSON, falling back to eval');

  eval('cjdroute = ' + cjdroute);
}

config.cjdns = cjdnsadmin;
config.cjdns.pubkey = cjdroute.publicKey;

var iptdServer = new iptd(config);

iptdServer.listen();
