// 几何图形具体实现
#ifndef GEOMETRY_H
#define GEOMETRY_H

#include "shape.h"
#include <cmath>

// 圆形类
class Circle : public Shape2D {
private:
    double radius;

public:
    Circle(const std::string& color, double r) 
        : Shape2D(color), radius(r) {}
    
    double calculateArea() override {
        area = M_PI * radius * radius;
        return area;
    }
    
    double calculatePerimeter() override {
        perimeter = 2 * M_PI * radius;
        return perimeter;
    }
    
    void draw() override {
        // 绘制圆形的实现
    }
    
    double getRadius() const { return radius; }
    void setRadius(double r) { radius = r; }
};

// 矩形类
class Rectangle : public Shape2D {
protected:
    double width;
    double height;

public:
    Rectangle(const std::string& color, double w, double h)
        : Shape2D(color), width(w), height(h) {}
    
    double calculateArea() override {
        area = width * height;
        return area;
    }
    
    double calculatePerimeter() override {
        perimeter = 2 * (width + height);
        return perimeter;
    }
    
    void draw() override {
        // 绘制矩形的实现
    }
    
    double getWidth() const { return width; }
    double getHeight() const { return height; }
    void setDimensions(double w, double h) { width = w; height = h; }
};

// 正方形类（继承自矩形）
class Square : public Rectangle {
public:
    Square(const std::string& color, double side)
        : Rectangle(color, side, side) {}
    
    void setSide(double side) {
        width = height = side;
    }
    
    double getSide() const { return width; }
};

#endif // GEOMETRY_H
