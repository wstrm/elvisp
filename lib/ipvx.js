const IPV6_MAX_BLOCK_SIZE = 65535;
const IPV4_MAX_BLOCK_SIZE = 255;

var IPv6 = {};
var IPv4 = {};
var IPvX = {};


//////////////////// IPvX ////////////////////
/*
 * Validate IPv6/IPv4 address
 * @param   address IPv6/IPv4 as string or array
 * @param   type    Number(4) for IPv4 and Number(6) for IPv6
 * @returns true or false
 */
IPvX.validAddress = function(address, type) {
  var regex;

  if (typeof address === 'object') {
    if(type === 4) {
      address = address.join('.');
    } else if (type === 6) {
      address = address.join(':');
    }
  }

  if (type === 4) { // IPv4
    regex = new RegExp(
      "^([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\." +
      "([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\." +
      "([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\." +
      "([01]?\\d\\d?|2[0-4]\\d|25[0-5])$", 'i'
    );
  } else if (type === 6) { // IPv6
    regex = new RegExp("^(::)?((?:[0-9a-f]+::?)+)?([0-9a-f]+)?(::)?$", 'i');
  }

  return (address.match(regex) ? true : false);
};


/*
 * Get the common prefix for two addresses (start and end address)
 * @param   start     The start address
 * @param   end       The end address
 * @param   callback  callback(result)
 * @returns function(result)
 */
IPvX.commonPrefix = function(start, end, callback) {
  var blockLen, bitLen;

  if (IPv6.validAddress(start) && IPv6.validAddress(end)) {
    blockLen = 8;
    bitLen = 4;
  } else if (IPv4.validAddress(start) && IPv4.validAddress(end)) {
    blockLen = 4;
    bitLen = 3;
  }

  var result = [];
  for (var block = 0; block < blockLen; block++) {
    for (var bit = 0; bit < bitLen; bit++) {

      if (start[block].charAt(bit) !== end[block].charAt(bit)) {
        return callback(result);
      }
    }

    result.push(start[block]);
  }
}


/*
 * Remove prefix from address
 * @param prefix  The prefix that should be removed
 * @param address The address that the prefix should be removed from
 * @returns       Suffix for the address
 */
IPvX.remPrefix = function(prefix, address) {
  address.splice(0, prefix.length);
  return address;
}



//////////////////// IPv6 ////////////////////

/* Add a bit to a array based IPv6 address
 * @param   address   Address to add one bit to
 * @param   callback  callback(newAddress)
 * @results function(newAddress)
 */
IPv6.addBit = function(address, callback) {
  var _this = this;

  this.toDec(address, function(err, dAddress) {
    for(var i = (dAddress.length - 1); i >= 0; i--) {
      if (dAddress[i] === IPV6_MAX_BLOCK_SIZE && dAddress[i - 1] < IPV6_MAX_BLOCK_SIZE) { // Filled
        dAddress[i - 1]++; // Add bit to above

        // Reset all the other ones behind
        for(var j = i; j < dAddress.length; j++) {
          dAddress[j] = 0;
        }

        break;
      } else if (address[i] < IPV6_MAX_BLOCK_SIZE) {
        dAddress[i]++;
        break;
      }
    }

    _this.toHex(dAddress, function(err, hAddress) {
      if (err) {
        return callback(err);
      }

      return callback(null, hAddress);
    });    

  });
};


/*
 * Parse hexadecimal IPv6 array to decimal
 * @param   address   Hexadecimal IPv6 array
 * @param   callback  callback(error, address)
 * @returns function(error, address)
 */
IPv6.toDec = function(address, callback) {
  if (typeof address !== 'object') {
    return callback('Address should be array');
  }

  for(var i = 0; i < address.length; i++) {
    address[i] = parseInt(address[i], 16);
  }

  return callback(null, address);
};


/*
 * Parse decimal IPv6 array to hexadecimal
 * @param   address   Decimal IPv6 array
 * @param   callback  callback(error, address)
 * @returns function(error, address)
 */
IPv6.toHex = function(address, callback) {
  if (typeof address !== 'object') {
    return callback('Address should be array');
  }

  for(var i = 0; i < address.length; i++) {
    address[i] = address[i].toString(16);
  }

  return callback(null, address);
};


/*
 * Expand a string IPv6 address and output as array splitted by ':'
 * @param address   IPv6 address string
 * @param callback  callback(error, finalAddress)
 * @returns         function(error, finalAddress)
 */
