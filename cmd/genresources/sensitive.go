package main

// sensitive lists plain (not write-only) string arguments whose values can
// hold a credential, keyed by path as in collections. They are marked
// Sensitive, so plans and output never print them. They stay ordinary
// arguments: the server stores and returns them, so they cannot be
// write-only.
//
// TestSensitiveEntriesExist fails on an entry that names no string argument.
var sensitive = map[string]string{
	"LibraryItem.spec.body": "an mcp item's body carries server headers, such as an Authorization " +
		"key, and a command or skill's text can quote a connection string; a body is also often set " +
		"from a sensitive variable, which without this plans as a sensitivity-only update.",
	"LibraryItem.spec.files[].body": "the files beside an item are free text, like its body.",
}
