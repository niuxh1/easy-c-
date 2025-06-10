// 简单的继承测试用例，用于验证线条间距优化
class BaseA {
public:
    virtual void methodA() = 0;
};

class BaseB {
public:
    virtual void methodB() = 0;
};

class BaseC {
public:
    virtual void methodC() = 0;
};

// 多个类继承自不同基类，用于测试线条分散效果
class Child1 : public BaseA {
public:
    void methodA() override;
};

class Child2 : public BaseA {
public:
    void methodA() override;
};

class Child3 : public BaseB {
public:
    void methodB() override;
};

class Child4 : public BaseB {
public:
    void methodB() override;
};

class Child5 : public BaseC {
public:
    void methodC() override;
};

// 多重继承测试
class MultiChild : public BaseA, public BaseB {
public:
    void methodA() override;
    void methodB() override;
};
