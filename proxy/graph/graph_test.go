package graph

import (
	"reflect"
	"testing"
)

func TestGenerateMermaid(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test",
			want: GenerateMermaid(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateMermaid(); got == "" {
				t.Errorf("GenerateMermaid() = %v, want %v", got, "")
			}
		})
	}
}

func TestNewNode(t *testing.T) {
	type args struct {
		id   int
		name string
		form string
	}
	tests := []struct {
		name string
		args args
		want *Node
	}{
		{
			name: "test",
			args: args{
				id:   1,
				name: "A1",
				form: "[Square Rect]",
			},
			want: NewNode(1, "A1", "[Square Rect]"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNode(tt.args.id, tt.args.name, tt.args.form); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_AddLink(t *testing.T) {
	type fields struct {
		ID      int
		Name    string
		Form    string
		Links   []*Node
		Visited bool
	}
	type args struct {
		node *Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test",
			fields: fields{
				ID:   1,
				Name: "A1",
				Form: "[Square Rect]",
			},
			args: args{
				node: NewNode(2, "A2", "[Square Rect]"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				ID:      tt.fields.ID,
				Name:    tt.fields.Name,
				Form:    tt.fields.Form,
				Links:   tt.fields.Links,
				Visited: tt.fields.Visited,
			}
			n.AddLink(tt.args.node)
		})
	}
}

func TestToMermaid(t *testing.T) {
	type args struct {
		nodes []*Node
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				nodes: []*Node{
					NewNode(1, "A1", "[Square Rect]"),
					NewNode(2, "A2", "[Square Rect]"),
				},
			},
			want: "graph LR\nA1[Square Rect]\nA2[Square Rect]\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToMermaid(tt.args.nodes); got == " " {
				t.Errorf("ToMermaid() = %v, want %v", got, " ")
			}
		})
	}
}

func Test_getRandomNodes(t *testing.T) {
	type args struct {
		nodes []*Node
		count int
	}
	tests := []struct {
		name string
		args args
		want []*Node
	}{
		{
			name: "test",
			args: args{
				nodes: []*Node{
					NewNode(1, "A1", "[Square Rect]"),
					NewNode(2, "A2", "[Square Rect]"),
				},
				count: 2,
			},
			want: []*Node{
				NewNode(1, "A1", "[Square Rect]"),
				NewNode(2, "A2", "[Square Rect]"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRandomNodes(tt.args.nodes, tt.args.count); len(got) != tt.args.count {
				t.Errorf("getRandomNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}
