import _ from 'lodash';
import { CamelCasedProperties } from 'type-fest';

export function convertSnakeToCamel<T extends Record<any, any>>(obj: T): CamelCasedProperties<T> {
  return Object.keys(obj).reduce(
    (result, key) => ({
      ...result,
      [_.camelCase(key)]: obj[key],
    }),
    {},
  ) as CamelCasedProperties<T>;
}
