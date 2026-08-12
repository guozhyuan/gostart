文档
https://golang.dbwu.tech/performance/slice_pre_alloc/


#### 方法 
    值接收者和指针接收者的区别
    1. 值接收者：函数内部修改的是副本，不会影响到调用者传入的参数。
    2. 指针接收者：函数内部修改的是传入的参数，会直接影响到调用者传入的参数。
    
    // type Human interface {
    //     ChangeName(s string)
    // }

    type Person struct {
        Name string
    }

    func (p *Person) ChangeName(s string) {
        p.Name = s
    }

    // func (p Person) ChangeName(s string) {
    // 	p.Name = s
    // }

    func Func() {
        p1 := Person{"初始名"}
        p1.ChangeName("P1")
        fmt.Println(p1.Name)

        p2 := &Person{"初始名"}
        p2.ChangeName("P2")
        fmt.Println(p2.Name)
    }

    // 输出 P1 P2


    如果把 func (p *Person) ChangeName(s string)
    改成 func  (p Person) ChangeName(s string)
    同样输出 P1 P1
    不同点在于(p *Person)会改变p的值，而(p Person)不会改变p的值，只是改变p的副本的值。
    
```