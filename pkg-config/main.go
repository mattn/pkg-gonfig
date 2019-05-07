package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func flags(root, pkg, name string) string {
	f, e := os.Open(filepath.Join(root, pkg+".pc"))
	if e != nil {
		println(e.Error())
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	re := regexp.MustCompile(`^(\w+)=(.*)$`)
	vv := make(map[string]string)

	for scanner.Scan() {
		l := scanner.Text()
		m := re.FindStringSubmatch(l)
		if m != nil && len(m) == 3 {
			key, val := m[1], m[2]
			for k, v := range vv {
				val = strings.Replace(val, "${"+k+"}", v, -1)
			}
			vv[key] = val
			continue
		}
		if strings.HasPrefix(l, name+":") {
			val := strings.TrimSpace(l[len(name)+1:])
			for k, v := range vv {
				val = strings.Replace(val, "${"+k+"}", v, -1)
			}
			return val
		}
	}
	if scanner.Err() != nil {
		println(scanner.Err().Error())
		return ""
	}
	return ""
}

func main() {
	root := os.Getenv("PKG_CONFIG_PATH")
	if len(root) == 0 {
		println("PKG_CONFIG_PATH is not set.")
		os.Exit(1)
	}

	var pkgs []string
	var fcflags, flibs bool
	for _, arg := range os.Args[1:] {
		if len(arg) < 2 || arg[:2] != "--" {
			pkgs = append(pkgs, arg)
		} else {
			if arg == "--cflags" {
				fcflags = true
			} else if arg == "--libs" {
				flibs = true
			} else if arg != "--" {
				println("invalid argument:", arg)
				os.Exit(2)
			}
		}
	}

	r := ""
	if fcflags {
		for _, pkg := range pkgs {
			r += " " + flags(root, pkg, "Cflags")
		}
	}
	if flibs {
		for _, pkg := range pkgs {
			r += " " + flags(root, pkg, "Libs")
		}
	}
	fmt.Println(r)
}
