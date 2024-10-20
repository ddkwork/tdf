/**
 *  BlazeObjectType.go
 *
 *  (c) Electronic Arts. All Rights Reserved.
 *
 */

package tdf

import (
	"fmt"
	"strconv"
	"strings"
)

// BlazeObjectType represents a particular Blaze type in the context of its owner component.
// A BlazeObjectType consists of a component Id and an entity type which is unique to the component itself.
type BlazeObjectType struct {
	ComponentID int
	TypeID      int
}

// NewBlazeObjectType creates a new BlazeObjectType instance with default values.
func NewBlazeObjectType() *BlazeObjectType {
	return &BlazeObjectType{}
}

// NewBlazeObjectTypeWithValues creates a new BlazeObjectType instance with given values.
func NewBlazeObjectTypeWithValues(componentID, typeID int) *BlazeObjectType {
	return &BlazeObjectType{
		ComponentID: componentID,
		TypeID:      typeID,
	}
}

// NewBlazeObjectTypeFromOther creates a new BlazeObjectType instance from another instance.
func NewBlazeObjectTypeFromOther(other *BlazeObjectType) *BlazeObjectType {
	return &BlazeObjectType{
		ComponentID: other.ComponentID,
		TypeID:      other.TypeID,
	}
}

// ToLong returns the BlazeObjectType as a long value.
func (bot *BlazeObjectType) ToLong() int64 {
	return int64((bot.ComponentID<<16)+bot.TypeID) & 0xffffffff
}

// Equals returns true if the given BlazeObjectType is equal to this instance.
func (bot *BlazeObjectType) Equals(other *BlazeObjectType) bool {
	return other != nil && bot.ComponentID == other.ComponentID && bot.TypeID == other.TypeID
}

// HashCode returns a hash code value for this instance.
func (bot *BlazeObjectType) HashCode() int {
	return bot.ComponentID ^ bot.TypeID
}

// CompareTo compares this instance with another instance and returns -1, 0, or 1 if this instance is less than, equal to, or greater than the other instance respectively.
func (bot *BlazeObjectType) CompareTo(other *BlazeObjectType) int {
	if bot.ToLong() < other.ToLong() {
		return -1
	} else if bot.ToLong() > other.ToLong() {
		return 1
	} else {
		return 0
	}
}

// ToString returns a string representation of this instance using the specified separator.
func (bot *BlazeObjectType) ToString(separator rune) string {
	return fmt.Sprintf("%d%c%d", bot.ComponentID, separator, bot.TypeID)
}

// ToString returns a string representation of this instance using the default separator.
func (bot *BlazeObjectType) String() string {
	return bot.ToString('/')
}

// FromString creates a BlazeObjectType instance from the given string.
func BlazeObjectTypeFromString(str string) *BlazeObjectType {
	tokens := strings.Split(str, "/")
	if len(tokens) == 2 {
		componentID, err1 := strconv.Atoi(tokens[0])
		typeID, err2 := strconv.Atoi(tokens[1])
		if err1 == nil && err2 == nil {
			return NewBlazeObjectTypeWithValues(componentID, typeID)
		}
	}
	return nil
}
