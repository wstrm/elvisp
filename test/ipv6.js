var IPv6 = require('../lib/ipv6.js');

var start = '2a03:b0c0:2:d0::1c0:f000',
    end = '2a03:b0c0:2:d0::1cf:f000';

function controller(start, end, callback) {

  // Expand and validate the addresses
  if (IPv6.validAddress(start) && IPv6.validAddress(end)) {
    for (var arg = 0; arg < 2; arg++) {

      IPv6.expandAddress(arguments[arg], function resultAddress (address) {
        if (arg === 0) { // Start address
          start = address;
        } else { // End address
          end = address;
        }
      });
    }
  } else {
    return callback('Invalid IPv6 address');
  }

  IPv6.commonPrefix(start, end, function resultPrefix (prefix) {

    var suffix = IPv6.remPrefix(prefix, start);
    console.log(suffix);
    var suffix = ['0x0001', '0xffff', '0xffff'];
    while(true) {
      IPv6.addBit(suffix, function resultAddress (address) {
        suffix = address;
      });

      console.log(suffix);
    }
    return callback(null, prefix);
  });

}

controller(start, end, function result(err, range) {
  if (err) {
    throw new Error(err);
  }
  
  console.log(range);
});
