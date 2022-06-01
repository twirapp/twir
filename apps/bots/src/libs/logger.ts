
import pc from 'picocolors';

export class ConsoleLogger {
  info(key: string, ...args: string[]) {
    console.log(
      `${pc.bgCyan(pc.black(`${key}:`))} ${args.join(' ')}`,
    );
  }
}