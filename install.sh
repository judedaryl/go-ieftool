if [[ $1 == "" ]]
then
    VER=v1.0.0
else
    VER=v$1
fi

echoerr() { echo "$@" 1>&2; }
if [[ ! ":$PATH:" == *":/usr/local/bin:"* ]]; then
    echoerr "Your path is missing /usr/local/bin, you need to add this to use this installer."
    exit 1
fi
if [ "$(uname)" == "Darwin" ]; then
    OS=darwin
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    OS=linux
else
    echoerr "This installer is only supported on Linux and MacOS"
    exit 1
fi

ARCH="$(uname -m)"
if [ "$ARCH" == "x86_64" ]; then
    ARCH=amd64
elif [[ "$ARCH" == aarch* ]]; then
    ARCH=arm
else
    echoerr "unsupported arch: $ARCH"
    exit 1
fi

DOWNLOAD_URL=https://github.com/judedaryl/go-ieftool/releases/download/$VER/ieftool-$OS-$ARCH

echo "Installing ieftool from $DOWNLOAD_URL"

if [ $(command -v curl) ]; then
curl -sL "$DOWNLOAD_URL" -o ieftool
else
wget -O- "$DOWNLOAD_URL"
fi

echo $ARCH
echo $OS