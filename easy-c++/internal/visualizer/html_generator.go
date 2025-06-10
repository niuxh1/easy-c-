package visualizer

import (
	"fmt"
	"strings"
)

// HTMLGenerator generates interactive HTML diagrams for class inheritance
type HTMLGenerator struct {
	classes map[string]*HTMLClass
}

// HTMLClass represents a class for HTML visualization
type HTMLClass struct {
	Name     string
	Members  []string
	Methods  []string
	Parents  []string
	Children []string
	Level    int
	FilePath string
}

// NewHTMLGenerator creates a new HTML generator
func NewHTMLGenerator() *HTMLGenerator {
	return &HTMLGenerator{
		classes: make(map[string]*HTMLClass),
	}
}

// AddClass adds a class to the HTML diagram
func (h *HTMLGenerator) AddClass(name string, members, methods, parents []string, filePath string) {
	h.classes[name] = &HTMLClass{
		Name:     name,
		Members:  members,
		Methods:  methods,
		Parents:  parents,
		Children: []string{},
		FilePath: filePath,
	}
}

// calculateRelationships builds parent-child relationships
func (h *HTMLGenerator) calculateRelationships() {
	// Build children relationships
	for className, class := range h.classes {
		for _, parentName := range class.Parents {
			if parent, exists := h.classes[parentName]; exists {
				parent.Children = append(parent.Children, className)
			}
		}
	}

	// Calculate inheritance levels
	h.calculateLevels()
}

// calculateLevels determines the inheritance level of each class
func (h *HTMLGenerator) calculateLevels() { // Find root classes (classes with no parents)
	for _, class := range h.classes {
		if len(class.Parents) == 0 {
			class.Level = 0
		} else {
			class.Level = -1 // Mark as unprocessed
		}
	}

	// Calculate levels for other classes
	changed := true
	for changed {
		changed = false
		for _, class := range h.classes {
			if class.Level >= 0 {
				continue
			}

			maxParentLevel := -1
			allParentsProcessed := true

			for _, parentName := range class.Parents {
				if parent, exists := h.classes[parentName]; exists {
					if parent.Level >= 0 {
						if parent.Level > maxParentLevel {
							maxParentLevel = parent.Level
						}
					} else {
						allParentsProcessed = false
						break
					}
				}
			}

			if allParentsProcessed && maxParentLevel >= 0 {
				class.Level = maxParentLevel + 1
				changed = true
			}
		}
	}

	// For classes still not processed (circular inheritance), assign them level 0
	for _, class := range h.classes {
		if class.Level < 0 {
			class.Level = 0
		}
	}
}

// GenerateHTML generates the complete interactive HTML content
func (h *HTMLGenerator) GenerateHTML() string {
	h.calculateRelationships()

	var sb strings.Builder

	// HTML header with embedded CSS and JavaScript
	sb.WriteString(h.generateHTMLHeader())

	// Generate inheritance tree view
	sb.WriteString(h.generateTreeView())

	// Generate detailed class cards
	sb.WriteString(h.generateClassCards())

	// HTML footer
	sb.WriteString(h.generateHTMLFooter())

	return sb.String()
}

