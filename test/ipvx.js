var assert = require('assert');


describe('IPvX', function() {
  describe('.IPv4', function() {
    var IPv4 = require('../lib/ipvx.js').IPv4;

    const mocks = {
      validAddress: {
        string: '192.168.1.43',
        array: ['192', '168', '1', '43']
      },
      invalidAddress: {
        string: '255.213.523.10',
        array: ['255', '213', '523', '10']
      },
      addressRange: {
        string: [
          '192.168.1.0',
          '192.168.1.255'
        ],
        array: [
          ['192', '168', '1', '0'],
          ['192', '168', '1', '255']
        ]
      },
      prefix: {
        string: '192.168.1',
        array: ['192', '168', '1']
      },
      suffix: {
        string: '43',
        array: ['43']
      }
    };


    describe('.validAddress(address)', function() {
      it('should validate an valid IPv4 address (string)', function() {
      
        assert.equal(IPv4.validAddress(mocks.validAddress.string), true);
      
      });
      
      it('should invalidate an invalid IPv4 address (string)', function() {
      
        assert.equal(IPv4.validAddress(mocks.invalidAddress.string), false);
      
      });
    });


    describe('.commonPrefix(first, second)', function() {
      it('should return the common prefix for two addresses (array)', function() {
        IPv4.commonPrefix(mocks.addressRange.array[0], mocks.addressRange.array[1], function(prefix) {
          
          assert.deepEqual(prefix, mocks.prefix.array);
        
        });
      });
    });

    
    describe('.remPrefix(first, second)', function() {
      it('should remove the prefix from an IPv4 address and return suffix (array)', function() {
        assert.deepEqual(IPv4.remPrefix(mocks.prefix.array, mocks.validAddress.array), mocks.suffix.array);
      });
    });
  });
  
  describe('.IPv6', function() {
    var IPv6 = require('../lib/ipvx.js').IPv6;

    const mocks = {
      validAddress: {
        compressed: {
          string: '2a03:b0c0:2:d0::1c0:f000',
          array: ['2a03', 'b0c0', '2', 'd0', '0', '0', '1c0', 'f000']
        },
        expanded: {
          string: '2a03:b0c0:0002:00d0:0000:0000:01c0:f000',
          array: [ '2a03', 'b0c0', '0002', '00d0', '0000', '0000', '01c0', 'f000' ]
        },
        decimal: [10755, 45248, 2, 208, 0, 0, 448, 61440]
      },
      invalidAddress: {
        compressed: {
          string: '2a03:b0c0:42322:d0::1c0:g000',
          array: ['2a03', 'b0c0', '42322', 'd0', '0', '0', '1c0', 'g000']
        },
        expanded: {
          string: '2a03:b0c0:42322:00d0:0000:0000:01c0:g000',
          array: [ '2a03', 'b0c0', '42322', '00d0', '0000', '0000', '01c0', 'g000' ]
        }
      },
      addressRange: {
        string: [
          '2a03:b0c0:2:d0::1c0:f000',
          '2a03:b0c0:2:d0::1cf:f000'
        ],
        array: [
          ['2a03', 'b0c0', '2', 'd0', '0', '0', '1c0', 'f000'],
          ['2a03', 'b0c0', '2', 'd0', '0', '0', '1cf', 'f000']
        ]
      },
      prefix: {
        string: '2a03:b0c0:2:d0',
        array: ['2a03', 'b0c0', '2', 'd0']
      },
      suffix: {
        string: '::1c0:f000',
        array: ['0', '0', '1c0', 'f000']
      }
    };
    
    
    describe('.validAddress(address)', function() {
      it('should validate an valid IPv6 address (string)', function() {
      
        assert.equal(IPv6.validAddress(mocks.validAddress.compressed.string), true);
      
      });
      
      it('should invalidate an invalid IPv6 address (string)', function() {
      
        assert.equal(IPv6.validAddress(mocks.invalidAddress.compressed.string), false);
      
      });
    });


    describe('.commonPrefix(first, second, callback(prefix))', function() {
      it('should return the common prefix for two addresses (array)', function() {
        IPv6.commonPrefix(mocks.addressRange.array[0], mocks.addressRange.array[1], function(prefix) {
          
          assert.deepEqual(prefix, mocks.prefix.array);
        
        });
      });
    });

    
    describe('.remPrefix(first, second)', function() {
      it('should remove the prefix from an IPv6 address and return suffix (array)', function() {
        assert.deepEqual(IPv6.remPrefix(mocks.prefix.array, mocks.validAddress.compressed.array), mocks.suffix.array);
      });
    });
    
    
    describe('.expandAddress(address, callback(error, resultAddress))', function() {
      it('should expand a compressed IPv6 address (string) and return result address (array)', function() {
        IPv6.expandAddress(mocks.validAddress.compressed.string, function(error, resultAddress) {

          assert.equal(error, null);
          assert.deepEqual(resultAddress, mocks.validAddress.expanded.array);
       
        });
      });
    });


    describe('.toDec(address, callback(error, hexAddress)', function() {
      it('shoud convert IPv6 address to decimal (array)', function() {
        IPv6.toDec(mocks.validAddress.expanded.array, function(error, decAddress) {

          assert.equal(error, null);
          assert.deepEqual(decAddress, mocks.validAddress.decimal);

        });
      });
    });


    describe('.toHex(address, callback(error, hexAddress)', function() {
      it('shoud convert IPv6 address to hexadecimal (array)', function() {
        IPv6.toHex(mocks.validAddress.decimal, function(error, hexAddress) {

          assert.equal(error, null);
          assert.deepEqual(hexAddress, mocks.validAddress.expanded.array);

        });
      });
    });


    describe('.addBit(address, callback(address)', function() {
      it('should add a bit to the address and return callback with address', function() {
        IPv6.addBit(mocks.validAddress.expanded.array, function(address) {

          assert.deepEqual(address, ['2a03', 'b0c0', '2', 'd0', '0', '0', '1c0', 'f001']);
       
        });
      });
    });
  });
});
