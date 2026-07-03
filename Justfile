run-whisper:
	./bin/whisper-server \
	-m ./models/ggml-small.en.bin \
	--host 127.0.0.1 \
	--port 8088
