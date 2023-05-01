package main

import (
	"BookStore/pb"
	"context"
	"testing"
)

func TestServer_ListBooks(t *testing.T) {
	db, _ := NewDB("test.db")
	s := server{bs: &bookstore{db: db}}
	req := &pb.ListBooksRequest{
		Shelf: 3,
	}
	res, err := s.ListBooks(context.Background(), req)
	if err != nil {
		t.Fatalf(" s.ListBooks failed, err:%v\n", err)
	}
	t.Logf("next_page_token:%v\n", res.GetNextPage())
	for i, book := range res.Books {
		t.Logf("%d:%#v\n", i, book)
	}
}
