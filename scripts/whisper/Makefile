ROOT := ../..

WHISPER_DIR := whisper.cpp
BIN_DIR := $(ROOT)/bin
MODEL_DIR := $(ROOT)/models

MODEL_NAME := small.en

CLI := $(BIN_DIR)/whisper-cli
SERVER := $(BIN_DIR)/whisper-server
MODEL := $(MODEL_DIR)/ggml-$(MODEL_NAME).bin

.PHONY: build clean rebuild

build: $(CLI) $(SERVER) $(MODEL)

$(WHISPER_DIR):
	git clone https://github.com/ggml-org/whisper.cpp $(WHISPER_DIR)

$(CLI) $(SERVER): $(WHISPER_DIR)
	cd $(WHISPER_DIR) && \
		cmake -B build && \
		cmake --build build -j

	mkdir -p $(BIN_DIR)

	cp $(WHISPER_DIR)/build/bin/whisper-cli $(CLI)
	cp $(WHISPER_DIR)/build/bin/whisper-server $(SERVER)

$(MODEL): $(WHISPER_DIR)
	mkdir -p $(MODEL_DIR)

	cd $(WHISPER_DIR) && \
		bash models/download-ggml-model.sh $(MODEL_NAME)

	cp $(WHISPER_DIR)/models/ggml-$(MODEL_NAME).bin $(MODEL)

clean:
	rm -rf $(WHISPER_DIR)
	rm -f $(CLI)
	rm -f $(SERVER)
	rm -f $(MODEL)

rebuild: clean build
