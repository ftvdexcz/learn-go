# I.Interface

Interface là 1 tập các function signature (k có cài đặt). Các struct `implement` interface bằng cách ‘cài đặt’ function trong interface đó. Khác với các ngôn ngữ khác, Go implement interface ngầm định (không thông qua keyword `implements` )

Nhiệm vụ chính của interface chỉ là cung cấp function signature: name, parameters, return type 

```go
type Person interface{
	greet(name string) string; 
}

type Student struct{
	name string
}

type Employee struct{
	name string 
}

func (s *Student) greet(name string) string{
	return name + " is a student"
}

func (e *Employee) greet(name string) string{
	return name + " is a employee"
}

func main(){
	var p Person 
	p = &Student{
		name: "Long",
	}

	fmt.Println(p.greet("Long"));

	p = &Employee{
		name: "Nam",
	}

	fmt.Println(p.greet("Nam"));
}
```

Go cho phép implement nhiều interface 

```go
type Animal interface {
	speak()
}

type Dog struct {}

// implements mà không cần khai báo tường minh như các ngôn ngữ khác
func (d Dog) speak() {
	fmt.Println(“woaww woaww”)	
}

type Movement interface {
	move()
}

func (d Dog) move() {
	fmt.Println(“Dog chạy bằng 4 chân”)
}

func main() {
	dog := Dog{}

	var m Movement = dog
	m.move()		//	Dog chạy bằng 4 chân

	var a animal = dog
	a.speak()		//	woaww woaww
}
```

## Interface extends Interface

Interface kế thừa interface khác. Trong Java có thể `extends` interface, Go không cho phép nhưng có thể đạt được bằng embbed interface. Một đối tượng muốn implement interface này cần implement đầy đủ các phương thức có trong các interface con

```go
type NextAnimal interface {
	Movement
	Animal
}

func main() {
	dog := Dog{} // Dog implement 2 phương thức speak() và move() có trong 2 interface con
	var na NextAnimal = dog 

	na.speak()		//	Dog chạy bằng 4 chân
	na.move()		// 	woaww woaww
}
```

## Empty interface

Tất cả type đều implement empty interface 

```go
type Body struct {
    Msg interface{}
}

func main() {
    b := Body{"Hello there"}
    fmt.Printf("%#v %T\n", b.Msg, b.Msg) // "Hello there" string

    b.Msg = 5
    fmt.Printf("%#v %T\n", b.Msg, b.Msg) // 5 int 
}
```

Có thể dùng empty interface làm tham số truyền vào hàm để hàm chấp nhận mọi kiểu dữ liệu 

```go
func GoOut(i interface{}){
	fmt.Println(i)
}

type data struct {
	val int
}

func main() {
	GoOut(10)		// 	10
	GoOut(12.5)		// 	12.5
	
	d := data{15}
	GoOut(d)		// 	{15}
}
```

## Type Assertion

Dùng để kiểm tra kiểu của biến trong run-time: `t := i.(T)`

Type assertion dùng với interface 

```go
greetingStr, ok := greeting.(string)
```

