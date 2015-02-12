const MAX_BLOCK_SIZE = 65535;

var IPv6 = {};


/* Add a bit to a array based IPv6 address
 * @param   address   Address to add one bit to
 * @param   callback  callback(newAddress)
 * @results function(newAddress)
 */
IPv6.addBit = function(address, callback) {
  var _this = this;

  this.toDec(address, function(err, dAddress) {
    for(var i = (dAddress.length - 1); i >= 0; i--) {
      if (dAddress[i] === MAX_BLOCK_SIZE && dAddress[i - 1] < MAX_BLOCK_SIZE) { // Filled
        dAddress[i - 1]++; // Add bit to above

        // Reset all the other ones behind
        for(var j = i; j < dAddress.length; j++) {
          dAddress[j] = 0;
        }
      
        break;
      } else if (address[i] < MAX_BLOCK_SIZE) {
        dAddress[i]++;
        break;
      }
    }

    _this.toHex(dAddress, function(err, hAddress) {
      return callback(address);
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
 * Get the common prefix for two addresses (start and end address)
 * @param   start     The start address
 * @param   end       The end address
 * @param   callback  callback(result)
 * @returns function(result)
 */
IPv6.commonPrefix = function(start, end, callback) {
  const blockLen = 8;
  const bitLen = 4;

  var result = []
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
IPv6.remPrefix = function(prefix, address) {
  address.splice(0, prefix.length)
  return address;
}


/*
 * Expand a string IPv6 address and output as array splitted by ':'
 * @param address   IPv6 address string
 * @param callback  callback(error)
 * @returns         function(error)
*/
IPv6.expandAddress = function(address, callback) {
  var finalAddress = [];
  var addressArray = address.split(':');
  var addressLen = addressArray.length;

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
      currentBlock += '0';
      blockLen = currentBlock.length;
    }

    finalAddress.push(currentBlock);
  }

  return callback(finalAddress);
}


/*
 * Validate IPv6 address
 * @param   address IPv6 as string
 * @returns true or false
 */
IPv6.validAddress = function(address) {
  var validator = new RegExp("^(::)?((?:[0-9a-f]+::?)+)?([0-9a-f]+)?(::)?$", 'i')
  return address.match(validator);
}


module.exports = IPv6;
