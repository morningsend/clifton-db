package cliftondbserver

import (
	"encoding/binary"
	"errors"
	"github.com/zl14917/MastersProject/api/internal_request"
	"sync/atomic"
	"unsafe"
)

type RequestIdGenerator interface {
	NextId() uint64
}

type requestIdGenerator struct {
	current uint64
}

func NewRequestIdGenerator(initial uint64) RequestIdGenerator {
	g := &requestIdGenerator{
		current: initial,
	}
	return g
}
func (g *requestIdGenerator) NextId() uint64 {
	v := atomic.AddUint64(&g.current, 1)
	return v
}

type RequestBuilder struct {
	idGenerator RequestIdGenerator
}

func NewRequestBuilder(gen RequestIdGenerator) *RequestBuilder {
	return &RequestBuilder{}
}

func (b *RequestBuilder) NewGetRequest(key []byte) *internal_request.InternalRequest {
	nextId := b.idGenerator.NextId()
	req := &internal_request.InternalRequest{
		Request: &internal_request.InternalRequest_GetReq{GetReq: &internal_request.GetReq{Key: key}},
		Header: &internal_request.RequestHeader{
			ID: nextId,
		},
	}
	return req
}

func (b *RequestBuilder) NewPutRequest(key []byte, value []byte) *internal_request.InternalRequest {
	nextId := b.idGenerator.NextId()
	req := &internal_request.InternalRequest{
		Request: &internal_request.InternalRequest_PutReq{&internal_request.PutReq{Key: key, Value: value}},
		Header:  &internal_request.RequestHeader{ID: nextId},
	}
	return req
}

func (b *RequestBuilder) NewDeleteRequest(key []byte) *internal_request.InternalRequest {
	nextId := b.idGenerator.NextId()
	req := &internal_request.InternalRequest{
		Request: &internal_request.InternalRequest_DeleteReq{&internal_request.DeleteReq{Key: key}},
		Header:  &internal_request.RequestHeader{ID: nextId},
	}
	return req
}

func IdAsBytes(id uint64, buffer []byte) {
	if len(buffer) < int(unsafe.Sizeof(id)) {
		panic(errors.New("buffer too small"))
	}

	binary.BigEndian.PutUint64(buffer, id)
}

func BytesAsId(buffer []byte) uint64 {
	return binary.BigEndian.Uint64(buffer)
}
