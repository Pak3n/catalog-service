package binding

import (
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

type defaultValidator struct {
	validate *validator.Validate
	once     sync.Once
}

var _ StructValidator = &defaultValidator{}

func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
	})
}

func (v *defaultValidator) Engine() any {
	v.lazyInit()
	return v.validate
}

func (v *defaultValidator) ValidateStruct(obj any) error {
	// 1. nil — нечего валидировать
	if obj == nil {
		return nil
	}

	// 2. reflect.ValueOf — получаем Value
	val := reflect.ValueOf(obj)

	// 3. Если указатель — проверяем на nil и разыменовываем через Elem()
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	// 4. Если не структура — нечего валидировать
	if val.Kind() != reflect.Struct {
		return nil
	}

	// 5. Ленивая инициализация и валидация
	v.lazyInit()
	return v.validate.Struct(obj)
}
