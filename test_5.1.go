package main

import (
  "errors"
  "testing"
)
// Тестиование нового IoC контейнера
func TestNewIoC(t *testing.T) {
  ioc := NewIoC()
  if ioc == nil {
    t.Error("Expected IoC container, got nil")
  }
}
// Тестиование регистрации
func TestIoC_Register(t *testing.T) {
  ioc := NewIoC()
  ioc.Register("create_string", func(args ...interface{}) (interface{}, error) {
    return "Hello IoC!", nil
  })

  _, err := ioc.Resolve("create_string")
  if err != nil {
    t.Errorf("Expected no error, got: %v", err)
  }
}

// Тестиование реализации
func TestIoC_Resolve(t *testing.T) {
  ioc := NewIoC()

  // Тестиование IoC.Register через IoC.Resolve
  ioc.Resolve("IoC.Register", "create_int", func(args ...interface{}) (interface{}, error) {
    return 42, nil
  })

  val, err := ioc.Resolve("create_int")
  if err != nil {
    t.Errorf("Expected no error, got: %v", err)
  }
  intVal, ok := val.(int)
  if !ok {
    t.Errorf("Expected integer, got: %T", val)
  }
  if intVal != 42 {
    t.Errorf("Expected 42, got: %d", intVal)
  }

  // Тестирование на отсутствие ключа
  _, err = ioc.Resolve("missing_key")
  if err == nil {
    t.Errorf("Expected error, got nil")
  }
  if !errors.Is(err, errors.New("key not found")) {
    t.Errorf("Expected 'key not found', got: %v", err)
  }

  // Тестирование недопустимых типы аргументов для  IoC.Register
  _, err = ioc.Resolve("IoC.Register", 42, "invalid")
  if err == nil || !errors.Is(err, errors.New("invalid argument types for IoC.Register")) {
    t.Errorf("Expected 'invalid argument types for IoC.Register', got: %v", err)
  }

  // Тестирование Scopes.New и Scopes.Current
  ioc.Resolve("Scopes.New", "myscope")
  val, err = ioc.Resolve("Scopes.Current", "myscope")
  if err != nil {
    t.Errorf("Expected no error, got: %v", err)
  }
  if val == nil {
    t.Errorf("Expected scope, got: %v", val)
  }
}