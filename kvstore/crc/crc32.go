package crc

type CRC interface {
	Update(data []byte)
	GetValue() uint32
}

type CRCStruct struct {
	value uint32
}

func NewCRC() CRC {
	return &CRCStruct{
		value: 0,
	}
}

func (c *CRCStruct) Update(data []byte) {

}

func (c *CRCStruct) GetValue() uint32 {
	return 0
}
