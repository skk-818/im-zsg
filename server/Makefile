api:
	goctl api go -api app/api/desc/api.api -dir app/api -style goZero

rpc:
	goctl rpc protoc app/rpc/pb/rpc.proto --go_out=app/rpc/ --go-grpc_out=app/rpc/ --zrpc_out=app/rpc/

user-model:
	goctl model mysql datasource -url="root:123456@tcp(localhost:3306)/zeroim" -table="im_users"  -dir="app/rpc/model" -cache=true --style=goZero

user-friend-model:
	goctl model mysql datasource -url="root:123456@tcp(localhost:3306)/zeroim" -table="im_user_friends"  -dir="app/rpc/model" -cache=true --style=goZero


edge:
	goctl rpc protoc app/edge/pb/edge.proto --go_out=app/edge/ --go-grpc_out=app/edge/ --zrpc_out=app/edge/