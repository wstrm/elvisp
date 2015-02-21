/*
 * THIS IS A PROOF OF CONCEPT
 */
var net   = require('net'),
    fs    = require('fs'),
    ipv6  = require('./ipvx').IPv6,
    ipv4  = require('./ipvx').IPv4,
    cjdns = require('cjdnsadmin');
  
const DB_DIR_MODE = 511;


/*
 * Load configuration, cjdnsadmin and initialize
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
    },
    db: {
      path: __dirname + '/../db'
    }
  };

  this.config = this.helpers().mergeObj(defaultConfig, config);
  this.cjdns = new cjdns(this.config.cjdns.config);

  // Get IP address version, probably going to change it to support both protocols at the same time.
  if (ipv6.validAddress(this.config.iptunnel.range[0]) && ipv6.validAddress(this.config.iptunnel.range[1])) {
    this.config.iptunnel.type = 6;
  } else if (ipv4.validAddress(this.config.iptunnel.range[0]) && ipv4.validAddress(this.config.iptunnel.range[1])) {
    this.config.iptunnel.type = 4;
  } else {
    throw new Error('Unable to get IP version');
  }

  // Initialize database
  var _this = this;
  fs.exists(_this.config.db.path, function(result) {

    if (!result) { // DB Path does not exist, looks like a dry run
      fs.mkdir(_this.config.db.path, DB_DIR_MODE, function(err) {
        if (err) {
          throw new Error(err);
        }

        _this.db().exists('last_ip', function(result) {
          
          if (!result) {
            _this.db().write('last_ip', { lastIP: _this.config.iptunnel.range[0] }, function(err) {
              if (err) {
                throw new Error(err);
              }
            });
          }

          return;
        });

        return;
      });
    }
    
    return;
  });

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
        
          connection.end();
        });
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

  function createPass(salt, callback) {
    return callback('password');
  }

  function addUser(address, password, callback) {
    _this.db().write(data.pubkey, {
      password: password,
      ip: address,
      misc: data.misc
    }, function(err) {
      if (err) {
        return callback(err);
      }

      // Update last IP in database
      _this.db().write('last_ip', { lastIP: address }, function(err) {
        if (err) {
          return callback(err);
        }

        return callback();
      });
    });
  }

  _this.db().exists(data.pubkey, function(result) {
    
    if (!result) {
      console.info('Info: Unique pubkey, will register', data.pubkey);

      _this.db().read('last_ip', function(err, data) {
        if (err) {
          return callback(err);
        }

        var lastIP = data.lastIP;

        if (_this.config.iptunnel.type === 6)   {
          ipv6.expandAddress(lastIP, function(err, address) {

            ipv6.addBit(address, function(err, address) {
              if (err) {
                return callback(err);
              }

              ipv6.toString(address, function(err, address) {
                if (err) {
                  return callback(err);
                }

                createPass(address, function(password) {
                  addUser(address, password, function(err) {
                    if (err) {
                      return callback(err);
                    }
                    
                    return callback(null, {
                      ip: address,
                      password: password
                    });
                  });
                });
              });
            });
          });

        } else if (_this.config.iptunnel.type === 4) {
          ipv4.expandAddress(lastIP, function(err, address) {

            ipv4.addBit(address, function(err, address) {
              if (err) {
                console.error(err);
                return callback(err);
              }

              ipv4.toString(address, function(err, address) {
                if (err) {
                  return callback(err);
                }

                addUser(address, password, function(err) {
                  if (err) {
                    return callback(err);
                  }

                  return callback(null, {
                    ip: address,
                    password: password
                  });
                });
              });
            });
          });
        }
      });
    
    } else {
      return callback('User already exist, will not add');
    }
  });

};

IPTD.prototype.db = function db () {
  var _this = this;

  return {
    read: function read(file, callback) {
      
      if (file === null) {
        var data = {};

        fs.readdir(_this.config.db.path, function(err, files) {
          if (err) {
            return callback(err);
          }

          for(var file = 0; file < files.length; file++) {
            fs.readFile(_this.config.db.path + '/' + files[file], function(err, buffer) {
              if (err) {
                return callback(err);
              }

              data[files[file]] = JSON.parse(buffer.toString());
            });
          }
          
          return callback(null, data);
        });
      
      } else {
        
        fs.readFile(_this.config.db.path + '/' + file, function(err, buffer) {
          if (err) {
            return callback(err);
          }
         
          try { 
            var data = JSON.parse(buffer.toString());

            return callback(null, data);
          } catch (err) {
            return callback(err);
          }
        });
      }
    },

    exists: function exists(file, callback) {
      fs.exists(_this.config.db.path + '/' + file, function(result) {
        return callback(result);
      });
    },
    
    write: function write(file, data, callback) {
      data = JSON.stringify(data);

      fs.writeFile(_this.config.db.path + '/' + file, data, function(err) {
        if (err) {
          return callback(err);
        }

        return callback();
      });
    }
  };
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
