"use strict";

var iptd = require('./lib/iptd.js'),
    log  = require('./lib/log.js'),
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
  log.warn('Failed to parse JSON, falling back to eval');

  eval('cjdroute = ' + cjdroute);
}

config.cjdns = cjdnsadmin;
config.cjdns.pubkey = cjdroute.publicKey;

var iptdServer = new iptd(config);
iptdServer.listen();

/*
 * Reload IPTd on SIGHUP
 * This is useful if CJDNS has crashed or restarted
 * and you want to load all the registered users into
 * CJDNS again.
 */
process.on('SIGHUP', function() {
  log.info('SIGHUP recieved, reloading...');
  iptdServer.reload(function(err, result) {
    if (err) {
      throw new Error(err);
    }
  });
});
