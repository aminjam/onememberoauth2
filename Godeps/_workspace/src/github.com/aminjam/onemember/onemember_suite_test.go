package onemember_test

import (
	"errors"
	"reflect"

	"github.com/aminjam/onemember"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var db onemember.DataConnector

type memoryDB struct {
	m   map[string]*onemember.Account
	seq int
}

func TestOnemember(t *testing.T) {
	RegisterFailHandler(Fail)
	db = &memoryDB{
		m:   make(map[string]*onemember.Account),
		seq: 0,
	}
	RunSpecs(t, "Onemember Suite")
}

func (db *memoryDB) OnememberCreate(in *onemember.Account) error {
	db.m[in.Username] = in
	return nil
}
func (db *memoryDB) OnememberRead(username string) (map[string]interface{}, error) {
	var item = db.m[username]
	if item == nil {
		return nil, errors.New("Not Found")
	}
	out := structToMap(item)
	return out, nil
}
func (db *memoryDB) OnememberUpdate(in *onemember.Account) error {
	db.m[in.Username] = in
	return nil
}

func structToMap(i interface{}) (values map[string]interface{}) {
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	values = make(map[string]interface{}, iVal.NumField())
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		// You ca use tags here...
		// tag := typ.Field(i).Tag.Get("tagname")
		// Convert each type into a string for the url.Values string map
		var v interface{}
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = f.Int()
		case uint, uint8, uint16, uint32, uint64:
			v = f.Uint()
		case float32:
			v = f.Float()
		case float64:
			v = f.Float()
		case []byte:
			v = f.Bytes()
		case string:
			v = f.String()
		}
		values[typ.Field(i).Name] = v
	}
	return
}
