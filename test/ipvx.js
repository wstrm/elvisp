var assert = require('assert');


describe('IPvX', function() {
  describe('.IPv4', function() {
    var IPv4 = require('../lib/ipvx.js').IPv4;

    
    describe('.expandAddress(address, callback(finalAddress)', function() {
      it('should expand an IPv4 address an return callback with finalAddress', function() {
        IPv4.expandAddress('192.168.1.43', function(finalAddress) {
          assert.deepEqual(['192', '168', '1', '43'], finalAddress);
        });
      });
    });


    describe('.validAddress(address)', function() {
      it('should validate an valid IPv4 address (string)', function() {
      
        assert.equal(IPv4.validAddress('192.168.1.43'), true);
      
      });
      
      it('should invalidate an invalid IPv4 address (string)', function() {
      
        assert.equal(IPv4.validAddress('255.213.523.10'), false);
      
      });
    });


    describe('.commonPrefix(first, second)', function() {
      it('should return the common prefix for two addresses (array)', function() {
        IPv4.commonPrefix(['192', '168', '1', '0'],  ['192', '168', '1', '255'], function(prefix) {
          
          assert.deepEqual(prefix, ['192', '168', '1']);
        
        });
      });
    });
    
    
    describe('.toStringArr(address, callback(error, address)', function() {
      it('should return callback with string based array (array)', function() {
        IPv4.toStringArr([192, 168, 1, 43], function(err, address) {
          
          assert.equal(err, null);
          assert.deepEqual(address, ['192', '168', '1', '43']);
        
        });
      });
    });
    
    
    describe('.toString(address, callback(error, address)', function() {
      it('should return callback with stringified address(string)', function() {
        IPv4.toString(['192', '168', '1', '43'], function(err, address) {
          
          assert.equal(err, null);
          assert.deepEqual(address, '192.168.1.43');
        
        });
      });
    });

    
    describe('.remPrefix(first, second)', function() {
      it('should remove the prefix from an IPv4 address and return suffix (array)', function() {
        assert.deepEqual(IPv4.remPrefix(['192', '168', '1'], ['192', '168', '1', '43']), ['43']);
      });
    });
  
  
    describe('.addBit(address, callback(address)', function() {
      it('should add a bit to 192.168.1.43 and return callback with address', function() {
        IPv4.addBit(['192', '168', '1', '43'], function(err, address) {

          assert.equal(err, null);
          assert.deepEqual(address, ['192', '168', '1', '44']);
       
        });
      });

      it('should add a bit to 192.168.1.255 and return callback with address', function() {
        IPv4.addBit(['192', '168', '1', '255'], function(err, address) {

          assert.equal(err, null);
          assert.deepEqual(address, ['192', '168', '2', '0']);
       
        });
      });
    });
  });
 


  describe('.IPv6', function() {
    var IPv6 = require('../lib/ipvx.js').IPv6;
   

    describe('.validAddress(address)', function() {
      it('should validate an valid IPv6 address (string)', function() {
      
        assert.equal(IPv6.validAddress('2a03:b0c0:2:d0::1c0:f000'), true);
      
      });
      
      it('should invalidate an invalid IPv6 address (string)', function() {
      
        assert.equal(IPv6.validAddress('2a03:b0c0:42322:d0::1c0:g000'), false);
      
      });
    });


    describe('.commonPrefix(first, second, callback(prefix))', function() {
      it('should return the common prefix for two addresses (array)', function() {
        IPv6.commonPrefix(['2a03', 'b0c0', '2', 'd0', '0', '0', '1c0', 'f000'], ['2a03', 'b0c0', '2', 'd0', '0', '0', '1cf', 'f000'], function(prefix) {
          
          assert.deepEqual(prefix, ['2a03', 'b0c0', '2', 'd0', '0', '0']);
        
        });
      });
    });

    
    describe('.remPrefix(first, second)', function() {
      it('should remove the prefix from an IPv6 address and return suffix (array)', function() {
        assert.deepEqual(IPv6.remPrefix(['2a03', 'b0c0', '2', 'd0'], ['2a03', 'b0c0', '2', 'd0', '0', '0', '1c0', 'f000']), ['0', '0', '1c0', 'f000']);
      });
    });
    
    
    describe('.expandAddress(address, callback(error, resultAddress))', function() {
      it('should expand a compressed IPv6 address (string) and return result address (array)', function() {
        IPv6.expandAddress('2a03:b0c0:2:d0::1c0:f000', function(error, resultAddress) {

          assert.equal(error, null);
          assert.deepEqual(resultAddress, [ '2a03', 'b0c0', '0002', '00d0', '0000', '0000', '01c0', 'f000' ]);
       
        });
      });
    });


    describe('.toHex(address, callback(error, hexAddress)', function() {
      it('shoud convert IPv6 address to hexadecimal (array)', function() {
        IPv6.toHex([10755, 45248, 2, 208, 0, 0, 448, 61440], function(error, hexAddress) {

          assert.equal(error, null);
          assert.deepEqual(hexAddress, ['2a03', 'b0c0', '2', 'd0', '0', '0', '1c0', 'f000']);

        });
      });
    });


    describe('.addBit(address, callback(address)', function() {
      it('should add a bit to 2a03:b0c0:2:d0::1c0:f000 and return callback with address', function() {
        IPv6.addBit([ '2a03', 'b0c0', '0002', '00d0', '0000', '0000', '01c0', 'f000' ], function(err, address) {

          assert.equal(err, null);
          assert.deepEqual(address, ['2a03', 'b0c0', '2', 'd0', '0', '0', '1c0', 'f001']);
       
        });
      });
      
      it('should add a bit to 2a03:b0c0:2:d0::1c0:ffff and return callback with address', function() {
        IPv6.addBit([ '2a03', 'b0c0', '0002', '00d0', '0000', '0000', '01c0', 'ffff' ], function(err, address) {

          assert.equal(err, null);
          assert.deepEqual(address, ['2a03', 'b0c0', '2', 'd0', '0', '0', '1c1', '0']);
       
        });
      });
    });


    describe('.toDec(address, callback(error, decAddress)', function() {
      it('shoud convert IPv6 address to decimal (array)', function() {
        IPv6.toDec([ '2a03', 'b0c0', '0002', '00d0', '0000', '0000', '01c0', 'f000' ], function(error, decAddress) {

          assert.equal(error, null);
          assert.deepEqual(decAddress, [10755, 45248, 2, 208, 0, 0, 448, 61440]);

        });
      });
    });
    
        
    describe('.toString(address, callback(error, stringAddress)', function() {
      it('should stringify IPv6 address', function() {
        IPv6.toString([ '2a03', 'b0c0', '0002', '00d0', '0000', '0000', '01c0', 'f000' ], function(error, stringAddress) {

          assert.equal(error, null);
          assert.deepEqual(stringAddress, '2a03:b0c0:0002:00d0:0000:0000:01c0:f000');

        });
      });
    });
  });
});
