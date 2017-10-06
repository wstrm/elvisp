go-cjdns/config
===============

Package config allows easy loading, manipulation, and saving of cjdns configuration files. It provides two functions for fetching the configuration data from the file:

* `LoadMinConfig()` returns a pointer to a Config struct consisting of only the basic fields required of a cjdns configuration file
* `LoadExtConfig()` returns a map[string]interface consisting of all the data found in the configuration file. 

The reason for the two different methods of accessing the data is because the cjdns configuration file, being JSON, allows you to add any number of arbitrary fields to it. This is useful for keeping connection details and passwords organized by adding "name" and "location" fields to them. However, since the content of the file beyond what cjdns expects is unknown it is impossible to put the custom fields in to a Go structure. You are therefore responsible for working your way through the map[string]interface if you want to work with custom fields. 

This package also contains a function `SaveConfig()` which will accept either a pointer to the minimal Config structure or the extended map[string]interface data and save it to a file. 

**NOTE:** cjdns supports C style comments in the JSON configuration file.  Comments that were present in it  will not be restored upon saving with this package.

### Example

     package main
	
     import (
          "github.com/ehmry/go-cjdns/config"
     )
	
     func main() {
	
		//Load the extended config map
		extConf, err := config.LoadExtConfig("/etc/cjdroute.conf")
		if err != nil {
               panic(err)
		} 
		
		// Do stuff
		    
		// Save changes
		err = config.SaveConfig("output.conf", extConf, 0666) 
		if err != nil {
               panic(err)
		}
		
		//Load the Config struct
		minConf, err := config.LoadMinConfig("/etc/cjdroute.conf")
		if err != nil {
               panic(err)
		}
		
		// Do stuff
		
		// Save changes
		err = config.SaveConfig("output.conf", *minConf, 0666)
		if err != nil {
               panic(err)
		}
	}
