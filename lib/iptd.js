"use strict";

/*
 * THIS IS A PROOF OF CONCEPT
 */
var net   = require('net'),
    fs    = require('fs'),
    log   = require('./log'),
    ipv6  = require('./ipvx').IPv6,
    ipv4  = require('./ipvx').IPv4,
    cjdns = require('cjdnsadmin');

const DB_DIR_MODE = 511;

/*
 * Load configuration, cjdnsadmin and initialize
 * @returns null
 */
var IPTD = function IPTD (config) {
  this.config = {
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

  this.config = this.helpers.mergeObj(this.config, config);
  this.cjdns = new cjdns(this.config.cjdns.config);

  // Get IP address version, probably going to change it to support both protocols at the same time.
  if (ipv6.validAddress(this.config.iptunnel.range[0]) && ipv6.validAddress(this.config.iptunnel.range[1])) {
    this.config.iptunnel.type = 6;

  } else if (ipv4.validAddress(this.config.iptunnel.range[0]) && ipv4.validAddress(this.config.iptunnel.range[1])) {
    throw new Error('IPv4 is not implemnted');
    this.config.iptunnel.type = 4;

  } else {
    throw new Error('Unable to get IP version');
  }

  // Initialize database
  let _this = this;
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
  let _this = this;

  if (optListen) {
    this.config.listen = optListen; // override port in config with optional
  }

  var server = net.createServer(function(connection) { //'connection' listener
    log.info('Client connected');

    connection.on('end', function() {
      log.info('Client disconnected');
    });

    connection.on('data', function(data) {

      try {
        data = JSON.parse(data);
      } catch (err) {
        log.warn('Invalid JSON');
       
        connection.write(JSON.stringify({
          error: 'Invalid JSON'
        }));
        connection.end();
        return;
      }

      if (data.password === _this.config.password) {
        log.info('Client authenticated');

        connection.write(JSON.stringify({
          error: null,
          status: 'Authenticated'
        }));

        _this.register(data, function(err, info) {
          if (err) {
            log.warn(err);
          }

          connection.write(JSON.stringify({
            error: err,
            data: info
          }));

          connection.end();
          return;
        });
      } else {
        log.warn('Invalid password');
       
        connection.write(JSON.stringify({
          error: 'Invalid password'
        }));
        connection.end();
        return;
      }
    });

    connection.write('Connection established\r\n');
    connection.pipe(connection);
  });

  server.listen(this.config.listen, function() { //'listening' listener
    log.info('Server bound');
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
  let range = this.config.iptunnel.range;
  let _this = this;

  // Load database to get current users and their ip's
  // Get IPv6/IPv4 from configurable range
  // Add user to CJDNS ip tunnel
  // Return success with clearnet ip or failure with error

  log.info('Info: Registration request for pubkey:', data.pubkey);

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
        throw new Error(err);
      }
    });
      
    // Update last IP in database
    _this.db().write('last_ip', { lastIP: address }, function(err) {
      if (err) {
        throw new Error(err);
      }

      return callback();
    });
  }

  if (!data.pubkey) {
    return callback('No pubkey defined');
  }

  _this.db().exists(data.pubkey, function(result) {

    if (!result) {
      log.info('Unique pubkey, will register ' + data.pubkey);

      _this.db().read('last_ip', function(err, data) {
        if (err) {
          throw new Error(err);
        }

        var lastIP = data.lastIP;
        log.info('Last IP registered ' + lastIP);

        if (_this.config.iptunnel.type === 6)   {
          log.info('Using IP type 6');

          ipv6.expandAddress(lastIP, function(err, address) {
            if (err) {
              throw new Error(err);
            }
            log.debug('Expanded address: ' + address);

            ipv6.addBit(address, function(err, address) {
              if (err) {
                throw new Error(err);
              }
              log.debug('Bit added to address: ' + address);

              ipv6.inRange(address, range, function(err, result) {
                if (err) {
                  throw new Error(err);
                }

                if (!result) {
                  return callback('No available IP addresses');
                }

                createPass(address, function(password) {
                  log.debug('Password created: ' + password);
                  
                  ipv6.toString(address, function(err, address) {
                    if (err) {
                      throw new Error(err);
                    }

                    addUser(address, password, function() {
                      log.info('Address: ' + address + ' added to leases');


                      return callback(null, {
                        ip: address,
                        password: password
                      });
                    });
                  });
                });
              });
            });
          });

        } else if (_this.config.iptunnel.type === 4) {

          if (lastIP === _this.config.iptunnel.range[1]) {
            return callback('No available IP addresses');
          }

          ipv4.addBit(lastIP, function(err, address) {
            if (err) {
              console.error(err);
              return callback(err);
            }

            addUser(address, password, function(err) {
              if (err) {
                return callback(err);
              }

              ipv4.toString(address, function(err, address) {
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
        }
      });

    } else {
      return callback('User already exist, will not add');
    }
  });

};

/*
 * Filesystem based database, it's a wrapper for fs that changes
 * root '/' to the database path.
 * @returns object containing read, exists and write functions
 */
IPTD.prototype.db = function db () {
  let _this = this;

  return {

    /*
     * Read from filesystem
     * @param   file      File to read from
     * @param   callback  callback(error, data)
     * @returns function(error, data)
     */
    read: function read(file, callback) {

      if (file === null) {
        let data = {};

        fs.readdir(_this.config.db.path, function(err, files) {
          if (err) {
            return callback(err);
          }

          for(let file = 0; file < files.length; file++) {
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
            let data = JSON.parse(buffer.toString());

            return callback(null, data);
          } catch (err) {
            return callback(err);
          }
        });
      }
    },

    /*
     * Check if file exists
     * @param   file      File to check if exists
     * @param   callback  callback(result)
     * @returns function(result)
     */
    exists: function exists(file, callback) {
      fs.exists(_this.config.db.path + '/' + file, function(result) {
        return callback(result);
      });
    },

      /*
       * Write to filesystem
       * @param   file      File to write to
       * @param   callback  callback(error)
       * @returns function(error)
       */
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
IPTD.prototype.helpers = {
  /*
   * Overwrites obj1's values with obj2's and adds obj2's if non existent in obj1
   * @param obj1
   * @param obj2
   * @returns obj3 a new object based on obj1 and obj2
   */
  mergeObj: function mergeObj(obj1,obj2){
    var obj3 = {};
    for (let attrname in obj1) { obj3[attrname] = obj1[attrname]; }
    for (let attrname in obj2) { obj3[attrname] = obj2[attrname]; }
    return obj3;
  }
};

module.exports = IPTD;
