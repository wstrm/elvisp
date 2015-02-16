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

  if (ipv6.validAddress(this.config.iptunnel.range[0]) && ipv6.validAddress(this.config.iptunnel.range[1])) {
    this.config.iptunnel.type = 6;
  } else if (ipv4.validAddress(this.config.iptunnel.range[0]) && ipv4.validAddress(this.config.iptunnel.range[1])) {
    this.config.iptunnel.type = 4;
  } else {
    throw new Error('Unable to get IP type');
  }

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
    console.info('Info: Client connected');
    
    connection.on('end', function() {
      console.info('Info: Client disconnected');
    });

    connection.on('data', function(data) {

      try {
        data = JSON.parse(data);
      } catch (err) {
        console.warn('Warning: Invalid JSON');
        connection.write(JSON.stringify({
          error: 'Invalid JSON'
        }));
        connection.end();
      }

      if (data.password === _this.config.password) {
        console.info('Info: Client authenticated');
        connection.write(JSON.stringify({
          error: null,
          status: 'Authenticated'
        }));
        _this.register(data, function(err, info) {
        
          if (err) {
            console.warn('Warning:', err);
          }

          connection.write(JSON.stringify({
            error: err,
            data: info
          }));
        });
        connection.end();
      } else {
        console.warn('Warning: Invalid password');
        connection.write(JSON.stringify({
          error: 'Invalid password'
        }));
        connection.end();
      }
    });

    connection.write('Connection established\r\n');
    connection.pipe(connection);
  });

  server.listen(this.config.listen, function() { //'listening' listener
    console.info('Info: Server bound');
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
IPTD.prototype.register = function register (data, callback) {
  var _this = this;
  
  // Load database to get current users and their ip's
  // Get IPv6/IPv4 from configurable range
  // Add user to CJDNS ip tunnel
  // Return success with clearnet ip or failure with error

  console.info('Info: Registration request for pubkey:', data.pubkey);

  /* 
   * DATABASE STUFF
   */
  mockDb = [
    {
      password: 'pass',
      pubkey: 'pub1.k',
      ip: '2a03:b0c0:2:d0::1c0:f003',
      misc: 'some random info'
    },
    {
      password: 'pass',
      pubkey: 'pub2.k',
      ip: '2a03:b0c0:2:d0::1c0:f004',
      misc: 'some random info'
    }
  ];

  function checkUnique(pubkey, callback) {
    for (var i = 0; i < mockDb.length; i++) {
      if (mockDb[i].pubkey === pubkey) {
        return callback(false);
      }
    }

    return callback(true);
  }

  function getLastUser() {
    return mockDb[(mockDb.length - 1)];
  }

  function createPass(salt, callback) {
    return callback('password');
  }

  function addUser(address, password) {
    mockDb[mockDb.length] = {
      password: password,
      pubkey: data.pubkey,
      ip: address,
      misc: data.misc
    };
  }

  checkUnique(data.pubkey, function(result) {
    if (result) {
      console.info('Info: Unique pubkey, will register', data.pubkey);

      if (_this.config.iptunnel.type === 6)   {
        ipv6.expandAddress(getLastUser().ip, function(err, address) {

          ipv6.addBit(address, function(err, address) {
            if (err) {
              return callback(err);
            }

            ipv6.toString(address, function(err, address) {
              if (err) {
                return callback(err);
              }

              createPass(address, function(password) {
                addUser(address, password);
              });
            });
          });
        });
      } else if (_this.config.iptunnel.type === 4) {
        ipv4.expandAddress(getLastUser().ip, function(err, address) {
        
          ipv4.addBit(address, function(err, address) {
            if (err) {
              console.error(err);
              return callback(err);
            }

            ipv4.toString(address, function(err, address) {
              if (err) {
                return callback(err);
              }

              createPass(address, function(password) {
                addUser(address, password);
              });
            });
          });
        });
      }

      console.log(mockDb);
      return callback(null, {
        ip: mockDb[mockDb.length - 1].ip
      });
    } else {
      return callback('User already exist, will not add');
    }
  });

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
