cd gazernode
call npm run build
cd ..
rem Xcopy gazernode\build ..\bin\www /E /I /Y
go-bindata -pkg web -o res.go gazernode\build/...
