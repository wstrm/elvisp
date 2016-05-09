/*
 * You may redistribute this program and/or modify it under the terms of
 * the GNU General Public License as published by the Free Software Foundation,
 * either version 3 of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

"use strict";

/*
 * THIS IS A PROOF OF CONCEPT
 */
var net   = require('net'),
    fs    = require('fs'),
    log   = require('./log'),
    ipv6  = require('./ipvx').IPv6,
    ipv4  = require('./ipvx').IPv4,
    cjdns = require('cjdns-admin'),
    db    = require('./db');

const DB_DIR_MODE = 511;

/**
 * Check the configuration object
 * @param   config  Object, configuration
 * @returns null
 */
function checkConfig(config) {
  
  if (!config || typeof config !== 'object')
    throw new Error('No configuration object specified');

  if (!config.password || typeof config.password !== 'string')
    throw new Error('Invalid password');

  if (config.iptunnel) {
    if (config.iptunnel.range[0].length !== 8 || config.iptunnel.range[0].length !== 8)
      throw new Error('Invalid IPTunnel range');

    if (typeof config.iptunnel.prefix !== 'number')
      throw new Error('Invalid IPTunnel prefix');

  } else
    throw new Error('Invalid IPTunnel config');

  if (typeof config.db !== 'string')
    throw new Error('Invalid database path');

  if (!config.listen || typeof config.listen !== 'number')
    throw new Error('Invalid listen address/port');

  return;
}


/**
 * Throw error
 * @param   err Error to throw, won't do anything if null
 * @returns null
 */
function throwErr(err) {
  if (err) {
    throw new Error(err);
  }

  return;
}

/*
 * Load configuration, cjdnsadmin and initialize
 * @returns null
 */
var IPTD = function IPTD (config) {
  checkConfig(config);

  this.config = config;
  this.cjdns  = cjdns.Admin(this.config.cjdns);
  this.db     = new db(this.config.db);

  // Get IP address version, probably going to change it to support both protocols at the same time.
  if (ipv6.validAddress(this.config.iptunnel.range[0]) &&
      ipv6.validAddress(this.config.iptunnel.range[1])) {
    this.config.iptunnel.type = 6;

  } else if (ipv4.validAddress(this.config.iptunnel.range[0]) && 
             ipv4.validAddress(this.config.iptunnel.range[1])) {
    throwErr('IPv4 is not implemnted');
    this.config.iptunnel.type = 4;

  } else {
    throwErr('Unable to get IP version');
  }

  // Initialize database
  let elvisp = this;
  if (!fs.existsSync(elvisp.config.db)) {
    fs.mkdirSync(elvisp.config.db, DB_DIR_MODE);
    
    this.lastIP = elvisp.config.iptunnel.range[0];
    elvisp.db.write('last_ip', { lastIP: this.lastIP }, throwErr);
  
  } else {
    elvisp.db.read('last_ip', function(err, data) {
      throwErr(err);
      IPTD.prototype.lastIP = data.lastIP;
    });
  }

  IPTD.prototype.getLastIP = function() {
    return this.lastIP;
  };

  elvisp.init();
};

/*
 * Initializes the TCP server and listens on port defined by config or argument
 * @param optListen optional argument with port number or string with socket path
 * @returns null
 */
IPTD.prototype.listen = function listen (optListen) {
  if (optListen) {
    this.config.listen = optListen; // override port in config with optional
  }
  
  let elvisp = this;

  var server = net.createServer(function(connection) { //'connection' listener
    log.info('Client connected');

    connection.on('end', function() {
      log.info('Client disconnected');
    });
  
    connection.on('error', function(err) {
      log.error(err);
    });

    connection.on('data', function(data) {

      try {
        data = JSON.parse(data);
      } catch (err) {
        log.warn('Invalid JSON');
       
        connection.write(JSON.stringify({
          error: 'Invalid JSON',
          status: 0
        }) + '\r\n');
       
        connection.end();
        return;
      }

      if (data.password === elvisp.config.password) {
        log.info('Client authenticated');

        elvisp.register(data, function(err, address) {
          if (err) {
            log.warn(err);
          }

          connection.write(JSON.stringify({
            error: err,
            data: {
              address: err ? undefined : address,
              pubkey: err ? undefined : elvisp.config.cjdns.pubkey
            },
            status: err ? 0 : 1
          }) + '\r\n');

          connection.end();
          return;
        });
      } else {
        log.warn('Invalid password');
       
        connection.write(JSON.stringify({
          error: 'Invalid password',
          status: 0
        }) + '\r\n');
        
        connection.end();
        return;
      }
    });
  });

  server.listen(elvisp.config.listen, function() { //'listening' listener
    log.info('Server bound');
  });

};


