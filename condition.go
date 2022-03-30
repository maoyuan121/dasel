package dasel

import (
	"reflect"
)

// Condition 定义了一个 Check 方法，在 dynamic selectors 中使用
type Condition interface {
	Check(other reflect.Value) (bool, error)
}
