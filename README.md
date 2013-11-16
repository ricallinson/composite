# fcomposite

[![Build Status](https://secure.travis-ci.org/ricallinson/fcomposite.png?branch=master)](http://travis-ci.org/ricallinson/fcomposite)

__WARNING: WORK IN PROGRESS__

## Testing

The following should all be executed from the `forgery` directory _$GOPATH/src/github.com/ricallinson/forgery/_.

#### Install

    go get github.com/ricallinson/simplebdd

#### Run

    go test

### Code Coverage

#### Install

    go get github.com/axw/gocov/gocov
    go get -u github.com/matm/gocov-html

#### Generate

    gocov test | gocov-html > ./reports/coverage.html