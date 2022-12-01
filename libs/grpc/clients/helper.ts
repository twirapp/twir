export const createClientAddr = (env: string, service: string, port: number): string => {
  let ip = service;
  if (env != 'production') {
    ip = '127.0.0.1';
  }

  return `${ip}:${port}`;
};

export const CLIENT_OPTIONS = {
  'grpc.lb_policy_name': 'round_robin',
  'grpc.service_config': JSON.stringify({ loadBalancingConfig: [{ round_robin: {} }] }),
};
