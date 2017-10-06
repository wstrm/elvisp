package admin

import "errors"

// Core_exit tells cjdns to shutdown
func (c *Conn) Core_exit() error {
	resp := new(struct{ Error string })

	pack, err := c.sendCmd(&request{AQ: "Core_exit"})
	if err == nil {
		err = pack.Decode(resp)
		if err == nil && resp.Error != "none" {
			err = errors.New(resp.Error)
		}
	}
	return err
}
