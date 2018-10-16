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
			&Key{"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"}, false},
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

func TestKey_ToPublic(t *testing.T) {
	pvtkey := "npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"
	pubkey := "npuba4jaftckeebzgm7usrcx9jxve8rhst5uejqqtzdtjvhdeswdyzvhn22k98kq25iaaaaaaaaaaaapqhv86syt9pwwpm97n5dgixcmr3sc7ai4km65t9r4wt4s4kywai6fkiae5jkc"
	type fields struct {
		Key string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Key
		wantErr bool
	}{
		{"pub from pvt", fields{pvtkey}, &Key{pubkey}, false},
		{"pub from pub returns itself", fields{pubkey}, &Key{pubkey}, false},
		{"pub from bad key fails", fields{pvtkey + "xxx"}, nil, true},
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
		{"simple private child /1",
			fields{"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"},
			args{1},
			&Key{"npvta8jaftcjec9a4wruhmhf2cgjja658nn8pd5njs73r73hh9gpw9qz53bs26p6wap98fgsaaaaagts2dwuzsn3dsv9mwqm3zsbrwrsxbwavxw36zwsyik47scskjtguawdzw2haetn"},
			false},
		{"simple public child /1",
			fields{"npuba4jaftckeebzgm7usrcx9jxve8rhst5uejqqtzdtjvhdeswdyzvhn22k98kq25iaaaaaaaaaaaapqhv86syt9pwwpm97n5dgixcmr3sc7ai4km65t9r4wt4s4kywai6fkiae5jkc"},
			args{1},
			&Key{"npuba4jaftckeebm6f2y7a4uinxpch96xbea65xdka5mppvwsdg94m4tpxcasqxu6y2bz92w4aaaaaa4gdaquk8bveqcr7qt3rg8af8t8cwgscqyvhu8uc3bmmyakbjge4iibhabh8s7"},
			false},
		{"private child diff index",
			fields{"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"},
			args{2},
			&Key{"npvta8jaftcjeah8avkvz3dt3ccm6h7eqanas7mkaumji87v48zh23bmycpxywj68ap98fgsaaaaaj8igvzve4kvt94g5pvf2t3ikeqaqgea35p9gwtjrwqyn5ip6ru7atiq7w8stysn"},
			false},
		{"private grandchild /1/1",
			fields{"npvta8jaftcjec9a4wruhmhf2cgjja658nn8pd5njs73r73hh9gpw9qz53bs26p6wap98fgsaaaaagts2dwuzsn3dsv9mwqm3zsbrwrsxbwavxw36zwsyik47scskjtguawdzw2haetn"},
			args{1},
			&Key{"npvta8jaftcjebrjxmai8eiyk2s6tyguqwxqbmcrcquqwksqv5wjx84yhpm9q2jykaup36ysaaaaaf3affifpc6x3xh2zs86fnetuta7kqsnyfecana3zzgf9miuv4shpk3ra2ebd25h"},
			false},
		{"public grandchild /1/1",
			fields{"npuba4jaftckeebm6f2y7a4uinxpch96xbea65xdka5mppvwsdg94m4tpxcasqxu6y2bz92w4aaaaaa4gdaquk8bveqcr7qt3rg8af8t8cwgscqyvhu8uc3bmmyakbjge4iibhabh8s7"},
			args{1},
			&Key{"npuba4jaftckeebhmv4ddfe8e9asadxxynr8jxdtk8224rprqbvjcysqmat7iim5hsscjzhu4aaaaaazeawxaxwmuzgw9c8d5sxsugkedxj4bu2wsibsdg862z7pckrka7trbgc8iaqz"},
			false},
		{"bad index",
			fields{"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"},
			args{-1}, nil, true},
		{"bad key",
			fields{"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99xxxmg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"},
			args{11}, nil, true},
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
			fields{"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"},
			args{1},
			&Key{"npvta8jaftcjebe8sxbuvxcmh6xw37wuam8y2tmzachdzqh24mhjyr4j5pxgzw93aap98fg2aaaaaenphqxyh7nh2zhjfugk3a9xvqwkcarfau8239ykec4h69kzkcs8dx2ek3uwekhx"},
			false},
		{"attempt to create hardened child from a public key should fail",
			fields{"npuba4jaftckeebzgm7usrcx9jxve8rhst5uejqqtzdtjvhdeswdyzvhn22k98kq25iaaaaaaaaaaaapqhv86syt9pwwpm97n5dgixcmr3sc7ai4km65t9r4wt4s4kywai6fkiae5jkc"},
			args{1}, nil, true},
		{"bad key should fail",
			fields{"npvt8aaaaaaaaaaaadmt69zefwr5pfdxxxk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"},
			args{1}, nil, true},
		{"bad index should fail",
			fields{"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"},
			args{-11}, nil, true},
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
			fields{"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"},
			args{"AQIDBA=="},
			&Signature{"gbcseiia598sbs4u8p76adr2cgbkhy679867sba4dsaggchk657yzg3f92waeiaaxz8qf7k46cnwt3g2inycttseh38bw5j7nac2jkdg7nywbe7zxi======"},
			false},
		{"public key should error",
			fields{"npuba4jaftckeebzgm7usrcx9jxve8rhst5uejqqtzdtjvhdeswdyzvhn22k98kq25iaaaaaaaaaaaapqhv86syt9pwwpm97n5dgixcmr3sc7ai4km65t9r4wt4s4kywai6fkiae5jkc"},
			args{"AQIDBA=="},
			nil,
			true},
		{"bad key should error",
			fields{"npvt8aaaaaaaaaaaadmt69zefwr5pfdk99mg23ufxxxiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t"},
			args{"AQIDBA=="},
			nil,
			true},
		{"different key should gen diff sig",
			fields{"npvta8jaftcjebe8sxbuvxcmh6xw37wuam8y2tmzachdzqh24mhjyr4j5pxgzw93aap98fg2aaaaaenphqxyh7nh2zhjfugk3a9xvqwkcarfau8239ykec4h69kzkcs8dx2ek3uwekhx"},
			args{"AQIDBA=="},
			&Signature{"gbcaeiapsjv4ry5qtjya4exuyue2fv4a7ju7ytbr563ka4hhjvm9qzvyv2bcau5d9a5qd7gs5i6ynjx92aan2d3crcwr6jd87fqaz4atd85vuv74"},
			false},
		{"bad decode should error",
			fields{"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"},
			args{"AQIDxBA=="},
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
				"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf",
				"nd",
			},
			&Address{"ndad79yux8we7vk7dgvkqjwnkdhme57piydekb9bkbc6r7uj"}, false},
		{"addr from corresponding public key is the same",
			fields{
				"npuba4jaftckeebzgm7usrcx9jxve8rhst5uejqqtzdtjvhdeswdyzvhn22k98kq25iaaaaaaaaaaaapqhv86syt9pwwpm97n5dgixcmr3sc7ai4km65t9r4wt4s4kywai6fkiae5jkc",
				"nd",
			},
			&Address{"ndad79yux8we7vk7dgvkqjwnkdhme57piydekb9bkbc6r7uj"}, false},
		{"addr from private key on testnet",
			fields{
				"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf",
				"tn",
			},
			&Address{"tnad79yux8we7vk7dgvkqjwnkdhme57piydekb9bkbc6rkuf"}, false},
		{"addr from corresponding public key on testnet is the same",
			fields{
				"npuba4jaftckeebzgm7usrcx9jxve8rhst5uejqqtzdtjvhdeswdyzvhn22k98kq25iaaaaaaaaaaaapqhv86syt9pwwpm97n5dgixcmr3sc7ai4km65t9r4wt4s4kywai6fkiae5jkc",
				"tn",
			},
			&Address{"tnad79yux8we7vk7dgvkqjwnkdhme57piydekb9bkbc6rkuf"}, false},
		{"addr from bad key fails",
			fields{
				"npvt8aaaaaaaaaaaadmt69zefwxr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t",
				"nd",
			},
			nil, true},
		{"addr from bad chain fails",
			fields{
				"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf",
				"qz",
			},
			nil, true},
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
		{"private", fields{"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"}, true, false},
		{"public", fields{"npuba4jaftckeebzgm7usrcx9jxve8rhst5uejqqtzdtjvhdeswdyzvhn22k98kq25iaaaaaaaaaaaapqhv86syt9pwwpm97n5dgixcmr3sc7ai4km65t9r4wt4s4kywai6fkiae5jkc"}, false, false},
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
	privateKey := "npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"
	publicKey := "npuba4jaftckeebzgm7usrcx9jxve8rhst5uejqqtzdtjvhdeswdyzvhn22k98kq25iaaaaaaaaaaaapqhv86syt9pwwpm97n5dgixcmr3sc7ai4km65t9r4wt4s4kywai6fkiae5jkc"
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

