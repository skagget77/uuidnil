package uuidnil

import (
	"errors"
	"log"
	"reflect"

	"github.com/google/uuid"
)

// assignFunc assigns the from value to the to value.
type assignFunc func(to reflect.Value, from reflect.Value) error

// proxyArray returns a proxy type and assign function for the specified array
// type.
func proxyArray(typ reflect.Type, opts Option) (reflect.Type, assignFunc) {
	proxy, assign := proxyType(typ.Elem(), opts)

	return reflect.ArrayOf(typ.Len(), proxy), func(to reflect.Value, from reflect.Value) error {
		if opts.trace() {
			log.Printf("[TRACE] assign: %s => %s\n", from.Type(), to.Type())
		}

		for i := 0; i < typ.Len(); i++ {
			if err := assign(to.Index(i), from.Index(i)); err != nil {
				return err
			}
		}

		return nil
	}
}

// proxyMap returns a proxy type and assign function for the specified map
// type.
func proxyMap(typ reflect.Type, opts Option) (reflect.Type, assignFunc) {
	proxyKey, assignKey := proxyType(typ.Key(), opts)
	proxyVal, assignVal := proxyType(typ.Elem(), opts)

	return reflect.MapOf(proxyKey, proxyVal), func(to reflect.Value, from reflect.Value) error {
		if opts.trace() {
			log.Printf("[TRACE] assign: %s => %s\n", from.Type(), to.Type())
		}

		if from.IsNil() {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}

		if to.IsNil() {
			to.Set(reflect.MakeMapWithSize(to.Type(), from.Len()))
		}

		iter := from.MapRange()
		for iter.Next() {
			key := reflect.Indirect(reflect.New(typ.Key()))
			if err := assignKey(key, iter.Key()); err != nil {
				return err
			}

			val := reflect.Indirect(reflect.New(typ.Elem()))
			if err := assignVal(val, iter.Value()); err != nil {
				return err
			}

			to.SetMapIndex(key, val)
		}

		return nil
	}
}

// proxyPtr returns a proxy type and assign function for the specified pointer
// type.
func proxyPtr(typ reflect.Type, opts Option) (reflect.Type, assignFunc) {
	proxy, assign := proxyType(typ.Elem(), opts)

	return reflect.PtrTo(proxy), func(to reflect.Value, from reflect.Value) error {
		if opts.trace() {
			log.Printf("[TRACE] assign: %s => %s\n", from.Type(), to.Type())
		}

		if from.IsNil() {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}

		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}

		return assign(to.Elem(), from.Elem())
	}
}

// proxySlice returns a proxy type and assign function for the specified slice
// type.
func proxySlice(typ reflect.Type, opts Option) (reflect.Type, assignFunc) {
	proxy, assign := proxyType(typ.Elem(), opts)

	return reflect.SliceOf(proxy), func(to reflect.Value, from reflect.Value) error {
		if opts.trace() {
			log.Printf("[TRACE] assign: %s => %s\n", from.Type(), to.Type())
		}

		if from.IsNil() {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}

		slice := reflect.MakeSlice(to.Type(), 0, from.Len())
		for i := 0; i < from.Len(); i++ {
			elem := reflect.Indirect(reflect.New(to.Type().Elem()))
			if err := assign(elem, from.Index(i)); err != nil {
				return err
			}

			slice = reflect.Append(slice, elem)
		}

		to.Set(slice)
		return nil
	}
}

