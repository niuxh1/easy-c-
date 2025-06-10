package analyzer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCppAnalyzer_AnalyzeFile 测试文件分析功能
func TestCppAnalyzer_AnalyzeFile(t *testing.T) {
	// 创建测试用的临时C++文件
	testCode := `
// 测试基类
class Animal {
private:
    std::string name;
    int age;

public:
    Animal(const std::string& n, int a);
    virtual ~Animal();
    virtual void speak() = 0;
    std::string getName() const;
};

// 测试派生类
class Dog : public Animal {
private:
    std::string breed;

public:
    Dog(const std::string& n, int a, const std::string& b);
    void speak() override;
    void bark();
};

// 测试多重继承
class WorkingDog : public Dog {
private:
    std::string jobType;

public:
    WorkingDog(const std::string& n, int a, const std::string& b, const std::string& job);
    void performJob();
};
`

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "test*.cpp")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(testCode); err != nil {
		t.Fatalf("写入测试代码失败: %v", err)
	}
	tmpFile.Close()

	// 测试分析器
	analyzer := NewCppAnalyzer()
	classes, err := analyzer.AnalyzeFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("分析文件失败: %v", err)
	}

	// 验证结果
	if len(classes) != 3 {
		t.Errorf("期望找到3个类，实际找到%d个", len(classes))
	}

	// 验证Animal类
	animal := findClassByName(classes, "Animal")
	if animal == nil {
		t.Error("未找到Animal类")
	} else {
		if len(animal.BaseClasses) != 0 {
			t.Errorf("Animal应该是根类，但找到基类: %v", animal.BaseClasses)
		}
		if len(animal.Members) < 1 {
			t.Errorf("Animal应该有至少1个成员变量，实际有%d个", len(animal.Members))
		}
	}

	// 验证Dog类
	dog := findClassByName(classes, "Dog")
	if dog == nil {
		t.Error("未找到Dog类")
	} else {
		if len(dog.BaseClasses) != 1 || dog.BaseClasses[0] != "Animal" {
			t.Errorf("Dog应该继承自Animal，实际继承自: %v", dog.BaseClasses)
		}
	}

	// 验证WorkingDog类
	workingDog := findClassByName(classes, "WorkingDog")
	if workingDog == nil {
		t.Error("未找到WorkingDog类")
	} else {
		if len(workingDog.BaseClasses) != 1 || workingDog.BaseClasses[0] != "Dog" {
			t.Errorf("WorkingDog应该继承自Dog，实际继承自: %v", workingDog.BaseClasses)
		}
	}
}

// TestCppAnalyzer_ParseInheritance 测试继承关系解析
func TestCppAnalyzer_ParseInheritance(t *testing.T) {
	analyzer := NewCppAnalyzer()

	tests := []struct {
		input    string
		expected []string
	}{
		{"public Animal", []string{"Animal"}},
		{"private Animal", []string{"Animal"}},
		{"protected Animal", []string{"Animal"}},
		{"public Animal, private Mammal", []string{"Animal", "Mammal"}},
		{"Animal", []string{"Animal"}},
	}

	for _, test := range tests {
		result := analyzer.parseInheritance(test.input)
		if !sliceEqual(result, test.expected) {
			t.Errorf("解析继承关系 '%s' 失败，期望: %v，实际: %v", test.input, test.expected, result)
		}
	}
}

// TestCppAnalyzer_RemoveComments 测试注释移除功能
func TestCppAnalyzer_RemoveComments(t *testing.T) {
	analyzer := NewCppAnalyzer()

	tests := []struct {
		input    string
		expected string
	}{
		{"class Test { // 这是注释", "class Test { "},
		{"/* 块注释 */ class Test", " class Test"},
		{"class Test /* 内联块注释 */ {", "class Test  {"},
		{"// 整行注释", ""},
		{"class Test;", "class Test;"},
	}

	for _, test := range tests {
		inBlockComment := false
		result := analyzer.removeComments(test.input, &inBlockComment)
		if strings.TrimSpace(result) != strings.TrimSpace(test.expected) {
			t.Errorf("移除注释失败，输入: '%s'，期望: '%s'，实际: '%s'",
				test.input, test.expected, result)
		}
	}
}

