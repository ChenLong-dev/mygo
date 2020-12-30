/*
 * @Description:
 * @Author: Chen Long
 * @Date: 2020-12-16 14:48:13
 * @LastEditTime: 2020-12-16 14:48:13
 * @LastEditors: Chen Long
 * @Reference:
 */

package mmgo

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-mgo/mgo"
	"reflect"
	"sort"
	"strings"
)

func AddBsonField(m bson.M, path string, val interface{}, cover bool) (bson.M, error) {
	if m == nil {
		m = bson.M{}
	}
	var im interface{}
	im = m
	var bm bson.M
	var ok bool
	keys := strings.Split(path, ".")
	nKeys := len(keys)
	if nKeys == 0 {
		return m, fmt.Errorf("no specify path")
	}
	for i, key := range keys[:nKeys-1] {
		if len(key) == 0 {
			return m, fmt.Errorf("illegal path")
		}
		if key[0] == '[' {
			return m, fmt.Errorf("not support slice") //暂时不支持切片
		} else {
			bm, ok = im.(bson.M)
			if !ok {
				bm, ok = im.(map[string]interface{})
			}
			if !ok {
				return m, errors.New(fmt.Sprintf("key(%d:%s) not bson.M in %s", i, key, JsonPrintPretty(im)))
			}
			if im, ok = bm[key]; !ok {
				im = bson.M{}
				bm[key] = im
			}
		}
	}
	lastKey := keys[nKeys-1]
	if len(lastKey) == 0 {
		return m, fmt.Errorf("illegal path")
	}
	bm, ok = im.(bson.M)
	if !ok {
		bm, ok = im.(map[string]interface{})
	}
	if !ok {
		return m, errors.New(fmt.Sprintf("key(%s) not bson.M in %s", lastKey, JsonPrintPretty(im)))
	}
	if im, ok = bm[lastKey]; !ok {
		bm[lastKey] = val
	} else {
		if !cover {
			return m, fmt.Errorf("%s already exist but not set cover", path)
		}
		bm[lastKey] = val
	}

	return m, nil
}
func DelBsonField(m bson.M, path string) {
	pm := m
	lastKey := path
	lastDot := strings.LastIndex(path, ".")
	if lastDot >= 0 {
		v, verr := GetBsonField(m, path[:lastDot])
		if verr != nil {
			return
		}
		var ok bool
		if pm, ok = v.(bson.M); !ok {
			if pm, ok = v.(map[string]interface{}); !ok {
				return
				//return fmt.Errorf("path(%s) not bson.M value", path[:lastDot])
			}
		}
		lastKey = path[lastDot+1:]
	}

	delete(pm, lastKey)
	return
}
func GetBsonField(m bson.M, path string) (v interface{}, err error) {
	if m == nil || path == "" {
		return nil, errors.New(fmt.Sprintf("m or path(%s) is nil", path))
	}

	v = m
	var mv bson.M
	var ok bool
	keys := strings.Split(path, ".")
	for i, key := range keys {
		if len(key) == 0 {
			continue
		}
		if key[0] == '[' {
			val := reflect.ValueOf(v)
			if val.Kind() != reflect.Slice {
				return nil, errors.New(fmt.Sprintf("key(%d:%s): not op slice v=%v", i, key, JsonPrintPretty(v)))
			}
			i := ParseUint64(key[1:len(key)-1], 0)
			if val.Len() <= int(i) {
				return nil, errors.New(fmt.Sprintf("key(%d:%s): out bound slice v=%v", i, key, JsonPrintPretty(v)))
			}
			v = val.Index(int(i)).Interface()
		} else {
			mv, ok = v.(bson.M)
			if !ok {
				mv, ok = v.(map[string]interface{})
			}
			if !ok {
				return nil, errors.New(fmt.Sprintf("key(%d:%s) not bson.M in %s", i, key, JsonPrintPretty(v)))
			}
			if v, ok = mv[key]; !ok {
				return nil, errors.New(fmt.Sprintf("key(%d:%s) not exist in %s", i, key, JsonPrintPretty(v)))
			}
		}
	}
	return v, nil
}

