test:
	gotestsum --format testname

test-watch:
	gotestsum --watch --format dots

examples.tity:
	find examples -depth 1 -exec bash -c "cd '{}' && go mod tidy" \;
