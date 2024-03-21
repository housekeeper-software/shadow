package db

import (
	"encoding/json"
	"errors"
	"strings"
)

/*
{
   "faces":[
      {
         "id":"1",
         "no":"000001C-01B-01U-01F-0101R",
         "name":"xxx",
         "faceFeature":"feature base64 encode",
         "startTime":"UTC time",
         "endTime":"UTC time"
      }
   ],
   "cards":[
      {
         "id":"1",
         "no":"000001C-01B-01U",
         "code":"12121212",
         "startTime":"UTC time",
         "endTime":"UTC time"
      }
   ]
}
*/

var (
	errNotFound = errors.New("not found key(id)")
)

type Document struct {
	docs map[string][]map[string]interface{}
}

func NewDocument() *Document {
	return &Document{
		docs: make(map[string][]map[string]interface{}),
	}
}

//Load from database
func (d *Document) Load(data []byte) error {
	return json.Unmarshal(data, &d.docs)
}

//AddItem add a new item to this document
func (d *Document) AddItem(name string, value []byte) error {
	//unmarshal first
	m := make(map[string]interface{})
	err := json.Unmarshal(value, &m)
	if err != nil {
		return err
	}
	//must exist a key[id]
	id, ok := m["id"]
	if !ok {
		return errNotFound
	}
	c, ok := d.docs[name]
	if !ok {
		//convert to slice
		items := make([]map[string]interface{}, 0, 1)
		items = append(items, m)
		//add to section
		d.docs[name] = items
	} else {
		found := false
		for k, v := range c {
			if strings.EqualFold(v["id"].(string), id.(string)) {
				c[k] = m //replace it directly
				found = true
				break
			}
		}
		if !found {
			//append to tail
			c = append(c, m)
		}
		d.docs[name] = c
	}
	return nil
}

func (d *Document) RemoveItem(name string, id string) {
	c, ok := d.docs[name]
	if !ok {
		//name not exist
		return
	}
	if len(id) < 1 {
		//delete section if id is empty
		delete(d.docs, name)
		return
	}
	for k, v := range c {
		if strings.EqualFold(v["id"].(string), id) {
			c = append(c[:k], c[k+1:]...)
			break
		}
	}
	d.docs[name] = c
}

//IsEmpty check document is empty
func (d *Document) IsEmpty() bool {
	for _, v := range d.docs {
		if len(v) > 0 {
			return false
		}
	}
	return true
}

func (d *Document) ToJson() ([]byte, error) {
	return json.Marshal(d.docs)
}
