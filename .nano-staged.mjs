export default {
    "*": (api) => `prettier -luw ${api.filenames.join(" ")}`,
    "*.go": () => ["go fmt", "go fix"],
};
