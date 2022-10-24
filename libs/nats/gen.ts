import { exec } from 'child_process';
import { appendFile, readdir } from 'fs/promises';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { promisify } from 'util';

import * as subjects from './subjects.js';

const files = await readdir('./protos');

const promisedExec = promisify(exec);

await Promise.all(
  files.map(async (proto) => {
    const [name] = proto.split('.');

    const requests = await Promise.all([
      promisedExec(
        `protoc --proto_path=./protos --experimental_allow_proto3_optional --go_out=./${name}/ --go_opt=paths=source_relative ${name}.proto`,
      ),
      promisedExec(
        `protoc --plugin=./node_modules/.bin/protoc-gen-ts --experimental_allow_proto3_optional --ts_out ./${name}/ --ts_opt long_type_string --proto_path ./protos --ts_opt generate_dependencies ${name}.proto`,
      ),
    ]);

    console.info(`✅ Genered ${name} proto definitions for go and ts.`);

    return requests;
  }),
);

const __dirname = dirname(fileURLToPath(import.meta.url));

await Promise.all(
  Object.keys(subjects).map(async (pkg) => {
    const pkgName = pkg.toLowerCase();
    const packageDir = resolve(__dirname, pkgName);
    const entries = Object.entries(subjects[pkg]);

    const jsString = `export const SUBJECTS = {${entries.map((e) => `${e[0]}: "${e[1]}"`)}}`;
    const goString = `
  const (
    ${entries.map((e) => `SUBJECTS_${e[0]} = "${e[1]}"`).join('\n')}
    )
    `;

    await Promise.all([
      appendFile(resolve(packageDir, `${pkgName}.ts`), jsString),
      appendFile(resolve(packageDir, `${pkgName}.pb.go`), goString),
    ]);

    console.log(`✅ ${pkg} constants generated.`);
  }),
);

await promisedExec('pnpm ts:build');
console.info(`✅ All done.`);
