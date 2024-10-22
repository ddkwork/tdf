package tdf

import (
	"encoding/binary"
	"io"
	"math"
	. "reflect"
	"slices"
	"strings"
	"time"

	"github.com/ddkwork/app/widget"
	"github.com/ddkwork/encoding/struct2table"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

func Marshal(message any) (buf *stream.Buffer) {
	buf = stream.NewBuffer("")
	root := struct2table.New(message)
	marshalStruct(buf, root)
	return
}

func Unmarshal(buf []byte) (message *widget.Node[struct2table.StructField]) {
	b := stream.NewBuffer(buf)
	message = widget.NewRoot(struct2table.StructField{})
	unmarshalStruct(b, message)
	return
}

func marshalStruct(b *stream.Buffer, parent *widget.Node[struct2table.StructField]) {
	b.Append(encodeTagAndWireType(parent.Data.Tag, StructType))
	defer b.WriteByte(ID_TERM)
	for _, child := range parent.Children {
		switch child.Data.Value.Kind() {
		case Bool:
			b.Append(marshalSingular(child.Data.Tag, child.Data.Value.Bool()))
		case Int, Int8, Int16, Int32, Int64:
			b.Append(marshalSingular(child.Data.Tag, child.Data.Value.Int()))
		case Uint, Uint8, Uint16, Uint32, Uint64:
			b.Append(marshalSingular(child.Data.Tag, child.Data.Value.Uint()))
		case Float32, Float64:
			b.Append(marshalSingular(child.Data.Tag, child.Data.Value.Float()))
		case String:
			b.Append(marshalSingular(child.Data.Tag, child.Data.Value.String()))
		case Slice, Array:
			if child.Data.Value.Elem().Kind() == Int8 { // todo remove listType ?
				b.Append(marshalSingular(child.Data.Tag, child.Data.Value.Bytes()))
				continue
			}
			marshalList(b, child)
		case Map:
			marshalMap(b, child)
		case Struct:
			marshalStruct(b, child)
		default:
			switch child.Data.ValueAssert.(type) {
			case []byte:
				marshalSingular(child.Data.Tag, child.Data.Value.Bytes())
			case *Union:
				marshalSingular(child.Data.Tag, child.Data.ValueAssert.(*Union))
			case *Variable:
				marshalSingular(child.Data.Tag, child.Data.ValueAssert.(*Variable))
			case *BlazeObjectType:
				marshalSingular(child.Data.Tag, child.Data.ValueAssert.(*BlazeObjectType))
			case *BlazeObjectId:
				marshalSingular(child.Data.Tag, child.Data.ValueAssert.(*BlazeObjectId))
			case time.Time:
				b.Append(marshalSingular(child.Data.Tag, child.Data.ValueAssert.(time.Time).UnixNano()))
			}
			mylog.Check("unsupported type")
		}
	}
}

func unmarshalStruct(b *stream.Buffer, parent *widget.Node[struct2table.StructField]) {
	tag, wireType := decodeTagAndWireType(b)
	switch wireType {
	// case BoolType:
	//	parent.Data.Value.SetBool(unmarshalSingular(b, tag).Bool())
	// case IntType:
	//	parent.Data.Value.SetInt(unmarshalSingular(b, tag).Int())
	// case UintType:
	//	parent.Data.Value.SetUint(unmarshalSingular(b, tag).Uint())
	// case FloatType:
	//	parent.Data.Value.SetFloat(unmarshalSingular(b, tag).Float())
	// case StringType:
	//	parent.Data.Value.SetString(unmarshalSingular(b, tag).String())
	// case BinaryType:
	//	parent.Data.Value.SetBytes(unmarshalSingular(b, tag).Bytes())
	case ListType:
		unmarshalList(b, parent)
	case MapType:
		unmarshalMap(b, parent)
	case StructType:
		container := widget.NewContainerNode(parent.Data.Name, struct2table.StructField{})
		parent.AddChild(container)
		unmarshalStruct(b, container)
	default:
		parent.AddChildByData(struct2table.StructField{
			Number:      0,
			Name:        "",
			Type:        nil,
			Kind:        0,
			Tag:         StructTag(tag),
			Value:       Value{},
			ValueAssert: nil,
		})
	}
}

func marshalList(b *stream.Buffer, parent *widget.Node[struct2table.StructField]) {
	if parent.Data.Value.Len() == 0 {
		return
	}
	for i := range parent.Data.Value.Len() {
		elemType := parent.Data.Value.Index(i).Type()
		//	elemSize := elemType.Size()
		elemValue := parent.Data.Value.Index(i).Interface()
		elemValueOf := ValueOf(elemValue)
		switch elemType.Kind() {
		case Int, Int8, Int16, Int32, Int64:
			b.Append(marshalSingular(parent.Data.Tag, elemValueOf.Int()))
		case Uint, Uint8, Uint16, Uint32, Uint64:
			b.Append(marshalSingular(parent.Data.Tag, elemValueOf.Uint()))
		case Float32, Float64:
			b.Append(marshalSingular(parent.Data.Tag, elemValueOf.Float()))
		case String:
			b.Append(marshalSingular(parent.Data.Tag, elemValueOf.String()))
		case Slice, Array:
			if elemType.Elem().Kind() == Int8 { // todo remove listType ?
				b.Append(marshalSingular(parent.Data.Tag, elemValueOf.Bytes()))
				continue
			}
			marshalList(b, parent.Children[i])
		case Map:
			marshalMap(b, parent.Children[i])
		case Struct:
			marshalStruct(b, parent.Children[i])
		default:
			switch parent.Data.ValueAssert.(type) {
			case []byte:
				marshalSingular(parent.Data.Tag, elemValueOf.Bytes())
			case *Union:
				b.Append(marshalSingular(parent.Data.Tag, elemValue.(*Union)))
			case *Variable:
				b.Append(marshalSingular(parent.Data.Tag, elemValue.(*Variable)))
			case *BlazeObjectType:
				b.Append(marshalSingular(parent.Data.Tag, elemValue.(*BlazeObjectType)))
			case *BlazeObjectId:
				b.Append(marshalSingular(parent.Data.Tag, elemValue.(*BlazeObjectId)))
			case time.Time:
				b.Append(marshalSingular(parent.Data.Tag, elemValue.(time.Time).UnixNano()))
			}
			mylog.Check("unsupported type")
		}
	}
}

func unmarshalList(b *stream.Buffer, parent *widget.Node[struct2table.StructField]) {
	listType := BaseType(mylog.Check2(b.ReadByte()))
	if listType != ListType {
		mylog.Check("invalid list type") // todo skipElement
	}
	length := decompressInteger(b)
	if length == 0 {
		return
	}
	/* todo 不要解码tag和类型
	boolean tmpDecodeHeader = mDecodeHeader;
			mDecodeHeader = false;
			memberHelper.initializeVector(value, mTempValue.get().intValue());
			memberHelper.visitMembers(this, rootTdf, parentTdf, tag, value, referenceValue);
			mDecodeHeader = tmpDecodeHeader;
	*/
}

func marshalMap(b *stream.Buffer, parent *widget.Node[struct2table.StructField]) {
	value := parent.Data.Value
	if value.Len() == 0 {
		return
	}
	//keys := value.MapKeys()
	//for i, key := range keys {
	//	mapElemValue := value.MapIndex(key)
	//
	//}
	//todo
	//if e.encodeType(mapHelper.GetKeyType()) &&
	//	e.encodeType(mapHelper.GetValueType()) &&
	//	e.encodeVarsizeInteger(int64(len(*value))) {
	//	tmpEncodeHeader := e.mEncodeHeader //elem不需要编解码tag和wireType
	//	e.mEncodeHeader = false
	//	mapHelper.VisitMembers(e, rootTdf, parentTdf, tag, value, referenceValue)
	//	e.mEncodeHeader = tmpEncodeHeader
	//}

	//	for k, v := range t.Value {
	//		var kBuffer []byte
	//		switch t.KeyType {
	//		case Integer:
	//			kBuffer = compressInteger(k.(int))
	//		case String:
	//			kBuffer = writeString(k.(string))
	//		default:
	//			return nil, &NotImplemented{Type: string(t.KeyType)}
	//		}
	//
	//		var vBuffer []byte
	//		switch t.ValueType {
	//		case Integer:
	//			vBuffer = compressInteger(v.(int))
	//		case String:
	//			vBuffer = writeString(v.(string))
	//		case Struct:
	//			// 处理结构体类型的写入
	//		default:
	//			return nil, &NotImplemented{Type: string(t.ValueType)}
	//		}
	//		buffer.Write(kBuffer)
	//		buffer.Write(vBuffer)
	//	}
}

func unmarshalMap(b *stream.Buffer, parent *widget.Node[struct2table.StructField]) {
	//func DictionaryRead(label []byte, stream io.Reader) Node {
	//	keyType := decompressInteger(stream)
	//	valueType := decompressInteger(stream)
	//	length := decompressInteger(stream)
	//	value := make(Dictionary)
	//	for i := 0; i < length; i++ {
	//		var key, value any
	//		switch BaseKind(keyType) {
	//		case Integer:
	//			key = decompressInteger(stream)
	//		case String:
	//			key = readString(stream)
	//		default:
	//			return Node{}, &NotImplemented{Type: string(keyType)}
	//		}
	//
	//		switch BaseKind(valueType) {
	//		case Integer:
	//			value = decompressInteger(stream)
	//		case String:
	//			value = readString(stream)
	//		case Struct:
	//			value = readStruct(stream)
	//		default:
	//			return Node{}, &NotImplemented{Type: string(valueType)}
	//		}
	//		value[key] = value
	//	}
	//	return Dictionary{Node: Node{Label: string(label), Type: Dictionary}, KeyType: BaseKind(keyType), ValueType: BaseKind(valueType), Length: length, Value: value}
	//}
	//
}

func encodeTagAndWireType[T string | StructTag](tag T, wireType BaseType) (b *stream.Buffer) {
	b = stream.NewBuffer("")
	b.Write(EncodeTag(string(tag)))
	b.WriteByte(byte(wireType))
	return
}

func decodeTagAndWireType(b *stream.Buffer) (tag string, wireType BaseType) {
	tagBuf := make([]byte, 3)
	mylog.Check2(b.Read(tagBuf))
	tag = DecodeTag(tagBuf)
	typeBuf := mylog.Check2(b.ReadByte())
	wireType = BaseType(typeBuf)
	return
}

func DecodeTag(tag []byte) string {
	if len(tag) != 3 {
		mylog.Check("tagBuf length must be 3")
	}
	tag = slices.Concat([]byte{0}, tag)
	encodedTag := binary.BigEndian.Uint32(tag)
	decodedTag := make([]byte, 4)
	for i := 3; i >= 0; i-- {
		sixbits := encodedTag & 0x3F
		if sixbits != 0 {
			decodedTag[i] = byte(sixbits + 32)
		} else {
			decodedTag[i] = ' '
		}
		encodedTag >>= 6
	}
	return string(decodedTag)
}

func EncodeTag(tag string) []byte {
	if len(tag) != 4 {
		mylog.Check("tag must be 4 characters long")
	}

	tag = strings.ToUpper(tag)
	var encodedTag uint32 = 0

	for i := 0; i < len(tag); i++ {
		char := tag[i]
		if char == ' ' {
			continue
		}
		encodedTag |= (uint32(char) - 32) << (6 * (3 - i))
	}

	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, encodedTag)
	return result[1:]
}

