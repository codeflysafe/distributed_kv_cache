package go_rpc

import (
	"encoding/json"
	"fmt"
	"go_rpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

// 魔数， caffbybe
const MagicNumber = 0xcaffbabe

type Option struct {
	MagicNumber int        // 魔数
	CodecType   codec.Type // 采用的序列化模式
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,    // MagicNumber marks
	CodecType:   codec.JsonType, // 默认采用 json 格式， 便于跨语言调用
}

// server represents an rpc server
type Server struct{}

func NewServer() *Server {
	return &Server{}
}

// 提供一个默认的server
var DefaultServer = NewServer()

func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc Server: accept error ", err)
			return
		}
		// 开一个协程出去处理它
		go s.ServerCoon(conn)
	}
}

func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

// 用 json.NewDecoder 反序列化得到 Option 实例，检查 MagicNumber 和 CodeType 的值是否正确。
// 然后根据 CodeType 得到对应的消息编解码器，接下来的处理交给 serverCodec。
func (s *Server) ServerCoon(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server options error: ", err)
		return
	}

	// 去读取传输结果
	cc, ok := codec.NewCodecFuncMap[opt.CodecType]
	if !ok {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	s.serveCodec(cc(conn))

}

type request struct {
	h            *codec.Header
	argv, replyv reflect.Value
}

func (s *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
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
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

func (s *Server) sendResonse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {

	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

func (s *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("server go rpc resp %d", req.h.Seq))
	s.sendResonse(cc, req.h, req.replyv.Interface(), sending)
}

// invalidRequest is a placeholder for response argv when error occurs
var invalidRequest = struct{}{}

// 1. readrequest 2. handle request 3. sendresponse
func (s *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := s.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			s.sendResonse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go s.handleRequest(cc, req, sending, wg)
	}

	wg.Wait()
	_ = cc.Close()
}
