# HTTP/2, Websocket

# HTTP/2

- HTTP/1: Mỗi 1 request-response trên 1 connection, mỗi khi đóng mở connection (TCP) lại phải thực hiện “bắt tay 3 bước”
- HTTP/1.1: Một kết nối có thể được dùng cho nhiều request-response (keep-alive). Có cơ chế pipeline
- HTTP/2: Một kết nối TCP được chia thành nhiều luồng data (stream data)

Đặc điểm http/2:

### Multiplexed

Giống pipeline của http/1.1 (có thể gửi nhiều request trước khi nhận response) nhưng response nào tới trước sẽ nhận trước, không phải đợi theo thứ tự gửi request

### Prioritization

HTTP/2 cho phép client cung cấp mức độ ưu tiên cho các luồng dữ liệu cụ thể (stream). Mặc phía server không bị ràng buộc phải tuân theo ưu tiên này, nhưng cơ chế này cho phép server tối ưu hóa việc phân bổ tài nguyên dựa trên các yêu cầu của người dùng cuối.

VD: thẻ `<script>` trong `<head>` của trang sẽ được tải trong Chrome ở mức độ ưu tiên High (thấp hơn CSS - Highest)

### Dữ liệu dạng binary

### Nén Header

### Server push

Trong HTTP/2 server có thể gởi về nhiều response với chỉ một request từ client. Cơ chế này gọi là Server Push, giúp trình duyệt tiết kiệm được các requests không cần thiết.

VD: client yêu cầu 1 page .html, server có thể gửi file .css, .js

# Websocket

## polling

Khi chưa có websocket, để thực hiện realtime như lấy dữ liệu mới nhất từ server → luôn luôn gửi request (sau vài giây) tới server xem có gì mới không

```jsx
setInterval(function () {
  send_request();
}, 1000);
```

## long-polling

1 request sẽ được giữ cho tới khi server trả response giúp gửi ít request tới server hơn cách trên

## websocket

Sử dụng HTTP để giao tiếp nhưng khác là truyền 2 chiều, server cũng có thể gửi cho client, dùng trên 1 kết nối TCP

Dùng cho những ứng dụng cần realtime

# gRPC

Ý tưởng về việc làm thế nào để các service giao tiếp với nhau với **tốc độ cao nhất**, **giảm tải encode/decode data** chính là lý do thúc đẩy **gRPC** ra đời.

gRPC hoạt động trên **HTTP/2**, **HTTP/2** sẽ hoạt động rất tốt với **binary** thay vì là text. Vì thế Google phát minh kiểu dữ liệu binary mới với tên gọi: **Protobuf** (tên đầy đủ là **Protocol Buffers**).

