package server

import (
	"github.com/Godyu97/geerpc/codec"
	"net"
	"github.com/Godyu97/geerpc/logger"
	"io"
	"encoding/json"
	"sync"
	"reflect"
	"fmt"
)

var DefaultOption = codec.Option{
	MagicNumber: codec.MagicNumber,
	CodecType:   codec.JsonType,
}

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			logger.Error("rpc server: accept error:", err)
			return
		}
		logger.Info("开启连接:conn")
		go s.ServeConn(conn)
	}
}

func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

func (s *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() {
		_ = conn.Close()
	}()
	var opt codec.Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		logger.Error("rpc server: options error: ", err)
		return
	}
	if opt.MagicNumber != codec.MagicNumber {
		logger.Error("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	f, ok := codec.NewCodecFuncMap[opt.CodecType]
	if ok == false {
		logger.Error("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	logger.Info("开启连接:", opt)
	s.serveCodec(f(conn))
}

// invalidRequest is a placeholder for response argv when error occurs
var invalidRequest = struct{}{}

func (s *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		logger.Info("begin readReq")
		req, err := s.readRequest(cc)
		logger.Info("end readReq")
		if err != nil {
			logger.Error("it's not possible to recover, so close the connection err:", err)
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			s.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go s.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

// request stores all information of a call
type request struct {
	h      *codec.Header // header of request
	argv   reflect.Value
	replyv reflect.Value // argv and replyv of request
}

func (s *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	err := cc.ReadHeader(&h)
	if err != nil {
		logger.Error("rpc server: read header error:", err)
		return nil, err
	}
	return &h, nil
}

func (s *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := s.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}

	// TODO: now we don't know the type of request argv
	// day 1, just suppose it's string
	req.argv = reflect.New(reflect.TypeOf(""))
	if err := cc.ReadBody(req.argv.Interface()); err != nil {
		logger.Error("rpc server: read argv err:", err)
	}
	logger.Info("readRequest req:", req)
	return req, nil
}

func (s *Server) sendResponse(cc codec.Codec, h *codec.Header, body any, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		logger.Error("rpc server: write response error:", err)
	}
}

func (s *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	// TODO, should call registered rpc methods to get the right replyv
	// day 1, just print argv and send a hello message
	defer wg.Done()
	logger.Info("处理请求:", req.h, req.argv.Elem())

	req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.h.Seq))
	s.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}
