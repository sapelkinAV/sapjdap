// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jdbg

import (
	"reflect"
	"sapelkinav/javadap/jdwp/jdwpclient"
)

// marshal transforms o into a client.Value.
func (j *JDbg) marshal(o interface{}) jdwpclient.Value {
	switch o := o.(type) {
	case bool, jdwpclient.Char, int, int8,
		int16, int32, int64, float32, float64:
		return o

	case nil:
		return jdwpclient.ObjectID(0)

	case jdwpclient.Object:
		return o.ID()

	case Value:
		return j.marshal(o.val)

	case string:
		id, err := j.conn.CreateString(o)
		if err != nil {
			j.fail("Failed to marshal string: %v", err)
		}
		return id

	case []byte:
		return j.newArray(j.cache.byteTy, o)

	case []interface{}:
		return j.newArray(j.cache.objTy, o)

	default:
		j.fail("Unhandled type %T", o)
		return nil
	}
}

// newArray creates a new array with element type elTy, filled with values.
func (j *JDbg) newArray(elTy Type, values interface{}) jdwpclient.ArrayID {
	array := j.ArrayOf(elTy).New(reflect.ValueOf(values).Len())
	array.SetArrayValues(values)
	return array.val.(jdwpclient.ArrayID)
}

// toObjectType returns the corresponding java.lang.Object type for ty, or nil
// if there is no corresponding object type.
func (j *JDbg) toObjectType(ty reflect.Type) *Class {
	switch ty.Kind() {
	case reflect.Int:
		return j.cache.intObjTy
	case reflect.Int8, reflect.Uint8:
		return j.cache.byteObjTy
	case reflect.Int16:
		return j.cache.shortObjTy
	case reflect.Int32:
		return j.cache.intObjTy
	case reflect.Int64:
		return j.cache.longObjTy
	case reflect.Float32:
		return j.cache.floatObjTy
	case reflect.Float64:
		return j.cache.doubleObjTy
	default:
		return nil
	}
}

// toObject returns the value of o transformed to a java.lang.Object held in a
// client.Value.
func (j *JDbg) toObject(o interface{}) jdwpclient.Value {
	if obj := j.toObjectType(reflect.TypeOf(o)); obj != nil {
		return j.marshal(obj.New(o))
	}

	switch o := j.marshal(o).(type) {
	case jdwpclient.ObjectID, jdwpclient.StringID, jdwpclient.ArrayID:
		return o
	case jdwpclient.TaggedObjectID:
		return o.Object
	default:
		j.fail("Cannot convert %v (%T) to Object", o, o)
		return nil
	}
}

func (j *JDbg) toObjects(l []interface{}) []interface{} {
	objects := make([]interface{}, len(l))
	for i, v := range l {
		objects[i] = j.toObject(v)
	}
	return objects
}

// unmarshal unboxes the client.Value into a corresponding golang value.
func (j *JDbg) unmarshal(v jdwpclient.Value) interface{} {
	switch v := v.(type) {
	case jdwpclient.StringID:
		str, err := j.conn.GetString(v)
		if err != nil {
			j.fail("Failed to unmarshal string")
		}
		return str
	default:
		return v
	}
}

func (j *JDbg) marshalN(v []interface{}) []jdwpclient.Value {
	out := make([]jdwpclient.Value, len(v))
	for i, v := range v {
		out[i] = j.marshal(v)
	}
	return out
}

func (j *JDbg) unmarshalN(v []jdwpclient.Value) []interface{} {
	out := make([]interface{}, len(v))
	for i, v := range v {
		out[i] = j.unmarshal(v)
	}
	return out
}
