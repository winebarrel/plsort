package plsort

type Options struct {
	Creds string `kong:"default='client_secret.json',help='OAuth client ID credentials JSON path.'"`
	SortOptions
}

type SortOptions struct {
	PlaylistId string `kong:"arg='',required,help='YouTube playlist ID.'"`
	Reverse    bool   `kong:"short='r',help='Sort in reverse order.'"`
}
