package tdf

import (
	"bufio"
	"bytes"
	"gioui.org/unit"
	"github.com/ddkwork/golibrary/mylog"
	"io"
	"slices"
)

type Node struct {
	Data         any
	parent       *Node
	Children     []*Node
	filteredRows []*Node
	buf          *bytes.Buffer
	tag          string
	wireType     BaseType
}

func (n *Node) AddChildByData(data any) { n.AddChild(NewNode(nil)) }

func NewNode(b []byte) *Node {
	n := &Node{
		Data:         nil,
		parent:       nil,
		Children:     nil,
		filteredRows: nil,
		buf:          bytes.NewBuffer(b),
		tag:          "",
		wireType:     0,
	}

	if b != nil {
		n.decodeTagAndWireType()
	}
	return n
}
func (n *Node) Marshal() {}
func (n *Node) encodeTagAndWireType(tag string, wireType BaseType) {
	n.buf.Write(EncodeTag(tag))
	n.buf.WriteByte(byte(wireType))
}
func (n *Node) decodeTagAndWireType() (tag string, wireType BaseType) {
	tagBuf := make([]byte, 3)
	mylog.Check2(n.buf.Read(tagBuf))
	tag = DecodeTag(tagBuf)
	typeBuf := mylog.Check2(n.buf.ReadByte())
	wireType = BaseType(typeBuf)
	n.tag = tag
	n.wireType = wireType
	return
}

func (n *Node) EncodeString(tag, value string) []byte {
	n.encodeTagAndWireType(tag, StringType)
	n.buf.Write(slices.Concat(compressInteger(uint32(len(value)+1)), []byte(value), []byte{0}))
	return n.buf.Bytes()
}

func (n *Node) DecodeString() string {
	length := decompressInteger(n.buf) //这里reader管理偏移，前面读掉tag4字节，来到这里第一个字节就是长度，后面是字符串，最后一个字节是结束符
	result := make([]byte, length-1)   //去掉结束符
	mylog.Check2(io.ReadFull(n.buf, result))
	mylog.Check2(n.buf.ReadByte()) //c格式字符串的结束符0
	if length == 1 {               //这里合适?感觉应该往上移，需要过更多的单元测试
		return ""
	}
	return string(result)
}

func (n *Node) EncodeInteger(tag string, v uint32) []byte {
	n.encodeTagAndWireType(tag, IntegerType)
	n.buf.Write(compressInteger(v))
	return n.buf.Bytes()
}
func (n *Node) DecodeInteger() uint32 {
	return decompressInteger(n.buf)
}

func (n *Node) EncodeBlob(tag string, v []byte) []byte {
	n.encodeTagAndWireType(tag, BinaryType)
	n.buf.Write(slices.Concat(compressInteger(uint32(len(v))), v))
	return n.buf.Bytes()
}
func (n *Node) DecodeBlob() []byte {
	length := decompressInteger(n.buf)
	result := make([]byte, length)
	mylog.Check2(io.ReadFull(n.buf, result))
	return result
}

func (n *Node) EncodeStruct(tag string, v *Struct) []byte {
	n.encodeTagAndWireType(tag, StructType)
	//n.buf.Write(slices.Concat(compressInteger(uint32(len(v.Fields))), v.Fields))
	//js版本是填充0结束的，有点奇怪,大小似乎没有填充
	return nil

}
func (n *Node) DecodeStruct() *Struct {
	return nil
}

type Struct struct {
	objectID int
	Fields   []*Node
	*bufio.ReadWriter
}

func (t *Struct) AddField(field *Node) {
	t.Fields = append(t.Fields, field)
}

//todo 两个问题 结构体嵌套是否可以使用n叉树，头部16字节算法
// 结构体解码应该返回树形结构才合理，js版本似乎只是返回不同类型的切片
// js版本的list返回相同类型切片是合理的

// 这里相当于 SetRootRows 解码填充n叉树所有节点
//func readStruct(stream *bufio.Reader) []Node {
//	var result []Node
//	for {
//		buffer := mylog.Check2(stream.Peek(1)) //todo rename as tag
//		if buffer[0] == 0 {
//			break
//		}
//		t := Node{}
//		result = append(result, *t.Read(stream))
//	}
//	return result
//}

