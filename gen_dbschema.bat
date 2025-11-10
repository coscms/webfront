go get github.com/webx-top/db
go install github.com/webx-top/db/cmd/dbgenerator
dbgenerator -d nging -p root -o ./dbschema -match "^official_(ad|common|customer|page|short_url|i18n)($|_)" -backup "./library/setup/install.sql"

pause