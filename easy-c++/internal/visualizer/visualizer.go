package visualizer

import (
	"fmt"
	"strings"

	"cpp-inheritance-analyzer/internal/analyzer"
)

// Visualizer 可视化生成器
type Visualizer struct{}

// NewVisualizer 创建新的可视化生成器
func NewVisualizer() *Visualizer {
	return &Visualizer{}
}

// GenerateTextTree 生成文本格式的继承树
func (v *Visualizer) GenerateTextTree(classes []*analyzer.CppClass) string {
	var sb strings.Builder

	// 构建继承树
	tree := make(map[string][]*analyzer.CppClass)
	roots := []*analyzer.CppClass{}

	// 找到所有根类和构建树结构
	for _, class := range classes {
		if len(class.BaseClasses) == 0 {
			roots = append(roots, class)
		}

		for _, baseClass := range class.BaseClasses {
			// 查找对应的基类对象
			for _, bc := range classes {
				if bc.Name == baseClass {
					tree[bc.Name] = append(tree[bc.Name], class)
					break
				}
			}
		}
	}

	sb.WriteString("C++ 类继承关系树:\n")
	sb.WriteString(strings.Repeat("=", 30) + "\n\n")

	// 递归打印每个根类的继承树
	for _, root := range roots {
		v.printClassTree(&sb, root, tree, 0)
	}

	return sb.String()
}

// printClassTree 递归打印类继承树
func (v *Visualizer) printClassTree(sb *strings.Builder, class *analyzer.CppClass, tree map[string][]*analyzer.CppClass, level int) {
	indent := strings.Repeat("  ", level)
	prefix := "├─"
	if level == 0 {
		prefix = "📦"
	}

	sb.WriteString(fmt.Sprintf("%s%s %s", indent, prefix, class.Name))

	// 添加成员和方法信息
	if len(class.Members) > 0 || len(class.Methods) > 0 {
		sb.WriteString(fmt.Sprintf(" [%d个成员, %d个方法]", len(class.Members), len(class.Methods)))
	}

	sb.WriteString("\n")

	// 递归打印子类
	if children, exists := tree[class.Name]; exists {
		for _, child := range children {
			v.printClassTree(sb, child, tree, level+1)
		}
	}
}

// GenerateStatistics 生成统计信息
func (v *Visualizer) GenerateStatistics(classes []*analyzer.CppClass) string {
	var sb strings.Builder

	totalClasses := len(classes)
	rootClasses := 0
	maxDepth := 0
	totalMembers := 0
	totalMethods := 0

	// 计算统计信息
	depthMap := make(map[string]int)

	// 找到根类
	for _, class := range classes {
		if len(class.BaseClasses) == 0 {
			rootClasses++
			depthMap[class.Name] = 0
		}
		totalMembers += len(class.Members)
		totalMethods += len(class.Methods)
	}

	// 计算继承深度
	changed := true
	for changed {
		changed = false
		for _, class := range classes {
			if _, exists := depthMap[class.Name]; exists {
				continue
			}

			maxParentDepth := -1
			allParentsProcessed := true

			for _, baseClass := range class.BaseClasses {
				if depth, exists := depthMap[baseClass]; exists {
					if depth > maxParentDepth {
						maxParentDepth = depth
					}
				} else {
					allParentsProcessed = false
					break
				}
			}

			if allParentsProcessed && maxParentDepth >= 0 {
				depthMap[class.Name] = maxParentDepth + 1
				if maxParentDepth+1 > maxDepth {
					maxDepth = maxParentDepth + 1
				}
				changed = true
			}
		}
	}

	sb.WriteString("📊 继承关系统计信息\n")
	sb.WriteString(strings.Repeat("=", 25) + "\n")
	sb.WriteString(fmt.Sprintf("总类数量: %d\n", totalClasses))
	sb.WriteString(fmt.Sprintf("根类数量: %d\n", rootClasses))
	sb.WriteString(fmt.Sprintf("最大继承深度: %d\n", maxDepth))
	sb.WriteString(fmt.Sprintf("总成员变量: %d\n", totalMembers))
	sb.WriteString(fmt.Sprintf("总成员方法: %d\n", totalMethods))

	if totalClasses > 0 {
		sb.WriteString(fmt.Sprintf("平均成员变量/类: %.1f\n", float64(totalMembers)/float64(totalClasses)))
		sb.WriteString(fmt.Sprintf("平均成员方法/类: %.1f\n", float64(totalMethods)/float64(totalClasses)))
	}
	return sb.String()
}
