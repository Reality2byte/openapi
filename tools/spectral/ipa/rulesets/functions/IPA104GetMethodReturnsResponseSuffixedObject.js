import {
  isSingleResourceIdentifier,
  isSingletonResource,
  getResourcePathItems,
  isResourceCollectionIdentifier,
} from './utils/resourceEvaluation.js';
import { resolveObject } from './utils/componentUtils.js';
import { hasException } from './utils/exceptions.js';
import { collectAdoption, collectAndReturnViolation, collectException } from './utils/collectionUtils.js';
import { checkSchemaRefSuffixAndReturnErrors } from './utils/validations/checkSchemaRefSuffixAndReturnErrors.js';

const RULE_NAME = 'xgen-IPA-104-get-method-returns-response-suffixed-object';

export default (input, _, { path, documentInventory }) => {
  const resourcePath = path[1];
  const responseCode = path[4];
  const oas = documentInventory.unresolved;
  const resourcePaths = getResourcePathItems(resourcePath, oas.paths);
  const contentPerMediaType = resolveObject(oas, path);

  if (
    !responseCode.startsWith('2') ||
    !contentPerMediaType ||
    !contentPerMediaType.schema ||
    !input.endsWith('json') ||
    (!isSingleResourceIdentifier(resourcePath) &&
      !(isResourceCollectionIdentifier(resourcePath) && isSingletonResource(resourcePaths)))
  ) {
    return;
  }

  if (hasException(contentPerMediaType, RULE_NAME)) {
    collectException(contentPerMediaType, RULE_NAME, path);
    return;
  }

  const errors = checkSchemaRefSuffixAndReturnErrors(path, contentPerMediaType, 'Response', RULE_NAME);

  if (errors.length !== 0) {
    return collectAndReturnViolation(path, RULE_NAME, errors);
  }
  collectAdoption(path, RULE_NAME);
};
