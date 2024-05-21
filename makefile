push-git: mod
	git add .
	git commit -m "Update rush_pkg"
	git push

update-all-packets:
	go get -u ./...
	go mod tidy

mod:
	go mod tidy && go mod vendor
