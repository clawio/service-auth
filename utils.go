package main

import (
	"code.google.com/p/go-uuid/uuid"
	"golang.org/x/net/context"
	metadata "google.golang.org/grpc/metadata"
)

func newGRPCTraceContext(ctx context.Context, trace string) context.Context {
	md := metadata.Pairs("trace", trace)
	ctx = metadata.NewContext(ctx, md)
	return ctx

}

func getGRPCTraceID(ctx context.Context) string {

	md, ok := metadata.FromContext(ctx)
	if !ok {
		return uuid.New()

	}

	tokens := md["trace"]
	if len(tokens) == 0 {
		return uuid.New()

	}

	if tokens[0] != "" {
		return tokens[0]

	}

	return uuid.New()

}
