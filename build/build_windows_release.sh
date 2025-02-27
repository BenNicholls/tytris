#!/bin/bash

echo "BUILDING RELEASE FOR WINDOWS"
env CGO_ENABLED="1" CC="/usr/bin/x86_64-w64-mingw32-gcc" GOOS="windows" CGO_LDFLAGS="-lmingw32 -lSDL2" CGO_CFLAGS="-D_REENTRANT" go build -C ../game/ -ldflags  "-H=windowsgui" -tags audio

if [ $? -eq "0" ]
then 
echo "BUILD COMPLETE"
mv ../game/game.exe tytris_release.exe
else echo "POBLEMS, BUILD DID NOT GO WELL"
fi