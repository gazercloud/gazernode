rem cd gazernode
rem call npm run build
rem cd ..
rem Xcopy gazernode\build ..\bin\www /E /I /Y
go-bindata -pkg web -o res.go gazernode/...
