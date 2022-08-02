PKG_LST=$(go list -f '{{.Dir}}'/... -m | grep -v config)
go test $PKG_LST}
gocov test ${PKG_LST} | gocov report