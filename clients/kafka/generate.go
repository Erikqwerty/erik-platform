package kafka

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate $HOME/go/bin/minimock -i Producer  -o ./mocks/ -s "_minimock.go"
//go:generate $HOME/go/bin/minimock -i Consumer  -o ./mocks/ -s "_minimock.go"