func GetBsonFieldMultiName(m bson.M, paths ...string) (v interface{}, err error) {
	for _, path := range paths {
		v, err = GetBsonField(m, path)
		if err == nil {
			return v, err
		}
	}
	if err == nil {
		return v, errors.New(fmt.Sprintf("paths[%v] not find!", paths))
	}
	return v, err
}
func GetBsonM(m bson.M, path string) bson.M {
	if v, err := GetBsonField(m, path); err == nil {
		if r, ok := v.(bson.M); ok {
			return r
		}
	}
	return nil
}
func GetBsonInt(m bson.M, path string, defs ...int) int {
	def := 0
	if len(defs) > 0 {
		def = defs[0]
	}
	if v, err := GetBsonField(m, path); err == nil {
		return int(ToInt64(v, int64(def)))
	}
	return def
}
func GetBsonInt64(m bson.M, path string, defs ...int64) int64 {
	def := int64(0)
	if len(defs) > 0 {
		def = defs[0]
	}
	if v, err := GetBsonField(m, path); err == nil {
		return ToInt64(v, def)
	}
	return def
}
func GetBsonUint64(m bson.M, path string, defs ...uint64) uint64 {
	def := uint64(0)
	if len(defs) > 0 {
		def = defs[0]
	}
	if v, err := GetBsonField(m, path); err == nil {
		return uint64(ToInt64(v, int64(def)))
	}
	return def
}
func GetBsonFloat64(m bson.M, path string, defs ...float64) float64 {
	def := float64(0.0)
	if len(defs) > 0 {
		def = defs[0]
	}
	if v, err := GetBsonField(m, path); err == nil {
		return ToFloat64(v, def)
	}
	return def
}
func GetBsonString(m bson.M, path string, defs ...string) string {
	def := ""
	if len(defs) > 0 {
		def = defs[0]
	}
	if v, err := GetBsonField(m, path); err == nil {
		return ToString(v, def)
	}
	return def
}
func GetBsonBool(m bson.M, path string, defs ...bool) bool {
	def := false
	if len(defs) > 0 {
		def = defs[0]
	}
	if v, err := GetBsonField(m, path); err == nil {
		if r, ok := v.(bool); ok {
			return r
		}
		ri := ToInt64(v, 0)
		if ri != 0 {
			return true
		}
	}
	return def
}

func bsonM2D(bs bson.M, prefix string) (rd bson.D) {
	for key, val := range bs {
		name := prefix + key

		if subBs, ok := val.(bson.M); !ok {
			rd = append(rd, bson.DocElem{Name: name, Value: val})
		} else {
			if len(subBs) == 0 {
				rd = append(rd, bson.DocElem{Name: name, Value: val})
			} else {
				srd := bsonM2D(subBs, name+".")
				rd = append(rd, srd...)
				//snames := walkBsonName(subBs, name+".")
				//fields = append(fields, snames...)
			}
		}
	}
	return rd
}

func CompareBsonM(lm bson.M, rm bson.M) bool {
	ld := bsonM2D(lm, "")
	sort.Slice(ld, func(i, j int) bool {
		return ld[i].Name < ld[j].Name
	})
	ldata, _ := bson.Marshal(ld)

	rd := bsonM2D(rm, "")
	sort.Slice(rd, func(i, j int) bool {
		return rd[i].Name < rd[j].Name
	})
	rdata, _ := bson.Marshal(rd)

	return len(ldata) == len(rdata) && bytes.Compare(ldata, rdata) == 0
}

func GetAllBsonFieldName(bs bson.M, prefix string) (fields []string) {
	for key, val := range bs {
		name := prefix + key

		if subBs, ok := val.(bson.M); !ok {
			fields = append(fields, name)
		} else {
			if len(subBs) == 0 {
				fields = append(fields, name)
			} else {
				snames := GetAllBsonFieldName(subBs, name+".")
				fields = append(fields, snames...)
			}
		}
	}
	return fields
}
func checkFilterBson(filter bson.M) (need bool, err error) {
	hasNeed := false
	hasNotNeed := false
	var subFilter bson.M
	var ok bool
	for _, val := range filter {
		subFilter, ok = val.(bson.M)
		if !ok {
			subFilter, ok = val.(map[string]interface{})
		}
		if !ok {
			if ToInt64(val, 0) != 0 {
				hasNeed = true
			} else {
				hasNotNeed = true
			}
		} else {
			need, err = checkFilterBson(subFilter)
			if err != nil {
				return need, err
			}
			if need {
				hasNeed = true
			} else {
				hasNotNeed = true
			}
		}
	}
	if hasNeed && hasNotNeed {
		return false, fmt.Errorf("cannot have a mix of inclusion and exclusion")
	}
	return hasNeed, nil
}
func filtBsonM(bs bson.M, need bool, filterFields []string) (rbs bson.M) {
	if need {
		for _, name := range filterFields {
			if field, err := GetBsonField(bs, name); err == nil {
				rbs, _ = AddBsonField(rbs, name, field, true)
			}
		}
		return rbs
	} else {
		for _, name := range filterFields {
			DelBsonField(bs, name)
		}
		return bs
	}
}
func FiltBsonM(m bson.M, filter bson.M) (bson.M, error) {
	if len(filter) == 0 {
		return m, nil
	}
	need, err := checkFilterBson(filter)
	if err != nil {
		return m, err
	}
	filterFields := GetAllBsonFieldName(filter, "")
	return filtBsonM(m, need, filterFields), nil
}
