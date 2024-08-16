package common

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"strconv"
	"strings"
)

// UID is method to generate an virtual unique identifier for whole system
// its structure contains 62 bits: LocalID - ObjectType - ShardID
// 32 bits for LocalID
// 10 bits for ObjectType
// 18 bits for ShardID

type UID struct {
	localID    uint32
	objectType int
	shardID    uint32
}

func NewUID(localUID uint32, objectType int, shardID uint32) UID {
	return UID{
		localID:    localUID,
		objectType: objectType,
		shardID:    shardID,
	}
}

func (uid UID) String() string {
	val := uint64(uid.localID)<<28 | uint64(uid.objectType)<<18 | uint64(uid.shardID)
	return base58.Encode([]byte(fmt.Sprintf("%v", val)))
}

func (uid UID) GetLocalID() uint32 {
	return uid.localID
}

func (uid UID) GetObjectType() int {
	return uid.objectType
}

func (uid UID) GetShardID() uint32 {
	return uid.shardID
}

func DecomposeUID(s string) (UID, error) {
	uid, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		return UID{}, err
	}

	if (1 << 18) > uid {
		return UID{}, errors.New("wrong uid")
	}
	u := UID{
		localID:    uint32(uid >> 28),
		objectType: int(uid >> 18 & 0x3FF),
		shardID:    uint32(uid & 0x3FFFF),
	}
	return u, nil
}

func FromBase58(s string) (UID, error) {
	return DecomposeUID(string(base58.Decode(s)))
}

func (uid *UID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", uid.String())), nil
}

func (uid *UID) UnmmarshalJSON(data []byte) error {
	decodedUID, err := FromBase58(strings.Replace(string(data), "\"", "", -1))

	if err != nil {
		return err
	}

	uid.localID = decodedUID.localID
	uid.shardID = decodedUID.shardID
	uid.objectType = decodedUID.objectType

	return nil
}

func (uid *UID) Value() (driver.Value, error) {
	if uid == nil {
		return nil, nil
	}
	return int64(uid.localID), nil
}

func (uid *UID) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var i uint32

	switch t := value.(type) {
	case int:
		i = uint32(t)
	case int8:
		i = uint32(t)
	case int16:
		i = uint32(t)
	case int32:
		i = uint32(t)
	case int64:
		i = uint32(t)
	case uint8:
		i = uint32(t)
	case uint16:
		i = uint32(t)
	case uint32:
		i = t
	case uint64:
		i = uint32(t)
	case []byte:
		a, err := strconv.ParseUint(string(t), 10, 64)
		if err != nil {
			return err
		}
		i = uint32(a)
	default:
		return errors.New("invalid scan source")
	}
	*uid = NewUID(i, 0, 1)
	return nil
}
