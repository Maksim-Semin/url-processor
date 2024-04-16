package tests

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "main/pkg/api/proto"
	"os"
	"reflect"
	"strings"
	"testing"
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
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewApiClient(conn)
	sites := readURLsFromFile()

	for i, site := range sites {
		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			response, err := client.ChangeURL(context.Background(), &pb.URLRequest{URL: site})
			if err != nil {
				t.Fatalf("error calling ChangeURL RPC: %v", err)
			}

			code := response.URL

			newResponse, err := client.GetSourceURL(context.Background(), &pb.URLRequest{URL: code})
			if err != nil {
				t.Fatalf("error calling GetSourceURL RPC: %v", err)
			}

			if !reflect.DeepEqual(site, newResponse.URL) {
				t.Errorf("testcase %d failed: expected [%v] actual [%v]", i, site, newResponse.URL)
			}
		})
	}
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

	for _, site := range links {
		response, err := client.ChangeURL(context.Background(), &pb.URLRequest{URL: site})
		if err != nil {
			t.Errorf("error calling ChangeURL RPC: %v", err)
			return
		}
		code := response.URL
		data[code] = site
		shortLinks = append(shortLinks, code)
	}

	for _, shortLink := range shortLinks {
		resp, err := client.GetSourceURL(context.Background(), &pb.URLRequest{URL: shortLink})
		if err != nil {
			t.Errorf("error calling GetSourceURL RPC: %v", err)
			return
		}
		site := resp.URL
		expected, _ := data[shortLink]
		if !reflect.DeepEqual(expected, site) {
			t.Errorf("failed: expected [%v] actual [%v]", expected, site)
		}
	}
}

func TestGRPCRequestsEmptyValue(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewApiClient(conn)

	for i := 0; i < 5; i++ {
		response, err := client.ChangeURL(context.Background(), &pb.URLRequest{
			URL: "",
		})
		if err != nil {
			t.Errorf("error calling RPC: %v", err)
			return
		}
		if response.URL == "" && response.Error == "field URL should not be empty" {
			t.Logf("test %d on empty value success", i+1)
		} else {
			t.Errorf("test %d on empty value failed", i+1)
		}
	}
}