IPv6.expandAddress = function(address, callback) {
  var finalAddress = [];
  var addressArray = address.split(':');
  var addressLen = addressArray.length;

  if (IPv6.validAddress(address)) {

    try {
      // We want those redudant 0's
      for (var block = 0; addressLen > block; block++) {
        var currentBlock = addressArray[block];
        var blockLen = currentBlock.length;

        // If we got an empty block, fill up the array with the missing ones
        if (currentBlock === '') {
          var numMissing = 8 - addressLen;

          // Splice array with '0000' for every missing block
          for (var i = 0; i < numMissing; i++) {
            addressArray.splice(block, 0, '0000');
          }

          addressLen = addressArray.length; // Update address length
        }

        // Add missing 0's
        while (blockLen < 4) {
          currentBlock = '0' + currentBlock;
          blockLen = currentBlock.length;
        }

        finalAddress.push(currentBlock);
      }

      return callback(null, finalAddress);
    } catch (err) {
      return callback(err, null);
    }

  } else {
    return callback('Invalid IPv6 address');
  }
};


/*
 * Validate IPv6 address
 * @param   address IPv6 as string
 * @returns true or false
 */
IPv6.validAddress = function(address) {
  return IPvX.validAddress(address, 6);
};


/*
 * Stringify an array based address
 * @param   address   The array based address
 * @param   callback  callback(error, result)
 * @returns function(error, result)
 */
IPv6.toString = function(address, callback) {
  if (typeof address !== 'object') {
    return callback('Address should be array');
  }
  
  address = address.join(':');

  return callback(null, address); 
};


IPv6.commonPrefix = IPvX.commonPrefix;
IPv6.remPrefix = IPvX.remPrefix;



//////////////////// IPv4 ////////////////////

/* Add a bit to a array based IPv4 address
 * @param   address   Address to add one bit to
 * @param   callback  callback(newAddress)
 * @returns function(newAddress)
 */
IPv4.addBit = function(address, callback) {
  var _this = this;

  this.toNum(address, function(err, nAddress) {
    for(var i = (nAddress.length - 1); i >= 0; i--) {
      if (nAddress[i] === IPV4_MAX_BLOCK_SIZE && nAddress[i - 1] < IPV4_MAX_BLOCK_SIZE) { // Filled
        nAddress[i - 1]++; // Add bit to above

        // Reset all the other ones behind
        for(var j = i; j < nAddress.length; j++) {
          nAddress[j] = 0;
        }

        break;
      } else if (nAddress[i] < IPV4_MAX_BLOCK_SIZE) {
        nAddress[i]++;
        break;
      }
    }

    _this.toStringArr(nAddress, function(err, sAddress) {
      if (err) {
        return callback(err);
      }

      return callback(null, sAddress);
    });    
  });
};


/*
 * Convert string based IPv4 array to number based array
 * @param   address Address to convert to numbers
 * @param   callback  callback(error, address)
 * @returns function(error, address)
 */
IPv4.toNum = function(address, callback) {
  if (typeof address !== 'object') {
    return callback('Address should be array');
  }
  
  for(var i = 0; i < address.length; i++) {
    address[i] = parseInt(address[i]);
  }

  return callback(null, address);
};


/*
 * Convert number based IPv4 array to string based array
 * @param   address Address to convert to string
 * @param   callback  callback(error, address)
 * @returns function(error, address)
 */
IPv4.toStringArr = function(address, callback) {
  if (typeof address !== 'object') {
    return callback('Address should be array');
  }
  
  for(var i = 0; i < address.length; i++) {
    address[i] = address[i].toString();
  }

  return callback(null, address);
};


/*
 * Stringify an array based address
 * @param   address   The array based address
 * @param   callback  callback(error, result)
 * @returns function(error, result)
 */
IPv4.toString = function(address, callback) {
  if (typeof address !== 'object') {
    return callback('Address should be array');
  }
  
  address = address.join('.');

  return callback(null, address); 
};


/*
 * Expand a string IPv4 address and output as array splitted by '.'
 * @param   address   IPv4 address string
 * @param   callback  callback(error, finalAddress)
 * @returns function(error, finalAddress)
 */
IPv4.expandAddress = function(address, callback) {
  var finalAddress = address.split('.');

  return callback(finalAddress);
};


/*
 * Validate IPv4 address
 * @param   address IPv4 as string
 * @returns true or false
 */
IPv4.validAddress = function(address) {
  return IPvX.validAddress(address, 4);
};


IPv4.commonPrefix = IPvX.commonPrefix;
IPv4.remPrefix = IPvX.remPrefix;



IPvX.IPv6 = IPv6;
IPvX.IPv4 = IPv4;

module.exports = IPvX;
