package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sort"
	"time"

	"github.com/KangkanaBaruah/calci-Assignment/calcipb"
	"google.golang.org/grpc"
)

type server struct {
	calcipb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calcipb.SumRequest) (resp *calcipb.SumResponse, err error) {
	fmt.Println("Sum function is invoked")

	num1 := req.FirstNum
	num2 := req.SecondNum

	res := num1 + num2
	resp = &calcipb.SumResponse{
		Result: res,
	}
	return resp, nil
}

func (*server) Prime(req *calcipb.PrimeRequest, resp calcipb.CalculatorService_PrimeServer) error {
	fmt.Println("Prime function is invoked")
	num := req.Num

	isprime := func(n int) bool {
		if n <= 1 {
			return false
		}

		for i := 2; i < n; i++ {
			if n%i == 0 {
				return false
			}
		}

		return true
	}

	for i := 2; i <= int(num); i++ {
		if isprime(i) {
			res := calcipb.PrimeResponse{
				Result: int32(i),
			}
			time.Sleep(100 * time.Microsecond)
			resp.Send(&res)
		}
	}
	return nil

}

func (*server) Average(stream calcipb.CalculatorService_AverageServer) error {
	fmt.Println("Average function is invoked")
	var sum int32
	var i int32
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calcipb.AverageResponse{
				Result: float32(sum) / float32(i),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream : %v", err)
		}
		sum += msg.Num
		i++

	}

}

func (*server) Maxnumber(stream calcipb.CalculatorService_MaxnumberServer) error {
	fmt.Println("Maxnumber function is invoked")
	sl := []int{}
	maxNow := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("error while receiving data from Maxnumber client : %v", err)
			return err
		}
		num := req.Num
		sl = append(sl, int(num))
		sort.Ints(sl)
		if sl[len(sl)-1] > maxNow {
			maxNow = sl[len(sl)-1]
			sendErr := stream.Send(&calcipb.MaxnumberResponse{
				Result: int32(maxNow),
			})
			if sendErr != nil {
				log.Fatalf("error while sending response to Maxnumber Client : %v", err)
				return err
			}
		}

	}
}

func main() {
	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	calcipb.RegisterCalculatorServiceServer(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
