type Obj = Record<string, any>

export function merge<S extends Obj, T extends Obj>(source: S, target: T): S & T {
  return { ...structuredClone(source), ...structuredClone(target) };
}
