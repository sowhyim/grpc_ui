package method

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestGetMethod(t *testing.T) {
	GetMethod(nil)
}

func TestAllMethodsViaReflection(t *testing.T) {
	ctx := context.Background()
	cc, err := grpc.Dial("192.168.1.94:8061", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial failed, err: %v", err)
	}
	log.Println(AllMethodsViaReflection(ctx, cc))
}
