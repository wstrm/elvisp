/*
 * THIS IS A PROOF OF CONCEPT
 */
var net   = require('net'),
    ipv6  = require('./ipvx').IPv6,
    ipv4  = require('./ipvx').IPv4,
    cjdns = require('cjdnsadmin');


/*
 * Load configuration and cjdnsadmin
 * @returns null
 */
var IPTD = function IPTD (config) {
  var defaultConfig = {
    listen: 4132,
    password: 'examplePassword',
    iptunnel: {
      range: [
        '2a03:b0c0:2:d0::1c0:f000',
        '2a03:b0c0:2:d0::1c0:f00f'
      ]
    },
    cjdns: {
      config: process.env['HOME'] + '/.cjdnsadmin'
    }
  };

  this.config = this.helpers().mergeObj(defaultConfig, config);
  this.cjdns = new cjdns(this.config.cjdns.config);

};

/*
 * Initializes the TCP server and listens on port defined by config or argument
 * @param optListen optional argument with port number or string with socket path
 * @returns null
 */
IPTD.prototype.listen = function listen (optListen) {
  var _this = this;
  
  if (optListen) {
    this.config.listen = optListen; // override port in config with optional
  }

  var server = net.createServer(function(connection) { //'connection' listener
    console.log('Client connected');
    
    connection.on('end', function() {
      console.log('Client disconnected');
    });

    connection.on('data', function(data) {

      try {
        data = JSON.parse(data);
      } catch (err) {
        return new Error(err);
      }

      if (data.password === _this.config.password) {
        console.log('Client authenticated');
        _this.register(data);
      }
    });

    connection.write('Connection established\r\n');
    connection.pipe(connection);
  });

  server.listen(this.config.listen, function() { //'listening' listener
    console.info('Server bound');
  });
  
};

/*
 * Validate and handle new user registrations
 * @param data  object with the following data (misc is optional)
 *              {
 *                password: 'password',
 *                pubkey: 'pubkey.k',
 *                misc: 'some random information'
 *              }
 */
IPTD.prototype.register = function register (data) {
  var _this = this;
  
  // Load database to get current users and their ip's
  // Get IPv6/IPv4 from configurable range
  // Add user to CJDNS ip tunnel
  // Return success with clearnet ip or failure with error

  console.log('Registration request for pubkey:', data.pubkey);
};

/*
 * Helper functions
 * @returns objects with helper functions
 */
IPTD.prototype.helpers = function helpers () {
  return {
    /*
     * Overwrites obj1's values with obj2's and adds obj2's if non existent in obj1
     * @param obj1
     * @param obj2
     * @returns obj3 a new object based on obj1 and obj2
     */
    mergeObj: function mergeObj(obj1,obj2){
      var obj3 = {};
      for (var attrname in obj1) { obj3[attrname] = obj1[attrname]; }
      for (var attrname in obj2) { obj3[attrname] = obj2[attrname]; }
      return obj3;
    }
  };
};

module.exports = IPTD;
