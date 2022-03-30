package dasel

import (
	"fmt"
	"github.com/tomwright/dasel/storage"
	"io"
	"os"
	"reflect"
	"regexp"
)

// Selector 表示一个 Node 的选择去
type Selector struct {
	// Raw 是完整的 selector
	Raw string `json:"raw"`
	// Current 是与当前节点一起使用的选择器。
	Current string `json:"current"`
	// Remaining 是 Raw 选择器的剩余部分
	Remaining string `json:"remaining"`
	// Type 是 selector 类型
	Type string `json:"type"`
	// Property 是此选择器目标属性的名称(如果适用的话)。
	Property string `json:"property,omitempty"`
	// Index 是要使用的索引(如果适用的话)。
	Index int `json:"index,omitempty"`
	// Conditions 包含一组可选的匹配目标的条件。
	Conditions []Condition `json:"conditions,omitempty"`
}

// Copy 返回一个 selector 的拷贝
func (s Selector) Copy() Selector {
	return Selector{
		Raw:        s.Raw,
		Current:    s.Current,
		Remaining:  s.Remaining,
		Type:       s.Type,
		Property:   s.Property,
		Index:      s.Index,
		Conditions: s.Conditions,
	}
}

// Node表示选择器的节点链中的单个节点。
type Node struct {
	// Previous 是链中前一个节点
	Previous *Node `json:"-"`
	// Next 是链中的下一个节点
	// 与 Query 和 Put 请求一起使用。
	Next *Node `json:"next,omitempty"`
	// NextMultiple 包含链中的下一个节点
	// 与 QueryMultiple和  PutMultiple 请求一起使用的。
	// 当发生大版本变更时，将完全替换 Next。
	NextMultiple []*Node `json:"nextMultiple,omitempty"`
	// OriginalValue 是解析器返回的值。
	// 在大多数情况下，这与 Value 相同，但对于 thr YAML 解析器不同
	// 因为它包含原始文档的信息。
	OriginalValue interface{} `json:"-"`
	// Value 是当前节点的值
	Value reflect.Value `json:"value"`
	// Selector 是当前节点的选择器
	Selector       Selector `json:"selector"`
	wasInitialised bool
}

// String 以字符串的形式返回节点的值
// 没有格式，得到的是 raw value
func (n *Node) String() string {
	return fmt.Sprint(n.InterfaceValue())
}

// InterfaceValue 以 interface{} 的形式返回存储在节点中的值。
func (n *Node) InterfaceValue() interface{} {
	// We shouldn't be able to get here but this will stop a panic if we do.
	if !n.Value.IsValid() {
		return nil
	}
	return n.Value.Interface()
}

const (
	propertySelector = `(?P<property>[a-zA-Z\-_]+)`
	indexSelector    = `\[(?P<index>[0-9a-zA-Z\*]*?)\]`
)

var (
	propertyRegexp   = regexp.MustCompile(fmt.Sprintf("^\\.?%s", propertySelector))
	indexRegexp      = regexp.MustCompile(fmt.Sprintf("^\\.?%s", indexSelector))
	newDynamicRegexp = regexp.MustCompile(fmt.Sprintf("^\\.?((?:\\(.*\\))+)"))
)

func isValid(value reflect.Value) bool {
	return value.IsValid() && !safeIsNil(value)
}

func safeIsNil(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer,
		reflect.Interface, reflect.Slice:
		return value.IsNil()
	}
	return false
}

func nilValue() reflect.Value {
	return reflect.ValueOf(nil)
}

// 如果是接口类型，那么返回其 Elem，否则直接返回 value
func unwrapValue(value reflect.Value) reflect.Value {
	// value = reflect.Indirect(value)
	if value.Kind() == reflect.Interface {
		return value.Elem()
	}
	return value
}

// New 用给定的值返回一个新的根节点
func New(value interface{}) *Node {
	rootNode := &Node{
		Previous:     nil,
		Next:         nil,
		NextMultiple: nil,
		Selector: Selector{
			Raw:       ".",
			Current:   ".",
			Remaining: "",
			Type:      "ROOT",
			Property:  "",
		},
	}
	rootNode.setRealValue(value)
	return rootNode
}

// NewFromFile 通过使用指定的读解析器解析文件，返回一个新的根节点。
func NewFromFile(filename, parser string) (*Node, error) {
	readParser, err := storage.NewReadParserFromString(parser)
	if err != nil {
		return nil, err
	}

	data, err := storage.LoadFromFile(filename, readParser)
	if err != nil {
		return nil, err
	}

	return New(data), nil
}

// NewFromReader 通过使用指定的读解析器从 Reader中解析，返回一个新的根节点。
func NewFromReader(reader io.Reader, parser string) (*Node, error) {
	readParser, err := storage.NewReadParserFromString(parser)
	if err != nil {
		return nil, err
	}

	data, err := storage.Load(readParser, reader)
	if err != nil {
		return nil, err
	}

	return New(data), nil
}

// WriteToFile 使用指定的选项将数据写入给定的文件。
func (n *Node) WriteToFile(filename, parser string, writeOptions []storage.ReadWriteOption) error {
	f, err := os.Create(filename)

	if err != nil {
		return err
	}

	// https://www.joeshaw.org/dont-defer-close-on-writable-files/
	if err = n.Write(f, parser, writeOptions); err != nil {
		_ = f.Close()
		return err
	}

	return f.Close()
}

// Write 使用指定的写解析器和选项将写数据写入写入器。
func (n *Node) Write(writer io.Writer, parser string, writeOptions []storage.ReadWriteOption) error {
	writeParser, err := storage.NewWriteParserFromString(parser)
	if err != nil {
		return err
	}

	value := n.InterfaceValue()
	originalValue := n.OriginalValue

	if err := storage.Write(writeParser, value, originalValue, writer, writeOptions...); err != nil {
		return err
	}

	return nil
}

func (n *Node) setValue(newValue interface{}) {
	n.Value = reflect.ValueOf(newValue)
	if n.Selector.Type == "ROOT" {
		n.OriginalValue = newValue
	}
}

func (n *Node) setRealValue(newValue interface{}) {
	switch typed := newValue.(type) {
	case storage.RealValue:
		n.Value = reflect.ValueOf(typed.RealValue())
	default:
		n.Value = reflect.ValueOf(typed)
	}
	if n.Selector.Type == "ROOT" {
		n.OriginalValue = newValue
	}
}

func (n *Node) setReflectValue(newValue reflect.Value) {
	n.Value = newValue
	if n.Selector.Type == "ROOT" {
		n.OriginalValue = unwrapValue(newValue).Interface()
	}
}

func (n *Node) setRealReflectValue(newValue reflect.Value) {
	val := unwrapValue(newValue).Interface()
	switch typed := val.(type) {
	case storage.RealValue:
		n.OriginalValue = typed
		n.Value = reflect.ValueOf(typed.RealValue())
	default:
		n.Value = newValue
	}
	if n.Selector.Type == "ROOT" {
		n.OriginalValue = val
	}
}
