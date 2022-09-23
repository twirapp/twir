// "parser:gen:go": "protoc --proto_path=./protos --go_out=./parser/ --go_opt=paths=source_relative parser.proto",
// "parser:gen:ts": "protoc --plugin=./node_modules/.bin/protoc-gen-ts --ts_out ./parser/ --ts_opt long_type_string --proto_path ./protos  --ts_opt generate_dependencies parser.proto",
// "bots:gen:go": "protoc --proto_path=./protos --go_out=./bots/ --go_opt=paths=source_relative bots.proto",
// "bots:gen:ts": "protoc --plugin=./node_modules/.bin/protoc-gen-ts --ts_out ./bots/ --ts_opt long_type_string --proto_path ./protos  --ts_opt generate_dependencies bots.proto",
// "build": "rm -rf dist; pnpm parser:gen:go && pnpm parser:gen:ts && pnpm bots:gen:go && pnpm bots:gen:ts &&

import { exec } from 'child_process';
import { readdir } from 'fs/promises';
import { promisify } from 'util';

const files = await readdir('./protos');

const promisedExec = promisify(exec);

await Promise.all(
  files.map(async (proto) => {
    const [name, extension] = proto.split('.');

    const requests = await Promise.all([
      promisedExec(
        `protoc --proto_path=./protos --go_out=./${name}/ --go_opt=paths=source_relative ${name}.proto`,
      ),
      promisedExec(
        `protoc --plugin=./node_modules/.bin/protoc-gen-ts --ts_out ./${name}/ --ts_opt long_type_string --proto_path ./protos --ts_opt generate_dependencies ${name}.proto`,
      ),
    ]);

    console.info(`✅ Genered ${name} proto definitions for go and ts.`);

    return requests;
  }),
);

await promisedExec('pnpm ts:build');
