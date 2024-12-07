#!/bin/bash

rm /usr/bin/ccommits-cli
go build -o ccommits-cli
mv ccommits-cli /usr/bin