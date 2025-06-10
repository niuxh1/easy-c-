// 3D形状扩展
#ifndef SHAPES_3D_H
#define SHAPES_3D_H

#include "shape.h"
#include "geometry.h"
#include <cmath>

// 3D形状基类
class Shape3D : public Shape {
protected:
    double volume;
    double surfaceArea;

public:
    Shape3D(const std::string& color) 
        : Shape(color), volume(0.0), surfaceArea(0.0) {}
    virtual ~Shape3D() {}
    
    virtual double calculateVolume() = 0;
    virtual double calculateSurfaceArea() = 0;
    
    double getVolume() const { return volume; }
    double getSurfaceArea() const { return surfaceArea; }
    
    // 实现基类的纯虚函数
    double calculateArea() override {
        return calculateSurfaceArea();
    }
};

// 球体类
class Sphere : public Shape3D {
private:
    double radius;

public:
    Sphere(const std::string& color, double r)
        : Shape3D(color), radius(r) {}
    
    double calculateVolume() override {
        volume = (4.0/3.0) * M_PI * radius * radius * radius;
        return volume;
    }
    
    double calculateSurfaceArea() override {
        surfaceArea = 4 * M_PI * radius * radius;
        return surfaceArea;
    }
    
    void draw() override {
        // 绘制球体的实现
    }
    
    double getRadius() const { return radius; }
    void setRadius(double r) { radius = r; }
};

// 立方体类
class Cube : public Shape3D {
private:
    double side;

public:
    Cube(const std::string& color, double s)
        : Shape3D(color), side(s) {}
    
    double calculateVolume() override {
        volume = side * side * side;
        return volume;
    }
    
    double calculateSurfaceArea() override {
        surfaceArea = 6 * side * side;
        return surfaceArea;
    }
    
    void draw() override {
        // 绘制立方体的实现
    }
    
    double getSide() const { return side; }
    void setSide(double s) { side = s; }
};

// 圆柱体类（多重继承示例）
class Cylinder : public Shape3D {
private:
    double radius;
    double height;

public:
    Cylinder(const std::string& color, double r, double h)
        : Shape3D(color), radius(r), height(h) {}
    
    double calculateVolume() override {
        volume = M_PI * radius * radius * height;
        return volume;
    }
    
    double calculateSurfaceArea() override {
        surfaceArea = 2 * M_PI * radius * (radius + height);
        return surfaceArea;
    }
    
    void draw() override {
        // 绘制圆柱体的实现
    }
    
    double getRadius() const { return radius; }
    double getHeight() const { return height; }
    void setDimensions(double r, double h) { radius = r; height = h; }
};

#endif // SHAPES_3D_H