// TestGetInheritanceTree 测试继承树构建
func TestGetInheritanceTree(t *testing.T) {
	classes := []*CppClass{
		{Name: "Animal", BaseClasses: []string{}},
		{Name: "Dog", BaseClasses: []string{"Animal"}},
		{Name: "Cat", BaseClasses: []string{"Animal"}},
		{Name: "WorkingDog", BaseClasses: []string{"Dog"}},
	}

	tree := GetInheritanceTree(classes)

	// 验证Animal的子类
	if len(tree["Animal"]) != 2 {
		t.Errorf("Animal应该有2个子类，实际有%d个", len(tree["Animal"]))
	}

	// 验证Dog的子类
	if len(tree["Dog"]) != 1 {
		t.Errorf("Dog应该有1个子类，实际有%d个", len(tree["Dog"]))
	}

	// 验证WorkingDog是Dog的子类
	found := false
	for _, child := range tree["Dog"] {
		if child.Name == "WorkingDog" {
			found = true
			break
		}
	}
	if !found {
		t.Error("WorkingDog应该是Dog的子类")
	}
}

// TestFindRootClasses 测试根类查找
func TestFindRootClasses(t *testing.T) {
	classes := []*CppClass{
		{Name: "Animal", BaseClasses: []string{}},
		{Name: "Vehicle", BaseClasses: []string{}},
		{Name: "Dog", BaseClasses: []string{"Animal"}},
		{Name: "Car", BaseClasses: []string{"Vehicle"}},
	}

	rootClasses := FindRootClasses(classes)

	if len(rootClasses) != 2 {
		t.Errorf("应该找到2个根类，实际找到%d个", len(rootClasses))
	}

	expectedRoots := map[string]bool{"Animal": true, "Vehicle": true}
	for _, root := range rootClasses {
		if !expectedRoots[root.Name] {
			t.Errorf("意外的根类: %s", root.Name)
		}
	}
}

// TestAnalyzerWithComplexInheritance 测试复杂继承关系
func TestAnalyzerWithComplexInheritance(t *testing.T) {
	testCode := `
class Base1 {
public:
    virtual void method1() = 0;
};

class Base2 {
public:
    virtual void method2() = 0;
};

// 多重继承
class Derived : public Base1, public Base2 {
public:
    void method1() override;
    void method2() override;
    void additionalMethod();
};

// 虚继承
class VirtualBase {
public:
    int value;
};

class VirtualDerived1 : virtual public VirtualBase {
public:
    void func1();
};

class VirtualDerived2 : virtual public VirtualBase {
public:
    void func2();
};

class DiamondInheritance : public VirtualDerived1, public VirtualDerived2 {
public:
    void finalMethod();
};
`

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "complex*.cpp")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(testCode); err != nil {
		t.Fatalf("写入测试代码失败: %v", err)
	}
	tmpFile.Close()

	// 分析文件
	analyzer := NewCppAnalyzer()
	classes, err := analyzer.AnalyzeFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("分析复杂继承失败: %v", err)
	}
	// 验证类的数量 (期望找到7个类，包括Pet模板类)
	expectedCount := 7
	if len(classes) != expectedCount {
		t.Logf("找到的类:")
		for i, class := range classes {
			t.Logf("  %d. %s", i+1, class.Name)
		}
		t.Errorf("期望找到%d个类，实际找到%d个", expectedCount, len(classes))
	}

	// 验证多重继承
	derived := findClassByName(classes, "Derived")
	if derived == nil {
		t.Error("未找到Derived类")
	} else {
		expectedBases := []string{"Base1", "Base2"}
		if !sliceEqual(derived.BaseClasses, expectedBases) {
			t.Errorf("Derived类的基类不正确，期望: %v，实际: %v", expectedBases, derived.BaseClasses)
		}
	}

	// 验证菱形继承
	diamond := findClassByName(classes, "DiamondInheritance")
	if diamond == nil {
		t.Error("未找到DiamondInheritance类")
	} else {
		expectedBases := []string{"VirtualDerived1", "VirtualDerived2"}
		if !sliceEqual(diamond.BaseClasses, expectedBases) {
			t.Errorf("DiamondInheritance类的基类不正确，期望: %v，实际: %v", expectedBases, diamond.BaseClasses)
		}
	}
}

