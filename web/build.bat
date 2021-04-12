cd gazernode
call npm run build
cd ..
Xcopy gazernode\build ..\bin\www /E /I /Y
rem go-bindata -pkg httpdata -o res.go www/...

