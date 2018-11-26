This little program lists all files in a folder (in this case, the folder´s name is "folder")
It displays:
ID: MD5
Name: File´s name.
Size: Size in bytes
Modified date: The file´s modification time.

Your output should be:
[{"Id":"1593371f4e183af281e836868dbff1e7","Name":"consultas.txt","Size":189,"Modified":"2018-11-24T20:25:48.3378683-04:00"},{"Id":"65c7f06b4e0f8b05d4de6640e5bc29e7","Name":"tarea.jpg","Size":73572,"Modified":"2018-11-25T21:31:57.3363495-04:00"}]

If you have issues when trying to compile it, please run the following commands:

go get -u github.com/golang/dep/cmd/dep
go get github.com/gorilla/mux

(I am assuming that GO is installed on your OS)