export type ResolvableValue<T> = T | (() => T | Promise<T>);