// Write 将结构体类型的 Node 写入流,编码的话是递归遍历所有节点序列化并组合每个节点的序列化buffer
//func (t *Struct) Write() []byte {
//	buffer := bytes.Buffer{}
//	buffer.Write(t.Tag.Marshal())
//	for _,  := range t.Fields { //todo 改为walk root
//		data := .Write()
//		buffer.Write(data)
//	}
//	buffer.WriteByte(0)
//	return buffer.Bytes()
//}

func (n *Node) EncodeList(tag string, v []any) []byte {
	n.encodeTagAndWireType(tag, ListType)
	return nil
}
func (n *Node) DecodeList() *List {
	return nil
}
func (n *Node) EncodeDictionary(tag string, v map[any]any) []byte {
	n.encodeTagAndWireType(tag, MapType)
	return nil
}
func (n *Node) DecodeDictionary() map[any]any {
	return nil
}
func (n *Node) EncodeUnion(tag string, v any) []byte {
	n.encodeTagAndWireType(tag, UnionType)
	return nil
}
func (n *Node) DecodeUnion() *Union {
	return nil
}
func (n *Node) EncodeIntegerList(tag string, v []uint32) []byte {
	n.encodeTagAndWireType(tag, ListType)
	return nil
}
func (n *Node) DecodeIntegerList() []uint32 {
	return nil
}
func (n *Node) EncodeIntVector2(tag string, v any) []byte {
	n.encodeTagAndWireType(tag, BlazeObjectTypeType)
	return nil
}
func (n *Node) DecodeIntVector2() {
	return
}
func (n *Node) EncodeIntVector3(tag string, v any) []byte {
	n.encodeTagAndWireType(tag, BlazeObjectIdType)
	return nil
}
func (n *Node) DecodeIntVector3() {
	return
}

//func (n *Node) UnMarshal() {
//	label := mylog.Check2(n.Reader.Peek(3))
//	BaseKind := BaseKind(mylog.Check2(n.Reader.Peek(1))[0])
//	switch BaseKind { //类型应该是在tag后面，第四字节?
//	case Integer:
//		return IntegerRead(label, stream)
//	case String:
//		return StringRead(label, stream)
//	case Blob:
//		return BlobRead(label, stream)
//	case Struct:
//		return StructRead(label, stream)
//	case List:
//		return ListRead(label, stream)
//	case Dictionary:
//		return DictionaryRead(label, stream)
//	case Union:
//		return UnionRead(label, stream)
//	case IntegerList:
//		return IntegerListRead(label, stream)
//	case IntVector2:
//		return IntVector2Read(label, stream)
//	case IntVector3:
//		return IntVector3Read(label, stream)
//	default:
//		return Node{}, &UnknownBaseKind{Type: string(typeByte)}
//	}
//}

// List 同类型切片，结构体是不同类型切片+嵌套
type List struct {
	Value   any //todo 限制类型
	Subtype BaseType
	Length  int
	*bufio.ReadWriter
}

// ListRead 从流中读取一个列表类型的 Node
//func ListRead(label []byte, stream io.Reader) Node {
//	subtype := decompressInteger(stream) // 需要实现的函数
//	length := decompressInteger(stream)  // 需要实现的函数
//	var value any
//	switch BaseKind(subtype) {
//	case Integer:
//		value = decodeIntegerList(stream, length) // 需要实现的函数
//	case String:
//		value = decodeStringList(stream, length) // 需要实现的函数
//	case Struct:
//		value = readStructList(stream, length) // 需要实现的函数
//	default:
//		return Node{}, &NotImplemented{Type: string(subtype)}
//	}
//	return List{Node: Node{Label: string(label), Type: List}, Value: value, Subtype: BaseKind(subtype), Length: length}
//}

