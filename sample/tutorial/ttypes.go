// Autogenerated by Thrift Compiler (0.9.2)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package tutorial

import (
	"bytes"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/allenma/gosoa/sample/shared"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

var _ = shared.GoUnusedProtection__
var GoUnusedProtection__ int

//You can define enums, which are just 32 bit integers. Values are optional
//and start at 1 if not supplied, C style again.
type Operation int64

const (
	Operation_ADD      Operation = 1
	Operation_SUBTRACT Operation = 2
	Operation_MULTIPLY Operation = 3
	Operation_DIVIDE   Operation = 4
)

func (p Operation) String() string {
	switch p {
	case Operation_ADD:
		return "Operation_ADD"
	case Operation_SUBTRACT:
		return "Operation_SUBTRACT"
	case Operation_MULTIPLY:
		return "Operation_MULTIPLY"
	case Operation_DIVIDE:
		return "Operation_DIVIDE"
	}
	return "<UNSET>"
}

func OperationFromString(s string) (Operation, error) {
	switch s {
	case "Operation_ADD":
		return Operation_ADD, nil
	case "Operation_SUBTRACT":
		return Operation_SUBTRACT, nil
	case "Operation_MULTIPLY":
		return Operation_MULTIPLY, nil
	case "Operation_DIVIDE":
		return Operation_DIVIDE, nil
	}
	return Operation(0), fmt.Errorf("not a valid Operation string")
}

func OperationPtr(v Operation) *Operation { return &v }

//Thrift lets you do typedefs to get pretty names for your types. Standard
//C style here.
type MyInteger int32

func MyIntegerPtr(v MyInteger) *MyInteger { return &v }

type Work struct {
	Num1    int32     `thrift:"num1,1" json:"num1"`
	Num2    int32     `thrift:"num2,2" json:"num2"`
	Op      Operation `thrift:"op,3" json:"op"`
	Comment *string   `thrift:"comment,4" json:"comment"`
}

func NewWork() *Work {
	return &Work{}
}

func (p *Work) GetNum1() int32 {
	return p.Num1
}

func (p *Work) GetNum2() int32 {
	return p.Num2
}

func (p *Work) GetOp() Operation {
	return p.Op
}

var Work_Comment_DEFAULT string

func (p *Work) GetComment() string {
	if !p.IsSetComment() {
		return Work_Comment_DEFAULT
	}
	return *p.Comment
}
func (p *Work) IsSetComment() bool {
	return p.Comment != nil
}

func (p *Work) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.ReadField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.ReadField4(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *Work) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.Num1 = v
	}
	return nil
}

func (p *Work) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.Num2 = v
	}
	return nil
}

func (p *Work) ReadField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return fmt.Errorf("error reading field 3: %s", err)
	} else {
		temp := Operation(v)
		p.Op = temp
	}
	return nil
}

func (p *Work) ReadField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 4: %s", err)
	} else {
		p.Comment = &v
	}
	return nil
}

func (p *Work) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Work"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *Work) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("num1", thrift.I32, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:num1: %s", p, err)
	}
	if err := oprot.WriteI32(int32(p.Num1)); err != nil {
		return fmt.Errorf("%T.num1 (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:num1: %s", p, err)
	}
	return err
}

func (p *Work) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("num2", thrift.I32, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:num2: %s", p, err)
	}
	if err := oprot.WriteI32(int32(p.Num2)); err != nil {
		return fmt.Errorf("%T.num2 (2) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:num2: %s", p, err)
	}
	return err
}

func (p *Work) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("op", thrift.I32, 3); err != nil {
		return fmt.Errorf("%T write field begin error 3:op: %s", p, err)
	}
	if err := oprot.WriteI32(int32(p.Op)); err != nil {
		return fmt.Errorf("%T.op (3) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 3:op: %s", p, err)
	}
	return err
}

func (p *Work) writeField4(oprot thrift.TProtocol) (err error) {
	if p.IsSetComment() {
		if err := oprot.WriteFieldBegin("comment", thrift.STRING, 4); err != nil {
			return fmt.Errorf("%T write field begin error 4:comment: %s", p, err)
		}
		if err := oprot.WriteString(string(*p.Comment)); err != nil {
			return fmt.Errorf("%T.comment (4) field write error: %s", p, err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 4:comment: %s", p, err)
		}
	}
	return err
}

func (p *Work) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Work(%+v)", *p)
}

type InvalidOperation struct {
	WhatOp int32  `thrift:"whatOp,1" json:"whatOp"`
	Why    string `thrift:"why,2" json:"why"`
}

func NewInvalidOperation() *InvalidOperation {
	return &InvalidOperation{}
}

func (p *InvalidOperation) GetWhatOp() int32 {
	return p.WhatOp
}

func (p *InvalidOperation) GetWhy() string {
	return p.Why
}
func (p *InvalidOperation) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error: %s", p, err)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.ReadField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.ReadField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *InvalidOperation) ReadField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return fmt.Errorf("error reading field 1: %s", err)
	} else {
		p.WhatOp = v
	}
	return nil
}

func (p *InvalidOperation) ReadField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s", err)
	} else {
		p.Why = v
	}
	return nil
}

func (p *InvalidOperation) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("InvalidOperation"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("write struct stop error: %s", err)
	}
	return nil
}

func (p *InvalidOperation) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("whatOp", thrift.I32, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:whatOp: %s", p, err)
	}
	if err := oprot.WriteI32(int32(p.WhatOp)); err != nil {
		return fmt.Errorf("%T.whatOp (1) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:whatOp: %s", p, err)
	}
	return err
}

func (p *InvalidOperation) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("why", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:why: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Why)); err != nil {
		return fmt.Errorf("%T.why (2) field write error: %s", p, err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:why: %s", p, err)
	}
	return err
}

func (p *InvalidOperation) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("InvalidOperation(%+v)", *p)
}

func (p *InvalidOperation) Error() string {
	return p.String()
}
