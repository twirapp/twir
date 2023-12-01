import fs from 'fs';
import path from 'path';

const esmPath = path.join('./dist/esm', 'package.json');
const cjsPath = path.join('./dist/cjs', 'package.json');

fs.writeFileSync(esmPath, JSON.stringify({ type: 'module' }, null, 2));
fs.writeFileSync(cjsPath, JSON.stringify({ type: 'commonjs' }, null, 2));
