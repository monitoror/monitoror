export enum ConfigErrorId {
  ConfigErrorConfigNotFound = 'ERROR_CONFIG_NOT_FOUND',
  ConfigErrorConfigVersionTooOld = 'ERROR_CONFIG_VERSION_TOO_OLD',
  ConfigErrorInvalidFieldValue = 'ERROR_INVALID_FIELD_VALUE',
  ConfigErrorMissingRequiredField = 'ERROR_MISSING_REQUIRED_FIELD',
  ConfigErrorUnauthorizedField = 'ERROR_UNAUTHORIZED_FIELD',
  ConfigErrorUnauthorizedSubtileType = 'ERROR_UNAUTHORIZED_SUBTILE_TYPE',
  ConfigErrorUnableToHydrate = 'ERROR_UNABLE_TO_HYDRATE',
  ConfigErrorUnknownTileType = 'ERROR_UNKNOWN_TILE_TYPE',
  ConfigErrorUnknownVariant = 'ERROR_UNKNOWN_VARIANT',
  ConfigErrorUnsupportedVersion = 'ERROR_UNSUPPORTED_VERSION',
}

export default ConfigErrorId