// TestAnalyzeProject 测试项目目录分析功能
func TestAnalyzeProject(t *testing.T) {
	analyzer := NewCppAnalyzer()

	// 创建临时测试目录
	tempDir := t.TempDir()

	// 创建测试文件1
	file1Content := `
class BaseClass {
public:
    int baseValue;
    virtual void baseMethod() {}
};`
	file1Path := filepath.Join(tempDir, "base.h")
	err := os.WriteFile(file1Path, []byte(file1Content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建测试文件2
	file2Content := `
#include "base.h"
class DerivedClass : public BaseClass {
private:
    double derivedValue;
public:
    void derivedMethod() {}
};`
	file2Path := filepath.Join(tempDir, "derived.cpp")
	err = os.WriteFile(file2Path, []byte(file2Content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 分析项目
	classes, err := analyzer.AnalyzeProject(tempDir)
	if err != nil {
		t.Fatalf("项目分析失败: %v", err)
	}

	// 验证结果
	if len(classes) != 2 {
		t.Errorf("期望找到2个类，实际找到%d个", len(classes))
	}

	// 验证类名和继承关系
	classMap := make(map[string]*CppClass)
	for _, class := range classes {
		classMap[class.Name] = class
	}

	if baseClass, exists := classMap["BaseClass"]; !exists {
		t.Error("未找到BaseClass")
	} else {
		if len(baseClass.BaseClasses) != 0 {
			t.Error("BaseClass不应该有基类")
		}
		if baseClass.FilePath == "" {
			t.Error("BaseClass应该有文件路径")
		}
	}

	if derivedClass, exists := classMap["DerivedClass"]; !exists {
		t.Error("未找到DerivedClass")
	} else {
		if len(derivedClass.BaseClasses) != 1 || derivedClass.BaseClasses[0] != "BaseClass" {
			t.Errorf("DerivedClass应该继承自BaseClass，实际: %v", derivedClass.BaseClasses)
		}
		if derivedClass.FilePath == "" {
			t.Error("DerivedClass应该有文件路径")
		}
	}
}

// TestAnalyzeFiles 测试多文件分析功能
func TestAnalyzeFiles(t *testing.T) {
	analyzer := NewCppAnalyzer()

	// 创建临时测试文件
	tempDir := t.TempDir()

	files := []string{
		filepath.Join(tempDir, "class1.h"),
		filepath.Join(tempDir, "class2.h"),
	}

	contents := []string{
		`class Shape {
public:
    virtual void draw() = 0;
};`,
		`class Circle : public Shape {
private:
    double radius;
public:
    void draw() override {}
};`,
	}

	// 创建测试文件
	for i, content := range contents {
		err := os.WriteFile(files[i], []byte(content), 0644)
		if err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}
	}

	// 分析多个文件
	classes, err := analyzer.AnalyzeFiles(files)
	if err != nil {
		t.Fatalf("多文件分析失败: %v", err)
	}

	// 验证结果
	if len(classes) != 2 {
		t.Errorf("期望找到2个类，实际找到%d个", len(classes))
	}

	// 验证继承关系解析
	circleClass := classes[1] // Circle应该是第二个
	if len(circleClass.BaseClasses) != 1 || circleClass.BaseClasses[0] != "Shape" {
		t.Errorf("Circle应该继承自Shape，实际: %v", circleClass.BaseClasses)
	}
}

// TestComplexInheritance 测试复杂继承关系
func TestComplexInheritance(t *testing.T) {
	analyzer := NewCppAnalyzer()

	// 创建包含多重继承的测试代码
	content := `
class Interface1 {
public:
    virtual void method1() = 0;
};

class Interface2 {
public:
    virtual void method2() = 0;
};

class BaseClass {
protected:
    int baseValue;
public:
    virtual void baseMethod() {}
};

class MultiInheritance : public BaseClass, public Interface1, public Interface2 {
private:
    double value;
public:
    void method1() override {}
    void method2() override {}
    void multiMethod() {}
};`

	tempFile := filepath.Join(t.TempDir(), "multi.cpp")
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	classes, err := analyzer.AnalyzeFile(tempFile)
	if err != nil {
		t.Fatalf("分析失败: %v", err)
	}

	// 验证类数量
	if len(classes) != 4 {
		t.Errorf("期望找到4个类，实际找到%d个", len(classes))
	}

	// 查找多重继承类
	var multiClass *CppClass
	for _, class := range classes {
		if class.Name == "MultiInheritance" {
			multiClass = class
			break
		}
	}

	if multiClass == nil {
		t.Fatal("未找到MultiInheritance类")
	}

	// 验证多重继承
	expectedBases := []string{"BaseClass", "Interface1", "Interface2"}
	if len(multiClass.BaseClasses) != len(expectedBases) {
		t.Errorf("期望继承%d个基类，实际继承%d个", len(expectedBases), len(multiClass.BaseClasses))
	}

	for _, expected := range expectedBases {
		found := false
		for _, actual := range multiClass.BaseClasses {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("未找到期望的基类: %s", expected)
		}
	}
}

// 辅助函数：根据名称查找类
func findClassByName(classes []*CppClass, name string) *CppClass {
	for _, class := range classes {
		if class.Name == name {
			return class
		}
	}
	return nil
}

// 辅助函数：比较两个字符串切片是否相等
func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