func decompressInteger(b *stream.Buffer) uint32 {
	var result uint64 = 0
	var currentShift uint64 = 6
	buffer := make([]byte, 1)
	mylog.Check2(b.Read(buffer))
	result += uint64(buffer[0]) & 0x3F
	for buffer[0]&0x80 != 0 {
		mylog.Check2(b.Read(buffer))
		// todo test  Check2 是否可以工作

		result |= (uint64(buffer[0]) & 0x7F) << currentShift
		currentShift += 7
	}
	return uint32(result)
}

func compressInteger[T int64 | uint64](value T) []byte {
	var result []byte
	if value < 0x40 {
		result = append(result, byte(value))
	} else {
		currentByte := (value & 0x3F) | 0x80
		result = append(result, byte(currentByte))
		currentShift := value >> 6
		for currentShift >= 0x80 {
			currentByte = (currentShift & 0x7F) | 0x80
			currentShift >>= 7
			result = append(result, byte(currentByte))
		}
		result = append(result, byte(currentShift))
	}
	return result
}

// 任何协议的编解码我们只处理指定的结构体字段类型，map，切片，结构体等复合类型的字段也只支持这些单一里类型
type singularType interface {
	~bool |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~string |
		~[]byte | // 其余类型的切片一般情况下不会存在二维字节切片的字段，所以把一维字节切片视为单一类型
		time.Time |
		*BlazeObjectType | *BlazeObjectId |
		*Union | *Variable | *Enum
}

