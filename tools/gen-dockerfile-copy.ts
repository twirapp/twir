import fs from 'fs';
import path from 'path';

function findPackageJsonFiles(startPath: string) {
	const result: string[] = [];
	if (!fs.existsSync(startPath)) {
		console.log('Directory path does not exist: ', startPath);
		return result;
	}

	const files = fs.readdirSync(startPath);

	for (const file of files) {
		const fullPath = path.join(startPath, file);
		if (!fs.statSync(fullPath).isDirectory()) {
			continue;
		}

		const hasPackage = fs.existsSync(path.join(fullPath, 'package.json'));

		if (hasPackage) {
			result.push(fullPath);
		}
	}

	return result;
}


const paths = [
	...findPackageJsonFiles('./libs'),
	...findPackageJsonFiles('./apps'),
	...findPackageJsonFiles('./frontend'),
];

const generatedStrings = paths.map((p) => `COPY ${p}/package.json ${p}/package.json`).join('\n');

console.log(generatedStrings);
