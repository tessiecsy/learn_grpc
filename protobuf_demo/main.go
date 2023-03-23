package main

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"protobuf_demo/api"
)

type Book struct {
	Price int64 // 如何区分默认值和0
	// go语言中两种方法
	//Price sql.NullInt64  // 自定义结构体， 字段 Valid bool：Valid is true if Int64 is not NULL
	//Price *int64  // 指针方式
}

func foo() {
	//var book Book
	/*	if book.Price == nil {
			// 没有赋值 （如果指针就是nil）
		} else {
			// 赋值
		}*/
	//book = Book{Price: 0}
}

func oneofDemo() {
	// client
	req1 := &api.NoticeReaderRequest{
		Msg: "博客更新了",
		NoticeWay: &api.NoticeReaderRequest_Email{
			Email: "123@xxx.com",
		},
	}
	/*	req2 := &api.NoticeReaderRequest{
		Msg: "博客更新了",
		NoticeWay: &api.NoticeReaderRequest_Phone{
			Phone: "123456",
		},
	}*/

	// server
	// 怎么知道客户端拿的是oneof的哪一个呢? 使用类型断言，根据不同的类型，执行不同的操作
	req := req1
	switch v := req.NoticeWay.(type) {
	case *api.NoticeReaderRequest_Email:
		noticeWithEmail(v) // 发邮件通知
	case *api.NoticeReaderRequest_Phone:
		noticeWithPhone(v)
	}

}

// 使用google/protobuf/wrappers.proto中定义的WrapValue，本质上就是使用自定义message代替基本类型。
// 判断是不是有没有设置值
/*func wrapValueDemo() {
	// client
	book := &api.Book{
		Title: "《hello world》",
		Price: &wrapperspb.Int64Value{Value: 9900},
		Memo:  &wrapperspb.StringValue{Value: "hhhhh"},
	}
	// server
	if book.GetPrice() == nil {
		// 说明没有给price赋值
		fmt.Println("没有设置price")
	} else {
		// 赋值了就可以去用
		fmt.Println(book.GetPrice().GetValue())
	}
	if book.GetMemo() != nil {
		fmt.Println(book.GetMemo().GetValue())
	}
}*/

/*func optionalDemo() {
	// client
	book := api.Book{
		Title: "《hello》",
		Price: proto.Int64(9900), // 返回的是一个指针，就可以判断是不是nil
	}
	// server
	// 如何判断book.Price有没有赋值呢？
	// 注意是book.Price,不是GetPrice,因为如果没赋值，GetPrice返回的是int64的0,没法和nil比较
	if book.Price == nil {
		fmt.Println("no price")
	} else {
		fmt.Printf("book with price:%v\n", book.GetPrice())
	}
}*/

// fieldMaskDemo 实现部分更新
func fieldMaskDemo() {
	// client
	paths := []string{"price", "info.b", "author"} // 更新的字段信息
	req := &api.UpdateBookRequest{
		Op: "csy",
		Book: &api.Book{
			Author: "zhangyf",
			Price:  &wrapperspb.Int64Value{Value: 8800},
			Info: &api.Book_Info{
				B: "bbbb",
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: paths}, // 提供哪些字段更新了，所以更新的信息一定要写入paths
	}

	// server
	// UpdateMask是更新的字段信息，放到mask里去
	mask, _ := fieldmask_utils.MaskFromProtoFieldMask(req.UpdateMask, generator.CamelCase)
	var bookDst = make(map[string]interface{}) // 这是更新的信息
	// 将数据读取到map[string]interface{}
	// fieldmask-utils支持读取到结构体等，更多用法可查看文档。
	// 把req.Book的数据取到bookDst里
	fieldmask_utils.StructToMap(mask, req.Book, bookDst)
	// do update with bookDst
	fmt.Printf("bookDst:%#v\n", bookDst)
}

// 发送通知相关的功能函数
func noticeWithEmail(in *api.NoticeReaderRequest_Email) {
	fmt.Printf("notice reader by email:%v\n", in.Email)
}

func noticeWithPhone(in *api.NoticeReaderRequest_Phone) {
	fmt.Printf("notice reader by phone:%v\n", in.Phone)
}

func main() {
	oneofDemo()
	//wrapValueDemo()
	fieldMaskDemo()
}
