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
		s    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"basic", args{"en", "AAECAwQFBgcICQoLDA0ODw=="},
			"abandon amount liar amount expire adjust cage candy arch gather drum bundle", false},
		{"minor change gets slightly different words", args{"en", "AAECAwQFBgcIcQoLDA0ODw=="},
			"abandon amount liar amount expire adjust canyon candy arch gather drum business", false},
		{"generates an error", args{"foo", ""}, "", true},
		{"detects encoding error (bad length)", args{"en", "AAECAwQFBgcIcoLDA0ODw=="}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WordsFromBytes(tt.args.lang, tt.args.s)
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
		want    string
		wantErr bool
	}{
		{"basic", args{"en", "abandon amount liar amount expire adjust cage candy arch gather drum bundle"},
			"AAECAwQFBgcICQoLDA0ODw==", false},
		{"generates an error for bad words", args{"en", "abandon amount blah amount expire adjust cage candy arch gather drum bundle"},
			"", true},
		{"generates an error for language", args{"foo", "blah"}, "", true},
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

func TestWordsFromPrefix(t *testing.T) {
	type args struct {
		lang   string
		prefix string
		max    int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"two", args{"en", "riv", 0}, "rival river"},
		{"one", args{"en", "oxy", 0}, "oxygen"},
		{"five", args{"en", "dri", 0}, "drift drill drink drip drive"},
		{"five limit 3", args{"en", "dri", 3}, "drift drill drink"},
		{"none", args{"en", "zpj", 0}, ""},
		{"bad lang", args{"xx", "act", 0}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WordsFromPrefix(tt.args.lang, tt.args.prefix, tt.args.max); got != tt.want {
				t.Errorf("WordsFromPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKey(t *testing.T) {
	type args struct {
		seed string
	}
	tests := []struct {
		name    string
		args    args
		want    *Key
		wantErr bool
	}{
		{"generates a known key", args{"AAECAwQFBgcICQoLDA0ODw=="},
			&Key{"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"}, false},
		{"fails for bad encoding", args{"AAAwQFBgcICQoLDA0ODw=="}, nil, true},
		{"fails for too-short key", args{"AQIDBA=="}, nil, true},
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

func TestKey_Public(t *testing.T) {
	type fields struct {
		Key string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Key
		wantErr bool
	}{
		{"simple public", fields{"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"},
			&Key{"npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7q5"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Key{
				Key: tt.fields.Key,
			}
			got, err := k.ToPublic()
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.Public() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Key.Public() = %v, want %v", got, tt.want)
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
		msg string
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
			args{"AQIDBA=="},
			&Signature{"gbcseiia598sbs4u8p76adr2cgbkhy679867sba4dsaggchk657yzg3f92waeiaaxz8qf7k46cnwt3g2inycttseh38bw5j7nac2jkdg7nywbe7zxi======"},
			false},
		{"public key should error",
			fields{"npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7q5"},
			args{"AQIDBA=="},
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
		Key   string
		SKind string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Address
		wantErr bool
	}{
		{"addr from private key",
			fields{
				"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t",
				"nd",
			},
			&Address{"ndad79yux8we7vk7dgvkqjwnkdhme57piydekb9bkbc6r7uj"}, false},
		{"addr from corresponding public key is the same",
			fields{
				"npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7q5",
				"nd",
			},
			&Address{"ndad79yux8we7vk7dgvkqjwnkdhme57piydekb9bkbc6r7uj"}, false},
		{"addr from private key on testnet",
			fields{
				"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t",
				"tn",
			},
			&Address{"tnad79yux8we7vk7dgvkqjwnkdhme57piydekb9bkbc6rkuf"}, false},
		{"addr from corresponding public key on testnet is the same",
			fields{
				"npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7q5",
				"tn",
			},
			&Address{"tnad79yux8we7vk7dgvkqjwnkdhme57piydekb9bkbc6rkuf"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Key{
				Key: tt.fields.Key,
			}
			got, err := k.NdauAddress(tt.fields.SKind)
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

func TestKey_IsPrivate(t *testing.T) {
	type fields struct {
		Key string
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{"private", fields{"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"}, true, false},
		{"public", fields{"npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7q5"}, false, false},
		{"bad key", fields{"npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7xx"}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Key{
				Key: tt.fields.Key,
			}
			got, err := k.IsPrivate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.IsPrivate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Key.IsPrivate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	privateKey := "npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"
	publicKey := "npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7q5"
	badKey := "npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7xx"
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *Key
		wantErr bool
	}{
		{"private", args{privateKey}, &Key{privateKey}, false},
		{"public", args{publicKey}, &Key{publicKey}, false},
		{"bad key", args{badKey}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
