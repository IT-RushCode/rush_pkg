push-git: mod
	git add .
	git commit -m "Update rush_pkg"
	git push

update-mod:
	go get -u ./...

mod:
	go mod tidy && go mod vendor
