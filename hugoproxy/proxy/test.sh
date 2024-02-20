#!bash
 go test . -coverprofile cover.out.tmp &&                                               1 ✘   
cat cover.out.tmp | grep -v "_mock.go" > cover.out &&
rm cover.out.tmp &&
go tool cover -func cover.out

