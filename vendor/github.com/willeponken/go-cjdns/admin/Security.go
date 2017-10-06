package admin

// Security_setUser sets the user ID which cjdns is running under to a different user.
// This function allows cjdns to shed privileges after starting up.
func (c *Conn) Security_setUser(user string) error {
	_, err := c.sendCmd(&request{
		AQ: "Security_setUser",
		Args: &struct {
			User string `bencode:"user"`
		}{user}})
	return err
}
