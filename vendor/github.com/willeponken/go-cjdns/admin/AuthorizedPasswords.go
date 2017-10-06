package admin

import "fmt"

func IsPasswordAlreadyAdded(err error) bool { return err.Error() == "Password already added." }

// AuthorizedPasswords_add adds a password with will allow neighbors to make
// direct connections.
// Set authType to zero to invoke default.
func (c *Conn) AuthorizedPasswords_add(user, password string, authType int) error {
	var args = &struct {
		AuthType int    `bencode:"authType,omitempty"`
		Password string `bencode:"password"`
		User     string `bencode:"user"`
	}{authType, password, user}

	_, err := c.sendCmd(&request{AQ: "AuthorizedPasswords_add", Args: args})
	return err
}

// AuthorizedPasswords_list returns a list of users with passwords.
func (c *Conn) AuthorizedPasswords_list() (users []string, err error) {
	resp := new(struct {
		Total int
		Users []string
	})

	var pack *packet
	if pack, err = c.sendCmd(&request{AQ: "AuthorizedPasswords_list"}); err == nil {
		err = pack.Decode(resp)
	}
	if err == nil && len(resp.Users) != resp.Total {
		err = fmt.Errorf("users total reported as %d, but only unmarshaled %d", resp.Total, len(resp.Users))
	}
	return resp.Users, err
}

// AuthorizedPasswords_list removes a password for a given user.
func (c *Conn) AuthorizedPasswords_remove(user string) error {
	_, err := c.sendCmd(&request{
		AQ: "AuthorizedPasswords_remove",
		Args: &struct {
			User string `bencode:"user"`
		}{user},
	})
	return err
}