// SingularAssert 解码无法让泛型有用武之地,为了避免取值类型不匹配，取值之前需要判断baseType是否匹配
type SingularAssert struct {
	data any
}

func NewSingularAssert(data any) *SingularAssert {
	return &SingularAssert{data: data}
}

func (s *SingularAssert) String() string                   { return s.data.(string) }
func (s *SingularAssert) Int() int                         { return s.data.(int) }
func (s *SingularAssert) Int8() int8                       { return s.data.(int8) }
func (s *SingularAssert) Int16() int16                     { return s.data.(int16) }
func (s *SingularAssert) Int32() int32                     { return s.data.(int32) }
func (s *SingularAssert) Int64() int64                     { return s.data.(int64) }
func (s *SingularAssert) Uint() uint                       { return s.data.(uint) }
func (s *SingularAssert) Uint8() uint8                     { return s.data.(uint8) }
func (s *SingularAssert) Uint16() uint16                   { return s.data.(uint16) }
func (s *SingularAssert) Uint32() uint32                   { return s.data.(uint32) }
func (s *SingularAssert) Uint64() uint64                   { return s.data.(uint64) }
func (s *SingularAssert) Float32() float32                 { return s.data.(float32) }
func (s *SingularAssert) Float64() float64                 { return s.data.(float64) }
func (s *SingularAssert) Bytes() []byte                    { return s.data.([]byte) }
func (s *SingularAssert) Time() time.Time                  { return s.data.(time.Time) }
func (s *SingularAssert) BlazeObjectId() BlazeObjectId     { return s.data.(BlazeObjectId) }
func (s *SingularAssert) BlazeObjectType() BlazeObjectType { return s.data.(BlazeObjectType) }
func (s *SingularAssert) TimeValue() time.Time             { return s.data.(time.Time) }

