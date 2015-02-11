var net = require('net');

var password = 'xK43X68fP4SGS17DsaKV9HnA99JKPKz5';

var json = JSON.stringify({
  password: password,
  pubkey: 'testpubkey.k',
  misc: 'Cellphone 1231231232, Email john@doe'
});

var client = net.connect({ port: 4132 },
    function() { //'connect' listener
      console.log('connected to server!');
      client.write(json);
    });

client.on('data', function(data) {
  console.log(data.toString());
  client.end();
});

client.on('end', function() {
  console.log('disconnected from server');
});
