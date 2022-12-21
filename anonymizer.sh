#!/bin/sh
sed -ri 's/dogu/cmd/' main.go
sed -ri 's/(const app_.*string=).*/\1"none"/' config.go
