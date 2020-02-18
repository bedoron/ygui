package treeBuilder

import (
	"fmt"
	"github.com/gizak/termui/v3/widgets"
	yaml "gopkg.in/yaml.v3"
	"reflect"
	"sort"
)

type Builder struct {
	root yaml.Node
	deserialized map[interface{}]interface{}
	treeMapping []*widgets.TreeNode
}

func (b *Builder) Nodes() []*widgets.TreeNode {
	return b.treeMapping
}

func NewBuilder(yamlText []byte) (*Builder, error) {
	builder := &Builder{
		deserialized: make(map[interface{}]interface{}),
	}

	err := yaml.Unmarshal(yamlText, &builder.root)
	if err == nil {
		err = yaml.Unmarshal(yamlText, &builder.deserialized)
	}

	builder.treeMapping = buildGraphNodes(builder.deserialized)

	return builder, err
}

type nodeValue string

func (nv nodeValue) String() string {
	return string(nv)
}

func buildGraphNodesFromSlice(root interface{}) []*widgets.TreeNode {
	graph := root.([]interface{})
	childNodes := make([]*widgets.TreeNode, 0)
	for _, k := range graph {
		nodes := buildGraphNodes(k)
		if len(nodes) != 1 {
			fmt.Errorf("weird node data")
			continue
		}

		childNodes = append(childNodes, nodes[0])
	}

	return childNodes
}

func buildGraphNodesFromMap(root interface{}) []*widgets.TreeNode {
	original := reflect.ValueOf(root)

	keys := make([]string, 0)
	for _, key := range original.MapKeys() {
		str := key.Interface().(string)
		keys = append(keys, str)
	}
	sort.Strings(keys)

	childNodes := make([]*widgets.TreeNode, 0)
	for _, key := range keys {
		mapValue := original.MapIndex(reflect.ValueOf(key))
		value := mapValue.Interface()

		node := &widgets.TreeNode{Value: nodeValue(key)}
		node.Nodes = buildGraphNodes(value)
		childNodes = append(childNodes, node)
	}

	return childNodes
}

func buildGraphNodes(graph interface{}) []*widgets.TreeNode {
	if graph == nil {
		return []*widgets.TreeNode{
			{
				Value: nodeValue(fmt.Sprintf("nil")),
				Nodes: nil,
			},
		}
	}

	switch reflect.TypeOf(graph).Kind() {
	case reflect.Map:
		return buildGraphNodesFromMap(graph)
	case reflect.Slice:
		return buildGraphNodesFromSlice(graph)
	default:
		return []*widgets.TreeNode{
			{
				Value: nodeValue(fmt.Sprintf("%v", graph)),
				Nodes: nil,
			},
		}
	}
}
