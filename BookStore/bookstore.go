package main

import (
	"BookStore/pb"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"strconv"
	"time"
)

const (
	defaultCursor   = "0" // 默认游标，第一页
	defaultPageSize = 2   // 默认每页多少数量
)

// bookstore 的grpc服务

type server struct {
	pb.UnimplementedBookstoreServer

	bs *bookstore // data.go
}

// ListShelves 列出所有书架
func (s *server) ListShelves(ctx context.Context, in *emptypb.Empty) (*pb.ListShelvesResponse, error) {
	// 调用orm操作的方法
	sl, err := s.bs.ListShelves(ctx)
	if err == gorm.ErrEmptySlice { // 没有数据
		return &pb.ListShelvesResponse{}, nil
	}
	if err != nil { // 查询数据库失败
		return nil, status.Error(codes.Internal, "query failed")
	}
	// 封装返回数据,把数据库里的书架格式抓换为pb的格式
	nsl := make([]*pb.Shelf, 0, len(sl))
	for _, s := range sl {
		nsl = append(nsl, &pb.Shelf{
			Id:    s.ID,
			Theme: s.Theme,
			Size:  s.Size,
		})
	}
	return &pb.ListShelvesResponse{Shelves: nsl}, nil
}

// CreateShelf 创建书架
func (s *server) CreateShelf(ctx context.Context, in *pb.CreateShelfRequest) (*pb.Shelf, error) {
	// 参数检查
	if len(in.GetShelf().GetTheme()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid theme")
	}
	// 准备数据
	data := Shelf{
		Theme: in.GetShelf().GetTheme(),
		Size:  in.GetShelf().GetSize(),
	}
	// 去数据库创建
	ns, err := s.bs.CreateShelf(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, "create failed")
	}
	return &pb.Shelf{Id: ns.ID, Theme: ns.Theme, Size: ns.Size}, nil
}

// GetShelf 获取某个书架
func (s *server) GetShelf(ctx context.Context, in *pb.GetShelfRequest) (*pb.Shelf, error) {
	// 参数校验
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	// 去数据库查询
	sh, err := s.bs.GetShelf(ctx, in.GetShelf())

	if err != nil {
		return nil, status.Error(codes.Internal, "query failed")
	}
	return &pb.Shelf{Id: sh.ID, Theme: sh.Theme, Size: sh.Size}, nil
}

func (s *server) DeleteShelf(ctx context.Context, in *pb.DeleteShelfRequest) (*emptypb.Empty, error) {
	// 参数校验
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	// 去数据库删除
	err := s.bs.DeleteShelf(ctx, in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "delete failed")
	}
	return &emptypb.Empty{}, nil
}

// ListBooks 返回某个书架上图书的列表
func (s *server) ListBooks(ctx context.Context, in *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	// 参数校验
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	var (
		cursor   = defaultCursor
		pageSize = defaultPageSize
	)
	/*	if in.GetPageToken() == "" {
		// 没有分页token，默认第一页
	} else {*/
	if len(in.GetPageToken()) > 0 {
		// 有分页的话，先解析分页数据
		pageInfo := Token(in.GetPageToken()).Decode()
		// 再判断解析结果是否有效
		if pageInfo.InValid() {
			return nil, status.Error(codes.InvalidArgument, "invalid page_token")
		}
		cursor = pageInfo.NextID
		pageSize = int(pageInfo.PageSize)
	}
	// 查询数据库
	books, err := s.bs.ListBooks(ctx, in.GetShelf(), cursor, pageSize+1)
	if err != nil {
		fmt.Printf("Get ListBooks failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "query failed")
	}
	// 如果查询出来的结果比pageSize大，说明有下一页
	var (
		hasNextPage   bool
		nextPageToken string
		realSize      int = len(books)
	)
	// 当查询数据库的结果大于pageSize：（1）有下一页（2）格式化数据不用返回所有，只用返回用户需要的一页数据，即pageSize个
	if len(books) > pageSize {
		hasNextPage = true
		realSize = pageSize
	}
	// 封装返回的数据
	// 将 []*Book ----> []*pb.Book
	res := make([]*pb.Book, 0, len(books))
	//for _, b := range books {
	for i := 0; i < realSize; i++ {
		res = append(res, &pb.Book{
			Id:     books[i].ID,
			Title:  books[i].Title,
			Author: books[i].Author,
		})
	}
	// 如果有下一页，就要生下一页的page_token
	if hasNextPage {
		nextPageInfo := Page{
			NextID:        strconv.FormatInt(res[realSize-1].Id, 10), // res[realSize-1].Id最后一个返回结果的id
			NextTimeAtUTC: time.Now().Unix(),
			PageSize:      int64(pageSize),
		}
		nextPageToken = string(nextPageInfo.Encode())
	}
	return &pb.ListBooksResponse{Books: res, NextPage: nextPageToken}, nil
}

func (s *server) CreateBook(ctx context.Context, in *pb.CreateBookRequest) (*pb.Book, error) {
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelfID")
	}
	data := Book{
		Author:  in.GetBook().GetAuthor(),
		Title:   in.GetBook().GetTitle(),
		ShelfID: in.GetShelf(),
	}
	book, err := s.bs.CreateBook(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, "create failed")
	}
	return &pb.Book{Id: book.ID, Author: book.Author, Title: book.Title}, nil
}

func (s *server) DeleteBook(ctx context.Context, in *pb.DeleteBookRequest) (*emptypb.Empty, error) {
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelfID")
	}
	err := s.bs.DeleteBook(ctx, in.GetShelf(), in.GetBook())
	if err != nil {
		return nil, status.Error(codes.Internal, "delete failed")
	}
	return &emptypb.Empty{}, nil
}
