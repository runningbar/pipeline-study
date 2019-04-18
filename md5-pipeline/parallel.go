package main

import (
	"os"
	"fmt"
	"sort"
	"sync"
	"errors"
	"crypto/md5"
	"io/ioutil"
	"path/filepath"
)

type result struct {
	path string
	sum [md5.Size]byte
	err error
}

func sumFiles1(done <-chan struct{}, root string) (<-chan result, <-chan error) {
	c := make(chan result)
	errc := make(chan error, 1)

	// sumFiles1要立即返回outbound channel给下游，而实际的数据可以按需发送，所以这里单启协程操作数据
	go func() {
		var wg sync.WaitGroup
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			wg.Add(1)
			// 对每个文件单启协程去计算MD5码，可以并行计算
			go func() {
				data, err := ioutil.ReadFile(path)
				select {
				case c <-result{path, md5.Sum(data), err}:
				case <-done:
				}
				wg.Done()
			}()

			// 监听结束信号，有信号立马结束整个目录的walk，因为Walk是递归实现的，遇到错误，会快速跳出递归
			// 没信号就这个文件的遍历结束，继续下一个
			select {
			case <-done:
				return errors.New("walk canceled")
			default:
				return nil
			}
		})

		// 因为errc要实时监听walk的错误，所以这里要单启协程去wait wg
		// wg结束了，才能close outbound channel
		go func() {
			wg.Wait()
			close(c)
		}()
		
		errc <- err
	}()
	return c, errc
}

func main() {
	m, err := MD5All1(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x %s\n", m[path], path)
	}
}

func MD5All1(root string) (map[string][md5.Size]byte, error) {

	done := make(chan struct{})
	defer close(done)

	c, errc := sumFiles1(done, root)
	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
} 