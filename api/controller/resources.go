package controller

import (
	"fmt"
	"reflect"

	"github.com/dzahariev/solei/api/model"
)

// Resource holds information about a resource
type Resource struct {
	Name     string
	IsGlobal bool
	Type     reflect.Type
}

// ResourceFactory is used to hold information about supported resources
type ResourceFactory struct {
	resources map[string]Resource
}

// Register is used to register a resource type
func (s ResourceFactory) Register(object model.Object) {
	name := object.ResourceName()
	isGlobal := object.IsGlobal()
	objectType := reflect.TypeOf(object).Elem()
	s.resources[name] = Resource{
		Name:     name,
		IsGlobal: isGlobal,
		Type:     objectType,
	}
}

// Names returns the names of all registered resources
func (s ResourceFactory) Names() []string {
	names := make([]string, 0, len(s.resources))
	for name := range s.resources {
		names = append(names, name)
	}
	return names
}

// Names returns the names of all registered resources
func (s ResourceFactory) Resources() []Resource {
	resources := make([]Resource, 0, len(s.resources))
	for name := range s.resources {
		resources = append(resources, s.resources[name])
	}
	return resources
}

// New is used to create a new resource object
func (s ResourceFactory) New(name string) (model.Object, error) {
	t, ok := s.resources[name]
	if !ok {
		return nil, fmt.Errorf("unrecognized resource name: %s", name)
	}

	obj, ok := reflect.New(t.Type).Interface().(model.Object)
	if !ok {
		return nil, fmt.Errorf("type %s does not implement model.Object", t.Type)
	}
	return obj, nil
}

// NewResourceFactory is used to create a new ResourceFactory
func NewResourceFactory() *ResourceFactory {
	return &ResourceFactory{resources: map[string]Resource{}}
}
