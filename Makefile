test:
	gotestsum --format testname

test-watch:
	gotestsum --watch --format dots
