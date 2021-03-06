package main

import (
	"fmt"
	"time"
    "os"
    "io/ioutil"
    "log"
    "strings"
)

type NodeInfo struct {
	fb2pointer string
	weight     float32
}

type Node struct {
	text  string
	info  []*NodeInfo
	left  *Node
	right *Node
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Please specify search word")
        return
    }
    search := os.Args[1]
    root := dirTraverse("fb2docs")

	stime := time.Now().UnixNano()
	find(root, search)
	etime := time.Now().UnixNano()
	fmt.Println("time elapsed", (etime-stime)/1e3, "ms")
	fmt.Println("nodeCount", nodeCount)
}

func dirTraverse(path string) *Node {
	c := make(chan *Word)
	qchan := make(chan bool)
	var root *Node = new(Node)
    files, err := ioutil.ReadDir(path)

    if err != nil {
        log.Fatal(err)
    }

    fileCount := 0
    for _, fi := range files {
        if strings.HasSuffix(fi.Name(), ".fb2") {
            fileName := path + string(os.PathSeparator) + fi.Name()
            fmt.Println("Indexing", fileName)
	        go processBook(fileName, c, qchan)
            fileCount += 1
        }
    }

    if fileCount > 0 {
L:
        for {
            select {
            case word := <-c:
                searchNode(root, word)
            case <-qchan:
                fileCount -= 1
                if 0 == fileCount {
                    break L
                }
            }
        }
    }

    return root
}

var nodeCount int

func searchNode(node *Node, word *Word) {
	if node.text == "" {
		node.text = word.text
		node.info = make([]*NodeInfo, 1)
		node.info[0] = &NodeInfo{word.fb2pointer, word.weight}
		node.right = new(Node)
		node.left = new(Node)
		nodeCount += 1
	} else if word.text > node.text {
		searchNode(node.right, word)
	} else if word.text < node.text {
		searchNode(node.left, word)
	} else {
		node.info = append(node.info, &NodeInfo{word.fb2pointer, word.weight})
	}
}

func find(node *Node, word string) {
	if node == nil {
		fmt.Println("no such phrase")
	} else if word > node.text {
		find(node.right, word)
	} else if word < node.text {
		find(node.left, word)
	} else if word == node.text {
		fmt.Println("found", word, len(node.info), "times")
		for _, info := range node.info {
			fmt.Println(info.fb2pointer, info.weight)
		}
	}
}
