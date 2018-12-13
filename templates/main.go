package scaffold

var MainTemplate = `
package main

import "{{ .importpath }}"

func main() {
	cmd.Execute()
}
`
`