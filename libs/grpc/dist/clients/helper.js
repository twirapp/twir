"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.waitReady = exports.CLIENT_OPTIONS = exports.createClientAddr = void 0;
const nice_grpc_1 = require("nice-grpc");
const createClientAddr = (env, service, port) => {
    let ip = service;
    if (env != 'production') {
        ip = '127.0.0.1';
    }
    return `${ip}:${port}`;
};
exports.createClientAddr = createClientAddr;
exports.CLIENT_OPTIONS = {
    'grpc.lb_policy_name': 'round_robin',
    'grpc.service_config': JSON.stringify({ loadBalancingConfig: [{ round_robin: {} }] }),
};
const waitReady = (channel) => __awaiter(void 0, void 0, void 0, function* () {
    const deadline = new Date();
    deadline.setSeconds(deadline.getSeconds() + 15);
    yield (0, nice_grpc_1.waitForChannelReady)(channel, deadline);
});
exports.waitReady = waitReady;
//# sourceMappingURL=helper.js.map