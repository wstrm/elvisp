"use strict";

var elvisp = require('./lib/elvisp.js'),
    log  = require('./lib/log.js'),
    nconf = require('nconf'),
    fs   = require('fs');

function parseConfig(config) {

  nconf.use('memory');
  nconf.load();

  nconf.argv().env();
    
  nconf.file({
    file: './config.json'
  });

  nconf.defaults({
    "password": "testPassword",
    "listen": 4132,
    "iptunnel": {
      "range": [
        ["2a03", "b0c0", "0002", "00d0", "0000", "0000", "01c0", "f000"],
        ["2a03", "b0c0", "0002", "00d0", "0000", "0000", "01c0", "f00f"]
      ],
      "prefix": 0
    },
    "cjdns": {
      "password": "ycdzz73bn17k22c017xtdxgmq7kn7xq",
      "pubkey": "lpu15wrt3tb6d8vngq9yh3lr4gmnkuv0rgcd2jwl5rp5v0mhlg30.k",
      "port": 11234,
      "address": "127.0.0.1"
    },
    "db": "./db"
  });

  return nconf;
}

var config = parseConfig('./config.json');

var elvispServer = new elvisp(config);
elvispServer.listen();

/*
 * Reload Elvisp on SIGHUP
 * This is useful if cjdns has crashed or restarted
 * and you want to load all the registered users into
 * cjdns again.
 */
process.on('SIGHUP', function() {
  log.info('SIGHUP recieved, reloading...');
  elvispServer.reload(function(err, result) {
    if (err) {
      throw new Error(err);
    }
  });
});
