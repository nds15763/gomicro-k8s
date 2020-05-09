package v1

import (
	"context"
	"database/sql"
	"errors"
	pt "api\proto\v1"
	"github.com/coreos/etcd/error"
)

const (
	apiVersion = "1.0"
)

type toDoServiceServer struct {
	db *sql.DB
}

func NewToDoServiceServer(db) *toDoServiceServer {
	return &toDoServiceServer{db}
}

func (s *toDoServiceServer) checkApi(api string) error {
	if api != apiVersion {
		return errors.New("版本错误，传入版本为%s,实际版本为%s", api, apiVersion)
	}
	return nil
}

func (s *toDoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, errors.New("数据库连接错误:%s", error.Error())
	}
	return c, nil
}

func (s *toDoServiceServer) Create(ctx context.Context, req *pt.CreateRequest) (*pt.CreateResponse, error) {
	if err := s.checkApi(req.Api); err != nil {
		return nil, err
	}
	c,err:=s.connect(ctx)
	if err!=nil{
		return nil,err
	}
	defer c.close()
	res,err:=c.ExecContext(ctx,"Insert Into ToDo(`title`) Values (?)",req.ToDo.Title)
	if err!=nil{
		return nil,err//插入失败
	}
	id,err:= res.LastInsertId()
	if err!=nil{
		return nil,err//获取最新ID失败
	}
	return &pt.CreateResponse(Api:apiVersion,Id:id)
}
