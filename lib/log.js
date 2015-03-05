"use strict";

function timeNow() {
  let time = new Date().toISOString();

  return time;
}

var log = {
  error: function(msg) {
    console.error('%s - Error: %s', timeNow(), msg);
  },

  warn: function(msg) {
    console.warn('%s - Warning: %s', timeNow(), msg);
  },

  info: function(msg) {
    console.info('%s - Info: %s', timeNow(), msg);
  },

  debug: function(msg) {
    console.log(process.env['NODE_DEBUG']);
    if (process.env['NODE_DEBUG'] && (process.env['NODE_DEBUG'].split(',')).indexOf('iptd')) {
      console.log('%s - Debug: %s', timeNow(), msg);
    }
  }
};

module.exports = log;
