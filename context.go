package cli

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
)

const abortIndex int = math.MaxInt32

// Context context for command
type Context struct {
	app   *App                   // the app
	args  []string               // all raw args
	cmds  []*Command             // all command chain
	flags map[string]*Flag       // all flag map
	datas map[string]interface{} // dynamic datas
	index int                    // use for call command
}

func newContext(app *App, args []string, cmds []*Command, flags map[string]*Flag) *Context {
	return &Context{
		app:   app,
		args:  args,
		cmds:  cmds,
		flags: flags,
		index: -1,
	}
}

func (c *Context) App() *App {
	return c.app
}

func (c *Context) CommandList() []*Command {
	return c.cmds
}

// Get return dynamic data
func (c *Context) Get(key string) interface{} {
	if c.datas != nil {
		return c.datas[key]
	}

	return nil
}

// Set dynamic data
func (c *Context) Set(key string, val interface{}) {
	if c.datas == nil {
		c.datas = make(map[string]interface{})
	}

	c.datas[key] = val
}

// NArg number of the command line arguments
func (c *Context) NArg() int {
	return len(c.args)
}

// Arg return arg by index
func (c *Context) Arg(i int) string {
	// add offset?
	return c.args[i]
}

// ArgBool return bool arg by index
func (c *Context) ArgBool(i int) bool {
	return strToBool(c.Arg(i))
}

// ArgInt return integer arg by index
func (c *Context) ArgInt(i int) int {
	return strToInt(c.Arg(i))
}

// ArgUint return uint arg by index
func (c *Context) ArgUint(i int) uint {
	return strToUint(c.Arg(i))
}

// ArgF32 return float32 arg
func (c *Context) ArgF32(i int) float32 {
	return strToF32(c.Arg(i))
}

// ArgF64 return float64 arg
func (c *Context) ArgF64(i int) float64 {
	return strToF64(c.Arg(i))
}

//////////////////////////////////////////////
// Flag data
//////////////////////////////////////////////

// NFlag number of the flags
func (c *Context) NFlag() int {
	return len(c.flags)
}

// Flag get flag by key
func (c *Context) Flag(key string) *Flag {
	return c.flags[key]
}

// FlagStr return string flag
func (c *Context) FlagStr(key string) string {
	flag := c.Flag(key)
	if flag != nil {
		return flag.Get()
	}

	return ""
}

// FlagInt return int flag
func (c *Context) FlagInt(key string) int {
	return strToInt(c.FlagStr(key))
}

// FlagUint return uint flag
func (c *Context) FlagUint(key string) uint {
	return strToUint(c.FlagStr(key))
}

// FlagBool return bool flag
func (c *Context) FlagBool(key string) bool {
	return strToBool(c.FlagStr(key))
}

// FlagF32 return float32 flag
func (c *Context) FlagF32(key string) float32 {
	return strToF32(c.FlagStr(key))
}

// FlagF64 return float64 flag
func (c *Context) FlagF64(key string) float64 {
	return strToF64(c.FlagStr(key))
}

// FlagList return list flag
func (c *Context) FlagList(key string) []string {
	flag := c.Flag(key)
	if flag != nil {
		return flag.options
	}

	return nil
}

// Bind auto bind struct pointer, support basic type and slice and map field
// limit key and value of map must be basic type
// example:
// func onCommand(ctx *cli.Context) {
// 	 var flags struct {
// 	 	Type string `cli:"type"`
// 	 }
// 	 ctx.Bind(&flag)
// }
func (c *Context) Bind(flags interface{}) {
	value := reflect.ValueOf(flags)
	if value.Kind() != reflect.Ptr {
		panic(fmt.Errorf("bind must be pointer"))
	}

	if value.Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("bind just support struct"))
	}

	value = value.Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		vtype := value.Type().Field(i)
		name := ""
		tags := vtype.Tag.Get("cli")
		if tags != "" {
			name = tags
		} else {
			name = toKebabCase(vtype.Name)
		}

		f := c.Flag(name)
		if f == nil || f.Len() == 0 {
			continue
		}

		kind := field.Kind()
		if kind >= reflect.Bool && kind <= reflect.Float64 && kind != reflect.Uintptr {
			err := c.bindValue(f.Get(), field, kind)
			if err != nil {
				panic(fmt.Errorf("bind field[%+v] fail,%+v", vtype.Name, err))
			}

		} else if kind == reflect.Slice {
			elem := vtype.Type.Elem()
			slice := reflect.MakeSlice(elem, f.Len(), f.Len())
			for i := 0; i < f.Len(); i++ {
				val := slice.Index(i)
				str := f.GetAt(i)
				err := c.bindValue(str, val, elem.Kind())
				if err != nil {
					panic(fmt.Errorf("bind slice field[%+v] fail,%+v", vtype.Name, err))
				}
			}
		} else if kind == reflect.Map {
			// create map
			field.Set(reflect.MakeMap(vtype.Type))

			ft := vtype.Type
			re := regexp.MustCompile("(=|:)")
			for i := 0; i < f.Len(); i++ {
				str := f.GetAt(i)
				vk := re.Split(str, 2)
				if len(vk) != 2 {
					panic(fmt.Errorf("bind map field[%+v] fail, just support split by {=|:}", vtype.Name))
				}

				key := reflect.New(ft.Key()).Elem()
				val := reflect.New(ft.Elem()).Elem()
				var err error
				err = c.bindValue(vk[0], key, ft.Key().Kind())
				if err != nil {
					panic(fmt.Errorf("bind map field[%+v] fail,%+v", vtype.Name, err))
				}

				err = c.bindValue(vk[1], val, ft.Elem().Kind())
				if err != nil {
					panic(fmt.Errorf("bind map field[%+v] fail, %+v", vtype.Name, err))
				}

				field.SetMapIndex(key, val)
			}
		}
	}
}

// bindValue just bind the value of the basic type
func (c *Context) bindValue(str string, value reflect.Value, kind reflect.Kind) error {
	switch kind {
	case reflect.String:
		value.SetString(str)
	case reflect.Bool:
		v, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		value.SetBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		value.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		value.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		value.SetFloat(v)
	default:
		return fmt.Errorf("not support type:%+v", value.Kind())
	}

	return nil
}

//////////////////////////////////////////////
// Flow control
//////////////////////////////////////////////

// Next executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++

	for s := len(c.cmds); c.index < s; c.index++ {
		cmd := c.cmds[c.index]
		if cmd.Run != nil {
			cmd.Run(c)
		}
	}
}

// Abort stop process
func (c *Context) Abort() {
	c.index = abortIndex
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}
