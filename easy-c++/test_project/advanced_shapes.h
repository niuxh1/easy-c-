// 特殊形状和工具类
#ifndef ADVANCED_SHAPES_H
#define ADVANCED_SHAPES_H

#include "geometry.h"
#include "shapes_3d.h"

// 可绘制接口
class Drawable {
public:
    virtual ~Drawable() {}
    virtual void render() = 0;
    virtual void setRenderMode(int mode) = 0;
};

// 可变换接口
class Transformable {
public:
    virtual ~Transformable() {}
    virtual void translate(double x, double y, double z = 0) = 0;
    virtual void rotate(double angle) = 0;
    virtual void scale(double factor) = 0;
};

// 高级圆形（多重继承）
class AdvancedCircle : public Circle, public Drawable, public Transformable {
private:
    double posX, posY;
    double rotationAngle;
    double scaleFactor;
    int renderMode;

public:
    AdvancedCircle(const std::string& color, double radius, double x = 0, double y = 0)
        : Circle(color, radius), posX(x), posY(y), 
          rotationAngle(0), scaleFactor(1.0), renderMode(0) {}
    
    // Drawable接口实现
    void render() override {
        // 高级渲染实现
    }
    
    void setRenderMode(int mode) override {
        renderMode = mode;
    }
    
    // Transformable接口实现
    void translate(double x, double y, double z = 0) override {
        posX += x;
        posY += y;
    }
    
    void rotate(double angle) override {
        rotationAngle += angle;
    }
    
    void scale(double factor) override {
        scaleFactor *= factor;
    }
    
    // 获取位置信息
    double getX() const { return posX; }
    double getY() const { return posY; }
    double getRotation() const { return rotationAngle; }
    double getScaleFactor() const { return scaleFactor; }
};

// 复合形状管理器
class ShapeManager {
private:
    static ShapeManager* instance;
    int shapeCount;

protected:
    ShapeManager() : shapeCount(0) {}

public:
    static ShapeManager* getInstance() {
        if (instance == nullptr) {
            instance = new ShapeManager();
        }
        return instance;
    }
    
    void addShape() { shapeCount++; }
    void removeShape() { if (shapeCount > 0) shapeCount--; }
    int getShapeCount() const { return shapeCount; }
    
    virtual ~ShapeManager() {}
};

// 形状工厂（抽象工厂模式）
class ShapeFactory {
public:
    virtual ~ShapeFactory() {}
    virtual Shape2D* create2DShape(const std::string& type, const std::string& color) = 0;
    virtual Shape3D* create3DShape(const std::string& type, const std::string& color) = 0;
};

// 具体工厂实现
class ConcreteShapeFactory : public ShapeFactory {
public:
    Shape2D* create2DShape(const std::string& type, const std::string& color) override {
        if (type == "circle") {
            return new Circle(color, 1.0);
        } else if (type == "rectangle") {
            return new Rectangle(color, 1.0, 1.0);
        } else if (type == "square") {
            return new Square(color, 1.0);
        }
        return nullptr;
    }
    
    Shape3D* create3DShape(const std::string& type, const std::string& color) override {
        if (type == "sphere") {
            return new Sphere(color, 1.0);
        } else if (type == "cube") {
            return new Cube(color, 1.0);
        } else if (type == "cylinder") {
            return new Cylinder(color, 1.0, 1.0);
        }
        return nullptr;
    }
};

#endif // ADVANCED_SHAPES_H
