package simpleserialize

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

// convert data struct into []byte
func MarshalStruct(v interface{}) ([]byte, error) {
	obj := reflect.ValueOf(v)

	var out bytes.Buffer

	for i := 0; i < obj.NumField(); i += 1 {
		val := obj.Field(i)

		switch val.Kind() {
		case reflect.Ptr:
			if val.IsNil() {
				// ignore
			} else {
				data, err := valToBytes(val.Elem())
				if err != nil {
					return []byte{}, err
				}

				out.Write([]byte{byte(uint8(i))})  // Field Num
				out.Write([]byte{byte(len(data))}) // Field Size
				out.Write(data)                    // Field data itself
			}
		default:
			panic("non-pointer type, please check structure field types")
		}

	}

	return out.Bytes(), nil
}

func UnMarshalStruct(v interface{}, data []byte) error {
	dataLen := len(data)

	obj := reflect.ValueOf(v).Elem()
	typeOfStruct := obj.Type()

	cursor := 0
	for cursor < dataLen {
		fieldNum := uint8(data[cursor])
		dataSize := uint8(data[cursor+1])

		fieldData := data[cursor+2 : cursor+2+int(dataSize)]

		t := typeOfStruct.Field(int(fieldNum)).Type

		val := obj.Field(int(fieldNum))

		value, err := bytesToVal(t, fieldData)

		if err != nil {
			return err
		}

		val.Set(value)

		cursor += 2 + int(dataSize)
	}
	return nil
}

func valToBytes(val reflect.Value) ([]byte, error) {
	switch val.Kind() {
	case reflect.Int, reflect.Uint:
		panic("please always specify exact size, like int8 or uint32")

	case reflect.Slice:
		return val.Bytes(), nil

	case reflect.String:
		return []byte(val.String()), nil

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf := new(bytes.Buffer)
		var err error

		// TODO: refactor somehow
		switch val.Kind() {
		case reflect.Int8:
			err = binary.Write(buf, binary.LittleEndian, int8(val.Int()))
		case reflect.Int16:
			err = binary.Write(buf, binary.LittleEndian, int16(val.Int()))
		case reflect.Int32:
			err = binary.Write(buf, binary.LittleEndian, int32(val.Int()))
		case reflect.Int64:
			err = binary.Write(buf, binary.LittleEndian, val.Int())
		case reflect.Uint8:
			err = binary.Write(buf, binary.LittleEndian, uint8(val.Uint()))
		case reflect.Uint16:
			err = binary.Write(buf, binary.LittleEndian, uint16(val.Uint()))
		case reflect.Uint32:
			err = binary.Write(buf, binary.LittleEndian, uint32(val.Uint()))
		case reflect.Uint64:
			err = binary.Write(buf, binary.LittleEndian, val.Uint())
		default:
			panic("bug in the code")
		}
		return buf.Bytes(), err

	case reflect.Bool:
		value := val.Bool()
		if value {
			return []byte{1}, nil
		} else {
			return []byte{0}, nil
		}
	default:
		panic("Unsupported type, please check structure field types")
	}

	var data []byte
	data = make([]byte, 5, 5)

	return data, nil
}

func bytesToVal(val reflect.Type, data []byte) (reflect.Value, error) {
	var err error

	switch val.Elem().Kind() {
	case reflect.Int, reflect.Uint:
		panic("please always specify exact size, like int8 or uint32")

	case reflect.Slice:
		return reflect.ValueOf(&data), nil

	case reflect.String:
		str := string(data)
		return reflect.ValueOf(&str), nil

	case reflect.Bool:
		var value bool
		if data[0] == 1 {
			value = true
		} else if data[0] == 0 {
			value = false
		} else {
			err = errors.New("can not decode input data")
		}
		return reflect.ValueOf(&value), err

	// TODO: refactor somehow
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf := bytes.NewReader(data)
		switch val.Elem().Kind() {
		case reflect.Int8:
			var n int8
			err = binary.Read(buf, binary.LittleEndian, &n)
			return reflect.ValueOf(&n), err
		case reflect.Int16:
			var n int16
			err = binary.Read(buf, binary.LittleEndian, &n)
			return reflect.ValueOf(&n), err
		case reflect.Int32:
			var n int32
			err = binary.Read(buf, binary.LittleEndian, &n)
			return reflect.ValueOf(&n), err
		case reflect.Int64:
			var n int64
			err = binary.Read(buf, binary.LittleEndian, &n)
			return reflect.ValueOf(&n), err
		case reflect.Uint8:
			var n uint8
			err = binary.Read(buf, binary.LittleEndian, &n)
			return reflect.ValueOf(&n), err
		case reflect.Uint16:
			var n uint16
			err = binary.Read(buf, binary.LittleEndian, &n)
			return reflect.ValueOf(&n), err
		case reflect.Uint32:
			var n uint32
			err = binary.Read(buf, binary.LittleEndian, &n)
			return reflect.ValueOf(&n), err
		case reflect.Uint64:
			var n uint64
			err = binary.Read(buf, binary.LittleEndian, &n)
			return reflect.ValueOf(&n), err
		default:
			panic("bug in the code")
		}

	default:
		return reflect.ValueOf(nil), errors.New("unknown field type")
	}
}
