package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cpp-inheritance-analyzer/internal/analyzer"
	"cpp-inheritance-analyzer/internal/visualizer"
)

func main() {
	// 检查帮助选项
	if len(os.Args) >= 2 && (os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help") {
		showHelp()
		return
	}

	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	var classes []*analyzer.CppClass
	var err error
	outputFormat := "all" // 默认生成所有格式

	analyzer := analyzer.NewCppAnalyzer()

	// 解析命令行参数
	if os.Args[1] == "-project" {
		// 项目目录分析模式
		if len(os.Args) < 3 {
			fmt.Println("错误: 请指定项目目录路径")
			os.Exit(1)
		}
		projectPath := os.Args[2]
		if len(os.Args) >= 4 {
			outputFormat = os.Args[3]
		}

		fmt.Printf("正在分析C++项目目录: %s\n", projectPath)
		classes, err = analyzer.AnalyzeProject(projectPath)

	} else if os.Args[1] == "-files" {
		// 多文件分析模式
		if len(os.Args) < 3 {
			fmt.Println("错误: 请指定至少一个C++文件")
			os.Exit(1)
		}
		// 查找输出格式参数
		files := []string{}
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if arg == "text" || arg == "html" || arg == "interactive" || arg == "all" {
				outputFormat = arg
				break
			}
			files = append(files, arg)
		}

		fmt.Printf("正在分析C++文件: %v\n", files)
		classes, err = analyzer.AnalyzeFiles(files)

	} else {
		// 单文件分析模式
		filePath := os.Args[1]
		if len(os.Args) >= 3 {
			outputFormat = os.Args[2]
		}

		fmt.Printf("正在分析C++文件: %s\n", filePath)
		classes, err = analyzer.AnalyzeFile(filePath)

		// 为单文件分析添加文件路径
		for _, class := range classes {
			class.FilePath = filePath
		}
	}

	if err != nil {
		fmt.Printf("分析失败: %v\n", err)
		os.Exit(1)
	}

	if len(classes) == 0 {
		fmt.Println("未发现任何类定义")
		return
	}

	fmt.Printf("发现 %d 个类:\n", len(classes))
	for _, class := range classes {
		fmt.Printf("- %s", class.Name)
		if len(class.BaseClasses) > 0 {
			fmt.Printf(" (继承自: %v)", class.BaseClasses)
		}
		if class.FilePath != "" {
			fmt.Printf(" [文件: %s]", filepath.Base(class.FilePath))
		}
		fmt.Println()
	} // 生成可视化结果
	switch outputFormat {
	case "text":
		err = generateTextReport(classes, "inheritance_report.txt")
		if err != nil {
			fmt.Printf("生成文本报告失败: %v\n", err)
		} else {
			fmt.Println("继承关系报告已生成: inheritance_report.txt")
		}
	case "interactive", "interactive-html":
		err = generateInteractiveHTMLReport(classes, "inheritance_interactive.html")
		if err != nil {
			fmt.Printf("生成交互式HTML报告失败: %v\n", err)
		} else {
			fmt.Println("交互式继承关系报告已生成: inheritance_interactive.html")
		}
	case "html":
		err = generateHTMLReport(classes, "inheritance_report.html")
		if err != nil {
			fmt.Printf("生成HTML报告失败: %v\n", err)
		} else {
			fmt.Println("继承关系报告已生成: inheritance_report.html")
		}
	case "all":
		// 生成文本报告
		err = generateTextReport(classes, "inheritance_report.txt")
		if err != nil {
			fmt.Printf("生成文本报告失败: %v\n", err)
		} else {
			fmt.Println("继承关系报告已生成: inheritance_report.txt")
		}

		// 生成HTML报告
		err = generateHTMLReport(classes, "inheritance_report.html")
		if err != nil {
			fmt.Printf("生成HTML报告失败: %v\n", err)
		} else {
			fmt.Println("继承关系报告已生成: inheritance_report.html")
		}

		// 生成交互式HTML报告
		err = generateInteractiveHTMLReport(classes, "inheritance_interactive.html")
		if err != nil {
			fmt.Printf("生成交互式HTML报告失败: %v\n", err)
		} else {
			fmt.Println("交互式继承关系报告已生成: inheritance_interactive.html")
		}
	default:
		fmt.Printf("不支持的输出格式: %s\n", outputFormat)
		fmt.Println("支持的格式: text, html, interactive, all")
	}
}

