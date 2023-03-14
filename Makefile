# See: https://pkg.go.dev/golang.org/x/text/message#hdr-Translation_Pipeline
updatelocales:
	@echo "Updating locales.."
	pushd cmd/kallu; go generate; popd
	find cmd/kallu/locales -mindepth 1 -maxdepth 2 -type d -exec ./json_merge.py "{}/messages.gotext.json" "{}/../en-US/out.gotext.json" \;
