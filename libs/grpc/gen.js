import { exec } from 'node:child_process';
import { existsSync, mkdirSync, rmSync } from 'node:fs';
import { readdir } from 'node:fs/promises';
import { platform } from 'node:os';
import { resolve, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));

rmSync('generated', { recursive: true, force: true });

const workDir = resolve(__dirname);
const cwd = resolve(workDir, '..', '..');

const binSuffixExe = platform() === 'win32' ? '.exe' : '';
const binSuffixCmd = platform() === 'win32' ? '.CMD' : '';

const protocPath = resolve(
	cwd,
	'.bin',
	`protoc${binSuffixExe}`,
);
const protocGenGoPath = resolve(
	cwd,
	'.bin',
	`protoc-gen-go${binSuffixExe}`,
);
const protocGenGoGrpcPath = resolve(
	cwd,
	'.bin',
	`protoc-gen-go-grpc${binSuffixExe}`,
);
const protocGenTsPath = resolve(
	__dirname,
	'node_modules',
	'.bin',
	`protoc-gen-ts_proto${binSuffixCmd}`,
);

const pathName = platform() === 'win32' ? 'Path' : 'PATH';

const execOptions = {
	env: {
		...process.env,
		[pathName]: `${process.env[pathName]}:${cwd}/.bin`,
	},
};

const promisedExec = async (cmd) => {
	return new Promise((resolve, reject) => {
		exec(cmd, execOptions, (err, stdout, stderr) => {
			if (err) {
				console.error(stderr);
				reject(err);
			} else {
				resolve(stdout);
			}
		});
	});
};

(async () => {
	const files = await readdir(`${workDir}/protos`);

	const ignoredFiles = [
		'api',
	];
	await Promise.all(
		files
			.filter((n) => n != 'google')
			.map(async (proto) => {
				const [name] = proto.split('.');

				if (!existsSync(`${workDir}/generated/${name}`)) {
					mkdirSync(`${workDir}/generated/${name}`, { recursive: true });
				}

				if (ignoredFiles.includes(name)) {
					return;
				}

				const goBuildString = [
					protocPath,
					`--plugin="go=${protocGenGoPath}"`,
					`--plugin="go-grpc=${protocGenGoGrpcPath}"`,
					`--go_out=${workDir}/generated/${name}`,
					`--go_opt=paths=source_relative`,
					`--go-grpc_out=${workDir}/generated/${name}`,
					`--go-grpc_opt=paths=source_relative`,
					`--experimental_allow_proto3_optional`,
					`--proto_path=${workDir}/protos`,
					`${name}.proto`,
				].join(' ');

				const nodeBuildString = [
					protocPath,
					`--plugin="protoc-gen-ts_proto=${protocGenTsPath}"`,
					`--ts_proto_out=${workDir}/generated/${name}`,
					'--ts_proto_opt=outputServices=nice-grpc',
					'--ts_proto_opt=outputServices=generic-definitions',
					'--ts_proto_opt=useExactTypes=false',
					'--ts_proto_opt=esModuleInterop=true',
					'--ts_proto_opt=importSuffix=.js',
					'--experimental_allow_proto3_optional',
					`--proto_path=${workDir}/protos`,
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

	console.info(`✅ Generated api proto definitions for go and ts.`);
})();
