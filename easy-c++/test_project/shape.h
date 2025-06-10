// 基础类定义
#ifndef SHAPE_H
#define SHAPE_H

#include <string>

// 抽象基类
class Shape {
protected:
    std::string color;
    double area;

public:
    Shape(const std::string& c) : color(c), area(0.0) {}
    virtual ~Shape() {}
    
    // 纯虚函数
    virtual double calculateArea() = 0;
    virtual void draw() = 0;
    
    // 普通成员函数
    std::string getColor() const { return color; }
    void setColor(const std::string& c) { color = c; }
    double getArea() const { return area; }
};

// 2D形状基类
class Shape2D : public Shape {
protected:
    double perimeter;

public:
    Shape2D(const std::string& c) : Shape(c), perimeter(0.0) {}
    virtual ~Shape2D() {}
    
    virtual double calculatePerimeter() = 0;
    double getPerimeter() const { return perimeter; }
};

#endif // SHAPE_H
