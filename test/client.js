var net = require('net');


var cnt = 0;
setInterval(function() {
  var password = 'testPassword';

  var json = JSON.stringify({
    password: password,
    pubkey: 'lpu15wrt3tb6d8vngq9yh3lr4gmnkuv0rgcd2jwl5rp5v0mhlg3' + cnt + '.k',
    misc: 'Cellphone 1231231232, Email john@doe'
  });
  
  var client = net.connect({ port: 4132 },
      function() { //'connect' listener
        console.log('connected to server!');
        client.write(json);
      });

  client.on('data', function(data) {
    console.log(data.toString());
    //client.end();
  });

  client.on('end', function() {
    console.log('disconnected from server');
  });

  cnt++;
}, 1000);
