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
    if (process.env['NODE_DEBUG'] && (process.env['NODE_DEBUG'].split(',')).indexOf('iptd') > -1) {
      console.log('%s - Debug: %s', timeNow(), msg);
    }
  }
};

module.exports = log;
