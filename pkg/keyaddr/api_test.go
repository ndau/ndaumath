package keyaddr

import (
	"reflect"
	"testing"
)

// We don't need a lot of test cases as all the supporting functions have their own
// test suites; all we need to do is make sure our wrappers are doing the right thing.

func TestWordsFromBytes(t *testing.T) {
	type args struct {
		lang string
		b    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"basic", args{"en", []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
			"abandon amount liar amount expire adjust cage candy arch gather drum bundle", false},
		{"generates an error", args{"foo", []byte{}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WordsFromBytes(tt.args.lang, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("WordsFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WordsFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWordsToBytes(t *testing.T) {
	type args struct {
		lang string
		w    string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"basic", args{"en", "abandon amount liar amount expire adjust cage candy arch gather drum bundle"},
			[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, false},
		{"generates an error for bad words", args{"en", "abandon amount blah amount expire adjust cage candy arch gather drum bundle"},
			nil, true},
		{"generates an error for language", args{"foo", "blah"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WordsToBytes(tt.args.lang, tt.args.w)
			if (err != nil) != tt.wantErr {
				t.Errorf("WordsToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WordsToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKey(t *testing.T) {
	type args struct {
		seed []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Key
		wantErr bool
	}{
		{"generates a known key", args{[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
			&Key{"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewKey(tt.args.seed)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKey_Neuter(t *testing.T) {
	type fields struct {
		Key string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Key
		wantErr bool
	}{
		{"simple neuter", fields{"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"},
			&Key{"npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7q5"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Key{
				Key: tt.fields.Key,
			}
			got, err := k.Neuter()
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.Neuter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Key.Neuter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKey_Child(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		n int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Key
		wantErr bool
	}{
		{"simple private child",
			fields{"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"},
			args{1},
			&Key{"npvt8ap98fgsaaaaagts2dwuzsn3dsv9mwqm3zsbrwrsxbwavxw36zwsyik47scskjtguaf8bxi9eqyqmsenuub3z62364hy2vb5u95uqr8n5j87rzudbt253kvnjj7w"},
			false},
		{"simple public child",
			fields{"npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7q5"},
			args{1},
			&Key{"npubaap98fgsaaaaagts2dwuzsn3dsv9mwqm3zsbrwrsxbwavxw36zwsyik47scskjtguax8c6mqspjegkytd98ksuaqp4txapxyy34ibvr7f7iy4taihk3qmn3yu7dj"},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Key{
				Key: tt.fields.Key,
			}
			got, err := k.Child(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.Child() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Key.Child() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKey_HardenedChild(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		n int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Key
		wantErr bool
	}{
		{"simple private hardened child",
			fields{"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"},
			args{1},
			&Key{"npvt8ap98fg2aaaaaenphqxyh7nh2zhjfugk3a9xvqwkcarfau8239ykec4h69kzkcs8cacj7bkdfhkeyr3mjv5jeaz7ptczqaeqhq6rtwyqvn9wvy5kprj9ubsk7wza"},
			false},
		{"attempt to create hardened child from a public key should fail",
			fields{"npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7q5"},
			args{1}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Key{
				Key: tt.fields.Key,
			}
			got, err := k.HardenedChild(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.HardenedChild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Key.HardenedChild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKey_Sign(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		msg []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Signature
		wantErr bool
	}{
		{"basic",
			fields{"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"},
			args{[]byte{1, 2, 3, 4}},
			&Signature{"gbcseiia598sbs4u8p76adr2cgbkhy679867sba4dsaggchk657yzg3f92waeiaaxz8qf7k46cnwt3g2inycttseh38bw5j7nac2jkdg7nywbe7zxi======"},
			false},
		{"public key should error",
			fields{"npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7q5"},
			args{[]byte{1, 2, 3, 4}},
			nil,
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Key{
				Key: tt.fields.Key,
			}
			got, err := k.Sign(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Key.Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKey_NdauAddress(t *testing.T) {
	type fields struct {
		Key string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Address
		wantErr bool
	}{
		{"addr from private key",
			fields{"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"},
			&Address{"ndad79yux8we7vk7dgvkqjwnkdhme57piydekb9bkbc6r7uj"}, false},
		{"addr from corresponding public key is the same",
			fields{"npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7q5"},
			&Address{"ndad79yux8we7vk7dgvkqjwnkdhme57piydekb9bkbc6r7uj"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Key{
				Key: tt.fields.Key,
			}
			got, err := k.NdauAddress()
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.NdauAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Key.NdauAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