func (t *List) encode() []byte {
	//buffer := bytes.Buffer{}
	//buffer.Write(t.Tag.Marshal())
	//buffer.WriteByte(byte(t.Subtype))
	//buffer.Write(compressInteger(t.Length)) // 需要实现的函数

	// 处理每个子项的写入
	// 需根据具体类型处理，伪代码如下：
	// switch t.Subtype {
	// case Integer:
	//     for _, v := range t.Fields.([]int) {
	//         buffer.Write(compressInteger(v))
	//     }
	// ...
	// }
	//return buffer.Bytes()
	return nil
}

type DictionaryMap map[any]any

// Dictionary 表示字典类型的 Node
type Dictionary struct {
	Value     DictionaryMap
	KeyType   BaseType
	ValueType BaseType
	Length    int
	*bufio.ReadWriter
}

// DictionaryRead 从流中读取一个字典类型的 Node
//func DictionaryRead(label []byte, stream io.Reader) Node {
//	keyType := decompressInteger(stream)   // 需要实现的函数
//	valueType := decompressInteger(stream) // 需要实现的函数
//	length := decompressInteger(stream)    // 需要实现的函数
//	value := make(Dictionary)
//	for i := 0; i < length; i++ {
//		var key, value any
//		switch BaseKind(keyType) {
//		case Integer:
//			key = decompressInteger(stream) // 需要实现的函数
//		case String:
//			key = readString(stream) // 需要实现的函数
//		default:
//			return Node{}, &NotImplemented{Type: string(keyType)}
//		}
//
//		switch BaseKind(valueType) {
//		case Integer:
//			value = decompressInteger(stream) // 需要实现的函数
//		case String:
//			value = readString(stream) // 需要实现的函数
//		case Struct:
//			value = readStruct(stream) // 需要实现的函数
//		default:
//			return Node{}, &NotImplemented{Type: string(valueType)}
//		}
//		value[key] = value
//	}
//	return Dictionary{Node: Node{Label: string(label), Type: Dictionary}, KeyType: BaseKind(keyType), ValueType: BaseKind(valueType), Length: length, Value: value}
//}
//
//// Write 将字典类型的 Node 写入流
//func (t *Dictionary) Write() []byte {
//	buffer := bytes.Buffer{}
//	buffer.Write(t.Tag.Marshal())
//	buffer.WriteByte(byte(t.KeyType))
//	buffer.WriteByte(byte(t.ValueType))
//	buffer.Write(compressInteger(t.Length)) // 需要实现的函数
//
//	for k, v := range t.Value {
//		var kBuffer []byte
//		switch t.KeyType {
//		case Integer:
//			kBuffer = compressInteger(k.(int)) // 需要实现的函数
//		case String:
//			kBuffer = writeString(k.(string)) // 需要实现的函数
//		default:
//			return nil, &NotImplemented{Type: string(t.KeyType)}
//		}
//
//		var vBuffer []byte
//		switch t.ValueType {
//		case Integer:
//			vBuffer = compressInteger(v.(int)) // 需要实现的函数
//		case String:
//			vBuffer = writeString(v.(string)) // 需要实现的函数
//		case Struct:
//			// 处理结构体类型的写入
//		default:
//			return nil, &NotImplemented{Type: string(t.ValueType)}
//		}
//		buffer.Write(kBuffer)
//		buffer.Write(vBuffer)
//	}
//
//	return buffer.Bytes()
//}

// Union 表示联合类型的 Node

// UnionRead 从流中读取一个联合类型的 Node
func UnionRead(label []byte, stream io.Reader) *Union {
	//unionType := ReadByte
	//if _ := io.ReadFull(stream, unionType); err != nil {
	//	return Node{}
	//}
	//value := UmMarshal(stream)
	//return Union{Node: Node{Label: string(label), Type: Union}, UnionType: unionType[0], Value: value}
	return nil
}

// Write 将联合类型的 Node 写入流
//func (t *Union) Write() []byte {
//	buffer := bytes.Buffer{}
//	buffer.Write(t.Tag.Marshal())
//	buffer.WriteByte(t.UnionType)
//	valueData := t.Value.Write()
//	buffer.Write(valueData)
//	return buffer.Bytes()
//}

// IntegerList 表示整数列表类型的 Node
type IntegerList struct {
	Value []int
	*bufio.ReadWriter
}

