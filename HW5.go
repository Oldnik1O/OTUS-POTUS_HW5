//Разработана структура для контейнера который учитывает разные "скоупы" (в немсозданы ключи команд для управления функциямии) и набор функций для регистрации и разрешения зависимостей

go
package main

import (
  "errors"
  "fmt"
  "sync"
)

type ResolverFunc func(args ...interface{}) (interface{}, error)

type IoC struct {
  mu       sync.Mutex
  registry map[string]ResolverFunc
  scopes   map[string]map[string]ResolverFunc
}

func NewIoC() *IoC {
  return &IoC{
    registry: make(map[string]ResolverFunc),
    scopes:   make(map[string]map[string]ResolverFunc),
  }
}

func (ioc *IoC) Register(key string, resolver ResolverFunc) {
  ioc.mu.Lock()
  defer ioc.mu.Unlock()
  ioc.registry[key] = resolver
}

func (ioc *IoC) Resolve(key string, args ...interface{}) (interface{}, error) {
  ioc.mu.Lock()
  defer ioc.mu.Unlock()

  // Check if the key is special command
  switch key {
  case "IoC.Register":
    if len(args) < 2 {
      return nil, errors.New("invalid arguments for IoC.Register")
    }
    keyArg, ok1 := args[0].(string)
    resolverArg, ok2 := args[1].(ResolverFunc)
    if !ok1 || !ok2 {
      return nil, errors.New("invalid argument types for IoC.Register")
    }
    ioc.Register(keyArg, resolverArg)
    return nil, nil
  case "Scopes.New":
    if len(args) < 1 {
      return nil, errors.New("invalid arguments for Scopes.New")
    }
    scopeID, ok := args[0].(string)
    if !ok {
      return nil, errors.New("invalid argument type for Scopes.New")
    }
    ioc.scopes[scopeID] = make(map[string]ResolverFunc)
    return nil, nil
  case "Scopes.Current":
    if len(args) < 1 {
      return nil, errors.New("invalid arguments for Scopes.Current")
    }
    scopeID, ok := args[0].(string)
    if !ok {
      return nil, errors.New("invalid argument type for Scopes.Current")
    }
    if scope, exists := ioc.scopes[scopeID]; exists {
      return scope, nil
    } else {
      return nil, errors.New("scope not found")
    }
  }

  // Resolve from global registry
  if resolver, exists := ioc.registry[key]; exists {
    return resolver(args...)
  }

  return nil, errors.New("key not found")
}

func main() {
  ioc := NewIoC()

  // Register a simple dependency
  ioc.Register("create_string", func(args ...interface{}) (interface{}, error) {
    return "Hello IoC!", nil
  })

  result, err := ioc.Resolve("create_string")
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    fmt.Println(result)
  }

  // Register dependency using IoC.Resolve
  ioc.Resolve("IoC.Register", "create_int", func(args ...interface{}) (interface{}, error) {
    return 42, nil
  })
  intResult, err := ioc.Resolve("create_int")
  if err != nil {
    fmt.Println("Error:", err)
  } else {
    fmt.Println(intResult)
  }

  // Use scopes
  ioc.Resolve("Scopes.New", "myscope")
  scope, _ := ioc.Resolve("Scopes.Current", "myscope")
  fmt.Println(scope)
}

