// Bun embeds the .node file when require() is called with a literal path
// The correct platform-specific file is copied by scripts/copy-native.ts
require('./backend.node')
