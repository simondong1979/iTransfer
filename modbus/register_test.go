package modbus

import (
	"bytes"
	"reflect"
	"testing"
)

const (
	bitQuantity  = 16
	wordQuantity = 3
)

func newNodeReg() *NodeRegister {
	return &NodeRegister{
		slaveID:           0x01,
		coilsAddrStart:    0,
		coilsQuantity:     bitQuantity,
		coils:             []byte{0x55, 0xaa},
		discreteAddrStart: 0,
		discreteQuantity:  bitQuantity,
		discrete:          []byte{0xaa, 0x55},
		inputAddrStart:    0,
		input:             []uint16{0x9012, 0x1234, 0x5678},
		holdingAddrStart:  0,
		holding:           []uint16{0x1234, 0x5678, 0x9012},
	}
}

var readReg = newNodeReg()

func TestNewNodeRegister(t *testing.T) {
	type args struct {
		slaveID           byte
		coilsAddrStart    uint16
		coilsQuantity     uint16
		discreteAddrStart uint16
		discreteQuantity  uint16
		inputAddrStart    uint16
		inputQuantity     uint16
		holdingAddrStart  uint16
		holdingQuantity   uint16
	}
	tests := []struct {
		name string
		args args
		want *NodeRegister
	}{
		{"new node register", args{
			slaveID:           0x01,
			coilsAddrStart:    0,
			coilsQuantity:     10,
			discreteAddrStart: 0,
			discreteQuantity:  10,
			inputAddrStart:    0,
			inputQuantity:     10,
			holdingAddrStart:  0,
			holdingQuantity:   10,
		}, &NodeRegister{
			slaveID:           0x01,
			coilsAddrStart:    0,
			coilsQuantity:     10,
			coils:             make([]byte, 2),
			discreteAddrStart: 0,
			discreteQuantity:  10,
			discrete:          make([]byte, 2),
			inputAddrStart:    0,
			input:             make([]uint16, 10),
			holdingAddrStart:  0,
			holding:           make([]uint16, 10),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewNodeRegister(tt.args.slaveID,
				tt.args.coilsAddrStart, tt.args.coilsQuantity,
				tt.args.discreteAddrStart, tt.args.discreteQuantity,
				tt.args.inputAddrStart, tt.args.inputQuantity,
				tt.args.holdingAddrStart, tt.args.holdingQuantity)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNodeRegister() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeRegister_SlaveID(t *testing.T) {
	tests := []struct {
		name string
		this *NodeRegister
		want uint8
	}{
		{
			"slave ID same",
			&NodeRegister{slaveID: 0x01},
			0x01,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.SlaveID(); got != tt.want {
				t.Errorf("NodeRegister.SlaveID() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNodeRegister_SetSlaveID(t *testing.T) {
	tests := []struct {
		name string
		this *NodeRegister
		want byte
	}{
		{"", &NodeRegister{}, 0x02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.this.SetSlaveID(tt.want)
			if tt.this.slaveID != tt.want {
				t.Errorf("NodeRegister.SetSlaveID() = got %#v, want %#v", tt.this.slaveID, tt.want)
			}
		})
	}
}
func Test_getBits(t *testing.T) {
	type args struct {
		buf   []byte
		start uint16
		nBits uint16
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{"??????0-8???,???8???", args{[]byte{0xaa, 0x5}, 0, 8}, 0xaa},
		{"??????0-4???,???4???", args{[]byte{0xaa, 0x55}, 0, 4}, 0x0a},
		{"??????4-8???,???4???", args{[]byte{0xaa, 0x55}, 4, 4}, 0x0a},
		{"??????4-12???,???4???", args{[]byte{0xaa, 0x55}, 4, 8}, 0x5a},
		{"??????7-9???,???3???", args{[]byte{0xaa, 0x55}, 7, 3}, 0x03},
		{"??????9-16???,???7???", args{[]byte{0xaa, 0x55}, 9, 7}, 0x2a},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBits(tt.args.buf, tt.args.start, tt.args.nBits); got != tt.want {
				t.Errorf("getBits() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_setBits(t *testing.T) {
	type args struct {
		buf   []byte
		start uint16
		nBits uint16
		value byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"??????0-8???,???8???", args{[]byte{0x00, 0x00}, 0, 8, 0xaa}, []byte{0xaa, 0x00}},
		{"??????0-4???,???4???", args{[]byte{0x00, 0x00}, 0, 4, 0x0a}, []byte{0x0a, 0x00}},
		{"??????4-12???,???8???", args{[]byte{0x00, 0x00}, 4, 8, 0xaa}, []byte{0xa0, 0x0a}},
		{"??????1???,???1???", args{[]byte{0x00, 0x00}, 1, 1, 0xff}, []byte{0x02, 0x00}},
		{"??????9-16???,???7???", args{[]byte{0x00, 0x00}, 9, 7, 0xff}, []byte{0x00, 0xfe}},
		{"??????7-9???,???3???", args{[]byte{0x00, 0x00}, 7, 3, 0xff}, []byte{0x80, 0x03}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setBits(tt.args.buf, tt.args.start, tt.args.nBits, tt.args.value)
			if !bytes.Equal(tt.args.buf, tt.want) {
				t.Errorf("setBits() = %#v, want %#v", tt.args.buf, tt.want)
			}
		})
	}
}

func TestNodeRegister_WriteCoils(t *testing.T) {
	type args struct {
		address uint16
		quality uint16
		valBuf  []byte
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []byte
		wantErr bool
	}{
		{"?????????????????????", newNodeReg(), args{address: bitQuantity + 1}, nil, true},
		{"???????????????", newNodeReg(), args{quality: bitQuantity + 1}, nil, true},
		{"?????????????????????", newNodeReg(), args{address: 1, quality: bitQuantity}, nil, true},
		{"???8???", newNodeReg(),
			args{address: 4, quality: 8, valBuf: []byte{0xff}}, []byte{0xf5, 0xaf}, false},
		{"???10???", newNodeReg(),
			args{address: 4, quality: 10, valBuf: []byte{0xff, 0xff}}, []byte{0xf5, 0xbf}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.WriteCoils(tt.args.address, tt.args.quality, tt.args.valBuf); (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.WriteCoils() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.this.coils, tt.want) {
				t.Errorf("NodeRegister.WriteCoils() got = %#v, want %#v", tt.this.coils, tt.want)
			}
		})
	}
}

func TestNodeRegister_WriteSingleCoil(t *testing.T) {
	type args struct {
		address uint16
		val     bool
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []byte
		wantErr bool
	}{
		{"???false", newNodeReg(), args{2, false}, []byte{0x51, 0xaa}, false},
		{"???true", newNodeReg(), args{1, true}, []byte{0x57, 0xaa}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.WriteSingleCoil(tt.args.address, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.WriteSingleCoil() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.this.coils, tt.want) {
				t.Errorf("NodeRegister.WriteSingleCoil() got = %#v, want %#v", tt.this.coils, tt.want)
			}
		})
	}
}

func TestNodeRegister_ReadCoils(t *testing.T) {
	type args struct {
		address uint16
		quality uint16
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []byte
		wantErr bool
	}{
		{"?????????????????????", readReg, args{address: bitQuantity + 1}, nil, true},
		{"???????????????", readReg, args{quality: bitQuantity + 1}, nil, true},
		{"?????????????????????", readReg, args{address: 1, quality: bitQuantity}, nil, true},
		{"???8???", readReg, args{address: 4, quality: 8}, []byte{0xa5}, false},
		{"???10???", readReg, args{address: 4, quality: 10}, []byte{0xa5, 0x02}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.this.ReadCoils(tt.args.address, tt.args.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.ReadCoils() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeRegister.ReadCoils() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNodeRegister_ReadSingleCoil(t *testing.T) {
	type args struct {
		address uint16
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    bool
		wantErr bool
	}{
		{"???false", readReg, args{5}, false, false},
		{"???true", readReg, args{6}, true, false},
		{"?????????", readReg, args{bitQuantity}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.this.ReadSingleCoil(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.ReadSingleCoil() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NodeRegister.ReadSingleCoil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeRegister_WriteDiscretes(t *testing.T) {
	type args struct {
		address uint16
		quality uint16
		valBuf  []byte
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []byte
		wantErr bool
	}{
		{"?????????????????????", newNodeReg(), args{address: bitQuantity + 1}, nil, true},
		{"???????????????", newNodeReg(), args{quality: bitQuantity + 1}, nil, true},
		{"?????????????????????", newNodeReg(), args{address: 1, quality: bitQuantity}, nil, true},
		{"???8???", newNodeReg(),
			args{address: 4, quality: 8, valBuf: []byte{0xff}}, []byte{0xfa, 0x5f}, false},
		{"???10???", newNodeReg(),
			args{address: 4, quality: 10, valBuf: []byte{0xff, 0xff}}, []byte{0xfa, 0x7f}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.WriteDiscretes(tt.args.address, tt.args.quality, tt.args.valBuf); (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.WriteDiscretes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.this.discrete, tt.want) {
				t.Errorf("NodeRegister.WriteDiscretes() got = %#v, want %#v", tt.this.discrete, tt.want)
			}
		})
	}
}

func TestNodeRegister_WriteSingleDiscrete(t *testing.T) {
	type args struct {
		address uint16
		val     bool
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []byte
		wantErr bool
	}{
		{"???false", newNodeReg(), args{1, false}, []byte{0xa8, 0x55}, false},
		{"???true", newNodeReg(), args{2, true}, []byte{0xae, 0x55}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.WriteSingleDiscrete(tt.args.address, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.WriteSingleDiscrete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.this.discrete, tt.want) {
				t.Errorf("NodeRegister.WriteSingleCoil() got = %#v, want %#v", tt.this.discrete, tt.want)
			}
		})
	}
}

func TestNodeRegister_ReadDiscretes(t *testing.T) {
	type args struct {
		address uint16
		quality uint16
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []byte
		wantErr bool
	}{
		{"?????????????????????", readReg, args{address: bitQuantity + 1}, nil, true},
		{"???????????????", readReg, args{quality: bitQuantity + 1}, nil, true},
		{"?????????????????????", readReg, args{address: 1, quality: bitQuantity}, nil, true},
		{"???8???", readReg, args{address: 4, quality: 8}, []byte{0x5a}, false},
		{"???10???", readReg, args{address: 4, quality: 10}, []byte{0x5a, 0x01}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.this.ReadDiscretes(tt.args.address, tt.args.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.ReadDiscretes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeRegister.ReadDiscretes() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNodeRegister_ReadSingleDiscrete(t *testing.T) {
	type args struct {
		address uint16
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    bool
		wantErr bool
	}{
		{"???false", readReg, args{5}, true, false},
		{"???true", readReg, args{6}, false, false},
		{"?????????", readReg, args{bitQuantity}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.this.ReadSingleDiscrete(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.ReadSingleDiscrete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NodeRegister.ReadSingleDiscrete() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestNodeRegister_WriteHoldingsBytes(t *testing.T) {
	type args struct {
		address uint16
		quality uint16
		valBuf  []byte
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []uint16
		wantErr bool
	}{
		{"?????????????????????", newNodeReg(), args{address: wordQuantity + 1}, nil, true},
		{"???????????????", newNodeReg(), args{quality: wordQuantity + 1}, nil, true},
		{"?????????????????????", newNodeReg(), args{address: 1, quality: wordQuantity}, nil, true},
		{"???2????????????", newNodeReg(), args{address: 1, quality: 2, valBuf: []byte{0x11, 0x11, 0x22, 0x22}}, []uint16{0x1234, 0x1111, 0x2222}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.WriteHoldingsBytes(tt.args.address, tt.args.quality, tt.args.valBuf); (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.WriteHoldingsBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.this.holding, tt.want) {
				t.Errorf("NodeRegister.WriteHoldingsBytes() got = %#v, want %#v", tt.this.holding, tt.want)
			}
		})
	}
}

func TestNodeRegister_WriteHoldings(t *testing.T) {
	type args struct {
		address uint16
		valBuf  []uint16
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []uint16
		wantErr bool
	}{
		{"?????????????????????", newNodeReg(), args{address: wordQuantity + 1}, nil, true},
		{"???????????????", newNodeReg(), args{valBuf: make([]uint16, wordQuantity+1)}, nil, true},
		{"?????????????????????", newNodeReg(), args{address: 1, valBuf: make([]uint16, wordQuantity+1)}, nil, true},
		{"???2????????????", newNodeReg(), args{address: 1, valBuf: []uint16{0x1111, 0x2222}}, []uint16{0x1234, 0x1111, 0x2222}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.WriteHoldings(tt.args.address, tt.args.valBuf); (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.WriteHoldings() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.this.holding, tt.want) {
				t.Errorf("NodeRegister.WriteHoldings() got = %#v, want %#v", tt.this.holding, tt.want)
			}
		})
	}
}

func TestNodeRegister_ReadHoldingsBytes(t *testing.T) {
	type args struct {
		address uint16
		quality uint16
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []byte
		wantErr bool
	}{
		{"?????????????????????", readReg, args{address: wordQuantity + 1}, nil, true},
		{"???????????????", readReg, args{quality: wordQuantity + 1}, nil, true},
		{"?????????????????????", readReg, args{address: 1, quality: wordQuantity + 1}, nil, true},
		{"???2????????????", readReg, args{address: 1, quality: 2}, []byte{0x56, 0x78, 0x90, 0x12}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.this.ReadHoldingsBytes(tt.args.address, tt.args.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.ReadHoldingsBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeRegister.ReadHoldingsBytes() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNodeRegister_ReadHoldings(t *testing.T) {
	type args struct {
		address uint16
		quality uint16
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []uint16
		wantErr bool
	}{
		{"?????????????????????", readReg, args{address: wordQuantity + 1}, nil, true},
		{"???????????????", readReg, args{quality: wordQuantity + 1}, nil, true},
		{"?????????????????????", readReg, args{address: 1, quality: wordQuantity + 1}, nil, true},
		{"???2????????????", readReg, args{address: 1, quality: 2}, []uint16{0x5678, 0x9012}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.this.ReadHoldings(tt.args.address, tt.args.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.ReadHoldings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeRegister.ReadHoldings() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNodeRegister_WriteInputsBytes(t *testing.T) {
	type args struct {
		address uint16
		quality uint16
		regBuf  []byte
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []uint16
		wantErr bool
	}{
		{"?????????????????????", newNodeReg(), args{address: wordQuantity + 1}, nil, true},
		{"???????????????", newNodeReg(), args{quality: wordQuantity + 1}, nil, true},
		{"?????????????????????", newNodeReg(), args{address: 1, quality: wordQuantity}, nil, true},
		{
			"???2????????????", newNodeReg(),
			args{address: 1, quality: 2, regBuf: []byte{0x11, 0x11, 0x22, 0x22}},
			[]uint16{0x9012, 0x1111, 0x2222},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.WriteInputsBytes(tt.args.address, tt.args.quality, tt.args.regBuf); (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.WriteInputsBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.this.input, tt.want) {
				t.Errorf("NodeRegister.WriteInputsBytes() got = %#v, want %#v", tt.this.input, tt.want)
			}
		})
	}
}

func TestNodeRegister_WriteInputs(t *testing.T) {
	type args struct {
		address uint16
		valBuf  []uint16
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []uint16
		wantErr bool
	}{
		{"?????????????????????", newNodeReg(), args{address: wordQuantity + 1}, nil, true},
		{"???????????????", newNodeReg(), args{valBuf: make([]uint16, wordQuantity+1)}, nil, true},
		{"?????????????????????", newNodeReg(), args{address: 1, valBuf: make([]uint16, wordQuantity+1)}, nil, true},
		{"???2????????????", newNodeReg(), args{address: 1, valBuf: []uint16{0x1111, 0x2222}}, []uint16{0x9012, 0x1111, 0x2222}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.WriteInputs(tt.args.address, tt.args.valBuf); (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.WriteInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.this.input, tt.want) {
				t.Errorf("NodeRegister.WriteInputs() got = %#v, want %#v", tt.this.input, tt.want)
			}
		})
	}
}

func TestNodeRegister_ReadInputsBytes(t *testing.T) {
	type args struct {
		address uint16
		quality uint16
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []byte
		wantErr bool
	}{
		{"?????????????????????", readReg, args{address: wordQuantity + 1}, nil, true},
		{"???????????????", readReg, args{quality: wordQuantity + 1}, nil, true},
		{"?????????????????????", readReg, args{address: 1, quality: wordQuantity + 1}, nil, true},
		{"???2????????????", readReg, args{address: 1, quality: 2}, []byte{0x12, 0x34, 0x56, 0x78}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.this.ReadInputsBytes(tt.args.address, tt.args.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.ReadInputsBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeRegister.ReadInputsBytes() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNodeRegister_ReadInputs(t *testing.T) {
	type args struct {
		address uint16
		quality uint16
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    []uint16
		wantErr bool
	}{
		{"?????????????????????", readReg, args{address: wordQuantity + 1}, nil, true},
		{"???????????????", readReg, args{quality: wordQuantity + 1}, nil, true},
		{"?????????????????????", readReg, args{address: 1, quality: wordQuantity + 1}, nil, true},
		{"???2????????????", readReg, args{address: 1, quality: 2}, []uint16{0x1234, 0x5678}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.this.ReadInputs(tt.args.address, tt.args.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.ReadInputs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeRegister.ReadInputs() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNodeRegister_MaskWriteHolding(t *testing.T) {
	nodeReg := &NodeRegister{
		holdingAddrStart: 0,
		holding:          []uint16{0x0000, 0x0012, 0x0000},
	}
	type args struct {
		address uint16
		andMask uint16
		orMask  uint16
	}
	tests := []struct {
		name    string
		this    *NodeRegister
		args    args
		want    uint16
		wantErr bool
	}{
		{"??????", nodeReg, args{1, 0xf2, 0x25}, 0x0017, false},
		{"???????????????", nodeReg, args{address: wordQuantity + 1}, 0x0012, true},
		{"???????????????", nodeReg, args{address: wordQuantity}, 0x0012, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.this.MaskWriteHolding(tt.args.address, tt.args.andMask, tt.args.orMask); (err != nil) != tt.wantErr {
				t.Errorf("NodeRegister.MaskWriteHolding() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.this.holding[int(tt.args.address)] != tt.want {
				t.Errorf("NodeRegister.MaskWriteHolding() got = %#v, want %#v", tt.this.holding[tt.args.address], tt.want)
			}
		})
	}
}

func Benchmark_getBits(b *testing.B) {
	val := []byte{0x00, 0x02, 0x03, 0x04, 0x05}
	for i := 0; i < b.N; i++ {
		getBits(val, 1, 30)
	}
}

func Benchmark_setBits(b *testing.B) {
	val := []byte{0x00, 0x02, 0x03, 0x04, 0x05}
	for i := 0; i < b.N; i++ {
		setBits(val, 12, 8, 0xaa)
	}
}
