package admin

func (c *Conn) SearchRunner_showActiveSearch(number int) error {
	_, err := c.sendCmd(&request{AQ: "SearchRunner",
		Args: &struct {
			Number int `bencode:"number"`
		}{number}})
	return err
}