// IntegerListRead 从流中读取一个整数列表类型的 Node
//func IntegerListRead(label []byte, stream io.Reader) IntegerList {
//	length := decompressInteger(stream) // 需要实现的函数
//	values := make([]int, length)
//	for i := 0; i < length; i++ {
//		values[i] = decompressInteger(stream) // 需要实现的函数
//	}
//	return IntegerList{Node: Node{Label: string(label), Type: IntegerList}, Value: values}
//}
//
//// Write 将整数列表类型的 Node 写入流
//func (t *IntegerList) Write() []byte {
//	buffer := bytes.Buffer{}
//	buffer.Write(t.Tag.Marshal())
//	buffer.Write(compressInteger(len(t.Value))) // 需要实现的函数
//	for _, value := range t.Value {
//		buffer.Write(compressInteger(value)) // 需要实现的函数
//	}
//	return buffer.Bytes()
//}

// IntVector2 表示二维整数向量类型的 Node
//type IntVector2 struct {
//	Tag
//	Value []int
//	*bufio.ReadWriter
//}
//
//// IntVector2Read 从流中读取一个二维整数向量类型的 Node
//func IntVector2Read(label []byte, stream io.Reader) Node {
//	result := make([]int, 2)
//	for i := 0; i < 2; i++ {
//		value := decompressInteger(stream) // 需要实现的函数
//		result[i] = value
//	}
//	return IntVector2{Node: Node{Label: string(label), Type: IntVector2}, Value: result}
//}
//
//// Write 将二维整数向量类型的 Node 写入流
//func (t *IntVector2) Write() []byte {
//	buffer := bytes.Buffer{}
//	buffer.Write(t.Tag.Marshal())
//	buffer.Write(compressInteger(t.Value[0])) // 需要实现的函数
//	buffer.Write(compressInteger(t.Value[1])) // 需要实现的函数
//	return buffer.Bytes()
//}
//
//// IntVector3 表示三维整数向量类型的 Node
//type IntVector3 struct {
//	Tag
//	Value []int
//	*bufio.ReadWriter
//}
//
//// IntVector3Read 从流中读取一个三维整数向量类型的 Node
//func IntVector3Read(label []byte, stream io.Reader) Node {
//	result := make([]int, 3)
//	for i := 0; i < 3; i++ {
//		value := decompressInteger(stream) // 需要实现的函数
//		result[i] = value
//	}
//	return IntVector3{Node: Node{Label: string(label), Type: IntVector3}, Value: result}
//}
//
//// Write 将三维整数向量类型的 Node 写入流
//func (t *IntVector3) Write() []byte {
//	buffer := bytes.Buffer{}
//	buffer.Write(t.Tag.Marshal())
//	buffer.Write(compressInteger(t.Value[0])) // 需要实现的函数
//	buffer.Write(compressInteger(t.Value[1])) // 需要实现的函数
//	buffer.Write(compressInteger(t.Value[2])) // 需要实现的函数
//	return buffer.Bytes()
//}

func (n *Node) AddChild(child *Node) {
	child.parent = n
	n.Children = append(n.Children, child)
}

func (n *Node) Depth() unit.Dp {
	if n.parent != nil {
		return n.parent.Depth() + 1
	}
	return 1
}
func (n *Node) LenChildren() int {
	return len(n.Children)
}
func (n *Node) LastChild() (lastChild *Node) {
	//if n.IsRoot() {
	//	return n.Children[len(n.Children)-1]
	//}
	return n.parent.Children[len(n.parent.Children)-1]
}
func (n *Node) IsLastChild() bool {
	return n.LastChild() == n
}
func (n *Node) ResetChildren() {
	n.Children = nil
	n.filteredRows = nil
}

func (n *Node) Walk(callback func(node *Node)) {
	callback(n)
	for _, child := range n.Children {
		child.Walk(callback)
	}
}
func (n *Node) WalkQueue(callback func(node *Node)) {
	queue := []*Node{n}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		callback(node)
		for _, child := range node.Children {
			queue = append(queue, child)
		}
	}
}
