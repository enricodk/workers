import "../assets/polyfill_performance.js";
import "../assets/wasm_exec.js";
import mod from "./dist/app.wasm";

const go = new Go();

const load = WebAssembly.instantiate(mod, go.importObject).then((instance) => {
  go.run(instance);
  return instance;
});

export default {
  async fetch(req) {
    await load;
    return handleRequest(req);
  },
};
