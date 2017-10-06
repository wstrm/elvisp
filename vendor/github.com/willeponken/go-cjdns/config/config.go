// Package config allows easy loading, manipulation, and saving of cjdns
// configuration files.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
)

// Loads and parses the input file and returns a Config structure with the
// minimal cjdroute.conf file requirements.
func LoadMinConfig(filein string) (*Config, error) {

	//Load the raw JSON data from the file
	raw, err := loadJson(filein)
	if err != nil {
		return nil, err
	}

	// Parse the JSON in to our struct which supports all requried fields for
	// cjdns
	structured, err := parseJSONStruct(raw)
	if err != nil {
		// BUG(inhies): Find a better way of dealing with these errors.
		if e, ok := err.(*json.SyntaxError); ok {
			// BUG(inhies): Instead of printing x amount of characters, print
			// the previous and following 2 lines
			fmt.Println("Invalid JSON")
			fmt.Println("----------------------------------------")
			fmt.Println(string(raw[e.Offset-60 : e.Offset+60]))
			fmt.Println("----------------------------------------")
		} else if _, ok := err.(*json.InvalidUTF8Error); ok {
			fmt.Println("Invalid UTF-8")
		} else if e, ok := err.(*json.InvalidUnmarshalError); ok {
			fmt.Println("Invalid unmarshall type", e.Type)
			fmt.Println(err)
		} else if e, ok := err.(*json.UnmarshalFieldError); ok {
			fmt.Println("Invalid unmarshall field", e.Field, e.Key, e.Type)
		} else if e, ok := err.(*json.UnmarshalTypeError); ok {
			fmt.Println("Invalid JSON")
			fmt.Println("Expected", e.Type, "but received a", e.Value)
			fmt.Println("I apologize for not being more helpful")
		} else if e, ok := err.(*json.UnsupportedTypeError); ok {
			fmt.Println("Invalid JSON")
			fmt.Println("I am unable to utilize type", e.Type)
		} else if e, ok := err.(*json.UnsupportedValueError); ok {
			fmt.Println("Invalid JSON")
			fmt.Println("I am unable to utilize value", e.Value, e.Str)
		}
		return nil, err
	}

	//Parse the JSON in to an object to preserve non-standard fields
	object, err := parseJSONObject(raw)
	if err != nil {
		return nil, err
	}

	//Parse the odd security section of the config
	for _, value := range object["security"].([]interface{}) {
		v := reflect.ValueOf(value)
		if value == "nofiles" {
			structured.Security.NoFiles = 1
		} else if v.Kind() == reflect.Map {
			user := value.(map[string]interface{})
			structured.Security.SetUser = user["setuser"].(string)
		}
	}
	return &structured, nil
}

// Loads and parses the input file and returns a map with all data found in the
// config file, including non-standard fields.
func LoadExtConfig(filein string) (map[string]interface{}, error) {

	//Load the raw JSON data from the file
	raw, err := loadJson(filein)
	if err != nil {
		return nil, err
	}

	//Parse the JSON in to an object to preserve non-standard fields
	object, err := parseJSONObject(raw)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// Saves either of the two config types to the specified file with the
// specified permissions.
func SaveConfig(fileout string, config interface{}, perms os.FileMode) error {

	// check to see if we got a struct or a map (minimal or extended config,
	// respectively)
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Struct {
		config := config.(Config)

		// Parse the nicely formatted security section, and set the raw values
		// for JSON marshalling
		newSecurity := make([]interface{}, 0)
		if config.Security.NoFiles != 0 {
			newSecurity = append(newSecurity, "nofiles")
		}
		setuser := make(map[string]interface{})
		setuser["setuser"] = config.Security.SetUser
		newSecurity = append(newSecurity, setuser)
		config.RawSecurity = newSecurity

		jsonout, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			return err
		}
		return ioutil.WriteFile(fileout, jsonout, perms)
	} else if v.Kind() == reflect.Map {
		jsonout, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			return err
		}
		return ioutil.WriteFile(fileout, jsonout, perms)
	}
	return fmt.Errorf("Something very bad happened")
}

// Returns a []byte of raw JSON with comments removed.
func loadJson(filein string) ([]byte, error) {
	file, err := ioutil.ReadFile(filein)
	if err != nil {
		return nil, err
	}

	raw, err := stripComments(file)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// Returns a Config structure with the JSON unmarshalled in to it.
func parseJSONStruct(jsonIn []byte) (Config, error) {
	var structured Config
	err := json.Unmarshal(jsonIn, &structured)
	if err != nil {
		return Config{}, err
	}
	return structured, nil
}

// Returns a map with the JSON unmarshalled in to it.
func parseJSONObject(jsonIn []byte) (map[string]interface{}, error) {
	var object map[string]interface{}
	err := json.Unmarshal(jsonIn, &object)
	if err != nil {
		return nil, err
	}

	cleanObj := fixJSON(object)
	object = cleanObj.(map[string]interface{})

	return object, nil
}

// fixJSON works around the issue of empty []'s in the JSON
// https://codereview.appspot.com/7196050/
func fixJSON(in interface{}) interface{} {
	switch s := in.(type) {
	case map[string]interface{}:
		x := in.(map[string]interface{})
		for key, value := range s {
			x[key] = fixJSON(value)
		}
		in = x
	case []interface{}:
		x := in.([]interface{})
		// If the interface is nill, we need to initialize it
		// Otherwise it will be marshalled to 'null' in the JSON
		if x == nil {
			x = make([]interface{}, 0)
		} else {
			for i, value := range s {
				x[i] = fixJSON(value)
			}
		}
		in = x
	}
	return in
}

// Replaces all C-style comments (prefixed with "//" and inside "/* */") with
// empty strings. This is necessary in parsing JSON files that contain them.
// Returns b without comments. Credit to SashaCrofter, thanks!
func stripComments(b []byte) ([]byte, error) {
	regComment, err := regexp.Compile("(?s)//.*?\n|/\\*.*?\\*/")
	if err != nil {
		return nil, err
	}
	out := regComment.ReplaceAllLiteral(b, nil)
	return out, nil
}