// 更安全的做法是这里先解码tag和类型并传入类型，然后改成泛型函数，类型就安全了
// 在构造树节点的时候传入类型返回对应类型的值，然后填充节点元数据?
// 构造树类型安全之后，用户取值还是需要按类型取，感觉又回到原点了
// 测试一下root取节点元数据如何传入类型?
// 解码返回的只是一个树形，无法预先知道每个节点的元数据类型，因为我们不使用协议辅助文件系统

// 解码二进制预先是不知道类型的，所以这里不限制类型才合理？部分返回值需要断言看来是无法避免的了,或者需要一个统一的断言函数来确保解码安全
func unmarshalSingular(buf []byte) (tag string, wireType BaseType, data any) {
	b := stream.NewBuffer(buf)
	tag, wireType = decodeTagAndWireType(b)
	switch wireType {
	case IntegerType: ///很明显后期按类型取值的时候这里分不清是32位还是64位
		data = decompressInteger(b) // todo 如何确定是bool类型,不使用idt文件我看是确认不了类型的
	case StringType:
		length := decompressInteger(b)
		result := make([]byte, length-1)
		mylog.Check2(io.ReadFull(b, result))
		mylog.Check2(b.ReadByte())
		if length == 1 {
			data = ""
			return
		}
		data = string(result)
	case BinaryType:
		length := decompressInteger(b)
		metadataBuf := make([]byte, length)
		mylog.Check2(io.ReadFull(b, metadataBuf))
		data = metadataBuf
	case StructType:
	case ListType:
	case MapType: // 参考java，需要peek一下类型，然后再解码
	case UnionType:
		data = mylog.Check2(b.ReadByte())
	case VariableType:
	case BlazeObjectTypeType:
	case BlazeObjectIdType:
	case FloatType:
	case TimeValueType:
	case MaxType:
	case InvalidBaseType:
	}
	return
}

