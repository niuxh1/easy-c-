package analyzer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CppClass 表示一个C++类
type CppClass struct {
	Name        string   // 类名
	BaseClasses []string // 基类列表
	Members     []string // 成员变量
	Methods     []string // 成员方法
	LineNumber  int      // 类定义开始的行号
	FilePath    string   // 类定义所在的文件路径
}

// CppAnalyzer C++代码分析器
type CppAnalyzer struct {
	classRegex     *regexp.Regexp
	inheritRegex   *regexp.Regexp
	memberRegex    *regexp.Regexp
	methodRegex    *regexp.Regexp
	commentRegex   *regexp.Regexp
	blockCommentRe *regexp.Regexp
}

// NewCppAnalyzer 创建新的分析器实例
func NewCppAnalyzer() *CppAnalyzer {
	return &CppAnalyzer{
		// 匹配类定义 (支持继承)
		classRegex: regexp.MustCompile(`class\s+(\w+)(?:\s*:\s*(.+?))?\s*\{`),
		// 匹配继承关系
		inheritRegex: regexp.MustCompile(`(?:public|private|protected)?\s*(\w+)`),
		// 匹配成员变量
		memberRegex: regexp.MustCompile(`^\s*(?:private|public|protected)?\s*(\w+(?:\s*\*)?)\s+(\w+)(?:\[.*?\])?(?:\s*=.*?)?;`),
		// 匹配成员方法
		methodRegex: regexp.MustCompile(`^\s*(?:virtual\s+)?(?:static\s+)?(?:inline\s+)?(\w+(?:\s*\*)?)\s+(\w+)\s*\([^)]*\)(?:\s*const)?(?:\s*=\s*0)?(?:\s*override)?`),
		// 匹配单行注释
		commentRegex: regexp.MustCompile(`//.*$`),
		// 匹配块注释
		blockCommentRe: regexp.MustCompile(`/\*.*?\*/`),
	}
}

// AnalyzeFile 分析指定的C++文件
func (a *CppAnalyzer) AnalyzeFile(filePath string) ([]*CppClass, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件 %s: %v", filePath, err)
	}
	defer file.Close()

	var classes []*CppClass
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	inBlockComment := false

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// 处理块注释
		line = a.removeComments(line, &inBlockComment)

		// 跳过空行和纯注释行
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 查找类定义
		if matches := a.classRegex.FindStringSubmatch(line); matches != nil {
			class := &CppClass{
				Name:       matches[1],
				LineNumber: lineNumber,
			}

			// 解析继承关系
			if len(matches) > 2 && matches[2] != "" {
				class.BaseClasses = a.parseInheritance(matches[2])
			}

			// 读取类体内容
			a.parseClassBody(scanner, class, &lineNumber, &inBlockComment)
			classes = append(classes, class)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件时出错: %v", err)
	}

	return classes, nil
}

// removeComments 移除代码中的注释
func (a *CppAnalyzer) removeComments(line string, inBlockComment *bool) string {
	result := line

	// 处理块注释
	if *inBlockComment {
		if idx := strings.Index(result, "*/"); idx != -1 {
			*inBlockComment = false
			result = result[idx+2:]
		} else {
			return ""
		}
	}

	// 查找块注释开始
	for {
		startIdx := strings.Index(result, "/*")
		if startIdx == -1 {
			break
		}

		endIdx := strings.Index(result[startIdx:], "*/")
		if endIdx == -1 {
			*inBlockComment = true
			result = result[:startIdx]
			break
		} else {
			result = result[:startIdx] + result[startIdx+endIdx+2:]
		}
	}

	// 移除单行注释
	result = a.commentRegex.ReplaceAllString(result, "")

	return result
}

// parseInheritance 解析继承关系字符串
func (a *CppAnalyzer) parseInheritance(inheritanceStr string) []string {
	var baseClasses []string

	// 按逗号分割多重继承
	parts := strings.Split(inheritanceStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		// 使用正则表达式提取基类名
		if matches := a.inheritRegex.FindStringSubmatch(part); matches != nil {
			baseClasses = append(baseClasses, matches[1])
		}
	}

	return baseClasses
}