// generateTextReport 生成文本格式的报告
func generateTextReport(classes []*analyzer.CppClass, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入标题
	fmt.Fprintf(file, "C++ 类继承关系分析报告\n")
	fmt.Fprintf(file, "生成时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "%s\n\n", strings.Repeat("=", 50))

	// 概述
	fmt.Fprintf(file, "概述\n")
	fmt.Fprintf(file, "%s\n", strings.Repeat("-", 20))
	fmt.Fprintf(file, "总类数: %d\n", len(classes))

	rootClasses := analyzer.FindRootClasses(classes)
	fmt.Fprintf(file, "根类数: %d\n", len(rootClasses))
	fmt.Fprintf(file, "派生类数: %d\n\n", len(classes)-len(rootClasses))

	// 类详情
	fmt.Fprintf(file, "类详情\n")
	fmt.Fprintf(file, "%s\n", strings.Repeat("-", 20))

	for i, class := range classes {
		fmt.Fprintf(file, "%d. 类名: %s\n", i+1, class.Name)
		fmt.Fprintf(file, "   行号: %d\n", class.LineNumber)

		if len(class.BaseClasses) > 0 {
			fmt.Fprintf(file, "   继承自: %s\n", strings.Join(class.BaseClasses, ", "))
		} else {
			fmt.Fprintf(file, "   根类 (无继承)\n")
		}

		if len(class.Members) > 0 {
			fmt.Fprintf(file, "   成员变量 (%d):\n", len(class.Members))
			for _, member := range class.Members {
				fmt.Fprintf(file, "     - %s\n", member)
			}
		}

		if len(class.Methods) > 0 {
			fmt.Fprintf(file, "   成员方法 (%d):\n", len(class.Methods))
			for _, method := range class.Methods {
				fmt.Fprintf(file, "     - %s\n", method)
			}
		}

		fmt.Fprintf(file, "\n")
	}

	// 继承层次结构
	fmt.Fprintf(file, "继承层次结构\n")
	fmt.Fprintf(file, "%s\n", strings.Repeat("-", 20))

	tree := analyzer.GetInheritanceTree(classes)
	for _, rootClass := range rootClasses {
		printClassHierarchyText(file, rootClass, tree, 0)
	}

	return nil
}

// printClassHierarchyText 递归打印类层次结构到文本文件
func printClassHierarchyText(file *os.File, class *analyzer.CppClass, tree map[string][]*analyzer.CppClass, level int) {
	indent := strings.Repeat("  ", level)
	symbol := "+"
	if level > 0 {
		symbol = "├─"
	}

	fmt.Fprintf(file, "%s%s %s\n", indent, symbol, class.Name)

	if children, exists := tree[class.Name]; exists {
		for _, child := range children {
			printClassHierarchyText(file, child, tree, level+1)
		}
	}
}

// generateHTMLReport 生成HTML格式的报告
func generateHTMLReport(classes []*analyzer.CppClass, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// HTML头部
	fmt.Fprintf(file, `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>C++ 类继承关系分析报告</title>
    <style>
        body {
            font-family: 'Microsoft YaHei', Arial, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background-color: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0,0,0,0.1);
        }
        h1, h2, h3 {
            color: #333;
        }
        h1 {
            text-align: center;
            border-bottom: 3px solid #4CAF50;
            padding-bottom: 10px;
        }
        .overview {
            background-color: #e8f5e8;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .class-card {
            border: 1px solid #ddd;
            margin: 15px 0;
            padding: 15px;
            border-radius: 5px;
            background-color: #fafafa;
        }
        .class-name {
            font-size: 1.2em;
            font-weight: bold;
            color: #2196F3;
        }
        .inheritance {
            color: #FF9800;
            font-weight: bold;
        }
        .members, .methods {
            margin-top: 10px;
        }
        .members ul, .methods ul {
            list-style-type: none;
            padding-left: 20px;
        }
        .members li::before {
            content: "📝 ";
        }
        .methods li::before {
            content: "⚙️ ";
        }
        .hierarchy {
            background-color: #f0f0f0;
            padding: 15px;
            border-radius: 5px;
            font-family: monospace;
            white-space: pre;
        }
        .root-class {
            color: #4CAF50;
            font-weight: bold;
        }
        .derived-class {
            color: #2196F3;
        }
        .stats {
            display: flex;
            justify-content: space-around;
            margin: 20px 0;
        }
        .stat-item {
            text-align: center;
            padding: 15px;
            background-color: #e3f2fd;
            border-radius: 5px;
            min-width: 120px;
        }
        .stat-number {
            font-size: 2em;
            font-weight: bold;
            color: #1976D2;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>C++ 类继承关系分析报告</h1>
        <p style="text-align: center; color: #666;">生成时间: %s</p>
`, time.Now().Format("2006-01-02 15:04:05"))

	// 统计信息
	rootClasses := analyzer.FindRootClasses(classes)
	fmt.Fprintf(file, `
        <div class="overview">
            <h2>📊 统计概览</h2>
            <div class="stats">
                <div class="stat-item">
                    <div class="stat-number">%d</div>
                    <div>总类数</div>
                </div>
                <div class="stat-item">
                    <div class="stat-number">%d</div>
                    <div>根类数</div>
                </div>
                <div class="stat-item">
                    <div class="stat-number">%d</div>
                    <div>派生类数</div>
                </div>
            </div>
        </div>
`, len(classes), len(rootClasses), len(classes)-len(rootClasses))

	// 类详情
	fmt.Fprintf(file, `
        <h2>📚 类详情</h2>
`)

	for i, class := range classes {
		fmt.Fprintf(file, `
        <div class="class-card">
            <div class="class-name">%d. %s</div>
            <p><strong>定义位置:</strong> 第 %d 行</p>
`, i+1, htmlEscape(class.Name), class.LineNumber)

		if len(class.BaseClasses) > 0 {
			fmt.Fprintf(file, `            <p class="inheritance">🔗 继承自: %s</p>`, htmlEscape(strings.Join(class.BaseClasses, ", ")))
		} else {
			fmt.Fprintf(file, `            <p><em>🌳 根类 (无继承关系)</em></p>`)
		}

		if len(class.Members) > 0 {
			fmt.Fprintf(file, `
            <div class="members">
                <strong>成员变量 (%d个):</strong>
                <ul>
`, len(class.Members))
			for _, member := range class.Members {
				fmt.Fprintf(file, `                    <li>%s</li>`, htmlEscape(member))
			}
			fmt.Fprintf(file, `                </ul>
            </div>`)
		}

		if len(class.Methods) > 0 {
			fmt.Fprintf(file, `
            <div class="methods">
                <strong>成员方法 (%d个):</strong>
                <ul>
`, len(class.Methods))
			for _, method := range class.Methods {
				fmt.Fprintf(file, `                    <li>%s</li>`, htmlEscape(method))
			}
			fmt.Fprintf(file, `                </ul>
            </div>`)
		}

		fmt.Fprintf(file, `        </div>`)
	}

	// 继承层次结构
	fmt.Fprintf(file, `
        <h2>🌲 继承层次结构</h2>
        <div class="hierarchy">`)

	tree := analyzer.GetInheritanceTree(classes)
	for _, rootClass := range rootClasses {
		hierarchyHTML := buildClassHierarchyHTML(rootClass, tree, 0)
		fmt.Fprintf(file, "%s", hierarchyHTML)
	}

	fmt.Fprintf(file, `        </div>
    </div>
</body>
</html>`)

	return nil
}

// htmlEscape 转义HTML特殊字符
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// buildClassHierarchyHTML 构建HTML格式的类层次结构
func buildClassHierarchyHTML(class *analyzer.CppClass, tree map[string][]*analyzer.CppClass, level int) string {
	var result strings.Builder

	indent := strings.Repeat("  ", level)
	symbol := "+"
	if level > 0 {
		symbol = "├─"
	}

	className := htmlEscape(class.Name)
	if level == 0 {
		className = fmt.Sprintf(`<span class="root-class">%s</span>`, className)
	} else {
		className = fmt.Sprintf(`<span class="derived-class">%s</span>`, className)
	}

	result.WriteString(fmt.Sprintf("%s%s %s\n", indent, symbol, className))

	if children, exists := tree[class.Name]; exists {
		for _, child := range children {
			result.WriteString(buildClassHierarchyHTML(child, tree, level+1))
		}
	}
	return result.String()
}

// generateInteractiveHTMLReport 生成交互式HTML格式的报告
func generateInteractiveHTMLReport(classes []*analyzer.CppClass, outputPath string) error {
	htmlGen := visualizer.NewHTMLGenerator()

	// 添加所有类到HTML生成器
	for _, class := range classes {
		filePath := class.FilePath
		if filePath == "" {
			filePath = "未知文件"
		} else {
			filePath = filepath.Base(filePath)
		}

		htmlGen.AddClass(class.Name, class.Members, class.Methods, class.BaseClasses, filePath)
	}

	// 生成HTML内容
	htmlContent := htmlGen.GenerateHTML()

	// 写入文件
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(htmlContent)
	return err
}

// showHelp 显示命令行帮助信息
func showHelp() {
	fmt.Println("C++ 类继承关系分析器")
	fmt.Println("===================")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  go run main.go [选项] <文件/目录> [输出格式]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -h, --help     显示此帮助信息")
	fmt.Println("  -project <dir> 分析指定项目目录")
	fmt.Println("  -files <f1> <f2> ... 分析多个指定文件")
	fmt.Println()
	fmt.Println("输出格式:")
	fmt.Println("  text         纯文本报告 (inheritance_report.txt)")
	fmt.Println("  html         静态HTML报告 (inheritance_report.html)")
	fmt.Println("  interactive  交互式HTML报告 (inheritance_interactive.html)")
	fmt.Println("  all          生成所有格式 (默认)")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run main.go example.cpp")
	fmt.Println("  go run main.go example.cpp interactive")
	fmt.Println("  go run main.go -project ./test_project")
	fmt.Println("  go run main.go -project ./test_project interactive")
	fmt.Println("  go run main.go -files file1.cpp file2.h")
	fmt.Println()
	fmt.Println("支持的C++特性:")
	fmt.Println("  ✓ 类定义和继承关系")
	fmt.Println("  ✓ 单继承和多重继承")
	fmt.Println("  ✓ 虚函数和纯虚函数")
	fmt.Println("  ✓ 成员变量和方法")
	fmt.Println("  ✓ 访问修饰符 (public, private, protected)")
	fmt.Println("  ✓ 多文件项目分析")
	fmt.Println()
	fmt.Println("项目主页: https://github.com/yourusername/cpp-inheritance-analyzer")
}
