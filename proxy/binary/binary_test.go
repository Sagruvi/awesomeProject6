package binary

import (
	"reflect"
	"testing"
)

func TestAVLTree_Insert(t1 *testing.T) {
	type fields struct {
		Root *Node
	}
	type args struct {
		key int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test",
			fields: fields{
				Root: nil,
			},
			args: args{
				key: 1,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &AVLTree{
				Root: tt.fields.Root,
			}
			t.Insert(tt.args.key)
		})
	}
}

func TestAVLTree_ToMermaid(t1 *testing.T) {
	type fields struct {
		Root *Node
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test",
			fields: fields{
				Root: nil,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &AVLTree{
				Root: tt.fields.Root,
			}
			if got := t.ToMermaid(); got == "" {
				t1.Errorf("ToMermaid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateTree(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name string
		args args
		want *AVLTree
	}{
		{
			name: "test",
			args: args{
				count: 10,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateTree(tt.args.count); reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNode(t *testing.T) {
	type args struct {
		key int
	}
	tests := []struct {
		name string
		args args
		want *Node
	}{
		{
			name: "test",
			args: args{
				key: 1,
			},
			want: &Node{
				Key:    1,
				Height: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNode(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBalance(t *testing.T) {
	type args struct {
		node *Node
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test",
			args: args{
				node: &Node{
					Key:    1,
					Height: 1,
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBalance(tt.args.node); got != tt.want {
				t.Errorf("getBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_height(t *testing.T) {
	type args struct {
		node *Node
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test",
			args: args{
				node: &Node{
					Key:    1,
					Height: 1,
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := height(tt.args.node); got != tt.want {
				t.Errorf("height() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_insert(t *testing.T) {
	type args struct {
		node *Node
		key  int
	}
	tests := []struct {
		name string
		args args
		want *Node
	}{
		{
			name: "test",
			args: args{
				node: &Node{
					Key:    1,
					Height: 1,
				},
				key: 1,
			},
			want: &Node{
				Key:    1,
				Height: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := insert(tt.args.node, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_leftRotate(t *testing.T) {
	type args struct {
		x *Node
	}
	tests := []struct {
		name string
		args args
		want *Node
	}{
		{
			name: "test",
			args: args{
				x: &Node{
					Key:    1,
					Height: 1,
					Right: &Node{
						Key:    2,
						Height: 2,
						Left: &Node{
							Key:    3,
							Height: 3,
						},
					},
				},
			},
			want: &Node{
				Key:    2,
				Height: 5,
				Left: &Node{
					Key:    3,
					Height: 3,
					Right: &Node{
						Key:    2,
						Height: 2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := leftRotate(tt.args.x); !reflect.DeepEqual(got.Height, tt.want.Height) {
				t.Errorf("leftRotate() = %v, want %v", got.Height, tt.want.Height)
			}
		})
	}
}

func Test_max(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test",
			args: args{
				a: 1,
				b: 2,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeToMermaid(t *testing.T) {
	type args struct {
		node *Node
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				node: &Node{
					Key:    1,
					Height: 1,
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nodeToMermaid(tt.args.node); got != tt.want {
				t.Errorf("nodeToMermaid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateHeight(t *testing.T) {
	type args struct {
		node *Node
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				node: &Node{
					Key:    1,
					Height: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateHeight(tt.args.node)
		})
	}
}