// generateHTMLHeader generates the HTML header with styles and scripts
func (h *HTMLGenerator) generateHTMLHeader() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>C++ 类继承关系图</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        
        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: white;
            border-radius: 15px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(45deg, #2c3e50, #34495e);
            color: white;
            padding: 30px;
            text-align: center;
        }
        
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
        }
        
        .header p {
            font-size: 1.2em;
            opacity: 0.9;
        }
        
        .content {
            display: flex;
            min-height: 600px;
        }
        
        .tree-panel {
            flex: 1;
            background: #f8f9fa;
            border-right: 2px solid #e9ecef;
            padding: 20px;
            overflow-y: auto;
            max-height: 80vh;
        }
        
        .details-panel {
            flex: 1;
            padding: 20px;
            overflow-y: auto;
            max-height: 80vh;
        }
        
        .tree-title {
            font-size: 1.4em;
            font-weight: bold;
            color: #2c3e50;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid #3498db;
        }
        
        .level {
            margin-bottom: 25px;
        }
        
        .level-header {
            font-size: 1.1em;
            font-weight: bold;
            color: #34495e;
            margin-bottom: 12px;
            padding: 8px 15px;
            background: linear-gradient(45deg, #3498db, #2980b9);
            color: white;
            border-radius: 8px;
            text-align: center;
        }
        
        .class-node {
            background: white;
            border: 2px solid #bdc3c7;
            border-radius: 10px;
            margin: 8px 0;
            padding: 12px 15px;
            cursor: pointer;
            transition: all 0.3s ease;
            position: relative;
        }
        
        .class-node:hover {
            border-color: #3498db;
            transform: translateX(5px);
            box-shadow: 0 5px 15px rgba(52, 152, 219, 0.2);
        }
        
        .class-node.selected {
            border-color: #e74c3c;
            background: linear-gradient(45deg, #e74c3c, #c0392b);
            color: white;
        }
        
        .class-name {
            font-weight: bold;
            font-size: 1.1em;
            margin-bottom: 4px;
        }
        
        .class-file {
            font-size: 0.85em;
            color: #7f8c8d;
            font-style: italic;
        }
        
        .class-node.selected .class-file {
            color: rgba(255,255,255,0.8);
        }
        
        .inheritance-info {
            margin-top: 8px;
            font-size: 0.9em;
        }
        
        .parents {
            color: #27ae60;
        }
        
        .children {
            color: #8e44ad;
        }
        
        .class-card {
            background: white;
            border-radius: 12px;
            box-shadow: 0 8px 25px rgba(0,0,0,0.1);
            margin-bottom: 20px;
            overflow: hidden;
            border: 2px solid #ecf0f1;
            display: none;
        }
        
        .class-card.active {
            display: block;
            animation: slideIn 0.3s ease;
        }
        
        @keyframes slideIn {
            from {
                opacity: 0;
                transform: translateY(20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }
        
        .card-header {
            background: linear-gradient(45deg, #e74c3c, #c0392b);
            color: white;
            padding: 20px;
            font-size: 1.4em;
            font-weight: bold;
        }
        
        .card-content {
            padding: 20px;
        }
        
        .section {
            margin-bottom: 20px;
        }
        
        .section-title {
            font-size: 1.2em;
            font-weight: bold;
            color: #2c3e50;
            margin-bottom: 10px;
            padding-bottom: 5px;
            border-bottom: 2px solid #3498db;
        }
        
        .member-list {
            list-style: none;
        }
        
        .member-list li {
            padding: 5px 0;
            padding-left: 20px;
            position: relative;
        }
        
        .member-list li:before {
            content: "•";
            color: #3498db;
            font-weight: bold;
            position: absolute;
            left: 0;
        }
        
        .method-list li:before {
            content: "⚡";
        }
        
        .inheritance-list {
            display: flex;
            flex-wrap: wrap;
            gap: 8px;
        }
        
        .inheritance-item {
            background: #3498db;
            color: white;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 0.9em;
            cursor: pointer;
            transition: all 0.2s ease;
        }
        
        .inheritance-item:hover {
            background: #2980b9;
            transform: scale(1.05);
        }
        
        .children-item {
            background: #8e44ad;
        }
        
        .children-item:hover {
            background: #732d91;
        }
        
        .no-data {
            color: #7f8c8d;
            font-style: italic;
        }
        
        .stats {
            background: #ecf0f1;
            padding: 15px;
            border-radius: 8px;
            margin-top: 20px;
        }
        
        .stat-item {
            display: inline-block;
            margin-right: 20px;
            font-weight: bold;
        }
        
        .stat-number {
            color: #e74c3c;
            font-size: 1.2em;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🏗️ C++ 类继承关系分析</h1>
            <p>交互式继承关系图表 - 点击类名查看详细信息</p>
        </div>
        <div class="content">
            <div class="tree-panel">
                <div class="tree-title">📊 继承层级结构</div>
`
}

// generateTreeView generates the hierarchical tree view
func (h *HTMLGenerator) generateTreeView() string {
	var sb strings.Builder

	// Group classes by level
	levelGroups := make(map[int][]*HTMLClass)
	maxLevel := 0

	for _, class := range h.classes {
		levelGroups[class.Level] = append(levelGroups[class.Level], class)
		if class.Level > maxLevel {
			maxLevel = class.Level
		}
	}

	// Generate each level
	for level := 0; level <= maxLevel; level++ {
		if classes, exists := levelGroups[level]; exists {
			levelName := "根类"
			if level > 0 {
				levelName = fmt.Sprintf("第 %d 层", level)
			}

			sb.WriteString(fmt.Sprintf(`                <div class="level">
                    <div class="level-header">%s (%d 个类)</div>
`, levelName, len(classes)))

			for _, class := range classes {
				sb.WriteString(fmt.Sprintf(`                    <div class="class-node" onclick="showClassDetails('%s')">
                        <div class="class-name">%s</div>
                        <div class="class-file">📁 %s</div>
`, class.Name, class.Name, class.FilePath))

				if len(class.Parents) > 0 {
					sb.WriteString(fmt.Sprintf(`                        <div class="inheritance-info parents">
                            ⬆️ 继承自: %s
                        </div>
`, strings.Join(class.Parents, ", ")))
				}

				if len(class.Children) > 0 {
					sb.WriteString(fmt.Sprintf(`                        <div class="inheritance-info children">
                            ⬇️ 子类: %s
                        </div>
`, strings.Join(class.Children, ", ")))
				}

				sb.WriteString("                    </div>\n")
			}

			sb.WriteString("                </div>\n")
		}
	}

	// Add statistics
	totalClasses := len(h.classes)
	rootClasses := len(levelGroups[0])
	maxDepth := maxLevel

	sb.WriteString(fmt.Sprintf(`                <div class="stats">
                    <div class="stat-item">📈 总类数: <span class="stat-number">%d</span></div>
                    <div class="stat-item">🌳 根类数: <span class="stat-number">%d</span></div>
                    <div class="stat-item">📏 最大深度: <span class="stat-number">%d</span></div>
                </div>
`, totalClasses, rootClasses, maxDepth))

	sb.WriteString("            </div>\n            <div class=\"details-panel\">\n")

	return sb.String()
}

// generateClassCards generates detailed cards for each class
func (h *HTMLGenerator) generateClassCards() string {
	var sb strings.Builder

	sb.WriteString(`                <div id="welcome-message" class="class-card active">
                    <div class="card-header">👋 欢迎使用继承关系分析器</div>
                    <div class="card-content">
                        <p style="font-size: 1.1em; line-height: 1.6; color: #2c3e50;">
                            点击左侧的任意类名来查看详细的类信息，包括：
                        </p>
                        <ul style="margin: 20px 0; color: #34495e;">
                            <li>🔧 成员变量和方法</li>
                            <li>⬆️ 父类继承关系</li>
                            <li>⬇️ 子类派生关系</li>
                            <li>📁 所在文件位置</li>
                        </ul>
                        <p style="color: #7f8c8d;">
                            这个交互式图表让您能够更清晰地理解复杂的C++继承关系。
                        </p>
                    </div>
                </div>
`)

	for _, class := range h.classes {
		sb.WriteString(fmt.Sprintf(`                <div id="card-%s" class="class-card">
                    <div class="card-header">🎯 %s</div>
                    <div class="card-content">
`, class.Name, class.Name))

		// File info
		sb.WriteString(fmt.Sprintf(`                        <div class="section">
                            <div class="section-title">📁 文件位置</div>
                            <p>%s</p>
                        </div>
`, class.FilePath))

		// Parents
		sb.WriteString(`                        <div class="section">
                            <div class="section-title">⬆️ 继承关系</div>
`)
		if len(class.Parents) > 0 {
			sb.WriteString(`                            <div class="inheritance-list">
`)
			for _, parent := range class.Parents {
				sb.WriteString(fmt.Sprintf(`                                <span class="inheritance-item" onclick="showClassDetails('%s')">%s</span>
`, parent, parent))
			}
			sb.WriteString(`                            </div>
`)
		} else {
			sb.WriteString(`                            <p class="no-data">此类为根类，无父类</p>
`)
		}
		sb.WriteString(`                        </div>
`)

		// Children
		sb.WriteString(`                        <div class="section">
                            <div class="section-title">⬇️ 派生类</div>
`)
		if len(class.Children) > 0 {
			sb.WriteString(`                            <div class="inheritance-list">
`)
			for _, child := range class.Children {
				sb.WriteString(fmt.Sprintf(`                                <span class="inheritance-item children-item" onclick="showClassDetails('%s')">%s</span>
`, child, child))
			}
			sb.WriteString(`                            </div>
`)
		} else {
			sb.WriteString(`                            <p class="no-data">此类无派生类</p>
`)
		}
		sb.WriteString(`                        </div>
`)

		// Members
		sb.WriteString(`                        <div class="section">
                            <div class="section-title">🔧 成员变量</div>
`)
		if len(class.Members) > 0 {
			sb.WriteString(`                            <ul class="member-list">
`)
			for _, member := range class.Members {
				sb.WriteString(fmt.Sprintf(`                                <li>%s</li>
`, member))
			}
			sb.WriteString(`                            </ul>
`)
		} else {
			sb.WriteString(`                            <p class="no-data">未发现成员变量</p>
`)
		}
		sb.WriteString(`                        </div>
`)

		// Methods
		sb.WriteString(`                        <div class="section">
                            <div class="section-title">⚡ 成员方法</div>
`)
		if len(class.Methods) > 0 {
			sb.WriteString(`                            <ul class="member-list method-list">
`)
			for _, method := range class.Methods {
				sb.WriteString(fmt.Sprintf(`                                <li>%s</li>
`, method))
			}
			sb.WriteString(`                            </ul>
`)
		} else {
			sb.WriteString(`                            <p class="no-data">未发现成员方法</p>
`)
		}
		sb.WriteString(`                        </div>
`)

		sb.WriteString(`                    </div>
                </div>
`)
	}

	return sb.String()
}

// generateHTMLFooter generates the HTML footer with JavaScript
func (h *HTMLGenerator) generateHTMLFooter() string {
	return `            </div>
        </div>
    </div>
    
    <script>
        function showClassDetails(className) {
            // Hide all cards
            const cards = document.querySelectorAll('.class-card');
            cards.forEach(card => {
                card.classList.remove('active');
            });
            
            // Remove selection from all nodes
            const nodes = document.querySelectorAll('.class-node');
            nodes.forEach(node => {
                node.classList.remove('selected');
            });
            
            // Show selected card
            const targetCard = document.getElementById('card-' + className);
            if (targetCard) {
                targetCard.classList.add('active');
            }
            
            // Highlight selected node
            const targetNodes = document.querySelectorAll('.class-node');
            targetNodes.forEach(node => {
                if (node.querySelector('.class-name').textContent === className) {
                    node.classList.add('selected');
                }
            });
        }
        
        // Add some interactive effects
        document.addEventListener('DOMContentLoaded', function() {
            // Smooth scrolling for inheritance item clicks
            const inheritanceItems = document.querySelectorAll('.inheritance-item');
            inheritanceItems.forEach(item => {
                item.addEventListener('click', function(e) {
                    e.stopPropagation();
                    const className = this.textContent;
                    showClassDetails(className);
                    
                    // Find and scroll to the corresponding node in tree
                    const nodes = document.querySelectorAll('.class-node');
                    nodes.forEach(node => {
                        if (node.querySelector('.class-name').textContent === className) {
                            node.scrollIntoView({ behavior: 'smooth', block: 'center' });
                        }
                    });
                });
            });
        });
    </script>
</body>
</html>`
}
