package a

type Header map[string][]string

type AnotherHeader Header

type AliasHeader = Header // want "gostyle.typealiases"

func f() {
	type header map[string][]string

	type anotherHeader Header

	type aliasHeader = Header // want "gostyle.typealiases"
}