**gRPC** nên dùng để giao tiếp **backend to backend**. CPU không gánh nhiều cost cho **encode**/**decoding** mỗi đầu nữa. Tuy nhiên đặc tính mỗi đầu cần import file model chung (gen từ **protobuf** file) nên nếu update thì phải update đủ. Việc này vô tình tạo dependency cho các bên sử dụng

4 Loại API:

- Unary
- Client streaming
- Server streaming
- Bi-direction streaming

## Protobuf file

- Định nghĩa message (data, request, response), service (tên, endpoint) bằng Protocol Buffer
- gRPC sẽ gen code, chỉ cần viết phần xử lý

## Demo

Tham khảo: [https://www.youtube.com/playlist?list=PLC4c48H3oDRzLAn-YsHzY306qhuEvjhmh](https://www.youtube.com/playlist?list=PLC4c48H3oDRzLAn-YsHzY306qhuEvjhmh)

---

Download Protocol buffers:

[https://github.com/protocolbuffers/protobuf/releases/tag/v23.4](https://github.com/protocolbuffers/protobuf/releases/tag/v23.4)

Install protocol compiler

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

Install grpc package

```bash
go get -u google.golang.org/grpc
```

Tạo file proto

```protobuf
syntax = "proto3";

// package: fullpath
option go_package = "grpc-example-2/pb";

message SumRequest{
  int32 num1 = 1;
  int32 num2 = 2;
}

message SumResponse{
  int32 result = 1;
}

service Calculator{
  rpc Create(SumRequest) returns (SumResponse);
}
```

Gen code

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative .\calculator.pr
oto
```

! Khi service rỗng chạy lệnh protoc sẽ không gen được code https://github.com/golang/protobuf/issues/1334

### Unary API

> Client gửi 1 request và Server trả về 1 response

Tạo file server, tạo struct implement các method được định nghĩa trong proto

```go
package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpc-example-2/calculator/calculatorpb"
	"log"
	"net"
)

// server implements CalculatorServer ...
type server struct {
	pb.UnimplementedCalculatorServer
}

func (s *server) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	fmt.Println("call: Sum()")
	resp := &pb.SumResponse{
		Result: req.Num1 + req.Num2,
	}

	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:8089")
	if err != nil {
		log.Fatalf("create listen err %s", err)
	}

	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})
	fmt.Println("Server is running ...")

	// Serve accepts incoming connections on the listener lis
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("err while serve %s", err)
	}
}
```

Tạo file client

```go
package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpc-example-2/calculator/calculatorpb"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:8089", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("err while dial %v", err)
	}

	defer conn.Close()

	client := pb.NewCalculatorClient(conn)

	callSum(client)
}

func callSum(c pb.CalculatorClient) {
	resp, err := c.Sum(context.Background(), &pb.SumRequest{
		Num1: 5,
		Num2: 3,
	})

	if err != nil {
		log.Fatalf("call Sum err: %s", err)
	}

	fmt.Printf("Result: %d", resp.Result)
}
```

### Server Streaming API

> Client gửi 1 request và nhận nhiều response từ Server. Đây là loại API mới nhờ HTTP/2 (Server push)

Sửa file proto

```protobuf
message PNDRequest{
  int32 number = 1;
}

message PNDResponse{
  int32 result = 1;
}

service Calculator{
  rpc Sum(SumRequest) returns (SumResponse);
  rpc PrimeNumberDecomposition(PNDRequest) returns (stream PNDResponse);
}
```

Implement (server)

```go
func (s *server) PrimeNumberDecomposition(req *pb.PNDRequest, stream pb.Calculator_PrimeNumberDecompositionServer) error {
	fmt.Println("call: PrimeNumberDecomposition()")

	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k
			stream.Send(&pb.PNDResponse{
				Result: k,
			})
		} else {
			k++
			log.Printf("k increase to %v\n", k)
		}
	}

	return nil
}
```

Gọi hàm (client)

```go
func callPrimeNumberDecomposition(c pb.CalculatorClient) {
	stream, err := c.PrimeNumberDecomposition(context.Background(), &pb.PNDRequest{
		Number: 120,
	})

	if err != nil {
		log.Fatalf("err when decomposition: %s", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("server finish streaming")
			return
		}

		fmt.Printf("Result: %d\n", resp.Result)
	}
}
```

### Client Streaming API

> Client gửi nhiều request và nhận 1 response từ Server. Đây là loại API mới nhờ HTTP/2

File Proto

```protobuf
message AverageRequest{
  float num = 1;
}

message AverageResponse{
  float result = 1;
}

service Calculator{
  rpc Sum(SumRequest) returns (SumResponse);
  rpc PrimeNumberDecomposition(PNDRequest) returns (stream PNDResponse);
  rpc Average(stream AverageRequest) returns (AverageResponse);
}
```

Implement (server)

```go
func (s *server) Average(stream pb.Calculator_AverageServer) error {
	fmt.Println("call: Average()")
	var avg float32
	var count int

	for {
		req, err := stream.Recv()

		avg += req.Num
		count++

		if err == io.EOF {
			return stream.SendAndClose(&pb.AverageResponse{
				Result: avg / float32(count),
			})
		}
		if err != nil {
			log.Fatalf("err while recv average %v", err)
		}
	}
}
```

Gọi hàm (client)

```go
func callAverage(c pb.CalculatorClient) {
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("err when call average: %s", err)
	}

	err = stream.Send(&pb.AverageRequest{
		Num: 5,
	})
	if err != nil {
		log.Fatalf("err when send request: %s", err)
	}

	err = stream.Send(&pb.AverageRequest{
		Num: 2,
	})
	if err != nil {
		log.Fatalf("err when send request: %s", err)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("err when receive: %s", err)
	}

	log.Printf("average: %v", resp)
}
```

### Bi-direction Streaming API

> Đây là loại API mới nhờ HTTP/2 cho phép client gửi nhiều request và nhận nhiều response từ server

- Khi cả client và server cần gửi nhiều data cho nhau (async)
- Long running connection

File proto

```protobuf
message FindMaxRequest{
  int32 num = 1;
}

message FindMaxResponse{
  int32 result = 32;
}

service Calculator{
  rpc Sum(SumRequest) returns (SumResponse);
  rpc PrimeNumberDecomposition(PNDRequest) returns (stream PNDResponse);
  rpc Average(stream AverageRequest) returns (AverageResponse);
  rpc FindMax(stream FindMaxRequest) returns (stream FindMaxResponse);
}
```

Implement (server)

```go
func (s *server) FindMax(stream pb.Calculator_FindMaxServer) error {
	fmt.Println("call: FindMax()")
	max := int32(0)

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			log.Println("EOF ...")
			return nil
		}
		if err != nil {
			log.Fatalf("err while recv FindMax %v", err)
			return err
		}

		num := req.Num
		if num > max {
			max = num
		}
		err = stream.Send(&pb.FindMaxResponse{
			Result: max,
		})
		if err != nil {
			log.Fatalf("err while send %v", err)
			return err
		}
	}
}
```

Gọi hàm (client)

```go
func callFindMax(c pb.CalculatorClient) {
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("err when call findmax: %s", err)
	}

	listReq := []int32{
		7, 8, 4, -1, 10, 6, 5,
	}

	waitc := make(chan struct{})

	// go routine using to send requests
	go func() {
		for _, num := range listReq {
			err := stream.Send(&pb.FindMaxRequest{
				Num: num,
			})
			if err != nil {
				log.Fatalf("err when send request: %s", err)
			}
			time.Sleep(1 * time.Second)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("err when close: %s", err)
		}
	}()

	// go routine using to receive responses
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("EOF ...")
				break
			}

			if err != nil {
				log.Fatalf("err when recv: %s", err)
			}

			log.Printf("Result: %d\n", resp.Result)
		}
		close(waitc)
	}()

	_ = <-waitc
}
```
