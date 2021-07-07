package uuidnil

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestUnmarshalArray(t *testing.T) {
	// Without UUID.
	var array1 struct {
		Array [2]string
	}
	data1 := []byte(`{"array": ["string1", "string2"]}`)
	if err := json.Unmarshal(data1, Wrap(&array1, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := array1.Array[0]; v != "string1" {
		t.Errorf("invalid value: %v", v)
	}
	if v := array1.Array[1]; v != "string2" {
		t.Errorf("invalid value: %v", v)
	}

	// With empty UUID.
	var array2 struct {
		Array [2]uuid.UUID
	}
	data2 := []byte(`{"array": ["", "ed9a09dc-dd17-11eb-a1ae-305a3a7ae79b"]}`)
	if err := json.Unmarshal(data2, Wrap(&array2, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := array2.Array[0]; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
	if v := array2.Array[1]; v != uuid.MustParse("ed9a09dc-dd17-11eb-a1ae-305a3a7ae79b") {
		t.Errorf("invalid value: %v", v)
	}

	// With malformed UUID.
	var array3 struct {
		Array [2]uuid.UUID
	}
	data3 := []byte(`{"array": ["test", "ed9a09dc-dd17-11eb-a1ae-305a3a7ae79b"]}`)
	if err := json.Unmarshal(data3, Wrap(&array3, AllowInvalid)); err != nil {
		t.Fatal(err)
	}
	if v := array3.Array[0]; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
	if v := array3.Array[1]; v != uuid.MustParse("ed9a09dc-dd17-11eb-a1ae-305a3a7ae79b") {
		t.Errorf("invalid value: %v", v)
	}

	// With malformed UUID.
	var array4 struct {
		Array [2]uuid.UUID
	}
	data4 := []byte(`{"array": ["test", "ed9a09dc-dd17-11eb-a1ae-305a3a7ae79b"]}`)
	if err := json.Unmarshal(data4, Wrap(&array4, AllowEmpty)); err == nil {
		t.Fatal("expected test to fail")
	}
}

func TestUnmarshalMap(t *testing.T) {
	// Without UUID.
	var map1 struct {
		Map map[string]string
	}
	data1 := []byte(`{"map": {"key1": "string1", "key2": "string2"}}`)
	if err := json.Unmarshal(data1, Wrap(&map1, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := map1.Map["key1"]; v != "string1" {
		t.Errorf("invalid value: %v", v)
	}
	if v := map1.Map["key2"]; v != "string2" {
		t.Errorf("invalid value: %v", v)
	}

	// With empty key UUID.
	var map2 struct {
		Map map[uuid.UUID]string
	}
	data2 := []byte(`{"map": {"": "string1", "ad04dd68-dd65-11eb-bef6-305a3a7ae79b": "string2"}}`)
	if err := json.Unmarshal(data2, Wrap(&map2, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := map2.Map[uuid.Nil]; v != "string1" {
		t.Errorf("invalid value: %v", v)
	}
	if v := map2.Map[uuid.MustParse("ad04dd68-dd65-11eb-bef6-305a3a7ae79b")]; v != "string2" {
		t.Errorf("invalid value: %v", v)
	}

	// With empty value UUID.
	var map3 struct {
		Map map[string]uuid.UUID
	}
	data3 := []byte(`{"map": {"key1": "", "key2": "ad04dd68-dd65-11eb-bef6-305a3a7ae79b"}}`)
	if err := json.Unmarshal(data3, Wrap(&map3, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := map3.Map["key1"]; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
	if v := map3.Map["key2"]; v != uuid.MustParse("ad04dd68-dd65-11eb-bef6-305a3a7ae79b") {
		t.Errorf("invalid value: %v", v)
	}

	// With malformed key UUID.
	var map4 struct {
		Map map[uuid.UUID]string
	}
	data4 := []byte(`{"map": {"test": "string1", "ad04dd68-dd65-11eb-bef6-305a3a7ae79b": "string2"}}`)
	if err := json.Unmarshal(data4, Wrap(&map4, AllowInvalid)); err != nil {
		t.Fatal(err)
	}
	if v := map4.Map[uuid.Nil]; v != "string1" {
		t.Errorf("invalid value: %v", v)
	}
	if v := map4.Map[uuid.MustParse("ad04dd68-dd65-11eb-bef6-305a3a7ae79b")]; v != "string2" {
		t.Errorf("invalid value: %v", v)
	}

	// With malformed value UUID.
	var map5 struct {
		Map map[string]uuid.UUID
	}
	data5 := []byte(`{"map": {"key1": "test", "key2": "ad04dd68-dd65-11eb-bef6-305a3a7ae79b"}}`)
	if err := json.Unmarshal(data5, Wrap(&map5, AllowInvalid)); err != nil {
		t.Fatal(err)
	}
	if v := map5.Map["key1"]; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
	if v := map5.Map["key2"]; v != uuid.MustParse("ad04dd68-dd65-11eb-bef6-305a3a7ae79b") {
		t.Errorf("invalid value: %v", v)
	}

	// With malformed value UUID.
	var map6 struct {
		Map map[uuid.UUID]string
	}
	data6 := []byte(`{"map": {"test": "string1", "ad04dd68-dd65-11eb-bef6-305a3a7ae79b": "string2"}}`)
	if err := json.Unmarshal(data6, Wrap(&map6, AllowEmpty)); err == nil {
		t.Fatal("expected test to fail")
	}

	// With malformed key UUID.
	var map7 struct {
		Map map[string]uuid.UUID
	}
	data7 := []byte(`{"map": {"key1": "test", "key2": "ad04dd68-dd65-11eb-bef6-305a3a7ae79b"}}`)
	if err := json.Unmarshal(data7, Wrap(&map7, AllowEmpty)); err == nil {
		t.Fatal("expected test to fail")
	}
}

func TestUnmarshalPtr(t *testing.T) {
	// Without UUID.
	var ptr1 struct {
		Ptr *string
	}
	data1 := []byte(`{"ptr": "pointer"}`)
	if err := json.Unmarshal(data1, Wrap(&ptr1, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := *ptr1.Ptr; v != "pointer" {
		t.Errorf("invalid value: %v", v)
	}

	// With empty UUID.
	var ptr2 struct {
		Ptr *uuid.UUID
	}
	data2 := []byte(`{"ptr": ""}`)
	if err := json.Unmarshal(data2, Wrap(&ptr2, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := *ptr2.Ptr; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}

	// With malformed UUID.
	var ptr3 struct {
		Ptr *uuid.UUID
	}
	data3 := []byte(`{"ptr": "test"}`)
	if err := json.Unmarshal(data3, Wrap(&ptr3, AllowInvalid)); err != nil {
		t.Fatal(err)
	}
	if v := *ptr3.Ptr; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}

	// With malformed UUID.
	var ptr4 struct {
		Ptr *uuid.UUID
	}
	data4 := []byte(`{"ptr": "test"}`)
	if err := json.Unmarshal(data4, Wrap(&ptr4, AllowEmpty)); err == nil {
		t.Fatal("expected test to fail")
	}
}

func TestUnmarshalSlice(t *testing.T) {
	// Without UUID.
	var slice1 struct {
		Slice []string
	}
	data1 := []byte(`{"slice": ["string1", "string2"]}`)
	if err := json.Unmarshal(data1, Wrap(&slice1, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := slice1.Slice[0]; v != "string1" {
		t.Errorf("invalid value: %v", v)
	}
	if v := slice1.Slice[1]; v != "string2" {
		t.Errorf("invalid value: %v", v)
	}

	// With empty UUID.
	var slice2 struct {
		Slice []uuid.UUID
	}
	data2 := []byte(`{"slice": ["", "14720916-dd67-11eb-b5da-305a3a7ae79b"]}`)
	if err := json.Unmarshal(data2, Wrap(&slice2, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := slice2.Slice[0]; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
	if v := slice2.Slice[1]; v != uuid.MustParse("14720916-dd67-11eb-b5da-305a3a7ae79b") {
		t.Errorf("invalid value: %v", v)
	}

	// With malformed UUID.
	var slice3 struct {
		Slice []uuid.UUID
	}
	data3 := []byte(`{"slice": ["test", "14720916-dd67-11eb-b5da-305a3a7ae79b"]}`)
	if err := json.Unmarshal(data3, Wrap(&slice3, AllowInvalid)); err != nil {
		t.Fatal(err)
	}
	if v := slice3.Slice[0]; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
	if v := slice3.Slice[1]; v != uuid.MustParse("14720916-dd67-11eb-b5da-305a3a7ae79b") {
		t.Errorf("invalid value: %v", v)
	}

	// With malformed UUID.
	var slice4 struct {
		Slice []uuid.UUID
	}
	data4 := []byte(`{"slice": ["test", "14720916-dd67-11eb-b5da-305a3a7ae79b"]}`)
	if err := json.Unmarshal(data4, Wrap(&slice4, AllowEmpty)); err == nil {
		t.Fatal("expected test to fail")
	}
}

func TestUnmarshalStruct(t *testing.T) {
	// Without UUID.
	var struct1 struct {
		Struct struct {
			Key1 string
			Key2 string
		}
	}
	data1 := []byte(`{"struct": {"key1": "string1", "key2": "string2"}}`)
	if err := json.Unmarshal(data1, Wrap(&struct1, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := struct1.Struct.Key1; v != "string1" {
		t.Errorf("invalid value: %v", v)
	}
	if v := struct1.Struct.Key2; v != "string2" {
		t.Errorf("invalid value: %v", v)
	}

	// With empty UUID.
	var struct2 struct {
		Struct struct {
			Key1 uuid.UUID
			Key2 uuid.UUID
		}
	}
	data2 := []byte(`{"struct": {"key1": "", "key2": "14720916-dd67-11eb-b5da-305a3a7ae79b"}}`)
	if err := json.Unmarshal(data2, Wrap(&struct2, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := struct2.Struct.Key1; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
	if v := struct2.Struct.Key2; v != uuid.MustParse("14720916-dd67-11eb-b5da-305a3a7ae79b") {
		t.Errorf("invalid value: %v", v)
	}

	// With malformed UUID.
	var struct3 struct {
		Struct struct {
			Key1 uuid.UUID
			Key2 uuid.UUID
		}
	}
	data3 := []byte(`{"struct": {"key1": "test", "key2": "14720916-dd67-11eb-b5da-305a3a7ae79b"}}`)
	if err := json.Unmarshal(data3, Wrap(&struct3, AllowInvalid)); err != nil {
		t.Fatal(err)
	}
	if v := struct3.Struct.Key1; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
	if v := struct3.Struct.Key2; v != uuid.MustParse("14720916-dd67-11eb-b5da-305a3a7ae79b") {
		t.Errorf("invalid value: %v", v)
	}

	// With malformed UUID.
	var struct4 struct {
		Struct struct {
			Key1 uuid.UUID
			Key2 uuid.UUID
		}
	}
	data4 := []byte(`{"struct": {"key1": "test", "key2": "14720916-dd67-11eb-b5da-305a3a7ae79b"}}`)
	if err := json.Unmarshal(data4, Wrap(&struct4, AllowEmpty)); err == nil {
		t.Fatal("expected test to fail")
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var str string
	if err := json.Unmarshal([]byte("{}"), Wrap(str)); err == nil {
		t.Fatal("expected test to fail")
	}
}
