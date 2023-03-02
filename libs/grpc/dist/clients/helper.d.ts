import { Channel } from '@grpc/grpc-js';
export declare const createClientAddr: (env: string, service: string, port: number) => string;
export declare const CLIENT_OPTIONS: {
    'grpc.lb_policy_name': string;
    'grpc.service_config': string;
};
export declare const waitReady: (channel: Channel) => Promise<void>;
//# sourceMappingURL=helper.d.ts.map