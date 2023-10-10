importScripts("wasm_exec.js");

const go = new Go();

let wasmReady = false;

async function initWasm() {
  const wasmResponse = await fetch("main.wasm");
  const wasmBytes = await wasmResponse.arrayBuffer();
  const { instance } = await WebAssembly.instantiate(
    wasmBytes,
    go.importObject
  );
  go.run(instance);
  wasmReady = true;
}

initWasm();

self.onmessage = function (event) {
  const callbackChunk = (...args) => {
    let [
      chunkType,
      offset,
      uncompressSized,
      filePath,
      content,
      packageName,
      prefix,
    ] = args;
    content = JSON.parse(content);
    delete content.Reader;
    delete content.Chunk;
    content = JSON.stringify(content, null, "    ");
    self.postMessage({
      action: "callbackChunk",
      args: [
        chunkType,
        offset,
        uncompressSized,
        filePath,
        content,
        packageName,
        prefix,
      ],
    });
  };

  const callbackBinary = (...args) => {
    self.postMessage({ action: "callbackBinary", args });
  };

  const { action, data } = event.data;

  if (!wasmReady) {
    self.postMessage({
      action: "error",
      message: "WASM is not initialized yet.",
    });
    return;
  }

  switch (action) {
    case "processFile":
      processFile(data.fileName, data.byteArray, callbackChunk, callbackBinary);
      break;
  }
};
