This is a library of Go code that is designed to be used in native apps for android and IOS.

To build it, you need [gomobile](https://godoc.org/golang.org/x/mobile/cmd/gomobile), which you can install with:

```sh
go get golang.org/x/mobile/cmd/gomobile
```

To build IOS, you need XCode -- and its command line tools -- installed, and then you need to do:

```sh
gomobile init
```

After that, you can do:

```sh
gomobile bind -target ios -v
```

This will build IOS and generate Keyaddr.framework.

For Android, you need a current JDK, a current Android Studio, and a current NDK. Install the JDK and Android Studio, then set the location for ANDROID_HOME:

```sh
export ANDROID_HOME=$HOME/Library/Android/sdk/
```

If your android home is not there, then set this appropriately.

Download [the ndk](https://developer.android.com/ndk/downloads/) and put it wherever you like (I put it in my Downloads folder), then tell gomobile how to find it:

```sh
gomobile init -ndk ~/Downloads/android-ndk-r17b/
```

Then you can build the android target like this:

```sh
gomobile bind -target android -v
```

This will generate keyaddr-sources.jar and keyaddr.aar.
