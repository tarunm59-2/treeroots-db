package utilities

import (
	"encoding/binary"
	"errors"
)

const (
	NODE_INTERNAL = 1
	NODE_LEAF     = 2

	HEADER       = 4
	PAGE_SIZE    = 4096
	MAX_KEY_SIZE = 1000
	MAX_VAL_SIZE = 3000
)

type TreeNode struct {
	data []byte
}

type Tree struct {
	root uint64
	get  func(uint64) TreeNode
	new  func(TreeNode) uint64
	del  func(uint64)
}

func (a TreeNode) NodeType() uint16 {
	return binary.LittleEndian.Uint16(a.data[0:2])
}
func (a TreeNode) KeyCount() uint16 {
	return binary.LittleEndian.Uint16(a.data[2:4])
}
func (a TreeNode) SetHead(nodetype uint16, keycount uint16) {
	binary.LittleEndian.PutUint16((a.data[0:2]), nodetype)
	binary.LittleEndian.PutUint16((a.data[2:4]), keycount)
}
func (a TreeNode) GetNthPointer(id uint16) (ptr uint64, e error) {
	if id >= a.KeyCount() {
		return 0, errors.New("not enough keys")
	}
	position_to_get := HEADER + 8*id
	ptr = binary.LittleEndian.Uint64(a.data[position_to_get:])
	return ptr, nil
}
func (a TreeNode) SetNthPointer(id uint16, value uint64) (prt uint64, e error) {
	if id >= a.KeyCount() {
		return 0, errors.New("not enough keys in the set nth pointer function")
	}
	position_to_set := HEADER + 8*id
	binary.LittleEndian.PutUint64(a.data[position_to_set:], value)
	return value, nil
}
func (a TreeNode) GetNthOffset(id uint16) (prt uint16, e error) {
	if id == 0 {
		return 0, nil
	}
	if id > a.KeyCount() || id < 0 {
		return 0, errors.New("Not enough keys, in the get nth offset function")
	}

	kv_pairoffset_position := HEADER + 8*a.KeyCount() + 2*(id-1)
	return binary.LittleEndian.Uint16(a.data[kv_pairoffset_position:]), nil
}
func (a TreeNode) SetNthOffset(id uint16, offset_value uint16) (prt uint16, e error) {
	if id == 0 {
		return 0, nil
	}
	if id > a.KeyCount() || id < 0 {
		return 0, errors.New("Not enough keys, in the set nth offset function")
	}

	kv_pairoffset_position := HEADER + 8*a.KeyCount() + 2*(id-1)
	binary.LittleEndian.PutUint16(a.data[kv_pairoffset_position:], offset_value)
	return offset_value, nil
}

func (a TreeNode) GetNthKvPos(id uint16) (pos uint16, e error) {
	if id < 0 || id > a.KeyCount() {
		return 0, errors.New("not enough keys at getnthkvpos")
	}
	offset, _ := a.GetNthOffset(id)
	ret := HEADER + 10*(a.KeyCount()) + offset
	return ret, nil

}

func (a TreeNode) RetrieveNthKey(id uint16) (keybytes []byte, e error) {
	if id < 0 || id > a.KeyCount() {
		return nil, errors.New("not enough keys at getnthkvpos")
	}
	key_start_pos, _ := a.GetNthKvPos(id)
	key_length := binary.LittleEndian.Uint16(a.data[key_start_pos:])
	key_origin := a.data[key_start_pos+4:]
	retkey := key_origin[:key_length]
	return retkey, nil
}
func (a TreeNode) RetrieveNthValue(id uint16) (keybytes []byte, e error) {
	if id < 0 || id > a.KeyCount() {
		return nil, errors.New("not enough keys at getnthkvpos")
	}
	key_start_pos, _ := a.GetNthKvPos(id)
	key_length := binary.LittleEndian.Uint16(a.data[key_start_pos:])
	value_length := binary.LittleEndian.Uint16(a.data[key_start_pos+2:])
	value_origin := a.data[key_start_pos+4+key_length:]
	retval := value_origin[:value_length]
	return retval, nil
}
