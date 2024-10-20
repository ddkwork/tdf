package tdf

import (
	"bytes"
	"encoding/binary"
	"github.com/ddkwork/app/widget"
	"github.com/ddkwork/encoding/struct2table"
	"github.com/ddkwork/golibrary/mylog"
	"io"
	"math/big"
)

type Heat2Decoder struct {
	mByteBuffer   *bytes.Buffer
	mTempValue    *big.Int
	mBuf          []byte
	mDecodeHeader bool
	mErrorCount   int

	tag         string
	wireType    BaseType
	metadataBuf []byte
	root        *widget.Node[struct2table.StructField]
}

func NewHeat2Decoder(b []byte) (d *Heat2Decoder) {
	d = &Heat2Decoder{
		mByteBuffer:   bytes.NewBuffer(b),
		mTempValue:    big.NewInt(0),
		mBuf:          make([]byte, HEADER_SIZE),
		mDecodeHeader: true,
		mErrorCount:   0,
		tag:           "",
		wireType:      0,
		metadataBuf:   nil,
	}
	d.tag, d.wireType = decodeTagAndWireType(d.mByteBuffer)
	switch d.wireType {
	case IntegerType:
	case StringType:
		length := decompressInteger(d.mByteBuffer)
		result := make([]byte, length-1)
		mylog.Check2(io.ReadFull(d.mByteBuffer, result))
		mylog.Check2(d.mByteBuffer.ReadByte())
		if length == 1 {
			d.metadataBuf = []byte("")
			return
		}
		d.metadataBuf = result
	case BinaryType: //DecodeBlob
		length := decompressInteger(d.mByteBuffer)
		d.metadataBuf = make([]byte, length)
		mylog.Check2(io.ReadFull(d.mByteBuffer, d.metadataBuf))
	case StructType:
	case ListType:
	case MapType:
	case UnionType:
	case VariableType:
	case BlazeObjectTypeType:
	case BlazeObjectIdType:
	case FloatType:
	case TimeValueType:
	case MaxType:
	case InvalidBaseType:
	}
	return d
}

func (d *Heat2Decoder) readHeader(tag int64, tdfType BaseType) bool {
	//    public final boolean hasRemaining() {
	//        return position < limit;
	//    }
	//if d.w == nil || !d.w.HasRemaining() {
	//	return false
	//}
	for d.mByteBuffer.Len() >= HEADER_SIZE {
		mylog.Check(
			// d.w.Mark()
			mylog.Check2(d.mByteBuffer.Read(d.mBuf)))

		if d.mBuf[0] == ID_TERM {
			d.mByteBuffer.Reset() // todo test
			return false
		}
		bufTag := int64(binary.BigEndian.Uint32(d.mBuf[:3])) << 8
		bufType := int(d.mBuf[HEADER_TYPE_OFFSET])

		if bufType >= int(MaxType) {
			d.mErrorCount++
			return false
		}
		if bufTag == tag {
			if bufType != int(tdfType) {
				d.mErrorCount++
				return false
			}
			return true
		}
		if bufTag > tag {
			d.mByteBuffer.Reset()
			return false
		}
		if !d.SkipElement(BaseType(bufType)) {
			d.mErrorCount++
			return false
		}
	}
	return false
}

func (d *Heat2Decoder) getStructTerminator() bool {
	if d.mByteBuffer == nil {
		return false
	}
	for d.mByteBuffer.Len() >= 1 {
		mylog.Check(mylog.Check2(d.mByteBuffer.Read(d.mBuf)))

		if d.mBuf[0] == ID_TERM {
			if d.mByteBuffer.Len() < 1 {
				return false
			}
			mylog.Check(mylog.Check2(d.mByteBuffer.Read(d.mBuf[:1])))

			return true
		}
		bufType := BaseType(d.mBuf[HEADER_TYPE_OFFSET])
		if bufType.Valid() { // todo
			if !d.SkipElement(bufType) {
				break
			}
		} else {
			d.mErrorCount++
			break
		}
	}
	return false
}

func (d *Heat2Decoder) SkipElement(bufType BaseType) bool {
	//if d.mByteBuffer == nil {
	//	return false
	//}
	//var rc bool
	//
	//switch bufType {
	//case TDF_TYPE_INTEGER:
	//	if !d.DecodeVarsizeInteger(d.mTempValue) {
	//		rc = false
	//	}
	//case TDF_TYPE_STRING, TDF_TYPE_BINARY:
	//	if !d.DecodeVarsizeInteger(d.mTempValue) || d.mTempValue.Int64() < 0 || int64(d.mByteBuffer.Len()-d.mByteBuffer.Position()) < d.mTempValue.Int64() {
	//		rc = false
	//	} else {
	//		d.mByteBuffer.Seek(d.mTempValue.Int64(), 1)
	//	}
	//case TDF_TYPE_STRUCT:
	//	rc = d.getStructTerminator()
	//case TDF_TYPE_LIST:
	//	if d.mByteBuffer.Len()-d.mByteBuffer.Position() < 1 {
	//		rc = false
	//	} else {
	//		mylog.Check(binary.Read(d.mByteBuffer, binary.BigEndian, d.mBuf[:1]))
	//		listType := BaseType(d.mBuf[0])
	//
	//		if !listType.Valid() {
	//			rc = false
	//		} else {
	//			if !d.DecodeVarsizeInteger(d.mTempValue) || d.mTempValue.Int64() < 0 {
	//				rc = false
	//			} else {
	//				listMemberCount := d.mTempValue.Int64()
	//				for idx := int64(0); idx < listMemberCount && rc; idx++ {
	//					rc = d.SkipElement(listType)
	//				}
	//			}
	//		}
	//	}
	//case TDF_TYPE_MAP:
	//	if int64(d.mByteBuffer.Len()-d.mByteBuffer.Position()) < 2 {
	//		rc = false
	//	} else {
	//		mylog.Check(binary.Read(d.mByteBuffer, binary.BigEndian, d.mBuf[:2]))
	//		keyType := BaseType(d.mBuf[0])
	//		valueType := BaseType(d.mBuf[1])
	//
	//		if keyType.Valid() && valueType.Valid() {
	//			if !d.DecodeVarsizeInteger(d.mTempValue) || d.mTempValue.Int64() < 0 {
	//				rc = false
	//			} else {
	//				mapLength := d.mTempValue.Int64()
	//				for idx := int64(0); idx < mapLength && rc; idx++ {
	//					rc = d.SkipElement(keyType)
	//					if rc {
	//						rc = d.SkipElement(valueType)
	//					}
	//				}
	//			}
	//		}
	//	}
	//default:
	//	rc = false
	//}

	return false
}