func ch(s string, n int32) string {
	pk, _ := FromString(s)
	ch, err := pk.Child(n)
	if err != nil {
		panic(err)
	}
	return ch.Key
}

func hch(s string, n int32) string {
	pk, _ := FromString(s)
	ch, err := pk.HardenedChild(n)
	if err != nil {
		panic(err)
	}
	return ch.Key
}

func pub(s string) string {
	pk, _ := FromString(s)
	ch, err := pk.ToPublic()
	if err != nil {
		panic(err)
	}
	return ch.Key
}

func TestDeriveFrom(t *testing.T) {
	// let's make sure that DeriveFrom creates keys with the right sequence of derivations
	privateKey1 := "npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"
	privateKey2 := ch(privateKey1, 456)
	privateKey3 := ch(privateKey2, 789)
	privateKey4 := ch(privateKey2, 1)
	publicKey1 := pub(privateKey1)
	publicKey2 := pub(privateKey2)
	hprivateKey1 := hch(privateKey1, 456)
	badKey := "npubaaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacga5vf83ihtk9w43urhv2i73cezhi5t2w3vtuikb5m3vynnfr9fhnpxzbg7xx"
	type args struct {
		parentKey  string
		parentPath string
		childPath  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Key
		wantErr bool
	}{
		{"private one level", args{privateKey1, "/123", "/123/456"}, &Key{privateKey2}, false},
		{"private next level", args{privateKey2, "/123/456", "/123/456/789"}, &Key{privateKey3}, false},
		{"private next level diff", args{privateKey2, "/123/456", "/123/456/1"}, &Key{privateKey4}, false},
		{"private two levels", args{privateKey1, "/123", "/123/456/789"}, &Key{privateKey3}, false},
		{"private two levels diff", args{privateKey1, "/123", "/123/456/1"}, &Key{privateKey4}, false},
		{"public one level", args{publicKey1, "/123", "/123/456"}, &Key{publicKey2}, false},
		{"hardened private one level", args{privateKey1, "/123", "/123/456'"}, &Key{hprivateKey1}, false},
		{"attempt to create hardened from pubkey fails", args{publicKey1, "/123", "/123/456'"}, nil, true},
		{"bad key fails", args{badKey, "/123", "/123/456"}, nil, true},
		{"unrelated paths fails", args{privateKey1, "/123", "/456/789"}, nil, true},
		{"same path fails", args{privateKey1, "/123", "/123"}, nil, true},
		{"bad path fails", args{privateKey1, "/123x", "/123/456"}, nil, true},
		{"bad path fails", args{privateKey1, "/123", "/123/456x"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeriveFrom(tt.args.parentKey, tt.args.parentPath, tt.args.childPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeriveFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeriveFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}
