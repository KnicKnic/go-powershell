# install
Windows - install dependencies - Use choco

1. `choco install mingw -y`

## hacks
1. `copy "C:\Program Files (x86)\Windows Kits\NETFXSDK\4.7.2\Include\um\metahost.h" .`
1. `copy "C:\Program Files (x86)\Windows Kits\NETFXSDK\4.7.2\lib\um\x64\mscoree.lib" .`


1. https://code.google.com/archive/p/lib2a/downloads
```

How it works

The conversion process is accomplished in several steps:
1) Copy your .LIB file and .DLL file into the "convert" folder.
2) Edit and replace the files names in the four first lines of the LIB2A.bat.
3) Run LIB2A.bat.

You can find your .A linker library into the "convert" folder.
```

# Docs
https://docs.microsoft.com/en-us/dotnet/framework/unmanaged-api/hosting/clrcreateinstance-function

https://grokbase.com/t/gg/golang-nuts/154m672a6t/go-nuts-linking-cgo-with-visual-studio-x64-release-libraries-on-windows
