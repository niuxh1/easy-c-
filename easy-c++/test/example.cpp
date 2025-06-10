// 示例C++代码 - 展示多层继承关系
#include <iostream>
#include <string>
#include <vector>

// 基类：动物
class Animal {
private:
    std::string name;
    int age;

public:
    Animal(const std::string& n, int a);
    virtual ~Animal();
    
    virtual void speak() = 0;  // 纯虚函数
    virtual void move();
    
    std::string getName() const;
    int getAge() const;
    void setName(const std::string& n);
};

// 派生类：哺乳动物
class Mammal : public Animal {
protected:
    bool hasFur;
    double bodyTemperature;

public:
    Mammal(const std::string& n, int a, bool fur);
    virtual ~Mammal();
    
    virtual void breathe();
    bool getHasFur() const;
    double getBodyTemperature() const;
};

// 派生类：鸟类
class Bird : public Animal {
protected:
    bool canFly;
    double wingSpan;

public:
    Bird(const std::string& n, int a, bool fly, double span);
    virtual ~Bird();
    
    virtual void layEggs();
    bool getCanFly() const;
    double getWingSpan() const;
    
    void speak() override;  // 实现纯虚函数
};

// 多重继承示例：蝙蝠
class Bat : public Mammal {
private:
    double flightSpeed;
    bool isNocturnal;

public:
    Bat(const std::string& n, int a);
    virtual ~Bat();
    
    void speak() override;
    void fly();
    void echolocate();
    
    double getFlightSpeed() const;
    bool getIsNocturnal() const;
};

// 具体哺乳动物：狗
class Dog : public Mammal {
private:
    std::string breed;
    bool isTrained;

public:
    Dog(const std::string& n, int a, const std::string& b);
    virtual ~Dog();
    
    void speak() override;
    void bark();
    void wagTail();
    
    std::string getBreed() const;
    bool getIsTrained() const;
    void train();
};

// 具体哺乳动物：猫
class Cat : public Mammal {
private:
    int livesRemaining;
    bool isIndoor;

public:
    Cat(const std::string& n, int a);
    virtual ~Cat();
    
    void speak() override;
    void purr();
    void climb();
    
    int getLivesRemaining() const;
    bool getIsIndoor() const;
    void setIndoor(bool indoor);
};

// 具体鸟类：鹰
class Eagle : public Bird {
private:
    double huntingRange;
    int preyCount;

public:
    Eagle(const std::string& n, int a, double range);
    virtual ~Eagle();
    
    void hunt();
    void soar();
    
    double getHuntingRange() const;
    int getPreyCount() const;
};

// 具体鸟类：企鹅
class Penguin : public Bird {
private:
    double swimSpeed;
    bool isEmperor;

public:
    Penguin(const std::string& n, int a, bool emperor);
    virtual ~Penguin();
    
    void swim();
    void slide();
    
    double getSwimSpeed() const;
    bool getIsEmperor() const;
};

// 工作犬：继承自Dog
class WorkingDog : public Dog {
private:
    std::string jobType;
    int experienceYears;

public:
    WorkingDog(const std::string& n, int a, const std::string& b, const std::string& job);
    virtual ~WorkingDog();
    
    void performJob();
    void receiveCommand(const std::string& command);
    
    std::string getJobType() const;
    int getExperienceYears() const;
};

// 宠物类：可以包装任何动物
template<typename T>
class Pet {
private:
    T* animal;
    std::string ownerName;
    bool isVaccinated;

public:
    Pet(T* a, const std::string& owner);
    virtual ~Pet();
    
    T* getAnimal();
    std::string getOwnerName() const;
    bool getIsVaccinated() const;
    void vaccinate();
};

int main() {
    // 示例使用
    Dog myDog("Buddy", 3, "Golden Retriever");
    Cat myCat("Whiskers", 2);
    Eagle wildEagle("Swift", 5, 10.0);
    
    std::vector<Animal*> animals;
    animals.push_back(&myDog);
    animals.push_back(&myCat);
    animals.push_back(&wildEagle);
    
    for (Animal* animal : animals) {
        std::cout << animal->getName() << " says: ";
        animal->speak();
        std::cout << std::endl;
    }
    
    return 0;
}