func (d *Heat2Decoder) DecodeVarsizeInteger(value *big.Int) bool {
	len := d.mByteBuffer.Len()
	if len == 0 {
		return false
	}
	buf := make([]byte, 1)
	mylog.Check(binary.Read(d.mByteBuffer, binary.BigEndian, buf))

	valueIsNegative := (buf[0] & VARSIZE_NEGATIVE) == VARSIZE_NEGATIVE
	hasMore := (buf[0] & VARSIZE_MORE) == VARSIZE_MORE
	v := big.NewInt(int64(buf[0] & (VARSIZE_NEGATIVE - 1)))
	if hasMore {
		shift := 6
		hasMore = false
		var bufVal *big.Int
		for d.mByteBuffer.Len() > 0 {
			mylog.Check(binary.Read(d.mByteBuffer, binary.BigEndian, buf))

			bufVal = big.NewInt(int64(buf[0]))
			bufVal.And(bufVal, big.NewInt(int64(VARSIZE_MORE-1)))
			bufVal.Lsh(bufVal, uint(shift))
			v.Or(v, bufVal)
			hasMore = (buf[0] & VARSIZE_MORE) != 0
			if !hasMore {
				break
			}
			shift += 7
		}
		if hasMore {
			d.mErrorCount++
			value.Set(big.NewInt(0))
			return false
		}
	}
	if valueIsNegative {
		value.Set(v.Neg(v))
	} else {
		value.Set(v)
	}
	return true
}

func (d *Heat2Decoder) DecodeNextHeaderAndInteger(tag int64, value *big.Int, defaultValue *big.Int) {
	value.Set(defaultValue)
	if !d.mDecodeHeader || d.readHeader(tag, IntegerType) {
		d.DecodeVarsizeInteger(value)
	}
}

func (d *Heat2Decoder) PeekHeader(tag *int64, ttype *int32) bool {
	if d.mByteBuffer == nil || d.mByteBuffer.Len() < HEADER_SIZE {
		return false
	}
	//mBuf := make([]byte, HEADER_SIZE)
	//mylog.Check(binary.Read(d.mByteBuffer, binary.BigEndian, mBuf))
	//if mBuf[0] == ID_TERM {
	//	return false
	//}
	//*tag = int64((BaseType(mBuf[0]) << 24) | (BaseType(mBuf[1]) << 16) | (BaseType(mBuf[2]) << 8))
	//*ttype = int32(BaseType(mBuf[HEADER_TYPE_OFFSET]))
	//d.mByteBuffer.Seek(-1*HEADER_SIZE, io.SeekCurrent)
	return true
}

func (d *Heat2Decoder) PeekListType(ttype *int32) bool {
	if d.mByteBuffer == nil || d.mByteBuffer.Len() < (HEADER_SIZE+1) {
		return false
	}
	//d.mByteBuffer.Seek(HEADER_SIZE, io.SeekCurrent)
	//mylog.Check(binary.Read(d.mByteBuffer, binary.BigEndian, ttype))
	//
	//d.mByteBuffer.Seek(-1*(HEADER_SIZE+1), io.SeekCurrent)
	return true
}

func (d *Heat2Decoder) PeekMapType(keyType, valueType *int32) bool {
	if d.mByteBuffer == nil || d.mByteBuffer.Len() < (HEADER_SIZE+2) {
		return false
	}
	//d.mByteBuffer.Seek(HEADER_SIZE, io.SeekCurrent)
	//mylog.Check(binary.Read(d.mByteBuffer, binary.BigEndian, keyType))
	//
	//mylog.Check(binary.Read(d.mByteBuffer, binary.BigEndian, valueType))
	//
	//d.mByteBuffer.Seek(-1*(HEADER_SIZE+2), io.SeekCurrent)
	return true
}

func (d *Heat2Decoder) Reset() {
	d.mErrorCount = 0
	d.mDecodeHeader = true
}

func (d *Heat2Decoder) DecodeString(value *string, defaultValue string) {
	if !d.mDecodeHeader || d.readHeader(0, StringType) {
		d.DecodeVarsizeInteger(d.mTempValue)
		if d.mTempValue.Int64() > int64(d.mByteBuffer.Len()) {
			d.mErrorCount++
			*value = ""
			return
		}
		buf := make([]byte, d.mTempValue.Int64())
		d.mByteBuffer.Read(buf)
		*value = string(buf)
	}
}
