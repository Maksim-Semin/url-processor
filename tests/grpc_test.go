package tests

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "main/pkg/api/proto"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func readURLsFromFile() []string {
	var urls []string

	file, err := os.Open("./sites.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		urls = append(urls, url)
	}

	if err = scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return urls
}

func TestGRPCRequests(t *testing.T) {
	var wg sync.WaitGroup
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewApiClient(conn)
	sites := readURLsFromFile()

	for i := 0; i < len(sites); i++ {
		wg.Add(1)
		go func(i int) {
			data := sites[i]
			response, err := client.ChangeURL(context.Background(), &pb.ChangeURLMessage{
				URL: data,
			})
			if err != nil {
				t.Errorf("error calling RPC: %v", err)
				return
			}

			code := response.URL

			newResponse, err := client.GetSourceURL(context.Background(), &pb.ChangeURLMessage{
				URL: code,
			})
			if err != nil {
				t.Errorf("error calling RPC: %v", err)
				return
			}
			if reflect.DeepEqual(data, newResponse.URL) {
				fmt.Printf("testcase %d was succesfull \n", i)
			} else {
				fmt.Printf("testcase %d failed: expected [%v] actual [%v] \n", i, data, newResponse.URL)
			}
			wg.Done()
		}(i)
		time.Sleep(50 * time.Millisecond)

	}
	wg.Wait()
}

func TestGRPCRequests_1(t *testing.T) {
	data := make(map[string]string)

	var shortLinks []string

	links := readURLsFromFile()

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewApiClient(conn)
	for i := 0; i < len(links); i++ {
		site := links[i]
		response, err := client.ChangeURL(context.Background(), &pb.ChangeURLMessage{
			URL: site,
		})
		if err != nil {
			t.Errorf("error calling RPC: %v", err)
			return
		}

		code := response.URL

		data[code] = site
		shortLinks = append(shortLinks, code)
	}
	var count int
	for len(shortLinks) != 0 {

		randInd := rand.Intn(len(shortLinks))

		shortLink := shortLinks[randInd]
		resp, err := client.GetSourceURL(context.Background(), &pb.ChangeURLMessage{
			URL: shortLink,
		})
		if err != nil {
			t.Errorf("error calling RPC: %v", err)
			return
		}
		site := resp.URL

		value, _ := data[shortLink]
		if value == site {
			count++
			fmt.Printf("testcase %d was succesfull \n", count)
			shortLinks = append(shortLinks[:randInd], shortLinks[randInd+1:]...)
		}
		time.Sleep(50 * time.Millisecond)

	}

}
