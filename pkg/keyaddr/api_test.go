package keyaddr

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/oneiro-ndev/ndaumath/pkg/key"
	"github.com/stretchr/testify/require"
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
			&Signature{"aujaftchgbcseiiay7mr4bc69rgj3g4dnf82ijp8t22rnstjerrbwujy5mkbp2382vmseiahffbgsf6aujtn6vs5m3jhnm7qt952f59xemarjrcycphty726ybcgx82x"},
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
			&Signature{"ayjaftcggbcaeibatkbh6whxa5qvica5jp7btt3v44z6mtzcyuwad8xifbsnb3rzp6bcan54rfdkmvif8x4z9p5gwrx24f4qx9ukfx2k5u8rjkfgq5fmpmparr9zkyw9"},
			false},
		{"bad decode should error",
			fields{"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf"},
			args{"AQIDxBA=="},
			nil,
			true},
		{"real life test",
			fields{"npvta8jaftcjed4m8juti4sdgku42kq4g6nckgsprvexjb29jrewmznzshyatbzp6ceik9yaaaaaahvccndjymyem3tknrbsikfwrhfiwr6xg4g5md3gicyzd7adcmn4yn94rgjrgp4r"},
			args{"bnB1YmE0amFmdGNrZWViNjJ1YWFqOG41eGh1YnJoYjY4dWM2ZXR4bnltYjJqZWNwdnZ4YXZnZGEzYjRhZm1lcm4yMmV1d3JtZ2FhYWFhYXZoN3p4Y2d6ZnA0ZTd3bnozMjI5N3I4emd5N2V1OXRubW15MmRyODhkNHd4ejJ0bWppeWZ5ZDl6aW5hOHYAAAAAAAAAAW5kYXFmZ2Y5bnNkNTZlaGV0OWNxemZneTQ3cHd1YjM0ZDNjYWJha2l6ejN6ejJ4Zm5wdWJhNGphZnRja2VlYnZqYTV5c2VtZmIyZm16cjRrMnFjd2NlZWo3M3ZrYmQ2dzc2M2E1djZjYzZ1N3VkZGNoM2lpdGJtODJhYWFhYWE4bmlqc3BnM3EydHJnZmp0NmdiYml5dDY2eGN0OXV3NWk1cG5yZTNhazQ2cndhbmpwdms5NnNrd3ZwbjNp"},
			&Signature{"ayjaftcggbcaeiatiugceq4kkh6fwmwrnrzgqqc78xecd9u95djmxwngmip2wm6ut2bca7ghwhh2x7ftdies7df6zrbz8dh7r9br3bzxu8kri34jtrjfuks5yutxr2c3"},
			false},
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
			fields{
				"npvta8jaftcjebc56pvxgs8w2448fibvc4yqeub8b49b7k4tdg7t5dsdhayzi569eaaaaaaaaaaaadmt69zefwr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxachhw8sfiuejtf",
			},
			&Address{"ndad79yux8we7vk7dgvkqjwnkdhme57piydekb9bkbc6r7uj"}, false},
		{"addr from corresponding public key is the same",
			fields{
				"npuba4jaftckeebzgm7usrcx9jxve8rhst5uejqqtzdtjvhdeswdyzvhn22k98kq25iaaaaaaaaaaaapqhv86syt9pwwpm97n5dgixcmr3sc7ai4km65t9r4wt4s4kywai6fkiae5jkc",
			},
			&Address{"ndad79yux8we7vk7dgvkqjwnkdhme57piydekb9bkbc6r7uj"}, false},
		{"addr from bad key fails",
			fields{
				"npvt8aaaaaaaaaaaadmt69zefwxr5pfdk99mg23ufiu58nazicguu9g6r58xeqwguxxacgacfz25hkpb7jtxx6ksdgfxn6jed6dx8d4xxcgp5dyhagqbpqtz38kcrgm4t",
			},
			nil, true},
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

func TestPublicKey_SignaturePackageRoundtrip(t *testing.T) {
	seed, err := key.GenerateSeed(key.RecommendedSeedLen)
	require.NoError(t, err)
	ekeyPriv, err := key.NewMaster(seed)
	require.NoError(t, err)
	ekeyPub, err := ekeyPriv.Public()
	require.NoError(t, err)

	ka1, err := KeyFromExtended(ekeyPub)
	require.NoError(t, err)

	ka1B := []byte(ka1.Key)

	sk1, err := ka1.ToPublicKey()
	require.NoError(t, err)
	sk1B, err := sk1.MarshalText()
	require.NoError(t, err)
	require.Equal(t, ka1B, sk1B)

	ka2, err := KeyFromPublic(sk1)
	require.NoError(t, err)
	ka2B := []byte(ka2.Key)

	require.Equal(t, ka1B, ka2B)
	require.Equal(t, ka1, ka2)

	sk2, err := ka2.ToPublicKey()
	require.NoError(t, err)
	sk2B, err := sk2.MarshalText()
	require.NoError(t, err)

	require.Equal(t, ka2B, sk2B)
	require.Equal(t, sk1B, sk2B)
	require.Equal(t, sk1, sk2)
}

func TestPrivateKey_SignaturePackageRoundtrip(t *testing.T) {
	seed, err := key.GenerateSeed(key.RecommendedSeedLen)
	require.NoError(t, err)
	ekeyPriv, err := key.NewMaster(seed)
	require.NoError(t, err)

	ka1, err := KeyFromExtended(ekeyPriv)
	require.NoError(t, err)

	ka1B := []byte(ka1.Key)

	sk1, err := ka1.ToPrivateKey()
	require.NoError(t, err)
	sk1B, err := sk1.MarshalText()
	require.NoError(t, err)
	require.Equal(t, ka1B, sk1B)

	ka2, err := KeyFromPrivate(sk1)
	require.NoError(t, err)
	ka2B := []byte(ka2.Key)

	require.Equal(t, ka1B, ka2B)
	require.Equal(t, ka1, ka2)

	sk2, err := ka2.ToPrivateKey()
	require.NoError(t, err)
	sk2B, err := sk2.MarshalText()
	require.NoError(t, err)

	require.Equal(t, ka2B, sk2B)
	require.Equal(t, sk1B, sk2B)
	require.Equal(t, sk1, sk2)
}

func TestDebugAPIProblem(t *testing.T) {
	// This makes sure that DeriveFrom creates lots of derived keys without failure.
	// The particular sequence of bytes below drove an issue with the Child function where sometimes
	// keys would be returned that were shorter than desired because the key was just representing
	// a large number, and the big number package was trimming leading zeros -- so 1/256 of the time
	// it would return 31 bytes instead of 32. (And 1/65536 times it would return 30 bytes...)
	// This byte pattern generates that key on the 25th attempt.
	k, err := NewKey("ZDZjNmM5ZmNmYWJiODdkNg==")
	require.Nil(t, err)
	rootPrivateKey1, err := k.ToPrivateKey()
	require.Nil(t, err)
	s, err := rootPrivateKey1.MarshalText()
	require.Nil(t, err)
	for i := 0; i < 500; i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			childpath := fmt.Sprintf("/44'/20036'/100/%d", i)
			derivedKey, err := DeriveFrom(string(s), "/", childpath)
			require.Nil(t, err)
			x, err := derivedKey.IsPrivate()
			require.Nil(t, err)
			require.True(t, x)
			addr, err := derivedKey.NdauAddress()
			require.Nil(t, err)
			require.True(t, strings.HasPrefix(addr.Address, "nda"))
		})
	}
}
