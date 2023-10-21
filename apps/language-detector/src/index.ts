import { PORTS } from '@twir/grpc/constants/constants';
import * as LanguageDetector from '@twir/grpc/generated/language-detector/language-detector';
import cld from 'cld';
import { createServer } from 'nice-grpc';

const service = {
	async detect({ text }) {
		const result = await cld.detect(text);

		return {
			languages: result.languages.map(l => ({ name: l.name, score: l.score, percent: l.percent })),
			chunks: result.chunks.map(c => ({ name: c.name, offset: c.offset, bytes: c.bytes })),
		};
	},
};

const server = createServer({
	'grpc.keepalive_time_ms': 60 * 1000,
});

server.add(LanguageDetector.LanguageDetectorDefinition, service);

await server.listen(`0.0.0.0:${PORTS.LANGUAGE_DETECTOR_SERVER_PORT}`);
console.log('Language detector microservice started');
