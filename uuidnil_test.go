package uuidnil

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestUnmarshalArray(t *testing.T) {
	// Without UUID.
	data1 := []byte(`{
		"array": ["string1", "string2"]
	}`)
	var array1 struct {
		Array [2]string
	}
	if err := json.Unmarshal(data1, Wrap(&array1, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := array1.Array[0]; v != "string1" {
		t.Errorf("invalid value: %v", v)
	}
	if v := array1.Array[1]; v != "string2" {
		t.Errorf("invalid value: %v", v)
	}

	// With UUID.
	data2 := []byte(`{
		"array": ["", "ed9a09dc-dd17-11eb-a1ae-305a3a7ae79b"]
	}`)
	var array2 struct {
		Array [2]uuid.UUID
	}
	if err := json.Unmarshal(data2, Wrap(&array2, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := array2.Array[0]; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
	if v := array2.Array[1]; v != uuid.MustParse("ed9a09dc-dd17-11eb-a1ae-305a3a7ae79b") {
		t.Errorf("invalid value: %v", v)
	}
}

func TestUnmarshalMap(t *testing.T) {
	// Without UUID.
	data1 := []byte(`{
		"map": {"key1": "string1", "key2": "string2"}
	}`)
	var map1 struct {
		Map map[string]string
	}
	if err := json.Unmarshal(data1, Wrap(&map1, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := map1.Map["key1"]; v != "string1" {
		t.Errorf("invalid value: %v", v)
	}
	if v := map1.Map["key2"]; v != "string2" {
		t.Errorf("invalid value: %v", v)
	}

	// With UUID as key.
	data2 := []byte(`{
		"map": {"": "string1", "ad04dd68-dd65-11eb-bef6-305a3a7ae79b": "string2"}
	}`)
	var map2 struct {
		Map map[uuid.UUID]string
	}
	if err := json.Unmarshal(data2, Wrap(&map2, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := map2.Map[uuid.Nil]; v != "string1" {
		t.Errorf("invalid value: %v", v)
	}
	if v := map2.Map[uuid.MustParse("ad04dd68-dd65-11eb-bef6-305a3a7ae79b")]; v != "string2" {
		t.Errorf("invalid value: %v", v)
	}

	// With UUID as value.
	data3 := []byte(`{
			"map": {"key1": "", "key2": "ad04dd68-dd65-11eb-bef6-305a3a7ae79b"}
		}`)
	var map3 struct {
		Map map[string]uuid.UUID
	}
	if err := json.Unmarshal(data3, Wrap(&map3, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := map3.Map["key1"]; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
	if v := map3.Map["key2"]; v != uuid.MustParse("ad04dd68-dd65-11eb-bef6-305a3a7ae79b") {
		t.Errorf("invalid value: %v", v)
	}
}

func TestUnmarshalPtr(t *testing.T) {
	// Without UUID.
	data1 := []byte(`{
		"ptr": "pointer"
	}`)
	var ptr1 struct {
		Ptr *string
	}
	if err := json.Unmarshal(data1, Wrap(&ptr1, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := *ptr1.Ptr; v != "pointer" {
		t.Errorf("invalid value: %v", v)
	}

	// With UUID.
	data2 := []byte(`{
		"ptr": ""
	}`)
	var ptr2 struct {
		Ptr *uuid.UUID
	}
	if err := json.Unmarshal(data2, Wrap(&ptr2, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := *ptr2.Ptr; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
}

func TestUnmarshalSlice(t *testing.T) {
	// Without UUID.
	data1 := []byte(`{
		"slice": ["string1", "string2"]
	}`)
	var slice1 struct {
		Slice []string
	}
	if err := json.Unmarshal(data1, Wrap(&slice1, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := slice1.Slice[0]; v != "string1" {
		t.Errorf("invalid value: %v", v)
	}
	if v := slice1.Slice[1]; v != "string2" {
		t.Errorf("invalid value: %v", v)
	}

	// With UUID.
	data2 := []byte(`{
		"slice": ["", "14720916-dd67-11eb-b5da-305a3a7ae79b"]
	}`)
	var slice2 struct {
		Slice []uuid.UUID
	}
	if err := json.Unmarshal(data2, Wrap(&slice2, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := slice2.Slice[0]; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
	if v := slice2.Slice[1]; v != uuid.MustParse("14720916-dd67-11eb-b5da-305a3a7ae79b") {
		t.Errorf("invalid value: %v", v)
	}
}

func TestUnmarshalStruct(t *testing.T) {
	// Without UUID.
	data1 := []byte(`{
		"struct": {"key1": "string1", "key2": "string2"}
	}`)
	var struct1 struct {
		Struct struct {
			Key1 string
			Key2 string
		}
	}
	if err := json.Unmarshal(data1, Wrap(&struct1, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := struct1.Struct.Key1; v != "string1" {
		t.Errorf("invalid value: %v", v)
	}
	if v := struct1.Struct.Key2; v != "string2" {
		t.Errorf("invalid value: %v", v)
	}

	// With UUID.
	data2 := []byte(`{
		"struct": {"key1": "", "key2": "14720916-dd67-11eb-b5da-305a3a7ae79b"}
	}`)
	var struct2 struct {
		Struct struct {
			Key1 uuid.UUID
			Key2 uuid.UUID
		}
	}
	if err := json.Unmarshal(data2, Wrap(&struct2, AllowEmpty)); err != nil {
		t.Fatal(err)
	}
	if v := struct2.Struct.Key1; v != uuid.Nil {
		t.Errorf("invalid value: %v", v)
	}
	if v := struct2.Struct.Key2; v != uuid.MustParse("14720916-dd67-11eb-b5da-305a3a7ae79b") {
		t.Errorf("invalid value: %v", v)
	}
}
