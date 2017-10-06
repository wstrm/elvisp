go-cjdns/admin
==============

Package admin provides methods to access a running cjdns instance via the admin tcp socket. It not only allows you to send any command and receive the response but it also provides convenience functions. It relies on go-cjdns/config for loading of the configuration data and getting the IP address, port, and passsword for the admin interface.