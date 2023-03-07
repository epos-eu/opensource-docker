: '
    EPOS Open Source - Local installation with Docker
    Copyright (C) 2022  EPOS ERIC

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.'

#!/bin/bash

#Linux

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o out/linux-arm64/epos-linux-arm64 -ldflags="-extldflags=-static" # linux, arm64 arch
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/linux-amd64/epos-linux-amd64 -ldflags="-extldflags=-static" # linux, amd64 arch

#MacOS

CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o out/darwin-arm64/epos-darwin-arm64 -ldflags="-extldflags=-static" # mac, arm64 arch
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o out/darwin-amd64/epos-darwin-amd64 -ldflags="-extldflags=-static" # mac, amd64 arch

#Windows

CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o out/windows-arm64/epos-windows-arm64 -ldflags="-extldflags=-static" # mac, arm64 arch
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o out/windows-amd64/epos-windows-amd64 -ldflags="-extldflags=-static" # windows, amd64 arch