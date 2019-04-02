package vector

import (
	"errors"
	"math/bits"
	"sync/atomic"
	"unsafe"
)

type descriptor struct {
	size      uint
	pendingOp writeDescriptor
}

const firstBucketSize uint = 8
const maxBucketLevel uint = 8

type vectorValue struct {
	data interface{}
}

func higestBit(v uint) int {
	zeros := bits.LeadingZeros(uint(v))
	return bits.Len(uint(v)) - zeros
}

type writeDescriptor struct {
	oldValue *vectorValue
	newValue *vectorValue
	location int
	pending  bool
}

type ConcurrentVector struct {
	descriptor *descriptor
	memory     [][]*vectorValue
}

func NewConcurrentVector(cap int) ConcurrentVector {
	data := make([][]*vectorValue, 0, cap)
	return ConcurrentVector{
		descriptor: &descriptor{
			size: 0,
		},
		memory: data,
	}
}

func (v *ConcurrentVector) Resize(size int) {

}

func (v *ConcurrentVector) at(idx int) unsafe.Pointer {
	return unsafe.Pointer(&(v.memory[idx]))
}
func (v *ConcurrentVector) Get(idx int) interface{} {
	return nil
}

func (v *ConcurrentVector) Set(idx int, data interface{}) {

}

func (v *ConcurrentVector) completeWrite(writeop writeDescriptor) {
	if writeop.pending {
		atomic.CompareAndSwapPointer(
			((*unsafe.Pointer))(v.at(writeop.location)),
			unsafe.Pointer(writeop.oldValue),
			unsafe.Pointer(writeop.newValue),
		)
	}
}

func (v *ConcurrentVector) Append(value interface{}) {
	var nextDesc *descriptor
	for {
		desc := v.descriptor
		v.completeWrite(desc.pendingOp)

		bucket := higestBit(desc.size+firstBucketSize) - higestBit(firstBucketSize)

		if v.memory[bucket] == nil {
			v.allocateBucket(uint(bucket))
		}

		writeOp := writeDescriptor{
			location: int(desc.size),
			pending:  true,
			oldValue: nil,
			newValue: &vectorValue{
				data: value,
			},
		}

		nextDesc = &descriptor{
			size:      desc.size + 1,
			pendingOp: writeOp,
		}

		if atomic.CompareAndSwapPointer(nil, nil, nil) {
			break
		}
	}

	v.completeWrite(nextDesc.pendingOp)
}

func (v *ConcurrentVector) Pop() (data interface{}, ok bool) {
	return nil, false
}

func (v *ConcurrentVector) Len() int {
	desc := v.descriptor
	size := desc.size
	if desc.pendingOp.pending {
		size -= 1
	}
	return int(size)
}

func (v *ConcurrentVector) allocateBucket(bucketLevel uint) {
	if bucketLevel > maxBucketLevel {
		panic(errors.New("cannot allocate vector memory, size exceeded"))
	}

	bucketSize := firstBucketSize
	for i := 0; i <= int(bucketLevel); i++ {
		bucketLevel *= firstBucketSize
	}

	newBucket := make([]*vectorValue, 0, bucketSize)

	atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&(v.memory[bucketLevel]))),
		nil,
		unsafe.Pointer(&newBucket),
	)
}
