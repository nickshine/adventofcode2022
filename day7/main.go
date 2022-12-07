package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

const (
	fsSize    = 70000000
	unusedMin = 30000000
)

type file struct {
	size int
	name string
}

type dir struct {
	name   string
	parent *dir
	dirs   map[string]*dir
	files  map[string]*file
}

func newDir(name string, parent *dir) *dir {
	return &dir{
		name:   name,
		parent: parent,
		dirs:   make(map[string]*dir),
		files:  make(map[string]*file),
	}
}

func (d *dir) display(indent int) {

	indention := strings.Repeat(" ", indent)
	fmt.Printf("%s- %s (dir)\n", indention, d.name)
	for _, d := range d.dirs {
		d.display(indent + 2)
	}

	for _, f := range d.files {
		fmt.Printf("%s  - %s (file, size=%d)\n", indention, f.name, f.size)
	}
}

func (d *dir) size() int {

	size := 0
	for _, subDir := range d.dirs {
		size += subDir.size()
	}

	for _, f := range d.files {
		size += f.size
	}

	return size
}

func buildFS(in string) *dir {
	lines := strings.Split(in, "\n")
	// di)
	var curr, root *dir
	root = newDir("/", nil)
	curr = root

	for i := 1; i < len(lines); i++ {
		parts := strings.Fields(lines[i])

		if parts[0] == "$" {
			if parts[1] == "cd" {
				name := parts[2]
				if name == ".." {
					fmt.Printf("Traversing up a directory from %s, to %s\n", curr.name, curr.parent.name)
					curr = curr.parent
				} else if _, ok := curr.dirs[name]; ok {
					fmt.Printf("Traversing down a directory from %s, to %s\n", curr.name, name)
					curr = curr.dirs[name]
				} else {
					fmt.Printf("New dir that hasn't been seen: %s\n", name)
					child := newDir(name, curr)
					curr = child
				}

			} else if parts[1] == "ls" {
				fmt.Printf("Listing dirs/files in directory %s\n", curr.name)
			}
		} else if parts[0] == "dir" {
			name := parts[1]
			fmt.Printf("ls dir %s in directory %s\n", name, curr.name)
			if _, ok := curr.dirs[name]; !ok {
				curr.dirs[name] = newDir(name, curr)
			}
		} else {
			// is file
			size, _ := strconv.Atoi(parts[0])
			name := parts[1]
			fmt.Printf("ls file %s in directory %s\n", name, curr.name)
			if _, ok := curr.files[name]; !ok {
				curr.files[name] = &file{size: size, name: name}
			}
		}
	}

	return root
}

func sizes(d *dir) []int {

	dirSizes := []int{}

	for _, subDir := range d.dirs {
		dirSizes = append(dirSizes, sizes(subDir)...)
	}

	return append(dirSizes, d.size())
}

func part1() int {

	root := buildFS(strings.TrimSpace(input))
	root.display(0)
	dirSizes := sizes(root)
	fmt.Println(dirSizes)
	sum := 0
	for _, s := range dirSizes {
		if s <= 100000 {
			sum += s
		}
	}

	return sum
}

func part2() int {

	root := buildFS(strings.TrimSpace(input))
	unused := fsSize - root.size()
	deleteMin := unusedMin - unused

	dirSizes := sizes(root)
	sort.Ints(dirSizes)
	fmt.Println(dirSizes)
	for _, s := range dirSizes {
		if s >= deleteMin {
			return s
		}
	}

	return -1
}

func main() {
	fmt.Printf("part 1: %d\n", part1())
	fmt.Printf("part 2: %d\n", part2())
}
