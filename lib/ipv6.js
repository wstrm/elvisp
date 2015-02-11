var IPv6 = {};


/*
 * FIXME:
 * This method is limited by the highest Integer/Number in JavaScript (2^53 while an IPv6 goes up to 2^128)
 */
/*
 * Add a bit to a array based IPv6 address
 * @param   address   Address to add one bit to
 * @param   callback  callback(newAddress)
 * @results function(newAddress)
 */
IPv6.addBit = function(address, callback) {
  if (typeof address === 'object') {
    address = address.join('');
  }
  
  var newAddress = (((parseInt(address, 16) + 1).toString(16)).match(/.{1,4}/g));
  return callback(newAddress);
}


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
