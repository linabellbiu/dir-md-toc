package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// 读取忽略列表文件
func readIgnoreFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("warning: " + err.Error())
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
func listFiles(path string, level int, result *string, ignoreList []string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}

	for _, fileInfo := range files {
		if isIgnored(fileInfo.Name(), ignoreList) {
			continue
		}
		if fileInfo.IsDir() && !strings.HasPrefix(fileInfo.Name(), ".") {
			*result += fmt.Sprintf("%s- [%s](%s)\n", strings.Repeat("  ", level), fileInfo.Name(), strings.TrimRight(path, "/")+"/"+fileInfo.Name())
			listFiles(strings.TrimRight(path, "/")+"/"+fileInfo.Name(), level+1, result, ignoreList)
		} else if !fileInfo.IsDir() && !strings.HasPrefix(fileInfo.Name(), ".") {
			*result += fmt.Sprintf("%s- [%s](%s)\n", strings.Repeat("  ", level), fileInfo.Name(), strings.TrimRight(path, "/")+"/"+fileInfo.Name())
		}
	}
}

func main() {
	dirPath := flag.String("dir", ".", "Directory to generate")
	//outPath := flag.String("out", ".", "TOC to output path")
	flag.Parse()

	ignoreFilePath := fmt.Sprintf(".tocignore")

	ignoreList, _ := readIgnoreFile(ignoreFilePath)

	// Read README.md file
	readmeContent, err := os.ReadFile("README.md")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Find the position of <!--ts--> and <!--ed-->
	startIndex := strings.Index(string(readmeContent), "<!--ts-->")
	endIndex := strings.Index(string(readmeContent), "<!--ed-->")

	if startIndex != -1 && endIndex != -1 {
		result := string(readmeContent[:startIndex+len("<!--ts-->")]) + "\n"

		listFiles(*dirPath, 1, &result, ignoreList)

		result += string(readmeContent[endIndex:])

		err = os.WriteFile("README.md", []byte(result), 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Directory structure successfully inserted into README.md")
	} else {
		fmt.Println("Could not find <!--ts--> and <!--ed--> markers in README.md")
	}
}