// parseClassBody 解析类体内容
func (a *CppAnalyzer) parseClassBody(scanner *bufio.Scanner, class *CppClass, lineNumber *int, inBlockComment *bool) {
	braceCount := 1 // 已经遇到了开始的 {

	for scanner.Scan() && braceCount > 0 {
		*lineNumber++
		line := scanner.Text()

		// 处理注释
		line = a.removeComments(line, inBlockComment)
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 计算大括号
		braceCount += strings.Count(line, "{")
		braceCount -= strings.Count(line, "}")

		if braceCount <= 0 {
			break
		}

		// 跳过访问修饰符行
		trimmed := strings.TrimSpace(line)
		if trimmed == "public:" || trimmed == "private:" || trimmed == "protected:" {
			continue
		}

		// 尝试匹配成员变量
		if matches := a.memberRegex.FindStringSubmatch(line); matches != nil {
			member := fmt.Sprintf("%s %s", matches[1], matches[2])
			class.Members = append(class.Members, member)
			continue
		}

		// 尝试匹配成员方法
		if matches := a.methodRegex.FindStringSubmatch(line); matches != nil {
			method := fmt.Sprintf("%s %s(...)", matches[1], matches[2])
			class.Methods = append(class.Methods, method)
		}
	}
}

// GetInheritanceTree 构建继承树
func GetInheritanceTree(classes []*CppClass) map[string][]*CppClass {
	tree := make(map[string][]*CppClass)

	// 为每个基类建立子类列表
	for _, class := range classes {
		for _, baseClass := range class.BaseClasses {
			tree[baseClass] = append(tree[baseClass], class)
		}
	}

	return tree
}

// FindRootClasses 找到所有根类(没有基类的类)
func FindRootClasses(classes []*CppClass) []*CppClass {
	var roots []*CppClass

	for _, class := range classes {
		if len(class.BaseClasses) == 0 {
			roots = append(roots, class)
		}
	}

	return roots
}

// AnalyzeProject 分析整个C++项目目录
func (a *CppAnalyzer) AnalyzeProject(projectPath string) ([]*CppClass, error) {
	var allClasses []*CppClass

	// 遍历项目目录查找C++文件
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否为C++文件
		if !info.IsDir() && a.isCppFile(path) {
			classes, err := a.AnalyzeFile(path)
			if err != nil {
				fmt.Printf("警告: 分析文件 %s 时出错: %v\n", path, err)
				return nil // 继续处理其他文件
			}

			// 为每个类添加文件路径信息
			for _, class := range classes {
				class.FilePath = path
			}

			allClasses = append(allClasses, classes...)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历项目目录时出错: %v", err)
	}

	// 解析跨文件的继承关系
	a.resolveInterFileInheritance(allClasses)

	return allClasses, nil
}

// isCppFile 检查文件是否为C++源文件
func (a *CppAnalyzer) isCppFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".cpp" || ext == ".cxx" || ext == ".cc" || ext == ".c++" ||
		ext == ".h" || ext == ".hpp" || ext == ".hxx" || ext == ".h++"
}

// resolveInterFileInheritance 解析跨文件的继承关系
func (a *CppAnalyzer) resolveInterFileInheritance(classes []*CppClass) {
	// 创建类名到类对象的映射
	classMap := make(map[string]*CppClass)
	for _, class := range classes {
		classMap[class.Name] = class
	}

	// 验证和修正继承关系
	for _, class := range classes {
		var validBaseClasses []string
		for _, baseName := range class.BaseClasses {
			if _, exists := classMap[baseName]; exists {
				validBaseClasses = append(validBaseClasses, baseName)
			} else {
				fmt.Printf("警告: 类 %s 的基类 %s 未找到定义\n", class.Name, baseName)
			}
		}
		class.BaseClasses = validBaseClasses
	}
}

// AnalyzeFiles 分析多个指定的C++文件
func (a *CppAnalyzer) AnalyzeFiles(filePaths []string) ([]*CppClass, error) {
	var allClasses []*CppClass

	for _, filePath := range filePaths {
		classes, err := a.AnalyzeFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("分析文件 %s 时出错: %v", filePath, err)
		}

		// 为每个类添加文件路径信息
		for _, class := range classes {
			class.FilePath = filePath
		}

		allClasses = append(allClasses, classes...)
	}

	// 解析跨文件的继承关系
	a.resolveInterFileInheritance(allClasses)

	return allClasses, nil
}