func marshalSingular[T singularType, TagType string | StructTag](tag TagType, value T) (b *stream.Buffer) {
	b = encodeTagAndWireType(tag, NativeTypeBind[TypeOf(value)])
	switch v := any(value).(type) {
	case bool:
		if v {
			b.WriteByte(1)
			return
		}
		b.WriteByte(0)
	case int, int8, int16, int32, int64:
		mylog.Check2(b.Write(compressInteger(ValueOf(v).Int())))
	case uint, uint8, uint16, uint32, uint64:
		mylog.Check2(b.Write(compressInteger(ValueOf(v).Uint())))
	case float32:
		mBuf := make([]byte, FLOAT_SIZE)
		binary.LittleEndian.PutUint32(mBuf, math.Float32bits(v))
		b.Write(slices.Concat(mBuf, []byte{0}))
	case float64:
		mBuf := make([]byte, FLOAT_SIZE)
		binary.LittleEndian.PutUint64(mBuf, math.Float64bits(v))
		b.Write(slices.Concat(mBuf, []byte{0}))
	case string:
		b.Write(slices.Concat(compressInteger(uint64(len(v)+1)), []byte(v), []byte{0}))
	case []byte: // blob
		b.Write(slices.Concat(compressInteger(uint64(len(v))), v))
	case Union: // todo add enum ?
		b.WriteByte(v.activeMember)
	case Variable: // VariableTdfContainer
		if v.valid {
			b.WriteByte(1)
		} else {
			b.WriteByte(0)
		}
		if v.valid {
			b.Write(compressInteger(v.id))
			// todo how to visit value?
		}
		b.WriteByte(ID_TERM)
	case BlazeObjectType:
		//	//if e.encodeVarsizeInteger(int64(value.GetComponentId())) {//todo
		//	//	e.encodeVarsizeInteger(int64(value.GetTypeId()))
		//	//}
		b.Write(slices.Concat(compressInteger(uint64(v.ComponentID)), compressInteger(uint64(v.TypeID))))
	case BlazeObjectId:
		//	//if e.encodeVarsizeInteger(int64(value.GetType().GetComponentId())) &&//todo
		//	//	e.encodeVarsizeInteger(int64(value.GetType().GetTypeId())) {
		//	//	e.encodeVarsizeInteger(value.GetEntityId())
		//	//}
		b.Write(slices.Concat(compressInteger(uint64(v.GetType().ComponentID)), compressInteger(uint64(v.GetType().TypeID)), compressInteger(v.GetEntityId())))
	case time.Time:
		// b.Write(compressInteger(uint64(v.Second())))
		b.Write(compressInteger(uint64(v.UnixMicro())))
	}
	return
}

func (b BaseType) Valid() bool { return b >= 0 && b < MaxType }

type (
	Enum  struct{}
	Union struct {
		activeMember byte
	}
	Variable struct {
		value int64
		id    int64
		valid bool
	}
)

var NativeTypeBind = map[Type]BaseType{
	TypeFor[bool]():             IntegerType,
	TypeFor[int8]():             IntegerType,
	TypeFor[int16]():            IntegerType,
	TypeFor[int32]():            IntegerType,
	TypeFor[int64]():            IntegerType,
	TypeFor[uint8]():            IntegerType,
	TypeFor[uint16]():           IntegerType,
	TypeFor[uint32]():           IntegerType,
	TypeFor[uint64]():           IntegerType,
	TypeFor[string]():           StringType,
	TypeFor[[]byte]():           BinaryType,   // blob
	TypeFor[struct{}]():         StructType,   // java显示tdf就是结构体类型,ID_TERM 0结尾
	TypeFor[[]any]():            ListType,     // java显示方法名称为Vector,类型，大小，遍历
	TypeFor[map[any]any]():      MapType,      // k v type,size,遍历
	TypeFor[*Union]():           UnionType,    // todo mock enum
	TypeFor[*Variable]():        VariableType, // VariableTdfContainer ID_TERM 0结尾
	TypeFor[float32]():          FloatType,    // floatToIntBits
	TypeFor[float64]():          FloatType,
	TypeFor[*BlazeObjectType](): BlazeObjectTypeType, // getComponentId and getTypeId  整型编解码,难道是树形的层级下标和孩子节点下标?
	TypeFor[*BlazeObjectId]():   BlazeObjectIdType,   // getComponentId , getTypeId and getEntityId 整型编解码
	TypeFor[time.Time]():        TimeValueType,
}
