package method

import (
	"context"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"log"
)

type TRemote struct {
	Remote string
	Port   int
}

func GetMethod(remotes []TRemote) {
	ctx := context.Background()
	cc, err := grpc.Dial("192.168.1.94:8061", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial failed, err: %v", err)
	}
	client := grpcreflect.NewClient(ctx, grpc_reflection_v1alpha.NewServerReflectionClient(cc))
	res, err := client.ListServices()
	log.Println("list services: ", res, "\n", err)
	r, err := client.ResolveService(res[0])
	log.Println("resolve service: ", r, "\n", err)
	log.Println(r.GetMethods()[2].GetName(), "\n", r.GetMethods()[2].GetInputType().GetFields()[0].AsFieldDescriptorProto().GetTypeName(), "\n", r.GetMethods()[2].GetInputType().GetFields())
	log.Println(r.GetMethods()[2].GetInputType().GetFields()[0].AsFieldDescriptorProto().GetTypeName()[1:])
	log.Println(client.FileContainingSymbol(r.GetMethods()[2].GetInputType().GetFields()[0].AsFieldDescriptorProto().GetTypeName()[1:]))
	log.Println(r.GetFile().FindEnum(r.GetMethods()[2].GetInputType().GetFields()[0].GetName()))
	re, err := client.ResolveEnum(r.GetFile().GetName())
	log.Println("resolve service: ", re, "\n", err)
}



// AllMethodsForServices returns a slice that contains the method descriptors
// for all methods in the given services.
func AllMethodsForServices(descs []*desc.ServiceDescriptor) []*desc.MethodDescriptor {
	seen := map[string]struct{}{}
	var allMethods []*desc.MethodDescriptor
	for _, sd := range descs {
		if _, ok := seen[sd.GetFullyQualifiedName()]; ok {
			// duplicate
			continue
		}
		seen[sd.GetFullyQualifiedName()] = struct{}{}
		allMethods = append(allMethods, sd.GetMethods()...)
	}
	return allMethods
}

// AllMethodsForServer returns a slice that contains the method descriptors for
// all methods exposed by the given gRPC server.
func AllMethodsForServer(svr *grpc.Server) ([]*desc.MethodDescriptor, error) {
	svcs, err := grpcreflect.LoadServiceDescriptors(svr)
	if err != nil {
		return nil, err
	}
	var descs []*desc.ServiceDescriptor
	for _, sd := range svcs {
		descs = append(descs, sd)
	}
	return AllMethodsForServices(descs), nil
}

// AllMethodsViaReflection returns a slice that contains the method descriptors
// for all methods exposed by the server on the other end of the given
// connection. This returns an error if the server does not support service
// reflection. (See "google.golang.org/grpc/reflection" for more on service
// reflection.)
// This automatically skips the reflection service, since it is assumed this is not
// a desired inclusion.
func AllMethodsViaReflection(ctx context.Context, cc grpc.ClientConnInterface) ([]*desc.MethodDescriptor, error) {
	stub := grpc_reflection_v1alpha.NewServerReflectionClient(cc)
	cli := grpcreflect.NewClient(ctx, stub)
	svcNames, err := cli.ListServices()
	if err != nil {
		return nil, err
	}
	var descs []*desc.ServiceDescriptor
	for _, svcName := range svcNames {
		sd, err := cli.ResolveService(svcName)
		if err != nil {
			return nil, err
		}
		if sd.GetFullyQualifiedName() == "grpc.reflection.v1alpha.ServerReflection" {
			continue // skip reflection service
		}
		descs = append(descs, sd)
	}
	return AllMethodsForServices(descs), nil
}

