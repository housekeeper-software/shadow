package db

import "encoding/json"

// Index
/*
[
	{"name":"000001C-01B","hash":"9b97f3aa22abd1e7125f01c8523c243c"},
	{"name":"guest","hash":"9d2fad9ce02b83c585c500c53c3bb23d"}
]
*/

type Index struct {
	m map[string]string
}

//Result output format
type Result struct {
	Name     string `json:"name"`
	HashCode string `json:"hash"`
}

func NewIndex() *Index {
	return &Index{make(map[string]string)}
}

func (i *Index) Add(name string, hasCode string) {
	//delete first if exist,otherwise do nothing
	delete(i.m, name)
	i.m[name] = hasCode
}

func (i *Index) Remove(name string) {
	delete(i.m, name)
}

func (i *Index) Get(name string) (string, bool) {
	v, ok := i.m[name]
	return v, ok
}

func (i *Index) IsEmpty() bool {
	return len(i.m) < 1
}

func (i *Index) ToJson() ([]byte, error) {
	result := make([]Result, 0, len(i.m))
	for k, v := range i.m {
		r := Result{Name: k, HashCode: v}
		result = append(result, r)
	}
	return json.Marshal(result)
}
