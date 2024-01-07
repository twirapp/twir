import { exec } from 'node:child_process';
import { existsSync, mkdirSync, rmSync } from 'node:fs';
import { readFile, unlink, readdir, writeFile } from 'node:fs/promises';
import { platform } from 'node:os';
import { resolve, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';
import { promisify } from 'util';

const promisedExec = promisify(exec);

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
const protocGenTwirpPath = resolve(
	cwd,
	'.bin',
	`protoc-gen-twirp${binSuffixExe}`,
);
const protocGenTsPath = resolve(
	__dirname,
	'node_modules',
	'.bin',
	`protoc-gen-ts_proto${binSuffixCmd}`,
);
console.dir({
	workDir,
	cwd,
	protocPath,
	protocGenGoPath,
	protocGenGoGrpcPath,
	protocGenTwirpPath,
	protocGenTsPath,
});

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
					`--plugin=go=${protocGenGoPath}`,
					`--plugin=go-grpc=${protocGenGoGrpcPath}`,
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
					`--plugin=protoc-gen-ts_proto=${protocGenTsPath}`,
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

	await promisedExec([
		protocPath,
		`--plugin=ts=${protocGenTwirpPath}`,
		`--ts_out=${workDir}/generated/api`,
		`--ts_opt=optimize_code_size`,
		`--ts_opt=generate_dependencies,eslint_disable`,
		`--experimental_allow_proto3_optional`,
		`--proto_path=${workDir}/protos`,
		`api.proto`,
	].join(' '));


	await promisedExec([
		protocPath,
		`--plugin=twirp=${protocGenTwirpPath}`,
		`--plugin=go=${protocGenGoPath}`,
		`--go_opt=paths=source_relative`,
		`--twirp_opt=paths=source_relative`,
		`--go_out=${workDir}/generated/api`,
		`--twirp_out=${workDir}/generated/api`,
		`	--proto_path=${workDir}/protos`,
		`api.proto`,
	].join(' '));

	const apiFiles = await readdir(`${workDir}/protos/api`);

	for (const file of apiFiles) {
		const name = file.replace('.proto', '');
		const directoryPath = `${workDir}/generated/api/${name}`;
		const filePath = `${directoryPath}/${name}.pb.go`;
		mkdirSync(directoryPath, { recursive: true });

		await promisedExec([
			protocPath,
			`--plugin=go=${protocGenGoPath}`,
			`--go_opt=paths=source_relative`,
			`--go_out=${directoryPath}`,
			`--experimental_allow_proto3_optional`,
			`--proto_path=${workDir}/protos`,
			`api/${file}`,
		].join(' '));

		// найс костыль кекв / cool crutch kekw
		const sourceFilePath = `${directoryPath}/api/${name}.pb.go`;

		try {
			const fileData = await readFile(sourceFilePath, 'utf8');
			await writeFile(filePath, fileData, 'utf8');
			await unlink(sourceFilePath);
		} catch {
			process.exit(1);
		}
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
