import {
  getResourcePaths,
  hasGetMethod,
  isChild,
  isCustomMethod,
  isSingletonResource,
} from './utils/resourceEvaluation.js';
import { hasException } from './utils/exceptions.js';
import { getAllSuccessfulGetResponseSchemas } from './utils/methodUtils.js';

const RULE_NAME = 'xgen-IPA-113-singleton-must-not-have-id';
const ERROR_MESSAGE = 'Singleton resources must not have a user-provided or system-generated ID.';

export default (input, opts, { path, documentInventory }) => {
  const resourcePath = path[1];

  if (isCustomMethod(resourcePath) || isChild(resourcePath)) {
    return;
  }

  if (hasException(input, RULE_NAME)) {
    return;
  }

  const oas = documentInventory.resolved;
  const resourcePaths = getResourcePaths(resourcePath, Object.keys(oas.paths));

  if (isSingletonResource(resourcePaths) && hasGetMethod(input)) {
    const resourceSchemas = getAllSuccessfulGetResponseSchemas(input);
    if (resourceSchemas.some((schema) => schemaHasIdProperty(schema))) {
      return [
        {
          message: ERROR_MESSAGE,
        },
      ];
    }
  }
};

function schemaHasIdProperty(schema) {
  if (Object.keys(schema).includes('properties')) {
    const propertyNames = Object.keys(schema['properties']);
    return propertyNames.includes('id') || propertyNames.includes('_id');
  }
  return false;
}