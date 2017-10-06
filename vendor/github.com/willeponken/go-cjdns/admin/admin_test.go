// +build admin

package admin

import (
	"bytes"
	"math"
	"testing"
)

var c *Conn

func TestPathMarshalingUnmarshaling(t *testing.T) {
	path := new(Path)
	if err := path.UnmarshalText([]byte("0000.0114.a785.58e3")); err != nil {
		t.Error("Failed to unmarshal Path,", err)
		return
	}
	if *path == 0 {
		t.Error("unmarshaled path was empty")
		return
	}

	test, err := path.MarshalText()
	if err != nil {
		t.Error("Failed to marshal Path,", err)
		return
	}
	if !bytes.Equal([]byte("0000.0114.a785.58e3"), test) {
		t.Errorf("Path marshal and unmarshal mismatch, wanted \"0000.0114.a785.58e3\", got %q", test)
	}
}

func TestConnect(t *testing.T) {
	var err error
	c, err = Connect(nil)
	if err != nil {
		t.Fatal("Failed to connect,", err)
	}
}

func TestPing(t *testing.T) {
	if err := c.Ping(); err != nil {
		t.Error("Failed to ping,", err)
	}
}

func BenchmarkPing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := c.Ping(); err != nil {
			b.Error("Failed to ping,", err)
		}
	}
}

func TestCookie(t *testing.T) {
	_, err := c.cookie()
	if err != nil {
		t.Error(err)
	}
}

func TestAuth(t *testing.T) {
	err := c.authedPing()
	if err != nil {
		t.Error(err)
	}
}

func TestAdmin_asyncEnabled(t *testing.T) {
	_, err := c.Admin_asyncEnabled()
	if err != nil {
		t.Error("Admin_asyncEnabled failed,", err)
	}
}

var table Routes

func TestNodeStor(t *testing.T) {
	var err error
	table, err = c.NodeStore_dumpTable()
	if err != nil {
		t.Error("NodeStore_dumpTable failed,", err)
	}

	for _, r := range table {
		n, err := c.NodeStore_nodeForAddr(r.IP.String())
		if err != nil {
			t.Fatal("NodeStore_nodeForAddr failed,", err)
		}

		for i := 0; i < n.LinkCount; i++ {
			_, err = c.NodeStore_getLink(r.IP.String(), i)
			if err != nil {
				t.Fatal("NodeStore_getLink failed,", err)
			}
			return
		}
	}
}

func TestLog2x64Algos(t *testing.T) {
	for _, r := range table {
		path := *r.Path
		testA := Path(math.Log2(float64(path)))
		var testB Path
		for path > 1 {
			path >>= 1
			testB++
		}
		if testA != testB {
			t.Error("not equal,", testA, testB)
		}
	}
}

func BenchmarkLog2x64Float(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := table[i%len(table)]
		_ = Path(math.Log2(float64(*r.Path)))
	}

}

func BenchmarkLog2x64Shift(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := table[i%len(table)]
		in := *r.Path
		var out Path
		for in > 1 {
			in >>= 1
			out++
		}
	}
}

func TestInterfaceController_peerStats(t *testing.T) {
	peerStats, err := c.InterfaceController_peerStats()
	if err != nil {
		t.Error("InterfaceController_peerStats failed,", err)
	}
	for _, p := range peerStats {
		if !p.PublicKey.Valid() {
			t.Error("Bad public key in peerStats,", p.PublicKey)
		}
		if p.SwitchLabel == nil || *p.SwitchLabel == 0 {
			t.Error("Bad SwitchLabel in peerStats,", p.SwitchLabel)
		}
	}
}

func TestAuthorizedPasswords(t *testing.T) {
	user := "someguy"
	pass := "hackme"

	if err := c.AuthorizedPasswords_add(user, pass, 0); err != nil {
		t.Error("failed to add password to cjdns,", err)
		return
	}
	if err := c.AuthorizedPasswords_add(user, pass, 0); err != nil && !IsPasswordAlreadyAdded(err) {
		t.Error("failed to add password to cjdns again,", err)
		return
	}

	users, err := c.AuthorizedPasswords_list()
	if err != nil {
		t.Error("failed to get list of password users from cjdns,", err)
		return
	}

	var found bool
	for _, u := range users {
		if u == user {
			found = true
			break
		}
	}
	if !found {
		t.Error("previously added user not found in users list")
		return
	}

	err = c.AuthorizedPasswords_remove(user)
	if err != nil {
		t.Error("failed to remove password for user,", err)
		return
	}

	users, err = c.AuthorizedPasswords_list()
	if err != nil {
		t.Error("failed to get list of password users from cjdns", err)
		return
	}

	found = false
	for _, u := range users {
		if u == user {
			found = true
			break
		}
	}
	if found {
		t.Error("previously removed user still found in users list")
	}
}
