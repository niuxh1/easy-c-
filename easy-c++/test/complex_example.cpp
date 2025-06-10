// 更复杂的C++继承示例，用于测试边界情况
#include <iostream>
#include <vector>
#include <memory>

// 抽象基类
class Shape {
protected:
    double x, y;  // 位置坐标
    static int count;  // 静态成员

public:
    Shape(double x = 0, double y = 0);
    virtual ~Shape() = default;
    
    // 纯虚函数
    virtual double area() const = 0;
    virtual double perimeter() const = 0;
    virtual void draw() const = 0;
    
    // 虚函数
    virtual void move(double dx, double dy);
    virtual std::string toString() const;
    
    // 普通成员函数
    double getX() const { return x; }
    double getY() const { return y; }
    static int getCount() { return count; }
};

// 接口类
class Drawable {
public:
    virtual ~Drawable() = default;
    virtual void render() = 0;
    virtual void setColor(const std::string& color) = 0;
};

// 另一个接口
class Serializable {
public:
    virtual ~Serializable() = default;
    virtual std::string serialize() const = 0;
    virtual void deserialize(const std::string& data) = 0;
};

// 矩形类 - 多重继承
class Rectangle : public Shape, public Drawable, public Serializable {
private:
    double width, height;
    std::string color;

public:
    Rectangle(double x, double y, double w, double h);
    virtual ~Rectangle();
    
    // 实现Shape的纯虚函数
    double area() const override;
    double perimeter() const override;
    void draw() const override;
    
    // 实现Drawable接口
    void render() override;
    void setColor(const std::string& c) override;
    
    // 实现Serializable接口
    std::string serialize() const override;
    void deserialize(const std::string& data) override;
    
    // Rectangle特有方法
    double getWidth() const { return width; }
    double getHeight() const { return height; }
    void resize(double w, double h);
};

// 圆形类
class Circle : public Shape, public Drawable {
private:
    double radius;
    std::string color;

public:
    Circle(double x, double y, double r);
    virtual ~Circle();
    
    double area() const override;
    double perimeter() const override;
    void draw() const override;
    
    void render() override;
    void setColor(const std::string& c) override;
    
    double getRadius() const { return radius; }
    void setRadius(double r) { radius = r; }
};

// 正方形 - 继承自Rectangle
class Square : public Rectangle {
public:
    Square(double x, double y, double side);
    virtual ~Square();
    
    void resize(double side);  // 重写resize方法
    std::string toString() const override;
};

// 椭圆类
class Ellipse : public Shape {
private:
    double majorAxis, minorAxis;

public:
    Ellipse(double x, double y, double major, double minor);
    virtual ~Ellipse();
    
    double area() const override;
    double perimeter() const override;
    void draw() const override;
    
    double getMajorAxis() const { return majorAxis; }
    double getMinorAxis() const { return minorAxis; }
};

// 可填充形状的mixin类
template<typename T>
class Fillable {
private:
    bool filled;
    std::string fillColor;

public:
    Fillable() : filled(false), fillColor("white") {}
    virtual ~Fillable() = default;
    
    void setFilled(bool f) { filled = f; }
    bool isFilled() const { return filled; }
    void setFillColor(const std::string& color) { fillColor = color; }
    std::string getFillColor() const { return fillColor; }
};

// 可填充的圆形
class FilledCircle : public Circle, public Fillable<Circle> {
public:
    FilledCircle(double x, double y, double r);
    virtual ~FilledCircle();
    
    void draw() const override;  // 重写绘制方法以支持填充
    void render() override;
};

// 3D形状基类
class Shape3D {
protected:
    double z;  // Z坐标

public:
    Shape3D(double z = 0) : z(z) {}
    virtual ~Shape3D() = default;
    
    virtual double volume() const = 0;
    virtual double surfaceArea() const = 0;
    
    double getZ() const { return z; }
    void setZ(double newZ) { z = newZ; }
};

// 立方体 - 继承自Square和Shape3D
class Cube : public Square, public Shape3D {
public:
    Cube(double x, double y, double z, double side);
    virtual ~Cube();
    
    double volume() const override;
    double surfaceArea() const override;
    
    // 重写2D方法
    double area() const override;  // 返回底面积
    void draw() const override;
};

// 球体 - 继承自Circle和Shape3D
class Sphere : public Circle, public Shape3D {
public:
    Sphere(double x, double y, double z, double r);
    virtual ~Sphere();
    
    double volume() const override;
    double surfaceArea() const override;
    
    double area() const override;  // 返回截面积
    void draw() const override;
};

// 形状管理器
class ShapeManager {
private:
    std::vector<std::unique_ptr<Shape>> shapes;
    std::vector<std::unique_ptr<Shape3D>> shapes3D;

public:
    ShapeManager() = default;
    virtual ~ShapeManager() = default;
    
    void addShape(std::unique_ptr<Shape> shape);
    void addShape3D(std::unique_ptr<Shape3D> shape);
    
    double getTotalArea() const;
    double getTotalVolume() const;
    
    void drawAll() const;
    void renderAll() const;
    
    size_t getShapeCount() const { return shapes.size(); }
    size_t getShape3DCount() const { return shapes3D.size(); }
};

// 抽象工厂模式
class ShapeFactory {
public:
    virtual ~ShapeFactory() = default;
    virtual std::unique_ptr<Shape> createRectangle(double x, double y, double w, double h) = 0;
    virtual std::unique_ptr<Shape> createCircle(double x, double y, double r) = 0;
    virtual std::unique_ptr<Shape> createSquare(double x, double y, double side) = 0;
};

// 具体工厂
class ConcreteShapeFactory : public ShapeFactory {
public:
    ConcreteShapeFactory() = default;
    virtual ~ConcreteShapeFactory() = default;
    
    std::unique_ptr<Shape> createRectangle(double x, double y, double w, double h) override;
    std::unique_ptr<Shape> createCircle(double x, double y, double r) override;
    std::unique_ptr<Shape> createSquare(double x, double y, double side) override;
};

int main() {
    // 使用示例
    ShapeManager manager;
    ConcreteShapeFactory factory;
    
    // 创建各种形状
    auto rect = factory.createRectangle(0, 0, 10, 5);
    auto circle = factory.createCircle(5, 5, 3);
    auto square = factory.createSquare(10, 10, 4);
    
    manager.addShape(std::move(rect));
    manager.addShape(std::move(circle));
    manager.addShape(std::move(square));
    
    // 创建3D形状
    auto cube = std::make_unique<Cube>(0, 0, 0, 3);
    auto sphere = std::make_unique<Sphere>(10, 10, 10, 2);
    
    manager.addShape3D(std::move(cube));
    manager.addShape3D(std::move(sphere));
    
    std::cout << "总面积: " << manager.getTotalArea() << std::endl;
    std::cout << "总体积: " << manager.getTotalVolume() << std::endl;
    
    manager.drawAll();
    manager.renderAll();
    
    return 0;
}
