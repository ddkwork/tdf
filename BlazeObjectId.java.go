package tdf

import (
	"strconv"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

// BlazeObjectId refers to a particular Blaze object in the context of it's owner component.
// A BlazeObjectId consists of a BlazeObjectType and a unique value identifying the instance
// of the given type.
type BlazeObjectId struct {
	// Represents an invalid object id.
	BLAZE_OBJECT_ID_INVALID *BlazeObjectId
	// The entity id.
	mEntityId int64
	// The entity type.
	mType BlazeObjectType
}

// Instantations a BlazeObjectId with all values set to the provided default.
func NewBlazeObjectId_mEntityId(defaultValues int64) *BlazeObjectId {
	blazeObjectId := &BlazeObjectId{}
	// blazeObjectId.mType = *NewBlazeObjectType(int(defaultValues), int(defaultValues))
	blazeObjectId.mEntityId = defaultValues
	return blazeObjectId
}

// Instantiates a new blaze object id.
func NewBlazeObjectId() *BlazeObjectId {
	blazeObjectId := &BlazeObjectId{}
	blazeObjectId.mType = *NewBlazeObjectType()
	blazeObjectId.mEntityId = 0
	return blazeObjectId
}

// Instantiates a new blaze object id.
func NewBlazeObjectIdFrom(rhs *BlazeObjectId) *BlazeObjectId {
	blazeObjectId := &BlazeObjectId{}
	// blazeObjectId.mType = *NewBlazeObjectTypeFrom(rhs.mType)
	blazeObjectId.mEntityId = rhs.mEntityId
	return blazeObjectId
}

// Instantiates a new blaze object id.
func NewBlazeObjectIdWith(componentId, entityType int, entityId int64) *BlazeObjectId {
	blazeObjectId := &BlazeObjectId{}
	// blazeObjectId.mType = *NewBlazeObjectType(componentId, entityType)
	blazeObjectId.mEntityId = entityId
	return blazeObjectId
}

// Instantiates a new blaze object id.
func NewBlazeObjectIdWithType(t *BlazeObjectType, entityId int64) *BlazeObjectId {
	blazeObjectId := &BlazeObjectId{}
	blazeObjectId.mEntityId = entityId
	// blazeObjectId.mType = *NewBlazeObjectTypeFrom(t)
	return blazeObjectId
}

// Sets the.
func (blazeObjectId *BlazeObjectId) Set(rhs *BlazeObjectId) {
	blazeObjectId.mEntityId = rhs.mEntityId
	// blazeObjectId.mType.Set(rhs.GetType())
}

// Sets the entity id.
func (blazeObjectId *BlazeObjectId) SetEntityId(entityId int64) {
	blazeObjectId.mEntityId = entityId
}

// Gets the entity id.
func (blazeObjectId *BlazeObjectId) GetEntityId() int64 {
	return blazeObjectId.mEntityId
}

// Sets the entity type.
func (blazeObjectId *BlazeObjectId) SetType(rhs *BlazeObjectType) {
	// blazeObjectId.mType.Set(rhs)
}

// Gets the entity type.
func (blazeObjectId *BlazeObjectId) GetType() *BlazeObjectType {
	return &blazeObjectId.mType
}

// Equals.
func (blazeObjectId *BlazeObjectId) Equals(rhs interface{}) bool {
	if rhs, ok := rhs.(*BlazeObjectId); ok {
		return blazeObjectId.mEntityId == rhs.mEntityId && blazeObjectId.mType.Equals(&rhs.mType)
	}
	return false
}

// Returns a hash code value for the object.
func (blazeObjectId *BlazeObjectId) HashCode() int {
	return int(blazeObjectId.mEntityId^(blazeObjectId.mEntityId>>32)) ^ blazeObjectId.mType.HashCode()
}

// To string.
func (blazeObjectId *BlazeObjectId) ToString(separator rune) string {
	return blazeObjectId.mType.ToString(separator) + string(separator) + strconv.FormatInt(blazeObjectId.mEntityId, 10)
}

// From string.
func BlazeObjectFromString(s string) *BlazeObjectId {
	tok := strings.Split(s, "/")
	if len(tok) == 3 {
		compIdStr := tok[0]
		typeIdStr := tok[1]
		entityIdStr := tok[2]
		componentId, err1 := strconv.Atoi(compIdStr)
		entityId, err2 := strconv.ParseInt(entityIdStr, 10, 64)
		if err1 == nil && err2 == nil {
			typeId := mylog.Check2(strconv.Atoi(typeIdStr))
			return NewBlazeObjectIdWith(componentId, typeId, entityId)
		}
	}
	return nil
}
