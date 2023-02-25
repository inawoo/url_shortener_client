
git config --global url."ssh://git@github.com".insteadOf "https://github.com"
git config --list --show-origin
go env -w GOPRIVATE="github.com/inawoo/url_shortener"