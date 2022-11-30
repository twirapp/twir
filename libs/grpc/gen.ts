import { exec } from 'node:child_process';
import { existsSync, mkdirSync, rmSync } from 'node:fs';
import { readdir } from 'node:fs/promises';
import { promisify } from 'util';

const files = await readdir('./protos');

const promisedExec = promisify(exec);

rmSync('generated', { recursive: true, force: true });

await Promise.all(
  files.map(async (proto) => {
    const [name] = proto.split('.');

    if (!existsSync(`generated/${name}`)) {
      mkdirSync(`generated/${name}`, { recursive: true });
    }

    const requests = await Promise.all([
      promisedExec(`
      protoc --go_out=./generated/${name} --go_opt=paths=source_relative --go-grpc_out=./generated/${name} --go-grpc_opt=paths=source_relative --proto_path=./protos ${name}.proto`),
      promisedExec(
        `./node_modules/.bin/grpc_tools_node_protoc --plugin=protoc-gen-ts_proto=./node_modules/.bin/protoc-gen-ts_proto --ts_proto_out=./generated/${name} --ts_proto_opt=outputServices=nice-grpc,outputServices=generic-definitions,useExactTypes=false,esModuleInterop=true --proto_path=./protos ${name}.proto`,
      ),
    ]);

    console.info(`✅ Genered ${name} proto definitions for go and ts.`);

    return requests;
  }),
);
