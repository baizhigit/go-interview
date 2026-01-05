// взять data[key] (срез строк по ключу),
// добавить в него значения из newValue,
// не допуская дубликатов
// классический Go-паттерн set через map

package main

import (
	"fmt"
)

func MergeToMap(data map[string][]string, key string, values []string) {
	// 1. паттерн map[T]struct{} как set
	set := make(map[string]struct{}, len(data[key]))

	// 2. Заполняем set существующими значениями
	for _, v := range data[key] {
		set[v] = struct{}{}
	}

	// 3. Обрабатываем new values
	for _, v := range values {
		if _, ok := set[v]; ok {
			continue
		}
		data[key] = append(data[key], v)
		set[v] = struct{}{}
	}

	fmt.Println("Set", set)
}

func main() {
	fmt.Println("main start")

	oldMap := map[string][]string{
		"group1": {"apple", "banana"},
		"group2": {"carrot"},
	}

	newValues := []string{"banana", "cherry", "cherry"}

	fmt.Println("Before", oldMap)
	MergeToMap(oldMap, "group1", newValues)
	fmt.Println("After", oldMap)
}
