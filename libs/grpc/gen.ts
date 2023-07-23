import { exec } from 'node:child_process';
import { existsSync, mkdirSync, rmSync } from 'node:fs';
import { readdir } from 'node:fs/promises';
import { platform } from 'node:os';
import { resolve } from 'node:path';
import { promisify } from 'util';

const promisedExec = promisify(exec);

rmSync('generated', { recursive: true, force: true });

(async () => {
  const files = await readdir('./protos');

	const ignoredFiles = [
		'api',
	];
  await Promise.all(
    files
      .filter((n) => n != 'google')
      .map(async (proto) => {
        const [name] = proto.split('.');

        if (!existsSync(`generated/${name}`)) {
          mkdirSync(`generated/${name}`, { recursive: true });
        }

				if (ignoredFiles.includes(name)) {
					return;
				}

        const protocPath = resolve(
          __dirname,
          'node_modules',
          '.bin',
          `protoc-gen-ts_proto${platform() === 'win32' ? '.CMD' : ''}`,
        );

				const goBuildString = [
					`protoc --go_out=./generated/${name}`,
					`--go_opt=paths=source_relative`,
					`--go-grpc_out=./generated/${name}`,
					`--go-grpc_opt=paths=source_relative`,
					`--experimental_allow_proto3_optional`,
					`--proto_path=./protos`,
					`${name}.proto`,
				].join(' ');

				const nodeBuildString = [
					'protoc',
					`--plugin=protoc-gen-ts_proto=${protocPath}`,
					`--ts_proto_out=./generated/${name}`,
					'--ts_proto_opt=outputServices=nice-grpc',
					'--ts_proto_opt=outputServices=generic-definitions',
					'--ts_proto_opt=useExactTypes=false',
					'--ts_proto_opt=esModuleInterop=true',
					'--ts_proto_opt=importSuffix=.js',
					'--experimental_allow_proto3_optional',
					'--proto_path=./protos',
					`${name}.proto`,
				].join(' ');

        const requests = await Promise.all([
          promisedExec(goBuildString),
          promisedExec(nodeBuildString),
        ]);

        console.info(`✅ Generated ${name} proto definitions for go and ts.`);
        return requests;
      }),
  );

	await promisedExec(`protoc --experimental_allow_proto3_optional --ts_out ./generated/api --ts_opt=generate_dependencies,eslint_disable --proto_path ./protos api.proto`);
	await promisedExec(`protoc --experimental_allow_proto3_optional --go_opt=paths=source_relative --twirp_opt=paths=source_relative --go_out=./generated/api --twirp_out=./generated/api --proto_path=./protos api.proto`);
	const apiFiles = await readdir('./protos/api');

	for (const file of apiFiles) {
		const name = file.replace('.proto', '');
		const directoryPath = `./generated/api/${name}`;
		const filePath = `${directoryPath}/${name}.pb.go`;
		mkdirSync(directoryPath, { recursive: true });
		await promisedExec(`protoc --experimental_allow_proto3_optional --go_opt=paths=source_relative --go_out=${directoryPath} --proto_path=./protos api/${file}`);
		// найс костыль кекв / cool crutch kekw
		await promisedExec(`mv ${directoryPath}/api/${name}.pb.go ${filePath}`);
	}

	console.info(`✅ Generated api proto definitions for go and ts.`);

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
