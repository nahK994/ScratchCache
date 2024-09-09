echo "1) Start server"
echo "2) Start client"
echo "3) Start playground"
echo "4) Kill process"

read -p "Type: " cmd
if [[ $cmd == 1 ]]; then
    go run cmd/server/main.go
elif [[ $cmd == 2 ]]; then
    go run cmd/client/main.go
elif [[ $cmd == 3 ]]; then
    go run cmd/playground/main.go
elif [[ $cmd == 4 ]]; then
    sudo kill -9 $(sudo lsof -t -i:8000)
fi
