package hessianserialparse

import (
	"fmt"

	"github.com/fairyming/hessian_serial_parse/utils"
)

type MapKeyValue struct {
	Key   interface{}
	Value interface{}
}

type HessianMap struct {
	ClassName string
	Maps      []MapKeyValue
}

func NewHessianMap(className string) *HessianMap {
	return &HessianMap{ClassName: className, Maps: make([]MapKeyValue, 0)}
}

type HessianList struct {
	ListType string
	Args     []interface{}
}

func NewHessianList(listType string) *HessianList {
	return &HessianList{ListType: listType, Args: make([]interface{}, 0)}
}

type HessianParse struct {
	reader *utils.DataReader
}

func NewHessianParse(data []byte) (*HessianParse, error) {
	reader, err := utils.NewDataReader(data)
	if err != nil {
		return nil, err
	}
	return &HessianParse{reader: reader}, nil
}

func (h *HessianParse) parseInt() (int, error) {
	if len, err := h.reader.ReadU32Be(); err == nil {
		return int(len), nil
	} else {
		return 0, err
	}
}

func (h *HessianParse) parseLong() (int64, error) {
	tmp, err := h.reader.Read(8)
	if err != nil {
		return 0, err
	}
	return int64(tmp[0])<<56 + int64(tmp[1])<<48 + int64(tmp[2])<<40 + int64(tmp[3])<<32 + int64(tmp[4])<<24 + int64(tmp[5])<<16 + int64(tmp[6])<<8 + int64(tmp[7]), nil
}

func (h HessianParse) parseChar(length int) (string, error) {
	var result string
	for i := 0; i < length; i++ {
		ch, err := h.reader.Read1()
		if err != nil {
			return "", err
		}
		if ch < 0x80 {
			result += string(ch)
		} else if (ch & 0xe0) == 0xc0 {
			ch1, err := h.reader.Read1()
			if err != nil {
				return "", err
			}
			result += string(((ch & 0x1f) << 6) + (ch1 & 0x3f))
		} else if (ch & 0xf0) == 0xe0 {
			ch1, err := h.reader.Read1()
			if err != nil {
				return "", err
			}
			ch2, err := h.reader.Read1()
			if err != nil {
				return "", err
			}
			result += string(rune((int(ch&0x0f) << 12) + int((ch1&0x3f)<<6) + int(ch2&0x3f)))
		}
	}
	return result, nil
}

func (h HessianParse) readChunkLength() (int, error) {
	if len, err := h.reader.ReadU16Be(); err == nil {
		return int(len), nil
	} else {
		return 0, err
	}
}

func (h HessianParse) readLength() (int, error) {
	code, err := h.reader.Read1()
	if err != nil {
		return -1, err
	}
	if code != 'l' {
		return -1, nil
	}
	if r, err := h.parseInt(); err == nil {
		return r, nil
	} else {
		return -1, err
	}
}

func (h HessianParse) readType() (string, error) {
	chr, err := h.reader.Read1()
	if err != nil {
		return "", err
	}
	if chr != 't' {
		h.reader.Back(1)
		return "", nil
	}

	chunkLength, err := h.readChunkLength()
	if err != nil {
		return "", err
	}
	return h.parseChar(chunkLength)
}

func (h HessianParse) Parse() (interface{}, error) {
	chr, err := h.reader.Read1()
	fmt.Println(chr)
	// fmt.Println(h.reader.Position())
	if err != nil {
		return nil, err
	}
	switch chr {
	case HS_NULL_TAG:
		return nil, nil
	case HS_TRUE_TAG:
		return true, nil
	case HS_FALSE_TAG:
		return false, nil
	case HS_INT_TAG:
		return h.parseInt()
	case HS_LONG_TAG, HS_DOUBLE_TAG, HS_DATE_TAG:
		return h.parseLong()
	case HS_STRING_TAG:
		chunkLen, err := h.readChunkLength()
		if err != nil {
			return nil, err
		}
		result, err := h.parseChar(chunkLen)
		if err != nil {
			return nil, err
		}
		nextChunk, err := h.Parse()
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("%v%v", result, nextChunk), nil
	case HS_STRING_LAST_CHUNK_TAG:
		chunkLen, err := h.readChunkLength()
		if err != nil {
			return nil, err
		}
		return h.parseChar(chunkLen)
	case HS_BYTE_TAG:
		return nil, nil
	case HS_BYTE_LAST_CHUNK_TAG:
		return nil, nil
	case HS_LIST_TAG:
		hsType, err := h.readType()
		if err != nil {
			return nil, err
		}
		length, err := h.readLength()
		if err != nil {
			return nil, err
		}
		hsList := NewHessianList(hsType)
		for i := 0; i < length; i++ {
			arg, err := h.Parse()
			fmt.Println(arg)
			if err != nil {
				return nil, err
			}
			hsList.Args = append(hsList.Args, arg)
		}
		chr, err = h.reader.Read1()
		if err != nil {
			return nil, err
		}
		if chr == HS_END {
			return hsList, nil
		} else {
			return hsList, fmt.Errorf("list not end")
		}
	case HS_MAP_TAG:
		hsType, err := h.readType()
		if err != nil {
			return nil, err
		}
		hsMap := NewHessianMap(hsType)
		for {
			chr, err = h.reader.Read1()
			if err != nil {
				return nil, err
			}
			if chr == HS_END {
				break
			} else {
				h.reader.Back(1)
			}

			key, err := h.Parse()
			if err != nil {
				return nil, err
			}
			value, err := h.Parse()
			if err != nil {
				return nil, err
			}
			hsMap.Maps = append(hsMap.Maps, MapKeyValue{Key: key, Value: value})
		}
		return hsMap, nil
	case HS_ELEMENT_DATA_TAG:
		return nil, nil
	case HS_RESOLVE_REMOTE_TAG:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown tag %v", chr)
	}
}
