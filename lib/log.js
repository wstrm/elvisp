"use strict";

var log = {
  error: function(msg) {
    console.error('%s - Error: %s', this.date(), msg);
  },

  warn: function(msg) {
    console.warn('%s - Warning: %s', this.date(), msg);
  },

  info: function(msg) {
    console.info('%s - Info: %s', this.date(), msg);
  },

  debug: function(msg) {
    console.log(process.env['NODE_DEBUG']);
    if (process.env['NODE_DEBUG'] && (process.env['NODE_DEBUG'].split(',')).indexOf('iptd')) {
      console.log('%s - Debug: %s', this.date(), msg);
    }
  },

  date: function() {
    let timeNow = new Date().toISOString();

    return timeNow;
  }
};

module.exports = log;
