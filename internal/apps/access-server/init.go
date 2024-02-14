package access_server

import _client "Aurora/internal/apps/access-server/pkg/client"

func Init() {
	_client.InitReader()
}
