package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var b bytes.Buffer
	var program string

	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		t := sc.Text()
		if strings.HasPrefix(t, "Usage of ") {
			program = strings.TrimPrefix(t, "Usage of ")
			program = strings.Trim(program, ":")
			fmt.Fprintln(&b, "#compdef "+program+"\n\n_arguments \\")
			continue
		}

		if strings.HasPrefix(t, "  -") {
			argname := strings.Trim(t, " ")
			argname = strings.Split(argname, " ")[0]
			argname = strings.Trim(argname, " ")
			fmt.Fprint(&b, "  '"+argname)
			continue
		}

		if strings.HasPrefix(t, "    \t") {
			descr := strings.Trim(t, " \t")
			if strings.Contains(descr, "FILE") {
				fmt.Fprintln(&b, "["+descr+"]:filename:_files' \\")
				continue
			}
			if strings.Contains(descr, "{") {
				s := strings.Split(descr, "{")[1]
				s = strings.Trim(s, "}")
				list := strings.Split(s, "|")
				args := "(" + strings.Join(list, " ") + ")"
				fmt.Fprintln(&b, "["+descr+"]:filename:"+args+"' \\")
				continue
			}

			fmt.Fprintln(&b, "["+descr+"]' \\")
			continue
		}
	}

	fname := os.Getenv("HOME") + "/.config/zsh_completion/_" + program
	if err := ioutil.WriteFile(fname, b.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}

	fmt.Println(b.String())
	fmt.Println("written to ", fname)
	fmt.Println("run: unfunction", "_"+program, "&& compinit")
}
