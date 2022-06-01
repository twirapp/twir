export class NotATesterError extends Error {
  constructor(message: string) {
    super(message);
  }
}