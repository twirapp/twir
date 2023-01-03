import { exec } from 'node:child_process';
import { existsSync, mkdirSync, readFileSync, rmSync, statSync } from 'node:fs';
import { readdir } from 'node:fs/promises';
import { platform } from 'node:os';
import { resolve } from 'node:path';
import { promisify } from 'util';

const promisedExec = promisify(exec);

rmSync('generated', { recursive: true, force: true });

(async () => {
  const files = await readdir('./protos');
  // const pluginDir = isDocker()
  //   ? '/app/libs/gprc/node_modules/.bin/protoc-gen-ts_proto'
  //   : './node_modules/.bin/protoc-gen-ts_proto';
  await Promise.all(
    files
      .filter((n) => n != 'google')
      .map(async (proto) => {
        const [name] = proto.split('.');

        if (!existsSync(`generated/${name}`)) {
          mkdirSync(`generated/${name}`, { recursive: true });
        }

        const protocPath = resolve(
          __dirname,
          'node_modules',
          '.bin',
          `protoc-gen-ts_proto${platform() === 'win32' ? '.CMD' : ''}`,
        );

        const requests = await Promise.all([
          promisedExec(`
        protoc --go_out=./generated/${name} --go_opt=paths=source_relative --go-grpc_out=./generated/${name} --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional --proto_path=./protos ${name}.proto`),
          promisedExec(
            `protoc --plugin=protoc-gen-ts_proto=${protocPath} --ts_proto_out=./generated/${name} --ts_proto_opt=outputServices=nice-grpc,outputServices=generic-definitions,useExactTypes=false,esModuleInterop=true --experimental_allow_proto3_optional --proto_path=./protos ${name}.proto`,
          ),
        ]);

        console.info(`✅ Genered ${name} proto definitions for go and ts.`);

        return requests;
      }),
  );
})();

// function hasDockerEnv() {
//   try {
//     statSync('/.dockerenv');
//     return true;
//   } catch {
//     return false;
//   }
// }

// function hasDockerCGroup() {
//   try {
//     return readFileSync('/proc/self/cgroup', 'utf8').includes('docker');
//   } catch {
//     return false;
//   }
// }

// export default function isDocker() {
//   return hasDockerEnv() || hasDockerCGroup();
// }
