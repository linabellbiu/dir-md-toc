package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// 读取忽略列表文件
func readIgnoreFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var ignoreList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") { // 忽略空行和注释行
			ignoreList = append(ignoreList, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ignoreList, nil
}

// 判断文件或文件夹是否在忽略列表中
func isIgnored(name string, ignoreList []string) bool {
	for _, ignore := range ignoreList {
		if name == ignore || strings.HasPrefix(name, ignore) {
			return true
		}
	}
	return false
}

// 列出文件并生成markdown
func listFiles(path string, level int, file *os.File, ignoreList []string) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, fileInfo := range files {
		if isIgnored(fileInfo.Name(), ignoreList) {
			continue
		}
		if fileInfo.IsDir() && !strings.HasPrefix(fileInfo.Name(), ".") {
			_, _ = fmt.Fprintln(file, strings.Repeat("  ", level)+"- ["+fileInfo.Name()+"]("+path+"/"+fileInfo.Name()+")")
			listFiles(path+"/"+fileInfo.Name(), level+1, file, ignoreList)
		} else if !fileInfo.IsDir() && !strings.HasPrefix(fileInfo.Name(), ".") {
			_, _ = fmt.Fprintln(file, strings.Repeat("  ", level)+"- ["+fileInfo.Name()+"]("+path+"/"+fileInfo.Name()+")")
		}
	}
}

func main() {
	dirPath := flag.String("dir", ".", "Directory to generate")
	outPath := flag.String("out", ".", "TOC to output path")
	flag.Parse()

	ignoreFilePath := fmt.Sprintf("%s/.ignore", *dirPath)

	// 读取忽略列表
	ignoreList, err := readIgnoreFile(ignoreFilePath)
	if err != nil {
		fmt.Println("Error reading ignore file:", err)
		return
	}

	f, err := os.Create(fmt.Sprintf("%s/README.md", *outPath))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_, _ = fmt.Fprintf(f, "# 目录\n\n")
	listFiles(*dirPath, 1, f, ignoreList)
}
