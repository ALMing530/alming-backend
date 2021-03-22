package engine

import "testing"

func Test_trie(t *testing.T) {
	trie := rootNode()
	trie.AddNode("/wxm/abc")
	trie.AddNode("/wxm/abc/:user")
	res1 := trie.search("/wxm/abc")
	if res1 != nil {
		if !(res1.pattern == "wxm/abc") {
			t.Error("error in get node")
		}
	}
	res2 := trie.search("/wxm/abc/Little Amy")
	if res2 != nil {
		if !(res2.pattern == "wxm/abc/:user") {
			t.Error("error in get node")
		}
	}
}
