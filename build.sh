if [[ -z "$1" ]]; then
    echo "Please provide a version number"
    exit 1
fi

go build -o bin/gotasks -ldflags "-X 'github.com/okira-e/gotasks/cmd.version=v$1'"