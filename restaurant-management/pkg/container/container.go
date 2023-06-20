package container

import (
	"fmt"
	"reflect"
	"unsafe"
)

type binding struct {
	resolver interface{} // resolver function
	instance interface{} // instance stored for singleton bindings
}

// arguments will return resolved arguments of the given function.
func arguments(function interface{}) []reflect.Value {
	functionTypeOf := reflect.TypeOf(function)
	argumentsCount := functionTypeOf.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := functionTypeOf.In(i)

		var instance interface{}

		if concrete, ok := container[abstraction]; ok {
			instance = concrete.resolve()
		} else {
			panic("no concrete found for the abstraction: " + abstraction.String())
		}

		arguments[i] = reflect.ValueOf(instance)
	}

	return arguments
}

// invoke will call the given function and return its returned value.
// It only works for functions that return a single value.
func invoke(function interface{}) interface{} {
	return reflect.ValueOf(function).Call(arguments(function))[0].Interface()
}

// resolve will return the concrete of related abstraction.
func (b binding) resolve() interface{} {
	if b.instance != nil {
		return b.instance
	}

	return invoke(b.resolver)
}

// container is the IoC container that will keep all of the bindings.
var container = map[reflect.Type]binding{}

// bind will map an abstraction to a concrete and set instance if it's a singleton binding.
func bind(resolver interface{}, singleton bool) {
	resolverTypeOf := reflect.TypeOf(resolver)
	if resolverTypeOf.Kind() != reflect.Func {
		panic("the resolver must be a function")
	}

	for i := 0; i < resolverTypeOf.NumOut(); i++ {
		var instance interface{}
		if singleton {
			instance = invoke(resolver)
		}

		container[resolverTypeOf.Out(i)] = binding{
			resolver: resolver,
			instance: instance,
		}
	}
}

// Register will bind an abstraction to a concrete for further singleton resolves.
// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
// The resolver function can have arguments of abstraction that have bound already in Container.
func Register(resolver interface{}) {
	bind(resolver, true)
}

// Resolver will resolve the dependency and return a appropriate concrete of the given abstraction.
// It can take an abstraction (interface reference) and fill it with the related implementation.
// It also can takes a function (receiver) with one or more arguments of the abstractions (interfaces) that need to be
// resolved, Container will invoke the receiver function and pass the related implementations.
func Resolver(receiver interface{}) {
	receiverTypeOf := reflect.TypeOf(receiver)
	if receiverTypeOf == nil {
		panic("cannot detect type of the receiver, make sure your are passing reference of the object")
	}

	if receiverTypeOf.Kind() == reflect.Ptr {
		abstraction := receiverTypeOf.Elem()

		if concrete, ok := container[abstraction]; ok {
			instance := concrete.resolve()
			reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))
			return
		}

		panic("no concrete found for the abstraction " + abstraction.String())
	}

	if receiverTypeOf.Kind() == reflect.Func {
		arguments := arguments(receiver)
		reflect.ValueOf(receiver).Call(arguments)
		return
	}

	panic("the receiver must be either a reference or a callback")
}

func ResolverNil(receiver interface{}) {
	receiverTypeOf := reflect.TypeOf(receiver)
	if receiverTypeOf == nil {
		fmt.Println("cannot detect type of the receiver, make sure your are passing reference of the object")
		return
	}

	if receiverTypeOf.Kind() == reflect.Ptr {
		abstraction := receiverTypeOf.Elem()

		if concrete, ok := container[abstraction]; ok {
			instance := concrete.resolve()
			reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))
			return
		}

		fmt.Println("no concrete found for the abstraction " + abstraction.String())
		return

	}

	if receiverTypeOf.Kind() == reflect.Func {
		arguments := arguments(receiver)
		reflect.ValueOf(receiver).Call(arguments)
		return
	}

	fmt.Println("the receiver must be either a reference or a callback")
}

func ResolverErr(receiver interface{}) error {
	receiverTypeOf := reflect.TypeOf(receiver)
	if receiverTypeOf == nil {
		return fmt.Errorf("cannot detect type of the receiver, make sure your are passing reference of the object")
	}

	if receiverTypeOf.Kind() == reflect.Ptr {
		abstraction := receiverTypeOf.Elem()

		if concrete, ok := container[abstraction]; ok {
			instance := concrete.resolve()
			reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))
			return nil
		}

		return fmt.Errorf("no concrete found for the abstraction " + abstraction.String())
	}

	if receiverTypeOf.Kind() == reflect.Func {
		arguments := arguments(receiver)
		reflect.ValueOf(receiver).Call(arguments)
		return nil
	}

	return fmt.Errorf("the receiver must be either a reference or a callback")
}

// Reset reset the container, remove all the bindings
func Reset() {
	container = map[reflect.Type]binding{}
}

// Fill takes a struct and fills the fields with the tag `di:"inject"`
func Fill(structure interface{}) {
	receiverType := reflect.TypeOf(structure)
	if receiverType == nil {
		panic("container: invalid structure")
	}

	if receiverType.Kind() == reflect.Ptr {
		elem := receiverType.Elem()
		if elem.Kind() == reflect.Struct {
			s := reflect.ValueOf(structure).Elem()

			for i := 0; i < s.NumField(); i++ {
				f := s.Field(i)

				if t, ok := s.Type().Field(i).Tag.Lookup("di"); ok && t == "inject" {
					if concrete, ok := container[f.Type()]; ok {
						instance := concrete.resolve()
						// if err != nil {
						//	return err
						// }

						ptr := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
						ptr.Set(reflect.ValueOf(instance))

						continue
					}

					panic(fmt.Sprintf("container: cannot resolve %v field", s.Type().Field(i).Name))
				}
			}

			return
		}
	}

	panic("container: invalid structure")
}
