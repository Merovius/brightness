// Copyright 2018 Axel Wagner
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(0)

	check := func(err error) {
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
	}

	// Can't use flag, because we want to accept -x% as an argument
	display, expr := parseArgs(check, os.Args)

	check = func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	matches, err := filepath.Glob("/sys/devices/pci*/*/drm/card*/*/*/brightness")
	check(err)

	if len(matches) == 0 {
		log.Fatal("no display found")
	}
	curf := matches[0]
	if len(matches) > 1 || display != "" {
		if display == "" {
			log.Fatal("more than one display found, use -display to specify")
		}
		disp := listDisplays(matches)
		if curf = disp[display]; curf == "" {
			var ds []string
			for k := range disp {
				ds = append(ds, k)
			}
			sort.Strings(ds)
			log.Fatalf("display %q not found. Available displays: %q", display, ds)
		}
	}
	maxf := filepath.Join(filepath.Dir(curf), "max_brightness")

	cur := readInt(check, curf)
	if expr == nil {
		fmt.Println(cur)
		return
	}
	max := readInt(check, maxf)
	writeInt(check, curf, expr(cur, max))
}

type fn func(cur, max int) int

func parseArgs(check func(error), args []string) (display string, expr fn) {
	usage := fmt.Errorf("usage: %s [-display=<disp>] <expr>", args[0])
	args = args[1:]
	if len(args) == 0 {
		return "", nil
	}
	if strings.HasPrefix(args[0], "-display") {
		a := args[0][len("-display"):]
		if a == "" {
			if len(args) == 1 {
				check(usage)
			}
			display = args[1]
			args = args[2:]
		} else if a[0] != '=' {
			check(usage)
		} else {
			display = a[1:]
			args = args[1:]
		}
	}
	if len(args) == 0 {
		return display, nil
	}
	if len(args) > 1 {
		check(usage)
	}
	return display, parseExpr(check, args[0])
}

func listDisplays(matches []string) map[string]string {
	out := make(map[string]string)
	for _, m := range matches {
		tmp := filepath.Dir(filepath.Dir(m))
		card := filepath.Base(filepath.Dir(tmp))
		disp := strings.TrimPrefix(filepath.Base(tmp), card)
		disp = strings.TrimPrefix(disp, "-")
		out[disp] = m
	}
	return out
}

func parseExpr(check func(error), e string) fn {
	var (
		v        int
		inc, dec bool
		percent  bool
	)
	switch {
	case strings.HasPrefix(e, "+"):
		inc = true
		e = e[1:]
	case strings.HasPrefix(e, "-"):
		dec = true
		e = e[1:]
	}
	if strings.HasSuffix(e, "%") {
		percent = true
		e = e[:len(e)-1]
	}
	v, err := strconv.Atoi(e)
	check(err)

	var f func(cur, max int) int
	switch {
	case inc && !percent:
		f = func(cur, max int) int { return cur + v }
	case dec && !percent:
		f = func(cur, max int) int { return cur - v }
	case inc && percent:
		f = func(cur, max int) int { return cur + int(float64(max)*float64(v)/100) }
	case dec && percent:
		f = func(cur, max int) int { return cur - int(float64(max)*float64(v)/100) }
	case percent:
		f = func(cur, max int) int { return int(float64(max) * float64(v) / 100) }
	default:
		f = func(cur, max int) int { return v }
	}
	return func(cur, max int) int {
		n := f(cur, max)
		if n < 0 {
			n = 0
		}
		if n > max {
			n = max
		}
		return n
	}
}

func readInt(check func(error), path string) int {
	buf, err := ioutil.ReadFile(path)
	check(err)
	buf = bytes.TrimSpace(buf)
	n, err := strconv.Atoi(string(buf))
	check(err)
	return n
}

func writeInt(check func(error), path string, n int) {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	check(err)
	_, err = f.WriteString(strconv.Itoa(n))
	check(err)
	check(f.Close())
}
