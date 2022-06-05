import * as fs from 'fs';
import { fileURLToPath } from 'node:url';
import * as path from 'path';

import glob from 'glob';

function createEntitiesIndex() {
  const src = path.resolve(fileURLToPath(import.meta.url), '..', '..', 'src', 'entities');
  if (!fs.existsSync(src)) {
    console.log(`App api cannot be found. Path not exist: ${src}`);
    process.exit(1);
  }
  const tmpFile = `${src}/tmp-index.ts`;
  const outFile = `${src}/index.ts`;

  for (const item of glob.sync(`${src}/**/*.ts`)) {
    const filePath = path.relative(src, item).replace(/\.ts$/, '');
    const data = `export * from './${filePath}.js';\n`;
    fs.writeFileSync(tmpFile, data, { flag: 'a+' });
  }
  if (fs.existsSync(outFile) && fs.existsSync(tmpFile)) {
    fs.unlinkSync(outFile);
    console.log(`Old file '${outFile}' removed`);
  }
  if (fs.existsSync(tmpFile)) {
    fs.renameSync(tmpFile, outFile);
    console.log(`New file ${outFile} saved`);
  }
}
createEntitiesIndex();