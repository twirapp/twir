import { exec } from 'node:child_process';
import { readdir } from 'node:fs/promises';
import { platform } from 'node:os';
import { dirname, resolve, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));

const currentDir = resolve(__dirname);
const projectRoot = resolve(currentDir, '..', '..');

const pathName = platform() === 'win32' ? 'Path' : 'PATH';
const binSuffixExe = platform() === 'win32' ? '.exe' : '';

const execOptions = {
	env: {
		...process.env,
		[pathName]: `${process.env[pathName]}:${projectRoot}/.bin`,
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

const protocPath = resolve(
	projectRoot,
	'.bin',
	`protoc${binSuffixExe}`,
);

const messagesDirPath = join(currentDir, 'messages');
const messagesDirs = await readdir(messagesDirPath);

await promisedExec([
	protocPath,
	`--go_opt=paths=source_relative`,
	`--twirp_opt=paths=source_relative`,
	`--go_out=.`,
	`--twirp_out=.`,
	`--experimental_allow_proto3_optional`,
	`--proto_path=.`,
	// `${messagesDirs.map(d => `messages/${d}/${d}.proto`).join(' ')} `,
	// `${currentDir}/protos/*.proto`,
	'*.proto',
].join(' '), (err) => {
	console.log(err);
	process.exit(1);
});


// const outDir = join(currentDir, 'messages');
//
// for (const file of messagesFiles) {
// 	const name = file.replace('.proto', '');
//
// 	mkdirSync(outDir, { recursive: true });
//
// 	await promisedExec([
// 		protocPath,
// 		`--go_opt=paths=source_relative`,
// 		`--go_out=.`,
// 		`--twirp_out=.`,
// 		`--experimental_allow_proto3_optional`,
// 		`--proto_path=${join(currentDir, 'protos')}`,
// 		`messages/alerts/alerts.proto`,
// 	].join(' '));
// }