[https://www.sohamkamani.com/golang/type-assertions-vs-type-conversions/](https://www.sohamkamani.com/golang/type-assertions-vs-type-conversions/)

As mentioned before, different types have different restrictions and methods defined on them, even though their data structure may be the same. When you convert from one type to another, you are *changing* what you can do with the type, rather than just *exposing* its underlying type, as is done in type assertions.

Type conversions also give you compilation errors if you try to convert to the wrong type, as opposed to runtime errors and optional `ok` return values that type assertions give.

## Kiểm tra value implements interface

```go
var _ Somether = (*MyType)(nil)
```

# II. Goroutines & Channel

Concurrency: Điều phối nhiều task chạy cùng lúc nhưng tại 1 thời điểm chỉ có 1 task được chạy 

Parallelism: Chạy nhiều task cùng 1 thời điểm (cpu 2 core trở lên) 

Process vs Thread: 1 Process đại diện 1 ứng dụng đang chạy, 1 Process chứa nhiều thread (các thread có bộ nhớ riêng) 

→ Go đưa ra khái niệm `goroutine` thay cho thread, giao tiếp giữa các goroutine bằng `channel` 

- Thread được quản lý bởi phần cứng, goroutine quản lý bởi go run time
- stack size của thread là 1MB fix cứng, goroutine là 8KB và max là 1GB
- Giao tiếp thread khó hơn goroutine (channel)
- Thread có định danh (TID), goroutine không có

Goroutine có thể chạy parallelism hoặc không, trong 1 chương trình go có ít nhất 1 goroutine là main function

```go
func main() {
    hello("Martin")
    hello("Lucia")
    hello("Michal")
    hello("Jozef")
    hello("Peter")
}

func hello(name string) {
    fmt.Printf("Hello %s!\n", name)
}
```

```go
func main() {
    go hello("Martin")
    go hello("Lucia")
    go hello("Michal")
    go hello("Jozef")
    go hello("Peter")
    fmt.Scanln() // đợi input từ user 
}

func hello(name string) {
    fmt.Printf("Hello %s!\n", name)
}
```

## Go sync.WaitGroup - synchronize

`sync.WaitGroup` đợi các goroutine hoàn thành 

```go
func main() {
    var wg sync.WaitGroup
    wg.Add(2) // đợi 2 goroutine hoàn thành 

    go func() {
        count("oranges")
        wg.Done() // gọi Done để thông báo hoàn thành 
    }()

    go func() {
        count("apples")
        wg.Done()
    }()

    wg.Wait() // đợi đến khi các goroutine hoàn thành 
}

func count(thing string) {
    for i := 0; i < 4; i++ {
        fmt.Printf("counting %s\n", thing)
        time.Sleep(time.Millisecond * 500)
    }
}
```

## Goroutine channels

Tham khảo: 

- [https://200lab.io/blog/golang-channel-la-gi/](https://200lab.io/blog/golang-channel-la-gi/)
- [Master Go Programming With These Concurrency Patterns (in 40 minutes) - YouTube](https://www.youtube.com/watch?v=qyM8Pi1KiiM)

Goroutines giao tiếp thông qua channels, cho phép gửi và nhận bằng toán tử `<-`

```go
c := make(chan string)
```

A new channel is created with the `make` function.

```go
c <- v    // send
v := <-c  // receive
```

### Lưu ý:

Toán tử  gửi và nhận dữ liệu channel `<-` block các câu lệnh phía sau cho tới khi data đó được đọc hoặc gửi 

```go
func main() {
	c := make(chan int)
	fmt.Println("Send")
	c <- 1 
	fmt.Println(<-c)
	fmt.Println("End")
}
```

Ở ví dụ trên `c <- 1` block goroutine main cho tới khi có goroutine khác đọc số 1 nhưng không có ai đọc nên bị deadlock 

```go
func main() {
	c := make(chan int)

	go func(){
		fmt.Println(<- c)
	}()
	fmt.Println("Send")
	c <- 1
	time.Sleep(2000 * time.Millisecond)
	fmt.Println("End")
}
```

Không bị block nữa do có goroutine anonymous đã đọc và in ra giá trị, nếu bỏ dòng Sleep có thể không thấy 1 được in ra do không bị block nên nó đến dòng in ra End và kết thúc main 

Cách làm trên chưa tốt vì không biết nên phải Sleep bao lâu, nên tạo 1 channel khác 

```go
func DemoChannel2() {
    pipe := make(chan string)
    done := make(chan bool)//kênh báo khi receiver nhận đủ
		go func() {
        for {
            receiver, more := <-pipe//khi không còn dữ liệu, more sẽ false
            fmt.Println(receiver)
            if !more {
                done <- true//không còn dữ liệu trong pipe, đã nhận đủ !
								return//thoát khỏi go routine
            }
        }
    }()

    pipe <- "water 1"
    pipe <- "water 2"
    pipe <- "water 3"
    pipe <- "water 4"
    close(pipe)
    <-done//Chờ khi nhận được tin báo từ receiver thì thoát
}
```

## Buffered Channel

```go
ch := make(chan int, 1) // create buffered channel size = 1 
```

Buffered channel như là nơi lưu tạm data được ghi vào cho tới khi đầy bộ đệm và không bị block câu lệnh để đợi data được đọc 

```go
func main() {
	c := make(chan int, 1)
	fmt.Println("Send")
	c <- 1
	fmt.Println(<-c)
	fmt.Println("End")
}
```

Lưu ý nếu size bộ đệm nhỏ hơn kích thước data gửi đến sẽ gây deadlock 

### Đọc từ buffered channel

```go
func main() {
	c := make(chan int, 10)
	fmt.Println("Send")
	c <- 1
	c <- 1
	c <- 2
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println("End")
}
```

Đoạn code trên chạy bình thường do có 3 data được gửi và đọc ra 3 lần, nhưng nếu có nhiều hơn và không biết được gửi bao nhiêu lần thì đọc ra kiểu gì → dùng range 

```go
func main() {
	c := make(chan int, 10)
	fmt.Println("Send")
	c <- 1
	c <- 1
	c <- 2
	close(c)
	for d := range c{
		fmt.Println(d)
	}
	fmt.Println("End")
}
```

Lưu ý đoạn `close(c)` nếu không có sẽ gây deadlock do ở vòng for range nó ngồi đợi data và block luồng 

## Select

Select block cho tới khi 1 trong các case được thực thi, nếu có nhiều case ready, nó sẽ chọn random 

```go
func main() {
	c1 := make(chan int)
	c2 := make(chan int)

	go func(){
		time.Sleep(time.Second * 2)
		c1 <- 1
	}()

	go func(){
		time.Sleep(time.Second * 3)
		c2 <- 2
	}()
	
	fmt.Println(<-c2)
	fmt.Println(<-c1)
}
```

```go
func main() {
	c1 := make(chan int)
	c2 := make(chan int)

	go func() {
		time.Sleep(time.Second * 2)
		c1 <- 1
	}()

	go func() {
		time.Sleep(time.Second * 1)
		c2 <- 2
	}()

	select {
		case x1 := <-c1:
			fmt.Println("x1")
			fmt.Println(x1)
		case x2 := <-c2:
			fmt.Println("x2")
			fmt.Println(x2)
	}

	fmt.Println("End")
}
```

Ở đoạn code trên khi dùng select, nó đợi cho tới khi 1 case được chạy nó sẽ kết thúc 

Để đọc được hết data mà không kết thúc khi vào 1 case (như ở trên muốn đọc cả x1) thì dùng `for` 

```go
func main() {
	c1 := make(chan int)
	c2 := make(chan int)
	q := make(chan int)
	

	go func() {
		time.Sleep(time.Second * 2)
		c1 <- 1
		q <- 0
	}()

	go func() {
		time.Sleep(time.Second * 1)
		c2 <- 2
	}()

	loop := true
	for loop{
		select {
		case x1 := <-c1:
			fmt.Println("x1")
			fmt.Println(x1)
		case x2 := <-c2:
			fmt.Println("x2")
			fmt.Println(x2)
		case <- q:
			loop = false
		}
	}

	fmt.Println("End")
}
```

## **Default Selection**

The `default` case in a `select` is run if no other case is ready.

Use a `default` case to try a send or receive without blocking:

```go
select {
	case i := <-c:
    // use i
	default:
    // receiving from c would block
}
```

## sync.Mutex

Data Racing - Khi 2 hoặc nhiều goroutine cùng đọc và ghi vào 1 biến → giá trị của biến không như mong muốn 

```go
for i := 0; i < 5; i++ {
	go func(){
		for j := 1; j <= 10000; j++{
			count += 1
		}
	}
}
```

Sử dụng `sync.RWMutext` để “chiếm lock”, các goroutine khác không được sử dụng biến này 

- RWLock: block tất cả goroutine dù đang là read hay write
- RLock: block tất cả goroutine write cho phép các goroutine read được phép truy xuất

## Truyền channel vào tham số của hàm

```go
func(ch chan <- string){ ... } // channel này chỉ cho phép ghi dữ liệu 

func(ch <- chan string){ ... } // channel này chỉ cho phép đọc dữ liệu 

func(ch chan <- string){ ... } // channel này cho phép cả đọc cả ghi  
```

## “Done” channel

Để gửi tín hiệu dừng goroutine từ goroutine khác (chẳng hạn main) sử dụng “done channel pattern” 

```go
func printSomething(done <-chan int){
	for {
		select{
			case <- done:
				return 
			default:
				fmt.Println("work")
		}
	}
}

func main(){
	done := make(chan int)

	go printSomething(done)

	time.Sleep(time.Second * 3)

	close(done)
}
```

# III.Context

Context dùng để truyền data giữa các function trong 1 request, hoặc cancel 1 chuỗi các function khi timeout 

## Tạo context và truyền data với context

Sử dụng context cho phép thêm data vào context, mỗi layer có thể thêm thông tin vào context và lấy ra data trong context

```go
func doSomething(ctx context.Context){
	fmt.Printf("mykey: %s\n", ctx.Value("mykey"))
}

func main(){
	ctx := context.Background()

	ctx = context.WithValue(ctx, "mykey", "value1")
	doSomething(ctx)
}
```

## Hủy context

Context cho phép gửi 1 signal tới các function cho biết context bị hủy. Ví dụ giả sử client đã thoát thì k cần xử lý request để gửi về nữa hoặc khi quá thời gian thực thi cho request 

Gọi hàm `cancel()` và sử dụng select để đọc từ `ctx.Done()`

```go
func doSomething(ctx context.Context) {
	ctx, cancelCtx := context.WithCancel(ctx)
	
	printCh := make(chan int)
	go doAnother(ctx, printCh)

	cancelCtx()

	for num := 1; num <= 3; num++ {
		printCh <- num
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("doSomething: finished\n")

	
}

func doAnother(ctx context.Context, printCh <-chan int) {
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				fmt.Printf("doAnother err: %s\n", err)
			}
			fmt.Printf("doAnother: finished\n")
			return
		case num := <-printCh:
			fmt.Printf("doAnother: %d\n", num)
		}
	}
}
```

## Context with deadline

Context với deadline sẽ tự gọi `cancel()` khi hết thời gian 

```go
func main() {
   start := time.Now()

   ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*10)
   defer cancel()

   s := calc(ctx, 20)
   fmt.Println(s)
   fmt.Println("Took ", time.Since(start))
}

// return -1 = cancel context
func calc(ctx context.Context, n int) int {
   ch := make(chan int)

   s := 0
   var a int

   go fibo(n, ch)

   for {
      select {
      case <-ctx.Done():
         fmt.Println("Time out")

         return s
      case a = <-ch:
         s += a
      }
   }

}

func fibo(n int, ch chan int) int {
   if n == 0 || n == 1 {
      ch <- n
      return n
   }
   return fibo(n-1, ch) + fibo(n-2, ch)
}
```

# IV.Sqlx
SQLx là một thư viện hỗ trợ truy vấn cơ sở dữ liệu trong Golang, nó cung cấp các công cụ để giúp đơn giản hóa việc tương tác với cơ sở dữ liệu. SQLx cho phép thực hiện các truy vấn SQL trong Golang một cách dễ dàng hơn thông qua việc cung cấp một loạt các hàm hỗ trợ như `Get`, `Select`, `NamedQuery`, `Exec`, vv.

SQLx cũng hỗ trợ các tính năng như thực hiện truy vấn trong giao dịch và caching kết quả truy vấn. Nó cung cấp một cách tiếp cận tương đối trực tiếp đến cơ sở dữ liệu và cho phép người dùng thực hiện các truy vấn phức tạp hơn.

## So sánh với GORM

GORM là một thư viện ORM (Object-Relational Mapping) trong Golang, nó cung cấp một cách tiếp cận trừu tượng hóa đến cơ sở dữ liệu. Với GORM, người dùng không cần phải viết mã SQL trực tiếp, mà thay vào đó có thể sử dụng các phương thức truy vấn được hỗ trợ sẵn như `Create`, `Find`, `Update`, vv.

Một trong những ưu điểm của GORM là nó cung cấp các tính năng như tự động tạo bảng, quản lý quan hệ giữa các bảng và hỗ trợ các kiểu dữ liệu phức tạp như json và array. Nó cũng hỗ trợ nhiều cơ sở dữ liệu như MySQL, PostgreSQL và SQLite.

Với SQLx, người dùng có thể viết mã SQL trực tiếp để tối ưu hiệu suất và sử dụng các tính năng như caching để giảm thiểu tốn kém cho truy vấn cơ sở dữ liệu.

---

**`sqlx`** is a package for Go which provides a set of extensions on top of the excellent built in **`database/sql`** package.

### Types

- `sqlx.DB` tương tự `sql.DB`
- `sqlx.Tx` tương tự `sql.Tx`
- `sqlx.Stmt` tương tự `sql.Stmt`
- `sqlx.NamedStmt`

Các kiểu trên nhúng (embed) kiểu trong `database/sql` tương đương ⇒ gọi `sqlx.DB.Query` giống gọi `sql.DB.Query` 

![Screenshot from 2023-07-06 14-25-49.png](Sqlx%209ed288dc3c77417b937c8b2a02d87699/Screenshot_from_2023-07-06_14-25-49.png)

Ngoài ra có 2 kiểu cursor 

- `sqlx.Rows` tương tự `sql.Rows`, return từ `Queryx`
- `sqlx.Row` tương tự `sql.Row`, return từ `QueryRowx`

---

The handle types in sqlx implement the same basic verbs for querying your database:

- **`Exec(...) (sql.Result, error)`** - unchanged from database/sql
- **`Query(...) (*sql.Rows, error)`** - unchanged from database/sql
- **`QueryRow(...) *sql.Row`** - unchanged from database/sql

These extensions to the built-in verbs:

- **`MustExec() sql.Result`** -- Exec, but panic on error
- **`Queryx(...) (*sqlx.Rows, error)`** - Query, but return an sqlx.Rows
- **`QueryRowx(...) *sqlx.Row`** -- QueryRow, but return an sqlx.Row

And these new semantics:

- **`Get(dest interface{}, ...) error`**
- **`Select(dest interface{}, ...) error`**

### Exec

`Exec` và `MustExec` lấy connection từ connection pool và thực thi query cho trước 

```go
schema := `CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer);`
 
// execute a query on the server
result, err := db.Exec(schema)
 
// or, you can use MustExec, which panics on error
cityState := `INSERT INTO place (country, telcode) VALUES (?, ?)`
countryCity := `INSERT INTO place (country, city, telcode) VALUES (?, ?, ?)`
db.MustExec(cityState, "Hong Kong", 852)
db.MustExec(cityState, "Singapore", 65)
db.MustExec(countryCity, "South Africa", "Johannesburg", 27)
```

### Query

```go
// fetch all places from the db
rows, err := db.Query("SELECT country, city, telcode FROM place")
 
// iterate over each row
for rows.Next() {
    var country string
    // note that city can be NULL, so we use the NullString type
    var city    sql.NullString
    var telcode int
    err = rows.Scan(&country, &city, &telcode)
}
// check the error from rows
err = rows.Err()
```

Query return `sql.Rows` có thể dùng để scan từng row một. Lưu ý nếu như không lặp qua hết row thì phải gọi `rows.Close()` để trả connection về cho pool 

### Queryx

Queryx giống Query nhưng return `sqlx.Rows` cho phép scan nâng cao hơn: 

```go
rows, err := db.Queryx("SELECT * FROM place")
for rows.Next() {
    var p Place
    err = rows.StructScan(&p)
}
```

Các field của struct phải viết hoa chữ đầu (export) để sqlx có thể ghi vào. Sử dụng `db` struct tag để ánh xạ field của struct với cột trong database, hoặc mặc định field sẽ chuyển dưới dạng viết thường để ánh xạ 

### QueryRow

QueryRow lấy ra 1 row, trả về `sql.Row` 

```go
row := db.QueryRow("SELECT * FROM place WHERE telcode=?", 852)
var telcode int
err = row.Scan(&telcode)
```

### Get và Select

**`Get`** and **`Select`** are time saving extensions to the handle types. They combine the execution of a query with flexible scanning semantics. To explain them clearly, we have to talk about what it means to be **`scannable`**:

- a value is scannable if it is not a struct, eg **`string`**, **`int`**
- a value is scannable if it implements **`sql.Scanner`**
- a value is scannable if it is a struct with no exported fields (eg. **`time.Time`**)

**`Get`** and **`Select`** use **`rows.Scan`** on scannable types and **`rows.StructScan`** on non-scannable types. They are roughly analagous to **`QueryRow`** and **`Query`**, where Get is useful for fetching a single result and scanning it, and Select is useful for fetching a slice of results

```go
p := Place{}
pp := []Place{}
 
// this will pull the first place directly into p
err = db.Get(&p, "SELECT * FROM place LIMIT 1")
 
// this will pull places with telcode > 50 into the slice pp
err = db.Select(&pp, "SELECT * FROM place WHERE telcode > ?", 50)
 
// they work with regular types as well
var id int
err = db.Get(&id, "SELECT count(*) FROM place")
 
// fetch at most 10 place names
var names []string
err = db.Select(&names, "SELECT name FROM place LIMIT 10")
```

# V.Fiber 

## Routing

### Handler

Syntax: 

```go
// Function signature
app.Method(path string, ...func(*fiber.Ctx) error)
```

Mỗi 1 route có thể có nhiều handler function được thực thi 

```go
// Simple GET handler
app.Get("/api/list", func(c *fiber.Ctx) error {
  return c.SendString("I'm a GET request!")
})
```

`Use` có thể dùng cho middleware và prefix match. Ví dụ `/john` match `/john/doe`, `/johnnnnn` etc

```go
// Match any request
app.Use(func(c *fiber.Ctx) error {
    return c.Next()
})

// Match request starting with /api
app.Use("/api", func(c *fiber.Ctx) error {
    return c.Next()
})

// Match requests starting with /api or /home (multiple-prefix support)
app.Use([]string{"/api", "/home"}, func(c *fiber.Ctx) error {
    return c.Next()
})

// Attach multiple handlers 
app.Use("/api", func(c *fiber.Ctx) error {
  c.Set("X-Custom-Header", random.String(32))
    return c.Next()
}, func(c *fiber.Ctx) error {
    return c.Next()
})
```

### Parameters

```go
// Parameters
app.Get("/user/:name/books/:title", func(c *fiber.Ctx) error {
    fmt.Fprintf(c, "%s\n", c.Params("name"))
    fmt.Fprintf(c, "%s\n", c.Params("title"))
    return nil
})

// Optional parameter
app.Get("/user/:name?", func(c *fiber.Ctx) error {
    return c.SendString(c.Params("name"))
})
```

### Constraint

[https://docs.gofiber.io/guide/routing#constraints](https://docs.gofiber.io/guide/routing#constraints)

Constraints aren't validation for parameters. If constraint aren't valid for parameter value, Fiber returns **404 handler**.

```go
app.Get("/:test<min(5)>", func(c *fiber.Ctx) error {
  return c.SendString(c.Params("test"))
})

// curl -X GET http://localhost:3000/12
// 12

// curl -X GET http://localhost:3000/1
// Cannot GET /1
```

### Middleware

1 Function được thực thi trong 1 request - response cycle. Gọi hàm `Next()` để thực thi middleware tiếp theo 

```go
app.Use(func(c *fiber.Ctx) error {
  // Set a custom header on all responses:
  c.Set("X-Custom-Header", "Hello, World")

  // Go to next middleware:
  return c.Next()
})

app.Get("/", func(c *fiber.Ctx) error {
  return c.SendString("Hello, World!")
})
```

## Group

Nhóm các route có chung prefix 

```go
func main() {
  app := fiber.New()

  api := app.Group("/api", middleware) // /api

  v1 := api.Group("/v1", middleware)   // /api/v1
  v1.Get("/list", handler)             // /api/v1/list
  v1.Get("/user", handler)             // /api/v1/user

  v2 := api.Group("/v2", middleware)   // /api/v2
  v2.Get("/list", handler)             // /api/v2/list
  v2.Get("/user", handler)             // /api/v2/user

  log.Fatal(app.Listen(":3000"))
}
```

Có thể đặt handler cho group nhưng phải có `Next()` để luồng tiếp tục 

```go
func main() {
    app := fiber.New()

    handler := func(c *fiber.Ctx) error {
        return c.SendStatus(fiber.StatusOK)
    }
    api := app.Group("/api") // /api

    v1 := api.Group("/v1", func(c *fiber.Ctx) error { // middleware for /api/v1
        c.Set("Version", "v1")
        return c.Next()
    })
    v1.Get("/list", handler) // /api/v1/list
    v1.Get("/user", handler) // /api/v1/user

    log.Fatal(app.Listen(":3000"))
}

## Error Handling

### Catching Error

Mặc định fiber sẽ handle lỗi xảy ra trong route handler và middleware. Lỗi phải được return ở trong handler function, nó sẽ được fiber xử lý 

```go
app.Get("/", func(c *fiber.Ctx) error {
    // Pass error to Fiber
    return c.SendFile("file-does-not-exist")
})
```

Panic mặc định không được handle bởi fiber, để fiber handle panic → dùng middleware `recover` 

```go
package main

import (
    "log"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
    app := fiber.New()

    app.Use(recover.New())

    app.Get("/", func(c *fiber.Ctx) error {
        panic("This panic is caught by fiber")
    })

    log.Fatal(app.Listen(":3000"))
}
```

### Fiber custom error

Có thể custom status code và message trả về sử dụng `fiber.NewError()` hoặc custom error của fiber

```go
app.Get("/", func(c *fiber.Ctx) error {
	// 503 Service Unavailable
	return fiber.ErrServiceUnavailable

	// 503 On vacation!
	return fiber.NewError(fiber.StatusServiceUnavailable, "On vacation!")
})
```

### Default error handle

Fiber cung cấp default error handle, với error thông thường (kiểu error) status code trả về 500, error kiểu fiber.Error, status code và message trả về tương ứng 

```go
// Default error handler
var DefaultErrorHandler = func(c *fiber.Ctx, err error) error {
    // Status code defaults to 500
    code := fiber.StatusInternalServerError

    // Retrieve the custom status code if it's a *fiber.Error
    var e *fiber.Error
    if errors.As(err, &e) {
        code = e.Code
    }

    // Set Content-Type: text/plain; charset=utf-8
    c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

    // Return status code with error message
    return c.Status(code).SendString(err.Error())
}
```

### **Custom Error Handler**
A custom error handler can be set using a [Config](https://docs.gofiber.io/api/fiber#config) when initializing a [Fiber instance](https://docs.gofiber.io/api/fiber#new): [https://docs.gofiber.io/guide/error-handling#custom-error-handler](https://docs.gofiber.io/guide/error-handling#custom-error-handler)

# VI.GORM 
https://tame-bison-f5f.notion.site/GORM-05b7e94b3cc24b898a494f7f960e1135
