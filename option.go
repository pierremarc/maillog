package main

// type mapFunc func(interface{}) interface{}

// type Option interface {
// 	Return(interface{}) Option
// 	Map(mapFunc) Option
// 	Fold(interface{}, mapFunc) interface{}
// }

// type None struct{}

// func (n None) Return(a interface{}) Option {
// 	return None{}
// }
// func (n None) Map(f mapFunc) Option {
// 	return None{}
// }
// func (n None) Fold(a interface{}, f mapFunc) interface{} {
// 	return a
// }

// func none() Option {
// 	return None{}
// }

// // string
type SomeString struct {
	v string
}

func (ss SomeString) Return(v string) SomeString {
	return SomeString{v}
}
func (ss SomeString) Map(f func(string) string) SomeString {
	return deriveFmap(f, ss.v)
}

// func (ss SomeString) Fold(s interface{}, f mapFunc) interface{} {
// 	return f(ss.v)
// }

// // int
// type SomeInt struct {
// 	v int
// }
// type intMapper func(int) interface{}

// func mapInt(f intMapper) mapFunc {
// 	return func(v interface{}) interface{} {
// 		return f(v.(int))
// 	}
// }

// func (si SomeInt) Return(i interface{}) Option {
// 	return SomeInt{i.(int)}
// }
// func (si SomeInt) Map(f mapFunc) Option {
// 	return si.Return(f(si.v))
// }
// func (si SomeInt) Fold(i interface{}, f mapFunc) interface{} {
// 	return f(si.v)
// }

// func option(v interface{}) Option {
// 	switch v.(type) {
// 	case string:
// 		return SomeString{v.(string)}
// 	case int:
// 		return SomeInt{v.(int)}
// 	}

// 	return None{}
// }

// // func foo() {
// // 	var o = option(3)
// // 	mi := func(a int) interface{} {
// // 		return fmt.Sprintf("Number is %v", a)
// // 	}

// // 	ms := func(a string) interface{} {
// // 		return fmt.Sprintf("Got %s", a)
// // 	}
// // 	rs := mapString(func(a string) interface{} { return 3 })
// // 	fi := mapInt(mi)
// // 	fs := mapString(ms)

// // 	r := o.Map(fi).Map(fs).Fold(1, rs)
// // }
