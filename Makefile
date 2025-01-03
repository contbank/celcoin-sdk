define get_version
$(shell cat current_version)
endef

VERSION=$(call get_version,)

install-autotag:
	wget -O autotag https://github.com/pantheon-systems/autotag/releases/download/1.1.1/Linux && sudo chmod +x autotag && sudo mv autotag /usr/local/bin/

set_version:
	autotag > current_version

run-rebase:
	git rebase -p master 2>/dev/null | grep "Your branch is up-to-date with 'origin/master'." || echo "\nPlease rebase your branch with master!"

run-tests:
	go test -failfast -cover ./...

build-package:
	go mod vendor

tag-version: set_version
	git tag -d v$(VERSION) && git tag v$(VERSION) && git push --tags

mockery:
ifeq (,$(wildcard $(GOPATH)/bin/mockery))
	echo "Installing mockery"
	curl https://github.com/vektra/mockery/releases/download/v2.14.0/mockery_2.14.0_$(shell uname -s)_x86_64.tar.gz \
	-o $(GOPATH)/bin/mockery
endif

mock: mockery
	mockery --disable-version-string --with-expecter --keeptree --all