// proxyStruct returns a proxy type and assign function for the specified
// struct type.
func proxyStruct(typ reflect.Type, opts Option) (reflect.Type, assignFunc) {
	n := typ.NumField()

	proxies := make([]reflect.StructField, 0, n)
	assigns := make([]assignFunc, 0, n)
	for i := 0; i < n; i++ {
		field := typ.Field(i)

		proxy, assign := proxyType(field.Type, opts)
		proxies = append(proxies, reflect.StructField{
			Name:      field.Name,
			Type:      proxy,
			Tag:       field.Tag,
			Index:     field.Index,
			Anonymous: false,
		})
		assigns = append(assigns, assign)
	}

	return reflect.StructOf(proxies), func(to reflect.Value, from reflect.Value) error {
		if opts.trace() {
			log.Printf("[TRACE] assign: %s => %s\n", from.Type(), to.Type())
		}

		for i := 0; i < n; i++ {
			if err := assigns[i](to.Field(i), from.Field(i)); err != nil {
				return err
			}
		}

		return nil
	}
}

// proxyUUID returns a proxy type (string) and assign function for the
// specified UUID type.
func proxyUUID(typ reflect.Type, opts Option) (reflect.Type, assignFunc) {
	return reflect.TypeOf(""), func(to reflect.Value, from reflect.Value) error {
		if opts.trace() {
			log.Printf("[TRACE] assign: %s => %s\n", from.Type(), to.Type())
		}

		// *UUID <=> *string.
		if from.Kind() == reflect.Ptr {
			if from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return nil
			}

			if from.Type().Elem().Kind() != reflect.String {
				id, ok := from.Interface().(*uuid.UUID)
				if !ok {
					return errors.New("uuidnil: internal inconsistency")
				}

				s := id.String()
				to.Set(reflect.ValueOf(&s))
				return nil
			}

			str := from.Elem().String()
			if str == "" && opts.empty() {
				to.Set(reflect.New(to.Elem().Type()))
				return nil
			}

			id, err := uuid.Parse(str)
			if err != nil && opts.invalid() {
				to.Set(reflect.New(to.Elem().Type()))
				return nil
			}
			if err != nil {
				return err
			}

			to.Set(reflect.ValueOf(&id))
			return nil
		}

		// UUID <=> string.
		if from.Type().Kind() != reflect.String {
			id, ok := from.Interface().(uuid.UUID)
			if !ok {
				return errors.New("uuidnil: internal inconsistency")
			}

			to.Set(reflect.ValueOf(id.String()))
			return nil
		}

		str := from.String()
		if str == "" && opts.empty() {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}

		id, err := uuid.Parse(str)
		if err != nil && opts.invalid() {
			to.Set(reflect.Zero(to.Type()))
			return nil
		}
		if err != nil {
			return err
		}

		to.Set(reflect.ValueOf(id))
		return nil
	}
}

// proxyType returns a proxy type and assign function for the specified type.
func proxyType(typ reflect.Type, opts Option) (reflect.Type, assignFunc) {
	if opts.trace() {
		log.Printf("[TRACE] proxy: %s", typ)
	}

	// UUID type.
	if typ.PkgPath() == "github.com/google/uuid" && typ.Name() == "UUID" {
		return proxyUUID(typ, opts)
	}

	// Composite types.
	switch typ.Kind() {
	case reflect.Array:
		return proxyArray(typ, opts)
	case reflect.Map:
		return proxyMap(typ, opts)
	case reflect.Ptr:
		return proxyPtr(typ, opts)
	case reflect.Slice:
		return proxySlice(typ, opts)
	case reflect.Struct:
		return proxyStruct(typ, opts)
	}

	// Simple type.
	return typ, func(to reflect.Value, from reflect.Value) error {
		if opts.trace() {
			log.Printf("[TRACE] assign: %s => %s\n", from.Type(), to.Type())
		}

		to.Set(from)
		return nil
	}
}

// proxyValue returns a proxy value and assign function for the specified
// value.
func proxyValue(value reflect.Value, opts Option) (reflect.Value, assignFunc) {
	proxy, assign := proxyType(value.Type().Elem(), opts)

	if opts.debug() {
		log.Printf("[DEBUG] value type: %s, proxy type: %s\n", value.Type().Elem(), proxy)
	}

	return reflect.New(proxy), assign
}