/**
 * Write user to the database/filesystem
 * @param   address Array, IP address
 * @param   pubkey  String, Public key for user
 * @param   misc    String, Optional - Misc information
 * @returns null
 */
IPTD.prototype.writeUser = function writeUser(address, pubkey, misc) {
  this.db.write(pubkey, {
    ip: address,
    misc: misc
  }, throwErr);

  // Update last IP in database and ram
  this.lastIP = address;
  this.db.write('last_ip', { lastIP: address }, function(err) {
    throwErr(err);

    log.info('Address: ' + address + ' added to leases');
  });

  return;
};


/**
 * Add user to CJDNS using cjdns admin API
 * @param   address   Array, IP address
 * @param   pubkey    String, Public key for user
 * @param   callback  Function, function(err, stringAddress)
 * @returns function(err, stringAddress)
 */
IPTD.prototype.addUser = function addUser(address, pubkey, callback) {
  let elvisp = this,
      cjdns  = elvisp.cjdns;

  ipv6.toString(address, function(err, stringAddress) {
    throwErr(err);

    var registration = cjdns.ipTunnel.allowConnection({ 
      ip6Address: stringAddress,
      ip6Prefix: elvisp.config.iptunnel.prefix,
      publicKeyOfAuthorizedNode: pubkey
    });

    cjdns.once(registration, function(res) {
      if (res.errors.length > 0) {
        return callback(res.errors);
      } else {
        return callback(null, stringAddress);
      }
    });
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
  let elvisp = this,
      range = elvisp.config.iptunnel.range;

  // Load database to get current users and their ip's
  // Get IPv6/IPv4 from configurable range
  // Add user to CJDNS ip tunnel
  // Return success with clearnet ip or failure with error

  log.info('Info: Registration request for pubkey:', data.pubkey);


  if (!data.pubkey) {
    return callback('No pubkey defined');
  }

  elvisp.db.exists(data.pubkey, function(result) {

    if (!result) {
      log.info('Unique pubkey, will register ' + data.pubkey);
      log.info('Last IP registered ' + elvisp.lastIP);

      if (elvisp.config.iptunnel.type === 6)   {
        log.info('Using IP type 6');

        ipv6.addBit(elvisp.getLastIP(), function(err, address) {
          throwErr(err);
          log.debug('Bit added to address: ' + address);

          ipv6.inRange(address, range, function(err, result) {
            throwErr(err);

            if (!result) {
              return callback('No available IP addresses');
            }
            
            elvisp.writeUser(address, data.pubkey, data.misc);
            elvisp.addUser(address, data.pubkey, callback);
          });
        });

      } /*else if (iptd.config.iptunnel.type === 4) {

        if (lastIP === iptd.config.iptunnel.range[1]) {
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
      }*/

    } else {
      return callback('User already exist, will not add');
    }
  });

};

/*
 * Initialize/reload all the registered users in the db, and add to CJDNS IPTunnel.
 * @returns null or error
 */
IPTD.prototype.init =
IPTD.prototype.reload = function init(callback) {
  let elvisp = this;

  elvisp.db.read(null, function(err, data) {
    if (err) {
      log.debug('Unable to load users:', err);
    }
    
    if (data && data.ip && data.pubkey) {
      log.debug('Loading user', data.pubkey, 'with IP', data.ip);
      elvisp.addUser(data.ip, data.pubkey);
    }
  });
};

module.exports = IPTD;
